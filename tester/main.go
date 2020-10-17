package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func stopForSyscall(duration int64, exitChan chan int) {
	if duration < 0 {
		return
	}
	time.Sleep(time.Second * time.Duration(duration))
	exitChan <- 0
}

func main() {
	var stopDuration int64
	var tickerDuration int64
	var addr string
	var exitCode int
	flag.Int64Var(&stopDuration, "s", -1, "stop for syscall")
	flag.Int64Var(&tickerDuration, "d", 5, "ticker")
	flag.StringVar(&addr, "a", "0.0.0.0:10000", "addr")
	flag.IntVar(&exitCode, "e", 0, "exit code")
	flag.Parse()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panicf("cannot listen %s :%s", addr, err.Error())
	}

	fmt.Println("START")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGTERM,
		syscall.SIGKILL)

	ticker := time.NewTicker(time.Second * time.Duration(tickerDuration))

	exitChan := make(chan int)

	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGTERM:
					fmt.Println("SIGTERM")
					go stopForSyscall(stopDuration, exitChan)
				case syscall.SIGKILL:
					fmt.Println("SIGKILL")
					go stopForSyscall(stopDuration, exitChan)
				}
			case <-ticker.C:
				fmt.Println(time.Now().Format("2006/01/02 15:04:05"))
			}
		}
	}()
	<-exitChan
	l.Close()
	os.Exit(exitCode)
}
