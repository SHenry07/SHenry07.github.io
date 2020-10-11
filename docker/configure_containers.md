## Use a restart policy

To configure the restart policy for a container, use the `--restart` flag when using the `docker run` command. The value of the `--restart` flag can be any of the following:

| Policy                     | Result                                                       |
| :------------------------- | :----------------------------------------------------------- |
| `no`                       | Do not automatically restart the container when it exits. This is the default. |
| `on-failure[:max-retries]` | Restart only if the container exits with a non-zero exit status. Optionally, limit the number of restart retries the Docker daemon attempts. Similar to `always`, except that when the container is stopped (manually or otherwise), it is not restarted even after Docker daemon restarts|
| `unless-stopped`           | Restart the container unless it is explicitly stopped or Docker itself is stopped or restarted.(See the second bullet listed in [restart policy details](https://docs.docker.com/config/containers/start-containers-automatically/#restart-policy-details)) |
| `always`                   | Always restart the container regardless of the exit status. When you specify always, the Docker daemon will try to restart the container indefinitely. The container will also always start on daemon startup, regardless of the current state of the container. |

```
/usr/bin/docker run --restart=on-failure:5 --env-file=/etc/etcd.env --net=host -v /etc/ssl/certs:/etc/ssl/certs:ro -v /etc/ssl/etcd/ssl:/etc/ssl/etcd/ssl:ro -v /var/lib/etcd:/var/lib/etcd:rw --memory=0 --blkio-weight=1000 --name=etcd1 quay.io/coreos/etcd:v3.2.26-arm64 /usr/local/bin/etcd
```

# 第一个启动的程序

A container’s main running process is the `ENTRYPOINT` and/or `CMD` at the end of the `Dockerfile`. It is generally recommended that you separate areas of concern by using one service per container. That service may fork into multiple processes (for example, Apache web server starts multiple worker processes). It’s ok to have multiple processes, but to get the most benefit out of Docker, avoid one container being responsible for multiple aspects of your overall application. You can connect multiple containers using user-defined networks and shared volumes.

The container’s main process is responsible for managing all processes that it starts. In some cases, the main process isn’t well-designed, and doesn’t handle “reaping” (stopping) child processes gracefully when the container exits. If your process falls into this category, you can use the `--init` option when you run the container. The `--init` flag inserts a tiny init-process into the container as the main process, and handles reaping of all processes when the container exits. Handling such processes this way is superior to using a full-fledged init process such as `sysvinit`, `upstart`, or `systemd` to handle process lifecycle within your container.