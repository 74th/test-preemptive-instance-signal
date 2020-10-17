package main

import (
	"bytes"
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type kubernetesClientConfig struct {
	UseKubeConfig bool
	Namespace     string
	Deployment    string
}

type kubernetesClient struct {
	client    *kubernetes.Clientset
	config    kubernetesClientConfig
	available bool
	condLock  sync.Mutex
	cond      *sync.Cond
}

func newKubernetesClient(ctx context.Context, config kubernetesClientConfig) *kubernetesClient {
	c := new(kubernetesClient)
	c.config = config
	c.cond = sync.NewCond(&c.condLock)

	var conf *rest.Config
	var err error
	if config.UseKubeConfig {
		conf, err = clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
		if err != nil {
			log.Panicf("cannot connect kubernetes : %s", err.Error())
		}
	} else {
		conf, err = rest.InClusterConfig()
		if err != nil {
			log.Panicf("cannot connect kubernetes : %s", err.Error())
		}
		b, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		if err != nil {
			log.Panicf("cannot get current namespace : %s", err.Error())
		}
		c.config.Namespace = string(bytes.TrimSpace(b))
	}
	c.client, err = kubernetes.NewForConfig(conf)
	if err != nil {
		log.Fatalf("cannot connect kubernetes : %s", err.Error())
	}

	return c
}

func (c *kubernetesClient) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
loop:
	for {
		select {
		case <-ticker.C:
			status, err := c.client.AppsV1().Deployments(c.config.Namespace).Get(ctx, c.config.Deployment, v1.GetOptions{})
			if err != nil {
				log.Fatalf("cannot get status : %s", err.Error())
			}
			log.Printf("%s %d/%d", c.config.Deployment, status.Status.ReadyReplicas, status.Status.Replicas)
		case <-ctx.Done():
			break loop
		}
	}
	ticker.Stop()
}

func main() {
	conf := kubernetesClientConfig{}
	flag.StringVar(&conf.Namespace, "n", "default", "namespace")
	flag.StringVar(&conf.Deployment, "d", "tester", "deployment name")
	flag.BoolVar(&conf.UseKubeConfig, "k", false, "using kube config")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGQUIT)
		<-signalChan
		cancel()
	}()

	cl := newKubernetesClient(ctx, conf)
	cl.Run(ctx)

}
