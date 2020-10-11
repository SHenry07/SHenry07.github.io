# docker

`IMAGE ID`应该是根据时间来算出来的



Yes, containerd sets localhost to default to HTTP https://github.com/containerd/containerd/blob/master/remotes/docker/resolver.go#L185 .

## cli

```shell
--link 共享网络空间
--link=""  : Add link to another container (<name or id>:alias or <name or id>)
其中，name和id是源容器的name和id，alias是源容器在link下的别名

-v some-docker-certs-client:/certs/client:ro  自动申请临时卷

--cap-add=NET_ADMIN # 限定上下文
--init		API 1.25+ Run an init inside the container that forwards signals and reaps processes
--dns 8.8.8.8
--rm
–-platform arm

manifest inspect --verbose 镜像名
```

**` docker run --init`专治crashloop起不来的容器**

```bash
# link怎么用
$ docker run --privileged --name some-docker -d -e DOCKER_TLS_CERTDIR=/certs     -v some-docker-certs-ca:/certs/ca     -v some-docker-certs-client:/certs/client     docker:dind

$ docker run --rm --link some-docker:docker  -e DOCKER_TLS_CERTDIR=/certs     -v some-docker-certs-client:/certs/client:ro     docker:latest version

# link 可能要被官方弃用，整体逻辑和k8s一样 就是环境变量
$ docker run --privileged --name some-docker -d \
    --network some-network --network-alias docker \
    -e DOCKER_TLS_CERTDIR=/certs \
    -v some-docker-certs-ca:/certs/ca \
    -v some-docker-certs-client:/certs/client \
    docker:dind
    
$ docker run --rm --network some-network \
    -e DOCKER_TLS_CERTDIR=/certs \
    -v some-docker-certs-client:/certs/client:ro \
    docker:latest version
```



# daemon configure

```bash
cat > /etc/docker/daemon.json <<EOF
{
  "registry-mirrors" : ["https://7bezldxe.mirror.aliyuncs.com"],
  "exec-opts": ["native.cgroupdriver=systemd"],
  "oom-score-adjust": -1000,
  "max-concurrent-downloads": 20,
  "live-restore": true,
  "max-concurrent-uploads": 10,
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "default-ipc-mode": "shareable",
  "storage-driver": "overlay2",
  "storage-opts": [
    "overlay2.override_kernel_check=true"
  ]
}
EOF
# 可选
  -H fd:/ ""
  "hosts": ["unix:///var/run/docker.sock", "tcp://0.0.0.0:2375"]
  "insecure-registries": ["10.13.30.108","172.22.8.103:5000"],
  "insecure-registries": ["127.0.0.1:5000"]
```



# 类似产品

[其他替换工具podman skopeo](https://zhuanlan.zhihu.com/p/77373246)

[skopeo](https://github.com/containers/skopeo)