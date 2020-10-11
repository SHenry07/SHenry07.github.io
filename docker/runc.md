# runc

https://github.com/opencontainers/runc

接下来，尝试使用runc启动一个容器。

第一步，执行如下命令将一个docker image文件解压，然后打包tar。我是在RHEL8上运行的，所以用了podman，用docker也是一样的。

```text
$docker run --rm -d --name ubuntu ubuntu:18.04 tail -f /dev/null
$ docker export ubuntu > rootfs.tar
$ docker kill ubuntu
```

第二步，将tar文件解压后制作成容器所需的Filesystem bundle，然后使用`runc spec`命令获得设置文件(config.json)。

```text
$ mkdir -p bundle/rootfs
$ tar xf rootfs.tar -C bundle/rootfs
$ runc spec -b bundle
```

第三步，运行`runc run`，Filesystem bundle作为参数，启动容器

```text
# sudo runc run --bundle bundle demo1
```

启动后的样子就像下面截图这样进入了容器中的一个ubuntu的shell环境，直接的感觉就和用cgroups和unshare命令启动的shell很类似，但是并不需要输入比较复杂的cgoups和unshare命令，而且容器的功能也更加完整，比如像截图中所示可以通过`ps -ef`查看process。







执行时（ runc） ）使用所谓的执行时根目录来储存和获取有关容器的资讯.在这个根目录下， runc   放置子目录（每个容器一个），每个子目录包含 state.json   档案，容器状态描述所在的位置。

执行时根目录的预设位置是` /run/runc`   （对于非无根容器）或 `$XDG_RUNTIME_DIR/runc`   （对于无根rootless容器） - 后者通常也指向 /run下的某个地方   （例如` /run/user/$UID/runc `）。

当容器引擎呼叫 runc时 ，它可能会覆盖预设的执行时根目录并指定自定义目录（ --root   runc的选择 ）. Docker使用这种可能性，例如 在我的docker上，它指定了 `/run/docker/runtime-runc/moby `作为执行时根。

那说，要做 runc list   看到你的Docker容器，你必须通过指定 --root将它指向Docker的执行时根目录。 此外，鉴于Docker容器预设情况下不是无根的，您将需要相应的权限来访问执行时根（例如，使用 sudo） ）。
所以，这应该是如何运作的：

```shell
$ docker run -d alpine sleep 1000
4acd4af5ba8da324b7a902618aeb3fd0b8fce39db5285546e1f80169f157fc69
$ sudo runc --root /run/docker/runtime-runc/moby/ list
ID                                                                 PID         STATUS      BUNDLE                                                                                                                               CREATED                          OWNER
4acd4af5ba8da324b7a902618aeb3fd0b8fce39db5285546e1f80169f157fc69   18372       running     /run/docker/containerd/daemon/io.containerd.runtime.v1.linux/moby/4acd4af5ba8da324b7a902618aeb3fd0b8fce39db5285546e1f80169f157fc69   2019-07-12T17:33:23.401746168Z   root

```

对于images，你不能运行 `runc `看到它们，因为它根本没有图像的概念 - 相反，它在捆绑上执行.建立包（例如基于图像）是呼叫者的责任（在您的情况下 - 容器）。



# cri-o

https://www.jianshu.com/p/5c7ffe9328e9

[知乎实践](https://zhuanlan.zhihu.com/p/133861092)