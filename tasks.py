from invoke import task

@task
def build(c):
    with c.cd("tester"):
        c.run("CGO_ENABLED=0 go build -o tester main.go")
    c.run("docker-compose build")
    c.run("docker-compose push")

@task
def deploy(c):
    c.run("kubectl apply -f manifest")

