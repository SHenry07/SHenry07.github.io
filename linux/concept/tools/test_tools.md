## 压测/检测工具
#### stress 

是一个 Linux 系统压力测试工具，这里我们用作异常进程模拟平均负载升高的场景。

>  iowait无法升高的问题，是因为案例中stress使用的是 sync() 系统调用，它的作用是刷新缓冲区内存到磁盘中。对于新安装的虚拟机，缓冲区可能比较小，无法产生大的IO压力，这样大部分就都是系统调用的消耗了。所以，你会看到只有系统CPU使用率升高。解决方法是使用stress的下一代stress-ng，它支持更丰富的选项，比如 stress-ng -i 1 --hdd 1 --timeout 600（--hdd表示读写临时文件）。

```
 stress --cpu 1 --timeout 600 // cpu
 stress -i 1 --timeout 600 // io
 stress -c 8  --timeout 600 // 大量进程
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
--connect-timeout <num> 超时时间 10
-s Silent or quiet mode. Don't show progress meter or error messages.  Makes Curl mute
--max-time 3 最多重试3次

```

### 硬盘/IO

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

### 网络IO

#### hping3 

是一个可以构造 TCP/IP 协议数据包的工具，可以对系统进行安全审计、防火墙测试等。
```
 -S参数表示设置TCP协议的SYN（同步序列号），
 -p表示目的端口为80
 -i u100表示每隔100微秒发送一个网络帧
   --flood no replies will be shown
 -c <num>
 -ltn 列出所有TCP端口
 -a 伪造IP模拟DDOS
 
 DUP! len=44 ip=10.13.30.109 ttl=58 DF id=0 sport=80 flags=SA seq=5 win=65535 rtt=3059.0 ms
# DUP 表示收到了重复包
```

#### pktgen

测试PPS的方法

```shell
$ modprobe pktgen
$ ps -ef | grep pktgen | grep -v grep
root     26384     2  0 06:17 ?        00:00:00 [kpktgend_0]
root     26385     2  0 06:17 ?        00:00:00 [kpktgend_1]
$ ls /proc/net/pktgen/
kpktgend_0  kpktgend_1  pgctrl

$ vim test.sh
# 定义一个工具函数，方便后面配置各种测试选项
pgset() {
    local result
    echo $1 > $PGDEV

    result=`cat $PGDEV | fgrep "Result: OK:"`
    if [ "$result" = "" ]; then
         cat $PGDEV | fgrep Result:
    fi
}

# 为0号线程绑定eth0网卡
PGDEV=/proc/net/pktgen/kpktgend_0
pgset "rem_device_all"   # 清空网卡绑定
pgset "add_device eth0"  # 添加eth0网卡

# 配置eth0网卡的测试选项
PGDEV=/proc/net/pktgen/eth0
pgset "count 1000000"    # 总发包数量
pgset "delay 5000"       # 不同包之间的发送延迟(单位纳秒)
pgset "clone_skb 0"      # SKB包复制
pgset "pkt_size 64"      # 网络包大小
pgset "dst 192.168.0.30" # 目的IP
pgset "dst_mac 11:11:11:11:11:11"  # 目的MAC

# 启动测试
PGDEV=/proc/net/pktgen/pgctrl
pgset "start"


$ cat /proc/net/pktgen/eth0
Params: count 1000000  min_pkt_size: 64  max_pkt_size: 64
     frags: 0  delay: 0  clone_skb: 0  ifname: eth0
     flows: 0 flowlen: 0
...
Current:
     pkts-sofar: 1000000  errors: 0
     started: 1534853256071us  stopped: 1534861576098us idle: 70673us
...
Result: OK: 8320027(c8249354+d70673) usec, 1000000 (64byte,0frags)
  120191pps 61Mb/sec (61537792bps) errors: 0
```

#### iperf3

测试TCP/UDP的性能的方法

```
# Ubuntu
apt-get install iperf3
# CentOS
yum install iperf3

# -s表示启动服务端，-i表示汇报间隔，-p表示监听端口
$ iperf3 -s -i 1 -p 10000

# -c表示启动客户端，192.168.0.30为目标服务器的IP
# -b表示目标带宽(单位是bits/s)
# -t表示测试时间
# -P表示并发数，-p表示目标服务器监听端口
$ iperf3 -c 192.168.0.30 -b 1G -t 15 -P 2 -p 10000
```

#### ab

测试HTTP
```
# -c表示并发请求数为1000，-n表示总的请求数为10000
$ ab -c 1000 -n 10000 http://192.168.0.30/```
```
#### wrk

应用负载测试

一个 HTTP 性能测试工具，内置了 LuaJIT，方便你根据实际需求，生成所需的请求负载，或者自定义响应的处理方法。

**安装**

```
$ git clone --depth=1 https://github.com/wg/wrk
$ cd wrk
$ apt-get install build-essential -y
$ make
$ sudo cp wrk /usr/local/bin/
```
**使用**

```
# 2 个线程、并发 1000 连接
# -c表示并发连接数1000，-t表示线程数为2
$ wrk -c 1000 -t 2 http://192.168.0.30/
Running 10s test @ http://192.168.0.30/
  2 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    65.83ms  174.06ms   1.99s    95.85%
    Req/Sec     4.87k   628.73     6.78k    69.00%
  96954 requests in 10.06s, 78.59MB read
  Socket errors: connect 0, read 0, write 0, timeout 179
Requests/sec:   9641.31
Transfer/sec:      7.82MB
# --latency TCP 延迟确认（Delayed ACK）的最小超时时间。
$ wrk --latency -c 100 -t 2 --timeout 2 http://192.168.0.30/
```

当然，wrk 最大的优势，是其内置的 LuaJIT，可以用来实现复杂场景的性能测试。wrk 在调用 Lua 脚本时，可以将 HTTP 请求分为三个阶段，即 setup、running、done，
如下图所示：![img](https://static001.geekbang.org/resource/image/d0/82/d02b845aa308b7a38a5735f3db8d9682.png)


比如，你可以在 setup 阶段，为请求设置认证参数[来自于 wrk 官方示例](https://github.com/wg/wrk/blob/master/scripts/auth.lua)：

```lua
-- example script that demonstrates response handling and
-- retrieving an authentication token to set on all future
-- requests

token = nil
path  = "/authenticate"

request = function()
   return wrk.format("GET", path)
end

response = function(status, headers, body)
   if not token and status == 200 then
      token = headers["X-Token"]
      path  = "/resource"
      wrk.headers["X-Token"] = token
   end
end
```

而在执行测试时，通过 -s 选项，执行脚本的路径：`$ wrk -c 1000 -t 2 -s auth.lua http://192.168.0.30/wrk `
需要你用 Lua 脚本，来构造请求负载。这对于大部分场景来说，可能已经足够了 。不过，它的缺点也正是，所有东西都需要代码来构造，并且工具本身不提供 GUI 环境。像 Jmeter 或者 LoadRunner（商业产品），则针对复杂场景提供了脚本录制、回放、GUI 等更丰富的功能，使用起来也更加方便。