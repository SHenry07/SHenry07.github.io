# context

有点类似namespace的概念，但是更偏向用户级别

## 如何创建context

- 本地
  默认已经创建

  `docker context create test01 --docker host=unix:///var/run/docker.sock`

- 远程ssh

  `docker context create remote ‐‐docker “host=ssh://user@remotemachine”`

  ssh 自测在buildx中不可用

- 远程tcp 端口

  `docker context create k8s-test      --docker host=tcp://10.14.41.57:2375`

| **Target Environment** | **Context name** | **API endpoint**            |
| ---------------------- | ---------------- | --------------------------- |
| localhost              | default          | unix:///var/run/docker.sock |
| Remote host            | remote           | ssh://user@remotemachine    |
| docker-in-docker       | dind             | tcp://127.0.0.1:2375        |

- 本地docker in docker

  `docker run ‐‐rm -d -p “2375:2375” ‐‐privileged -e “DOCKER_TLS_CERTDIR=” ‐‐name some-docker docker:19.03.3-dind`

  `docker context create dind ‐‐docker “host=tcp://127.0.0.1:2375” ‐‐default-stack-orchestrator swarm`

- k8s

`docker context create k8s-test --default-stack-orchestrator=kubernetes   --kubernetes config-file=/root/docker/default.kubeconfig  --docker host=unix:///var/run/docker.sock`

## Using DOCKER_HOST environment variable to set up the target engine

支持 [DOCKER_HOST environment variable](https://docs.docker.com/engine/reference/commandline/cli/#environment-variables) and [*-H, –host* command line option](https://docs.docker.com/compose/reference/overview/).

```
DOCKER_HOST=“ssh://user@remotehost” docker-compose up -d
```

## Using docker contexts

```
$ docker context ls
NAME   DESCRIPTION   DOCKER ENDPOINT   KUBERNETES ENDPOINT   ORCHESTRATOR
…
test              ssh://user@remotemachine
$ cd hello-docker
$ docker-compose ‐‐context test up -d
```

## Exporting a Kubernetes context

  You can export a Kubernetes context only if the context you are exporting has a Kubernetes endpoint configured. You cannot import a Kubernetes context using docker context import.

These steps will use the --kubeconfig flag to export only the Kubernetes elements of the existing k8s-test context to a file called “k8s-test.kubeconfig”. The cat command will then show that it’s exported as a valid kubeconfig file.

```bash
$ docker context export k8s-test --kubeconfig
Written file "k8s-test.kubeconfig"
```
