## 压测/检测工具
#### stress 

是一个 Linux 系统压力测试工具，这里我们用作异常进程模拟平均负载升高的场景。

>  iowait无法升高的问题，是因为案例中stress使用的是 sync() 系统调用，它的作用是刷新缓冲区内存到磁盘中。对于新安装的虚拟机，缓冲区可能比较小，无法产生大的IO压力，这样大部分就都是系统调用的消耗了。所以，你会看到只有系统CPU使用率升高。解决方法是使用stress的下一代stress-ng，它支持更丰富的选项，比如 stress-ng -i 1 --hdd 1 --timeout 600（--hdd表示读写临时文件）。

```
 stress --cpu 1 --timeout 600 // cpu
 stress -i 1 --timeout 600 // io
 stress -c 8  --timeout 600 // 大量进程
```
#### hping3 

是一个可以构造 TCP/IP 协议数据包的工具，可以对系统进行安全审计、防火墙测试等。
```
 -S参数表示设置TCP协议的SYN（同步序列号），
 -p表示目的端口为80
 -i u100表示每隔100微秒发送一个网络帧
 -c <num>
```
#### traceroute

```
-tcp 
-p 80 表示端口号
-n 表示不对结果中的IP地址执行反向域名解析
```



#### curl

```
-I,--head Fetch  the  HTTP-header  only!
-H, --header <header>
 	-H Host:httpbin.example.com 将 HTTP 头部参数 Host 设置为 “httpbin.example.com”
-x, --proxy <[protocol://][user:password@]proxyhost[:port]>
	Use the specified HTTP proxy. If the port number is not specified, it is assumed at port 1080.
-X,--request 指定请求方法 Specifies a custom request method to use when communicating with the HTTP server. 
-w表示只输出HTTP状态码及总时间，
-o Write to file instead of stdout
	 -w 'Http code: %{http_code}\nTotal time:%{time_total}s\n'
-O Write output to a file named as the remote file
```

#### fio

（Flexible I/O Tester）正是最常用的文件系统和磁盘 I/O 性能基准测试工具

```
# 随机读
fio -name=randread -direct=1 -iodepth=64 -rw=randread -ioengine=libaio -bs=4k -size=1G -numjobs=1 -runtime=1000 -group_reporting -filename=/dev/sdb

# 随机写
fio -name=randwrite -direct=1 -iodepth=64 -rw=randwrite -ioengine=libaio -bs=4k -size=1G -numjobs=1 -runtime=1000 -group_reporting -filename=/dev/sdb

# 顺序读
fio -name=read -direct=1 -iodepth=64 -rw=read -ioengine=libaio -bs=4k -size=1G -numjobs=1 -runtime=1000 -group_reporting -filename=/dev/sdb

# 顺序写
fio -name=write -direct=1 -iodepth=64 -rw=write -ioengine=libaio -bs=4k -size=1G -numjobs=1 -runtime=1000 -group_reporting -filename=/dev/sdb 
```

- direct，表示是否跳过系统缓存。上面示例中，我设置的 1 ，就表示跳过系统缓存。
- iodepth，表示使用异步 I/O（asynchronous I/O，简称 AIO）时，同时发出的 I/O 请求上限。在上面的示例中，我设置的是 64。
- rw，表示 I/O 模式。我的示例中， read/write 分别表示顺序读 / 写，而 randread/randwrite 则分别表示随机读 / 写。
- ioengine，表示 I/O 引擎，它支持同步（sync）、异步（libaio）、内存映射（mmap）、网络（net）等各种 I/O 引擎。上面示例中，我设置的 libaio 表示使用异步 I/O。
- bs，表示 I/O 的大小。示例中，我设置成了 4K（这也是默认值）。
- filename，表示文件路径，当然，它可以是磁盘路径（测试磁盘性能），也可以是文件路径（测试文件系统性能）。示例中，我把它设置成了磁盘 /dev/sdb。不过注意，用磁盘路径测试写，会破坏这个磁盘中的文件系统，所以在使用前，你一定要事先做好数据备份。

[磁盘优化](https://time.geekbang.org/column/article/79368)

针对磁盘和应用程序 I/O 模式的特征，我们可以选择最适合的 I/O 调度算法。比方说，SSD 和虚拟机中的磁盘，通常用的是 noop 调度算法。而数据库应用，我更推荐使用 deadline 算法。

```
如何修改磁盘的io调度算法哇？
/sys/block/{DEVICE-NAME}/queue/scheduler
```

#### smartctl

检查磁盘故障