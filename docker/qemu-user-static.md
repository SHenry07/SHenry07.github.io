# docker/binfmt

Register Arm executables
`docker run --rm --privileged docker/binfmt:820fdd95a9972a5308930a2bdfb8573dd4447ad3`

# multiarch/qemu-user-static

```bash
$ uname -m
x86_64

$ docker run --rm -t arm64v8/ubuntu uname -m
standard_init_linux.go:211: exec user process caused "exec format error"
或
standard_init_linux.go:211: exec user process caused "no such file or directory"
```

Here is the alternative workflow for your environment possibly.

`$ docker run --rm --privileged multiarch/qemu-user-static:register --reset`
or

`$ docker run --rm --privileged multiarch/qemu-user-static --reset -p yes`

`-p yes` (`--persistent yes`)  内核低于4.8的不支持，错误如下`sh: write error: Invalid argument`

Above commands are just to add `binfmt_misc` files to your host environment.

Then for example if you want to run arm64v8/ubuntu Ubuntu ARM 64-bit (= aarch64) container, it is like this with your container "my-aarch64-ubuntu".

```bash
$ docker build --rm -t my-aarch64-ubuntu -<<EOF
FROM multiarch/qemu-user-static:x86_64-aarch64 as qemu

FROM arm64v8/ubuntu
COPY --from=qemu /usr/bin/qemu-aarch64-static /usr/bin
EOF

$ docker run --rm -t my-aarch64-ubuntu uname -m
aarch64
```

## 更多常用选项

```
       --credential:  if yes, credential and security tokens are
                      calculated according to the binary to interpret
       --persistent:  if yes, the interpreter is loaded when binfmt is
                      configured and remains in memory. All future uses
                      are cloned from the open file.
```

# Error

- centos7/ubuntu 为什么没有正常工作

现象: `sh: write error: Invalid argument`

- 内核版本需要大于4.8 

- 自行安装`qemu-user-static`   

  > 此方法未实验 可以参考[1](https://github.com/docker/buildx/issues/138#issuecomment-664977087) 和[2](https://github.com/zalando/skipper/pull/1416/files)

- bin 文件手动挂载进容器

- 依赖docker的多级构建功能, 详见上面的dockerfile

That is to copy the binary qemu-*-static from the container to the host and use docker run -v host_dir:container_dir.

```bash
$ docker run --rm --privileged multiarch/qemu-user-static:register --reset

$ cat /proc/sys/fs/binfmt_misc/qemu-aarch64
enabled
interpreter /usr/bin/qemu-aarch64-static
flags: 
offset 0
magic 7f454c460201010000000000000000000200b700
mask ffffffffffffff00fffffffffffffffffeffffff

$ docker create -it --name dummy multiarch/qemu-user-static:x86_64-aarch64 bash

$ docker container ls -a
CONTAINER ID        IMAGE                                       COMMAND             CREATED             STATUS              PORTS               NAMES
6ab622a76dfa        multiarch/qemu-user-static:x86_64-aarch64   "bash"              3 minutes ago       Created                                 dummy

$ docker cp dummy:/usr/bin/qemu-aarch64-static qemu-aarch64-static

$ ls qemu-aarch64-static
qemu-aarch64-static*

$ docker rm -f dummy

$ docker container ls -a
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES

$ docker run --rm -t -v $(pwd)/qemu-aarch64-static:/usr/bin/qemu-aarch64-static arm64v8/ubuntu uname -m
aarch64
```

Because flags: `F` used for `-p yes`is relatively newer feature in binfmt_misc itself. I saw the option worked on Travis CI Ubuntu xenial, but did not work on Travis CI Ubuntu trusty.

Here is the document on kernel 4.10.
https://www.kernel.org/doc/html/v4.10/admin-guide/binfmt-misc.html
flags - F - fix binary is the used function for the -p yes.

Here is the binfmt_misc document for kernel 3.10. There is no F option of flags there.

https://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git/tree/Documentation/binfmt_misc.txt?h=v3.10#n37

# Reference

https://hub.docker.com/r/docker/binfmt/tags 

[multiarch/qemu-user-static](https://github.com/multiarch/qemu-user-static)

[qemu-user-static - part 2 register](http://junaruga.hatenablog.com/entry/2019/04/08/021753)

[multiarch/qemu-user-static#38](https://github.com/multiarch/qemu-user-static/issues/38) 

[multiarch/qemu-user-static#100](https://github.com/multiarch/qemu-user-static/issues/100)

