# 编译多平台docker镜像-buildx

buildx 是docker CLI plugin 是 `docker build` 的扩展。 `docker buildx install` 和`docker buildx uninstall`.  the full support of the features provided by [Moby BuildKit](https://github.com/moby/buildkit) builder toolkit

Buildx will always build using the BuildKit engine and does not require `DOCKER_BUILDKIT=1` environment variable for starting builds.

we support a "docker" driver that uses the BuildKit library bundled into the docker daemon binary, and a "docker-container" driver that automatically launches BuildKit inside a Docker container.

## 使用要求

### 多架构支持

#### qemu 虚拟化

For QEMU binaries registered with binfmt_misc on the host OS to work transparently inside containers they must be registed with the fix_binary flag. This requires a **kernel >= 4.8 and binfmt-support >= 2.1.7**. You can check for proper registration by checking if `F` is among the flags in `/proc/sys/fs/binfmt_misc/qemu-*`.

While Docker Desktop comes preconfigured with binfmt_misc support for additional platforms, for other installations it likely needs to be installed manually

Docker Desktop 内置了 binfmt_misc, 其他情况需要自己安装. [详见](./qemu-user-static.md)

#### 一个builder绑定多个实例(context)

To use a remote node you can specify the `DOCKER_HOST` or remote context name while creating the new builder(如果要使用远程node,可以设置DOCKER_HOST,或者远程context name，支持追加). After creating a new instance you can manage its lifecycle with the `inspect`, `stop` and `rm` commands and list all available builders with `ls`. After creating a new builder you can also append new nodes to it
生命周期: inspect stop rm ls
切换`docker buildx use <name>`

(使用多个不同架构的节点比交叉编译有更好的性能)Using multiple native nodes provides better support for more complicated cases not handled by QEMU and generally have better performance. Additional nodes can be added to the builder instance with `--append` flag.

```
# assuming contexts node-amd64 and node-arm64 exist in "docker context ls"
$ docker buildx create --use --name mybuild <amd64架构的节点>
mybuild
$ docker buildx create --append --name mybuild <追加arm64架构的节点>
$ docker buildx build --platform linux/amd64,linux/arm64 .
```

##### context

Docker 19.03 also features a new `docker context` command that can be used for giving names for remote Docker API endpoints. Buildx integrates with `docker context` so that all of your contexts automatically get a default builder instance. While creating a new builder instance or when adding a node to it you can also set the context name as the target.

通过`docker context`创建的context会自动出现在`docker buildx ls`, 但是不能修改
> context 支持unix:///sock,ssh,tcp等，[详见](./context.md)

### 语言支持交叉编译

Finally, depending on your project, the language that you use may have good support for cross-compilation. In that case, multi-stage builds in Dockerfiles can be effectively used to build binaries for the platform specified with `--platform` using the native architecture of the build node. List of build arguments like `BUILDPLATFORM` and `TARGETPLATFORM` are available automatically inside your Dockerfile and can be leveraged by the processes running as part of your build.

需要一个支持交叉编译的语言, build参数像 `BUILDPLATFORM` and `TARGETPLATFORM` 可以自动替换Dockerfile中的架构部分

### 多级构建

```
FROM --platform=$BUILDPLATFORM golang:alpine AS build
ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN echo "I am running on $BUILDPLATFORM, building for $TARGETPLATFORM" > /log
FROM alpine
COPY --from=build /log /log
```

`BUILDPLATFORM` 构架机器的架构

`TARGETPLATFORM` 想要构建成的架构

# 总结/实战

## 安装docker

```bash
# 安装19.0.3的docker
curl -s https://get.docker.com/ | bash - 
```

### with Docker 18.09+

```bash
## Download buildx
curl -s https://api.github.com/repos/docker/buildx/releases/latest \
| grep "browser_download_url.*buildx-*.*linux-amd64" \
| cut -d : -f 2,3 \
| tr -d \" \
| wget -qi -

## Confiure buildx
mkdir -p ~/.docker/cli-plugins
mv buildx* ~/.docker/cli-plugins/docker-buildx
chmod a+x ~/.docker/cli-plugins/docker-buildx
```

## 安装qemu-user-static

为了让在x86上可以运行arm64的docker镜像，这里需要安装[qemu-user-static](https://github.com/multiarch/qemu-user-static)，过程如下：

`docker run --rm --privileged multiarch/qemu-user-static --reset -p yes`

[如出现错误请参考](./qemu-user-static.md)

## 开启实验模式

```bash
export DOCKER_CLI_EXPERIMENTAL=enabled

mkdir -p $HOME/.docker
cat > $HOME/.docker/config.json  << EOF
{
  "experimental": "enabled",
  "debug": true
}
EOF
```

## 开始构建

```bash
# 创建构建器
docker buildx create --name actions_builder --use

docker buildx create --name actions_builder --use --config ./config.toml --driver-opt network=host node-amd64

# 启动构建器  
$ docker buildx inspect actions_builder --bootstrap

# 观察下当前使用的构建器及构建器支持的cpu架构，可以看到支持很多cpu架构
$ docker buildx ls

# 构建镜像
docker buildx build --platform=linux/amd64,linux/arm64,linux/arm/v7 -t <tag>  .
# 还有--push 的选项可直接到仓库
docker buildx build --push -f $DOCKERFILE_PATH -t $IMAGE_NAME --platform linux/amd64,linux/arm64 .


# 查看镜像信息
$ docker buildx imagetools inspect myusername/hello1 
# 暂未发现支持本地http仓库的方式
```

```
# syntax=docker/dockerfile:experimental 此选项在>=19.03后就不需要了

FROM --platform=$TARGETPLATFORM alpine

RUN uname -a > /os.txt

CMD cat /os.txt
```

## 编写脚本生成多平台docker镜像

假设有一个普通的golang程序源码，我们已经写好了Dockerfile生成其docker镜像，如下：

```dockerfile
# Start from the latest golang base image
FROM --platform=$BUILDPLATFORM golang:1.14.2-alpine as go-deps
# Add Maintainer Info
LABEL maintainer="Jeremy Xu <jeremyxu2010@gmail.com>"
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY ./cmd ./cmd
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -mod=readonly -o out/demo -ldflags "-s -w" ./cmd
# Start from the latest alpine base image
FROM alpine:latest
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy execute file from go-builder
COPY --from=go-builder /app/output/demo /app/demo
# Set docker image command
CMD [ "/app/demo" ]
```

```bash
# 生成linux/amd64 docker镜像
$ docker buildx build --rm -t go-mul-arch-build:latest-amd64 --platform=linux/amd64 --output=type=docker .
# 生成linux/arm64 docker镜像
$ docker buildx build --rm -t go-mul-arch-build:latest-arm64 --platform=linux/arm64 --output=type=docker .
```

## 如何push到私有仓库
### 如何推到本地http仓库

Use --config https://github.com/docker/buildx#--config-file with http: true set for the insecure registry.

1. 创建本地仓库的配置,用以覆盖buildKit的配置^Reference4^

  ```txt
  # registry configures a new Docker register used for cache import or output.
  [registry."10.13.30.108"]
    http = true
    insecure = true

  [registry."172.22.8.103"]
    http = true
    insecure = true
# 另外一个成功案例
# registry configures a new Docker register used for cache import or output. 
[registry."docker.io"]
  mirrors = ["wlzfs4t4.mirror.aliyuncs.com"]
  # http = true
  # insecure = true
  # ca=["/etc/config/myca.pem"]
  # [[registry."docker.io".keypair]]
  #   key="/etc/config/key.pem"
  #   cert="/etc/config/cert.pem"
  ```

2. 注意加上 本地网关。

`--driver docker-container --driver-opt image=moby/buildkit:master,network=host`

`--driver docker-container 是默认的驱动器`。netwerk=host 能被dns解析的remote registry经测试是可以省略的, local registry是**必需**的(也可手动进容器echo 进去 不过这只是临时方案)，/etc/hosts是不行的(原因在下面)，但是也可以通过dns来解析

  > Buildkit runs inside a container with buildx container driver and therefore will never access /etc/hosts inside your host machine. You can use buildkitd config with --config to configure registry mirrors and http/https. I assume mirrors = ["localhost:5000"] should work for you although it is a weird usage. ^[Reference1]^

3. 创建build KIt 默认是docker-container

  `docker buildx create --driver-opt network=host --use --config config.toml --name container-builder`

4. docker login

5. RUN 构建指令

> 所有配置都在buildx目录下可以尝试手动修改
> ```bash
> vim ~/.docker/buildx/instances/
> # 备注 但是按理来说 修改后是应该重启buildKit的
> ```

### 如何推到本地https仓库,主要解决自签证书的认证问题

- 基于moby/buildKit构建已经有ca证书的镜像
  driver-opt has a image option (for docker-container driver).

  > docker-container

  > - image=IMAGE - Sets the container image to be used for running buildkit.
    - network=NETMODE - Sets the network mode for running the buildkit container.
  - Example:
    `--driver docker-container --driver-opt image=moby/buildkit:master,network=host`

> I went with the option to have a 2 line Dockerfile that adds my internal CAs to moby/buildkit and use that image when creating the builder.

- 将ca证书复制进容器
A marginally more robust work-around, but still not pretty (no error checking etc):

```bash
BUILDER=$(sudo docker ps | grep buildkitd | cut -f1 -d' ')
sudo docker cp YOUR-CA.crt $BUILDER:/usr/local/share/ca-certificates/
sudo docker exec $BUILDER update-ca-certificates
sudo docker restart $BUILDER

#另外一个开启可行的办法, 创建私钥
A possible solution/suggestion would be to allow for something like:
docker secret create a-ca-secret my-ca-file.crt
docker buildx create ---name builder --driver-opt private-ca-cert=/run/secrets/a-ca-secret
```

...and then have the e.g. the docker-container builder pull in the cert and call update-ca-certificates when it starts.


## 更多存储方式

`-o, --output=[PATH,-,type=TYPE[,KEY=VALUE]`
Sets the export action for the build result. In docker build all builds finish by creating a container image and exporting it to docker images. buildx makes this step configurable allowing results to be exported directly to the client, oci image tarballs, registry etc.

Buildx with docker driver currently only supports local, tarball exporter and image exporter. docker-container driver supports all the exporters.

If just the path is specified as a value, buildx will use the local exporter with this path as the destination. If the value is “-”, buildx will use tar exporter and write to stdout.

Examples:

```
docker buildx build -o . .
docker buildx build -o outdir .
docker buildx build -o - - > out.tar
docker buildx build -o type=docker .
docker buildx build -o type=docker,dest=- . > myimage.tar
docker buildx build -t tonistiigi/foo -o type=registry 
```

Supported exported types are:

- local
  The local export type writes all result files to a directory on the client. The new files will be owned by the current user. On multi-platform builds, all results will be put in subdirectories by their platform.

  Attribute key:

  - dest - destination directory where files will be written

- tar
  The tar export type writes all result files as a single tarball on the client. On multi-platform builds all results will be put in subdirectories by their platform.

   Attribute key:

   - dest - destination path where tarball will be written. “-” writes to stdout.
  
- oci
  The oci export type writes the result image or manifest list as an OCI image layout tarball https://github.com/opencontainers/image-spec/blob/master/image-layout.md on the client.

    Attribute key:

    - dest - destination path where tarball will be written. “-” writes to stdout.

- docker
  The docker export type writes the single-platform result image as a Docker image specification tarball https://github.com/moby/moby/blob/master/image/spec/v1.2.md on the client. Tarballs created by this exporter are also OCI compatible.

  Currently, multi-platform images cannot be exported with the docker export type. The most common usecase for multi-platform images is to directly push to a registry (see registry)

  Attribute keys:
  - dest - destination path where tarball will be written. If not specified the tar will be loaded automatically to the current docker instance.
  - context - name for the docker context where to import the result

- image
  The image exporter writes the build result as an image or a manifest list. When using docker driver the image will appear in docker images. Optionally image can be automatically pushed to a registry by specifying attributes.
  
   Attribute keys:
  
  - name - name (references) for the new image
  - push - boolean to automatically push the image.
- registry
  The registry exporter is a shortcut for` type=image,push=true`
- --push
  Shorthand for `--output=type=registry`Will automatically push the build result to registry.
- --load
  Shorthand for `--output=type=docker`. 

## cache

`--cache-from=[NAME|type=TYPE[,KEY=VALUE]]`
Use an external cache source for a build. Supported types are registry and local. The registry source can import cache from a cache manifest or (special) image configuration on the registry. The local source can import cache from local files previously exported with --cache-to.

If no type is specified, registry exporter is used with a specified reference.

docker driver currently only supports importing build cache from the registry.

Examples:
```
docker buildx build --cache-from=user/app:cache .
docker buildx build --cache-from=user/app .
docker buildx build --cache-from=type=registry,ref=user/app .
docker buildx build --cache-from=type=local,src=path/to/cache .
```

`--cache-to=[NAME|type=TYPE[,KEY=VALUE]]`

xport build cache to an external cache destination. Supported types are registry, local and inline. Registry exports build cache to a cache manifest in the registry, local exports cache to a local directory on the client and inline writes the cache metadata into the image configuration.

docker driver currently only supports exporting inline cache metadata to image configuration. Alternatively, `--build-arg BUILDKIT_INLINE_CACHE=1` can be used to trigger inline cache exporter.

Attribute key:

- mode - Specifies how many layers are exported with the cache. “min” on only exports layers already in the final build build stage, “max” exports layers for all stages. Metadata is always exported for the whole build.

Examples:

```
docker buildx build --cache-to=user/app:cache .
docker buildx build --cache-to=type=inline .
docker buildx build --cache-to=type=registry,ref=user/app .
docker buildx build --cache-to=type=local,dest=path/to/cache .
```

## 卸装builder的某一个


```
docker buildx create --name mybuilder --node mybuilder0 --leave
```

# 常见错误

- `code = Unknown desc = failed to do request: Head https: ....`

见如何push到私有仓库

- `failed to solve: rpc error: code = Unknown desc = server message: insufficient_scope: authorization failed`

未登录, `docker login`

- unexpected status: 401 Unauthorized

现象: `failed to solve: rpc error: code = Unknown desc = failed commit on ref "manifest-sha256:1431e9cf96aaf4a236f16b8a7ec80abfa3be85a17c48b2b7cb21a2649709acef": unexpected status: 401 Unauthorized`

镜像构建缓存出了问题

- unexpected status: 415 Unsupported Media Type  

`failed to solve: rpc error: code = Unknown desc = failed commit on ref "index-sha256:79720c5cb034a2fdfa980db08b119c185645c8d38e5cf406dc17ab95c70fc8c2": unexpected status: 415 Unsupported Media Type`

原因不明，猜测是仓库的问题， 因为单独推送是可以的

- Building images for multi-arch with --load parameter fails

现象：
`failed to solve: rpc error: code = Unknown desc = docker exporter does not currently support exporting manifest lists`

`docker load` does not currently support loading manifest lists and images in `docker images` can only be for a single platform. This is documented in https://github.com/docker/buildx#docker

docker 19.0.3 版本中的registry不支持, 
The docker export type writes the single-platform result image as a Docker image specification tarball https://github.com/moby/moby/blob/master/image/spec/v1.2.md on the client. Tarballs created by this exporter are also OCI compatible.

Currently, multi-platform images cannot be exported with the docker export type. The most common usecase for multi-platform images is to directly push to a registry (see registry).
`use --push on multi-platform , use --load for single platform`

可以自行下载[moby](https://github.com/moby/moby/pull/38738),或导出到tar，推送到仓库

# 如何排错

cli 中加上`--progress=plain`

` docker buildx build --progress=plain -t reg.esgyn.cn/test/jenkins/jenkins:ownv9 --platform linux/arm64,linux/amd64 --push .  `
## TODO: 如何debug 调试/验证失败-需尝试升级buildx

`docker buildx build -f Dockerfile.buildx --target debug --platform linux/amd64,linux/arm64 -t localhost:5000/bmitch-public/golang-hello:buildx1 --output type=registry .`

`--target debug` 开启debug

 `--target frontend`

`--target string	Set the target build stage to build.`

# 参考

1.https://github.com/docker/buildx/issues/218
2.[在单机上运行buildx需要哪些准备](https://github.com/docker/buildx/issues/138)
3.[Custom registry, push error on self-signed cert](https://github.com/docker/buildx/issues/80) 
[custom registry HTTP](https://github.com/docker/buildx/issues/94)
4.[buildx调用的buildKit镜像模板](https://github.com/moby/buildkit/blob/master/docs/buildkitd.toml.md)
5.[别人封装好的构建工具](https://github.com/crazy-max/ghaction-docker-buildx/actions)
6.https://docs.docker.com/buildx/working-with-buildx/#build-multi-platform-images
7.https://github.com/docker/buildx
8.https://github.com/docker/buildx/issues/138
9.https://github.com/multiarch/qemu-user-static
10.[driver-opt 到底支持kubernetes吗？](https://github.com/docker/buildx/issues/342)