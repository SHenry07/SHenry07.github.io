# 名词解释

| 词                                         | 解释                                                         |
| ------------------------------------------ | ------------------------------------------------------------ |
| MTU(Maximum Transmission Unit)最大传输单元 | 规定了最大的 IP 包大小.在我们最常用的以太网中,MTU 默认值是 1500(这也是 Linux 的默认值) |
| TTL(time to live)                          | 生存时间,或者跳数                                            |
| MSL(Maximum Segment Lifetime)              | 报文最大生存时间, TCP的TIME_WAIT状态也称为2MSL等待状态,MSL要大于等于TTL |
| 往返延时 RTT（Round-Trip Time）            | 一个数据包往返所需时间,平均往返延迟                          |
| DMA(Direct Memory Access)                  | 直接存储器访问, 允许不同速度的硬件装置来沟通，而不需要依赖于 CPU 的大量中断负载。否则，CPU 需要从来源把每一片段的资料复制到[暂存器](https://baike.baidu.com/item/暂存器/4308343)，然后把它们再次写回到新的地方。 |
| ring buff 环形缓冲区                       | 由于需要 DMA 与网卡交互，理应属于网卡设备驱动的范围          |
| sk_buff 缓冲区                             | 一个维护网络帧结构的双向链表，链表中的每一个元素都是一个网络帧（Packet）。虽然 TCP/IP 协议栈分了好几层，但上下不同层之间的传递，实际上只需要操作这个数据结构中的指针，而无需进行数据复制 |
| skb_buff 套接字缓冲区                      | 允许应用程序，给每个套接字配置不同大小的接收或发送缓冲区。应用程序发送数据，实际上就是将数据写入缓冲区；而接收数据，其实就是从缓冲区中读取。至于缓冲区中数据的进一步处理，则由传输层的 TCP 或 UDP 协议来完成。 |
| 五元组                                     | 协议，源 IP、源端口、目的 IP、目的端口                       |

> 老师好，一直不太明白skb_buff和sk_buff的区别，这两者有关系吗
>
> sk_buff一般是说内核数据接口，而 skb则是套接字缓存(socket buffer)
>
> 实际上，sk_buff、套接字缓冲、连接跟踪等，都通过 slab 分配器来管理。你可以直接通过 /proc/slabinfo，来查看它们占用的内存大小。

# 网卡收发报文的过程：

<img src="../image/c7b5b16539f90caabb537362ee7c27ac.png" alt="network01" style="zoom:33%;" />

1. 内核分配一个主内存地址段（DMA缓冲区)，网卡设备可以在DMA缓冲区中读写数据
2. 当一个网络帧(Frame)到达网卡后，网卡会通过 DMA 方式，把这个网络包放到收包队列中；然后通过硬中断，告诉中断处理程序已经收到了网络包。
3. 硬中断处理程序锁定当前**DMA缓冲区**，为网络帧分配内核数据结构（sk_buff），并将其拷贝到 **sk_buff 缓冲区**中；清空并解锁当前DMA缓冲区, 然后再通过软中断，通知内核收到了新的网络帧。
4. 内核协议栈从缓冲区中取出网络帧，并通过网络协议栈，从下到上逐层处理这个网络帧,比如，
- 在链路层检查报文的合法性，找出上层协议的类型（比如 IPv4 还是 IPv6），再去掉帧头、帧尾，然后交给网络层。
- 网络层取出 IP 头，判断网络包下一步的走向，比如是交给上层处理还是转发。当网络层确认这个包是要发送到本机后，就会取出上层协议的类型（比如 TCP 还是 UDP），去掉 IP 头，再交给传输层处理。
- 传输层取出 TCP 头或者 UDP 头后，根据 < 源 IP、源端口、目的 IP、目的端口 > 四元组作为标识，找出对应的 Socket，并把数据拷贝到 Socket 的接收缓存中。
6. 当发送数据包时，与上述相反。应用程序调用 Socket API（比如 sendmsg）发送网络包。由于这是一个系统调用，所以会陷入到内核态的套接字层中。套接字层会把数据包放到 **Socket 发送缓冲区**中。网络协议栈从 Socket 发送缓冲区中，取出数据包；再按照 TCP/IP 栈，比如，传输层和网络层，分别为其增加 TCP 头和 IP 头，执行路由查找确认下一跳的 IP，并按照 MTU 大小进行分片。

7. 分片后的网络包，再送到网络接口层，进行物理地址寻址，以找到下一跳的 MAC 地址。然后添加帧头和帧尾，放到发包队列中。这一切完成后，会有软中断通知驱动程序：发包队列中有新的网络帧需要发送。最后，驱动程序通过 DMA ，从发包队列中读出网络帧，并通过物理网卡把它发送出去。

   > sk_buff、套接字缓冲、连接跟踪等，都通过 slab 分配器来管理。你可以直接通过 /proc/slabinfo，来查看它们占用的内存大小。
   >
   > sk_buff、套接字缓冲：这是两个不同的概念，具体到数据上，内核协议栈都是操作指针，并不会在不同协议层之间复制数据

![img](../image/network.png)

![img](../image/network2.png)

#  性能指标

实际上，我们通常用带宽、吞吐量、延时、PPS（Packet Per Second）等指标衡量网络的性能。

- **带宽**，表示链路的最大传输速率，单位通常为 b/s （比特 / 秒）。

- **吞吐量**，表示单位时间内成功传输的数据量，单位通常为 b/s（比特 / 秒）或者 B/s（字节 / 秒）。吞吐量受带宽限制，而吞吐量 / 带宽，也就是该网络的使用率。
- **延时**，表示从网络请求发出后，一直到收到远端响应，所需要的时间延迟。在不同场景中，这一指标可能会有不同含义。比如，它可以表示，建立连接需要的时间（比如 TCP 握手延时），或一个数据包往返所需的时间（比如 RTT）。
- **PPS**，是 Packet Per Second（包 / 秒）的缩写，表示以网络包为单位的传输速率。PPS 通常用来评估网络的转发能力，比如硬件交换机，通常可以达到线性转发（即 PPS 可以达到或者接近理论最大值）。而基于 Linux 服务器的**转发**，则容易受网络包大小的影响。

除了这些指标，**网络的可用性（网络能否正常通信）、并发连接数（TCP 连接数量）、丢包率（丢包百分比）、重传率（重新传输的网络包比例**）等也是常用的性能指标。

## 如何进行网络测试

- 应用层，我们关注的是应用程序的并发连接数、每秒请求数、处理延迟、错误数等，可以使用 wrk、JMeter 等工具，模拟用户的负载，得到想要的测试结果。 
- 传输层，我们关注的是 TCP、UDP 等传输层协议的工作状况，比如 TCP 连接数、 TCP 重传、TCP 错误数等。此时，你可以使用 iperf、netperf 等，来测试 TCP 或 UDP 的性能。
- 网络接口层和网络层，它们主要负责网络包的封装、寻址、路由以及发送和接收。我们关注的则是网络包的处理能力，即 PPS。**特别是 64B 小包的处理能力**，值得我们特别关注， 可选测试工具 hping3 和pktgen

  由于低层协议是高层协议的基础，底层是其上方各层的基础，底层性能也就决定了高层性能，所以一般情况下，我们所说的网络优化，实际上包含了整个网络协议栈的所有层的优化。当然，性能要求不同，具体需要优化的位置和目标并不完全相同。

# 套接字socket

套接字可以屏蔽掉 Linux 内核中不同协议的差异，为应用程序提供统一的访问接口。每个套接字，都有一个读写缓冲区。

- 读缓冲区，缓存了远端发过来的数据。如果读缓冲区已满，就不能再接收新的数据。
- 写缓冲区，缓存了要发出去的数据。如果写缓冲区已满，应用程序的写操作就会被阻塞。

接收队列（Recv-Q）和发送队列（Send-Q）需要你特别关注，它们通常应该是 0。当你发现它们不是 0 时，说明有网络包的堆积发生。当然还要注意，在不同套接字状态下，它们的含义不同。

- 当套接字处于连接状态（Established）时，

  - Recv-Q 表示套接字缓冲还没有被应用程序取走的字节数（即接收队列长度）
  - Send-Q 表示还没有被远端主机确认的字节数（即发送队列长度）
- 当套接字处于监听状态（Listening）时，
  - Recv-Q 表示全连接队列当前使用了多少,也就是全连接队列的当前长度。 
  - Send-Q 表示全连接队列的最大长度。

 >  所谓全连接，是指服务器收到了客户端的 ACK，完成了 TCP 三次握手，然后就会把这个连接挪到全连接队列中。这些全连接中的套接字，还需要被 accept() 系统调用取走，服务器才可以开始真正处理客户端的请求。

## 半连接/全连接

```
net.ipv4.tcp_max_syn_backlog 半连接容量
net.ipv4.tcp_synack_retries = 1 连接每个 SYN_RECV 时，如果失败的话，内核还会自动重试，centos默认的重试次数是 5 次。你可以执行下面的命令，将其减小为 1 次：
```

>  半连接状态不只这一个参数net.ipv4.tcp_max_syn_backlog 控制，实际是 这个，还有系统somaxcon，以及应用程序的backlog，三个一起控制的。具体可以看下相关的源码

TCP SYN Cookies 也是一种专门防御 SYN Flood 攻击的方法。SYN Cookies 基于连接信息（包括源地址、源端口、目的地址、目的端口等）以及一个加密种子（如系统启动时间），计算出一个哈希值（SHA1），这个哈希值称为 cookie

> 注意，开启 TCP syncookies 后，内核选项 net.ipv4.tcp_max_syn_backlog 也就无效了。

`net.ipv4.tcp_syncookies = 1`

**全连接队列的大小取决于：min(backlog, somaxconn) . backlog是在socket创建的时候传入的，somaxconn是一个os级别的系统参数**

**半连接队列的大小取决于：max(64, /proc/sys/net/ipv4/tcp_max_syn_backlog)。 不同版本的os会有些差异**

`$ netstat -s|egrep "listen|LISTEN"的溢出值一直在上升`

溢出值在升高，说明全连接队列满了，而全连接队列的长度是由backlog与somaxconn决定的，为min(backlog, somaxconn)。可以通过 cat /proc/sys/net/core/somaxconn 查看参数。

 >  `cat /proc/sys/net/ipv4/tcp_abort_on_overflow `
 >
 >  **tcp_abort_on_overflow 为0表示如果三次握手第三步的时候全连接队列满了那么server扔掉client 发过来的ack（在server端认为连接还没建立起来）**

>  与全连接队列相对应的，还有一个**半连接队列**。所谓半连接是指还没有完成 TCP 三次握手的连接，连接只进行了一半。服务器收到了客户端的 SYN 包后，就会把这个连接放到半连接队列中，然后再向客户端发送 SYN+ACK 包。
>
>  在linux 2.2以前，backlog大小包括了半连接状态和全连接状态两种队列大小。linux 2.2以后，分离为两个backlog来分别限制半连接SYN_RCVD状态的未完成连接队列大小跟全连接ESTABLISHED状态的已完成连接队列大小。互联网上常见的TCP SYN FLOOD恶意DOS攻击方式就是用/proc/sys/net/ipv4/tcp_max_syn_backlog来控制的，可参见《[TCP洪水攻击（SYN Flood）的诊断和处理](http://tech.uc.cn/?p=1790)》。
>
>  在使用listen函数时，内核会根据传入参数的backlog跟系统配置参数/proc/sys/net/core/somaxconn中，二者取最小值，作为“ESTABLISHED状态之后，完成TCP连接，等待服务程序ACCEPT”的队列大小。在kernel 2.4.25之前，是写死在代码常量SOMAXCONN，默认值是128。在kernel 2.4.25之后，在配置文件/proc/sys/net/core/somaxconn (即 /etc/sysctl.conf 之类 )中可以修改。我稍微整理了流程图，如下：
>
>  ![tcp-sync-queue-and-accept-queue-small](D:\Dropbox\linux\image\tcp_three_handshakes.jpg)
>
>   ss —lntp 这个 当session处于listening中 rec-q 确定是 syn的backlog吗？ 
>
>  作者回复: 是的
>
>  自己回复：不是,Recv-Q为全连接队列当前使用了多少,send-Q才是应用的backlog，表示全队列的最长长度

#### 全连接队列满了会影响半连接队列吗？

TCP三次握手第一步的时候如果全连接队列满了会影响第一步drop 半连接的发生。大概流程的如下：

```
tcp_v4_do_rcv->tcp_rcv_state_process->tcp_v4_conn_request
//如果accept backlog队列已满，且未超时的request socket的数量大于1，则丢弃当前请求  
  if(sk_acceptq_is_full(sk) && inet_csk_reqsk_queue_yong(sk)>1)
      goto drop;
```

### socket优化

为了提高网络的吞吐量，通常需要调整这些缓冲区的大小。``cat /proc/net/sockstat`比如：

增大每个套接字的缓冲区大小 net.core.optmem_max；

增大套接字接收缓冲区大小 net.core.rmem_max 和发送缓冲区大小 net.core.wmem_max；

增大 TCP 接收缓冲区大小 net.ipv4.tcp_rmem 和发送缓冲区大小 net.ipv4.tcp_wmem。

![img](../image/5f2d4957663dd8bf3410da8180ab18f0.png)

有几点需要你注意。

- tcp_rmem 和 tcp_wmem 的三个数值分别是 min，default，max，系统会根据这些设置，自动调整 TCP 接收 / 发送缓冲区的大小。
- udp_mem 的三个数值分别是 min，pressure，max，系统会根据这些设置，自动调整 UDP 发送缓冲区的大小。

当然，表格中的数值只提供参考价值，具体应该设置多少，还需要你根据实际的网络状况来确定。比如，发送缓冲区大小，**理想数值是吞吐量 * 延迟，**这样才可以达到最大网络利用率。
除此之外，套接字接口还提供了一些配置选项，用来修改网络连接的行为：-

- 为 TCP 连接设置 TCP_NODELAY 后，就可以禁用 Nagle 算法；

- 为 TCP 连接开启 TCP_CORK 后，可以让小包聚合成大包后再发送（注意会阻塞小包的发送）；
- 使用 SO_SNDBUF 和 SO_RCVBUF ，可以分别调整套接字发送缓冲区和接收缓冲区的大小。

# 协议栈统计信息

```
$ netstat -s
...
Tcp:
    3244906 active connection openings
    23143 passive connection openings
    115732 failed connection attempts
    2964 connection resets received
    1 connections established
    13025010 segments received
    17606946 segments sent out
    44438 segments retransmitted
    42 bad segments received
    5315 resets sent
    InCsumErrors: 42
...

$ ss -s
Total: 186 (kernel 1446)
TCP:   4 (estab 1, closed 0, orphaned 0, synrecv 0, timewait 0/0), ports 0

Transport Total     IP        IPv6
*    1446      -         -
RAW    2         1         1
UDP    2         2         0
TCP    4         3         1
...
```

这些协议栈的统计信息都很直观。ss 只显示已经连接、关闭、孤儿套接字等简要统计，而 netstat 则提供的是更详细的网络协议栈信息。

# 网络吞吐和 PPS(Packet Per Second)

```

# 数字1表示每隔1秒输出一组数据
$ sar -n DEV 1
Linux 4.15.0-1035-azure (ubuntu)   01/06/19   _x86_64_  (2 CPU)

13:21:40        IFACE   rxpck/s   txpck/s    rxkB/s    txkB/s   rxcmp/s   txcmp/s  rxmcst/s   %ifutil
13:21:41         eth0     18.00     20.00      5.79      4.25      0.00      0.00      0.00      0.00
13:21:41      docker0      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
13:21:41           lo      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
```

- rxpck/s 和 txpck/s 分别是接收和发送的 PPS，单位为包 / 秒。
- rxkB/s 和 txkB/s 分别是接收和发送的吞吐量，单位是 KB/ 秒。
- rxcmp/s 和 txcmp/s 分别是接收和发送的压缩数据包数，单位是包 / 秒。
- %ifutil 是网络接口的使用率，即半双工模式下为 (rxkB/s+txkB/s)/Bandwidth，而全双工模式下为 max(rxkB/s, txkB/s)/Bandwidth。

# C10K和C1000K

因为五元组的关系，客户端最大支持 65535 个连接，而服务器端可支持的连接数是海量的。当然，由于 Linux 协议栈本身的性能，以及各种物理和软件的资源限制等，这么大的连接数，还是远远达不到的（实际上，C10M 就已经很难了）

## I/O 模型优化

异步、非阻塞 I/O 的解决思路，你应该听说过，其实就是我们在网络编程中经常用到的 I/O 多路复用（I/O Multiplexing）。I/O 多路复用是什么意思呢？

两种 I/O 事件通知的方式：水平触发和边缘触发，它们常用在套接字接口的文件描述符中。

- **水平触发**(LT)：只要文件描述符可以非阻塞地执行 I/O ，就会触发通知。也就是说，应用程序可以随时检查文件描述符的状态，然后再根据状态，进行 I/O 操作。

- **边缘触发**(ET)：只有在文件描述符的状态发生改变（也就是 I/O 请求达到）时，才发送一次通知。这时候，应用程序需要尽可能多地执行 I/O，直到无法继续读写，才可以停止。如果 I/O 没执行完，或者因为某种原因没来得及处理，那么这次通知也就丢失了。

  > select/poll是LT模式，epoll缺省使用的也是水平触发模式（LT）。
  > 目前业界对于ET的最佳实践大概就是Nginx了，单线程redis也是使用的LT
  > 
  > LT:文件描述符准备就绪时（FD关联的读缓冲区不为空，可读。写缓冲区还没满，可写），触发通知。
  > 也就是你文中表述的"只要文件描述符可以非阻塞地执行 I/O ，就会触发通知..."
  > ET:当FD关联的缓冲区发生变化时（例如：读缓冲区由空变为非空，有新数据达到，可读。写缓冲区满变有空间了，有数据被发送走，可写），触发通知，仅此一次
  > 也就是你文中表述的"只有在文件描述符的状态发生改变（也就是 I/O 请求达到）时"

[详见](../io.md)

## 工作模型优化

### 主进程+多个worker子进程

- 主进程执行bind() + listen(), 然后创建多个子进程， 并管理子进程的生命周期

- 在每个子进程中，都通过accept() 或 epoll_wait(), 来处理相同的socket

- 如nginx

  <img src="https://static001.geekbang.org/resource/image/45/7e/451a24fb8f096729ed6822b1615b097e.png" alt="img" style="zoom:33%;" />

### 监听到相同端口的多进程模型

所有的进程都监听相同的接口，并且开启 SO_REUSEPORT 选项，由内核负责将请求负载均衡到这些监听进程中去

<img src="../image/90df0945f6ce5c910ae361bf2b135bbd.png" alt="img" style="zoom:33%;" />



## 应用层的网络协议优化

- 使用长连接取代短连接，可以显著降低TCP建立连接的成本
- 使用内存等方式，来缓存不长变化的数据，可以降低网络I/O次数，同时降快应用程序的响应速度
- 使用protocol buffer等序列化的方式，压缩网络I/O的数据量，可以提高应用程序的吞吐
- 使用DNS缓存，预取，HTTPDNS等方式，减少DNS解析的延迟，也可以提高网络I/O的整体速度

## C1000K C10M

多队列网卡、中断负载均衡、CPU 绑定、RPS/RFS（软中断负载均衡到多个 CPU 核上），以及将网络包的处理卸载（Offload）到网络设备（如 TSO/GSO、LRO/GRO、VXLAN OFFLOAD）等各种硬件和软件的优化。

跳过内核协议栈的冗长路径，把网络包直接送到要处理的应用程序那里去。这里有两种常见的机制，DPDK 和 XDP。

用 XDP 的方式，在内核协议栈之前处理网络包；或者用 DPDK 直接跳过网络协议栈，在用户空间通过轮询的方式直接处理网络包。

### DPDK

第一种机制，DPDK，是用户态网络的标准。它**跳过内核协议栈**，直接由用户态进程通过轮询的方式，来处理网络接收。

<img src="../image/998fd2f52f0a48a910517ada9f2bb23a.png" alt="img" style="zoom:33%;" />

### XDP

第二种机制，XDP（eXpress Data Path），则是 Linux 内核提供的一种高性能网络数据路径。它允许网络包，在进入内核协议栈之前，就进行处理，也可以带来更高的性能。XDP 底层跟我们之前用到的 bcc-tools 一样，都是基于 Linux 内核的 eBPF 机制实现的。

![img](../image/067ef9df4212cd4ede3cffcdac7001be.png)

基于 XDP 的应用程序通常是专用的网络应用，常见的有 IDS（入侵检测系统）、DDoS 防御、 [cilium](https://github.com/cilium/cilium) 容器网络插件等。

# QPS(Query per Second)

对 TCP 或者 Web 服务来说，更多会用并发连接数和每秒请求数（QPS，Query per Second）

web应用层，游戏通常会基于TCP或UDP

### DDOS

DDoS 的前身是 DoS（Denail of Service），即拒绝服务攻击

DDoS（Distributed Denial of Service） 则是在 DoS 的基础上，采用了分布式架构，利用多台主机同时攻击目标主机。


# 网络延迟

- 网络传输慢，导致延迟；

- Linux 内核协议栈报文处理慢，导致延迟；

- 应用程序数据处理慢，导致延迟等等。

  #### client端

  查询 TCP 文档（执行 man tcp），你就会发现，只有 TCP 套接字专门设置了 TCP_QUICKACK ，才会开启快速确认模式；否则，默认情况下，采用的就是延迟确认机制：

  ```
  TCP_QUICKACK (since Linux 2.4.4)
                Enable  quickack mode if set or disable quickack mode if cleared.  In quickack mode, acks are sent immediately, rather than delayed if needed in accordance to normal TCP operation.  This flag is  not  permanent,  it only enables a switch to or from quickack mode.  Subsequent operation of the TCP protocol will once again enter/leave quickack mode depending on internal  protocol  processing  and  factors  such  as delayed ack timeouts occurring and data transfer.  This option should not be used in code intended to be portable.
  ```

  #### server端

  Nagle 算法，是 TCP 协议中用于减少小包发送数量的一种优化算法，目的是为了提高实际带宽的利用率。
  举个例子，当有效负载只有 1 字节时，再加上 TCP 头部和 IP 头部分别占用的 20 字节，整个网络包就是 41 字节，这样实际带宽的利用率只有 2.4%（1/41）。往大了说，如果整个网络带宽都被这种小包占满，那整个网络的有效利用率就太低了。Nagle 算法正是为了解决这个问题。它通过合并 TCP 小包，提高网络带宽的利用率。Nagle 算法规定，一个 TCP 连接上，最多只能有一个未被确认的未完成分组；在收到这个分组的 ACK 前，不发送其他分组。这些小分组会被组合起来，并在收到 ACK 后，用同一个分组发送出去。

  ```shell
  $ man tcp 
  TCP_NODELAY
                If set, disable the Nagle algorithm.  This means that segments are always sent as soon as possible, even if there is only a small amount of data.  When not set, data is buffered until  there  is  a  sufficient amount  to  send out, thereby avoiding the frequent sending of small packets, which results in poor utilization of the network.  This option is overridden by TCP_CORK; however, setting this option forces  an explicit flush of pending output, even if TCP_CORK is currently set.
  ```

  > tcp_nodelay off; 默认nginx中是on

# NAT

NAT 基于Linux内核的连接追踪机制，实现了IP地址及端口号重写的功能。

主要目的，是实现地址转换。根据实现方式的不同，NAT 可以分为三类：

- 静态 NAT，即内网 IP 与公网 IP 是一对一的永久映射关系；
- 动态 NAT，即内网 IP 从公网 IP 池中，动态选择一个进行映射；
- 网络地址端口转换 NAPT（Network Address and Port Translation），即把内网 IP 映射到公网 IP 的不同端口上，让多个内网 IP 可以共享同一个公网 IP 地址。

MASQUERADE 是最常用的一种 SNAT 规则，常用来为多个内网 IP 地址提供共享的出口 IP。

当多个内网 IP 地址的端口号相同时，MASQUERADE 还可以正常工作吗？
如果内网 IP 地址数量或请求数比较多，这种方式有没有什么隐患呢？

问题1：Linux的NAT时给予内核的连接跟踪模块实现，保留了源IP、源端口、目的IP、目的端口之间的关系，多个内网IP地址的端口相同，但是IP不同，再nf_conntrack中对应不同的记录，所以MASQUERADE可以正常工作。

问题2：NAT方式所有流量都要经过NAT服务器，所以NAT服务器本身的软中断导致CPU负载、网络流量、文件句柄、端口号上限、nf_conntrack table full都可能是性能瓶颈

## DNAT工作流程

根据Netfilter 中，网络包的流向以及 NAT 的原理，要保证 NAT 正常工作，就至少需要两个步骤：

- 第一，利用 **Netfilter 中的钩子函数（Hook）**，修改源地址或者目的地址。
- 第二，利用**连接跟踪模块 conntrack** ，关联同一个连接的请求和响应。

```shell
#! /usr/bin/env stap

############################################################
# Dropwatch.stp
# Author: Neil Horman <nhorman@redhat.com>
# An example script to mimic the behavior of the dropwatch utility
# http://fedorahosted.org/dropwatch
############################################################

# Array to hold the list of drop points we find
global locations

# Note when we turn the monitor on and off
probe begin { printf("Monitoring for dropped packets\n") }
probe end { printf("Stopping dropped packet monitor\n") }

# increment a drop counter for every location we drop at
probe kernel.trace("kfree_skb") { locations[$location] <<< 1 }

# Every 5 seconds report our drop locations
probe timer.sec(5)
{
  printf("\n")
  foreach (l in locations-) {
    printf("%d packets dropped at %s\n",
           @count(locations[l]), symname(l))
  }
  delete locations
}
```

[这个脚本](https://sourceware.org/systemtap/SystemTap_Beginners_Guide/useful-systemtap-scripts.html#nettopsect)，跟踪内核函数 kfree_skb() 的调用，并统计丢包的位置。文件保存好后，执行下面的 stap 命令，就可以运行丢包跟踪脚本。这里的 stap，是 SystemTap 的命令行工具：

```
$ stap --all-modules dropwatch.stp
Monitoring for dropped packets
# 当你看到 probe begin 输出的 “Monitoring for dropped packets” 时，表明 SystemTap 已经将脚本编译为内核模块，并启动运行了。
```
```
$ sysctl -a | grep conntrack
net.netfilter.nf_conntrack_count = 180 表示当前连接追踪数
net.netfilter.nf_conntrack_max = 1000  表示最大连接追踪数
net.netfilter.nf_conntrack_buckets = 65536 表示连接追踪表的大小
net.netfilter.nf_conntrack_tcp_timeout_syn_recv = 60
net.netfilter.nf_conntrack_tcp_timeout_syn_sent = 120
net.netfilter.nf_conntrack_tcp_timeout_time_wait = 120
...

# 连接跟踪对象大小为376，链表项大小为16
nf_conntrack_max*连接跟踪对象大小+nf_conntrack_buckets*链表项大小 
= 1000*376+65536*16 B
= 1.4 MB
```

你可以用 conntrack 命令行工具，来查看连接跟踪表的内容。比如：
```
# -L表示列表，-o表示以扩展格式显示
$ conntrack -L -o extended | head
ipv4     2 tcp      6 7 TIME_WAIT src=192.168.0.2 dst=192.168.0.96 sport=51744 dport=8080 src=172.17.0.2 dst=192.168.0.2 sport=8080 dport=51744 [ASSURED] mark=0 use=1
ipv4     2 tcp      6 6 TIME_WAIT src=192.168.0.2 dst=192.168.0.96 sport=51524 dport=8080 src=172.17.0.2 dst=192.168.0.2 sport=8080 dport=51524 [ASSURED] mark=0 use=1
```
连接跟踪表里的对象，包括了协议、连接状态、源 IP、源端口、目的 IP、目的端口、跟踪状态等。由于这个格式是固定的，所以可以用 awk、sort 等工具，对其进行统计分析。

[nf_conntrack](https://www.kernel.org/doc/Documentation/networking/nf_conntrack-sysctl.txt)的文档 具有更多的配置选项，可以根据实际需求具体来配置

> 摘自评论区
>
> https://mp.weixin.qq.com/s/VYBs8iqf0HsNg9WAxktzYQ：（多个容器snat时因为搜索本地可用端口（都从1025开始，到找到可用端口并插入到conntrack表是一个非事务并且有时延--第二个插入会失败，进而导致第一个syn包被扔掉的错误，扔掉后重传找到新的可用端口，表现就是时延偶尔为1秒或者3秒）
>
> 这篇文章是我见过诊断NAT问题最专业的，大家要多学习一下里面的思路和手段
>
> 其实遇到很多问题的时候多看看内核日志就知道了，linux很智能的，很多报错信息都在日志里面，越遇到系统优化层面，就多要看看内核日志，我一般是使用journalctl -k -f来查看，有报错信息就Google，线上遇到nf_conntrack: table full，就是这样排查出来的，查看内核日志真的很重要，特别应用日志没看出什么来的时候
>
> ```
> $ dmesg | tail
> [104235.156774] nf_conntrack: nf_conntrack: table full, dropping packet
> [104243.800401] net_ratelimit: 3939 callbacks suppressed
> [104243.800401] nf_conntrack: nf_conntrack: table full, dropping packet
> [104262.962157] nf_conntrack: nf_conntrack: table full, dropping packet
> ```


```
# Linux的IP转发功能
$ sysctl -w net.ipv4.ip_forward=1
net.ipv4.ip_forward = 1
```
-----

# 如何优化/几个思路

就拿 NAT 网关来说，由于其直接影响整个数据中心的网络出入性能，所以 NAT 网关通常需要达到或接近线性转发，也就是说， PPS 是最主要的性能目标。
再如，对于数据库、缓存等系统，快速完成网络收发，即低延迟，是主要的性能目标。
而对于 Web 服务来说，则需要同时兼顾吞吐量和延迟。

## TCP

第一类，在请求数比较大的场景下，你可能会看到大量处于 TIME_WAIT 状态的连接，它们会占用大量内存和端口资源。这时，我们可以优化与 TIME_WAIT 状态相关的内核选项

1. 增大处于 TIME_WAIT 状态的连接数量 net.ipv4.tcp_max_tw_buckets ，并增大连接跟踪表的大小 net.netfilter.nf_conntrack_max。
2. 减小 net.ipv4.tcp_fin_timeout 和 net.netfilter.nf_conntrack_tcp_timeout_time_wait ，让系统尽快释放它们所占用的资源。
3. 开启端口复用 net.ipv4.tcp_tw_reuse。这样，被 TIME_WAIT 状态占用的端口，还能用到新建的连接中。
4. 增大本地端口的范围 net.ipv4.ip_local_port_range 。这样就可以支持更多连接，提高整体的并发能力。
5. 增加最大文件描述符的数量。你可以使用 fs.nr_open 和 fs.file-max ，分别增大**单个进程和系统**(下面图里错了)的最大文件描述符数；或在应用程序的 systemd 配置文件中，配置 LimitNOFILE ，设置应用程序的最大文件描述符数。

第二类，为了缓解 SYN FLOOD 等，利用 TCP 协议特点进行攻击而引发的性能问题，你可以考虑优化与 SYN 状态相关的内核选项

第三类，在长连接的场景中，通常使用 Keepalive 来检测 TCP 连接的状态，以便对端连接断开后，可以自动回收。但是，系统默认的 Keepalive 探测间隔和重试次数，一般都无法满足应用程序的性能要求。所以，这时候你需要优化与 Keepalive 相关的内核选项，比如：

- 缩短最后一次数据包到 Keepalive 探测包的间隔时间 net.ipv4.tcp_keepalive_time；
- 缩短发送 Keepalive 探测包的间隔时间 net.ipv4.tcp_keepalive_intvl；
- 减少 Keepalive 探测失败后，一直到通知应用程序前的重试次数 net.ipv4.tcp_keepalive_probes。

![img](https://static001.geekbang.org/resource/image/b0/e0/b07ea76a8737ed93395736795ede44e0.png)

## UDP

UDP 提供了面向数据报的网络协议，它不需要网络连接，也不提供可靠性保障。所以，UDP 优化，相对于 TCP 来说，要简单得多。这里我也总结了常见的几种优化方案。跟上篇套接字部分提到的一样，

增大套接字缓冲区大小以及 UDP 缓冲区范围；

跟前面 TCP 部分提到的一样，增大本地端口号的范围；

根据 MTU 大小，调整 UDP 数据包的大小，减少或者避免分片的发生。

## 网络层
接下来，我们再来看网络层的优化。
网络层，负责网络包的封装、寻址和路由，包括 IP、ICMP 等常见协议。在网络层，最主要的优化，其实就是对路由、 IP 分片以及 ICMP 等进行调优。

第一种，从路由和转发的角度出发，你可以调整下面的内核选项。
在需要转发的服务器中，比如

- 用作 NAT 网关的服务器或者使用 Docker 容器时，开启 IP 转发，即设置 net.ipv4.ip_forward = 1。

- 调整数据包的生存周期 TTL，比如设置 net.ipv4.ip_default_ttl = 64。注意，增大该值会降低系统性能。

- 开启数据包的反向地址校验，比如设置 net.ipv4.conf.eth0.rp_filter = 1。这样可以防止 IP 欺骗，并减少伪造 IP 带来的 DDoS 问题。


 第二种，从分片的角度出发，最主要的是调整 MTU（Maximum Transmission Unit）的大小。
    通常，MTU 的大小应该根据以太网的标准来设置。以太网标准规定，一个网络帧最大为 1518B，那么去掉以太网头部的 18B 后，剩余的 1500 就是以太网 MTU 的大小。
    在使用 VXLAN、GRE 等叠加网络技术时，要注意，网络叠加会使原来的网络包变大，导致 MTU 也需要调整。
    比如，就以 VXLAN 为例，它在原来报文的基础上，增加了 14B 的以太网头部、 8B 的 VXLAN 头部、8B 的 UDP 头部以及 20B 的 IP 头部。换句话说，每个包比原来增大了 50B。
    所以，我们就需要把交换机、路由器等的 MTU，增大到 1550， 或者把 VXLAN 封包前（比如虚拟化环境中的虚拟网卡）的 MTU 减小为 1450。
    另外，现在很多网络设备都支持巨帧，如果是这种环境，你还可以把 MTU 调大为 9000，以提高网络吞吐量。
    第三种，从 ICMP 的角度出发，为了避免 ICMP 主机探测、ICMP Flood 等各种网络问题，你可以通过内核选项，来限制 ICMP 的行为。比如，你可以禁止 ICMP 协议，即设置 `net.ipv4.icmp_echo_ignore_all = 1`。这样，外部主机就无法通过 ICMP 来探测主机。或者，你还可以禁止广播 ICMP，即设置 `net.ipv4.icmp_echo_ignore_broadcasts = 1`。

## 链路层

网络层的下面是链路层, 链路层负责网络包在物理网络中的传输，比如 MAC 寻址、错误侦测以及通过网卡传输网络帧等。自然，链路层的优化，也是围绕这些基本功能进行的。接下来，我们从不同的几个方面分别来看。
    由于网卡收包后调用的中断处理程序（特别是软中断），需要消耗大量的 CPU。所以，将这些中断处理程序调度到不同的 CPU 上执行，就可以显著提高网络吞吐量。这通常可以采用下面两种方法。
    比如，你可以为网卡硬中断配置 CPU 亲和性（smp_affinity），或者开启 irqbalance 服务。
    再如，你可以开启 RPS（Receive Packet Steering）和 RFS（Receive Flow Steering），将应用程序和软中断的处理，调度到相同 CPU 上，这样就可以增加 CPU 缓存命中率，减少网络延迟。
    另外，现在的网卡都有很丰富的功能，原来在内核中通过软件处理的功能，可以卸载到网卡中，通过硬件来执行。

- TSO（TCP Segmentation Offload）和 UFO（UDP Fragmentation Offload）：在 TCP/UDP 协议中直接发送大包；而 TCP 包的分段（按照 MSS 分段）和 UDP 的分片（按照 MTU 分片）功能，由网卡来完成-

- GSO（Generic Segmentation Offload）：在网卡不支持 TSO/UFO 时，将 TCP/UDP 包的分段，延迟到进入网卡前再执行。这样，不仅可以减少 CPU 的消耗，还可以在发生丢包时只重传分段后的包。

-  LRO（Large Receive Offload）：在接收 TCP 分段包时，由网卡将其组装合并后，再交给上层网络处理。不过要注意，在需要 IP 转发的情况下，不能开启 LRO，因为如果多个包的头部信息不一致，LRO 合并会导致网络包的校验错误。

-  GRO（Generic Receive Offload）：GRO 修复了 LRO 的缺陷，并且更为通用，同时支持 TCP 和 UDP。

-  RSS（Receive Side Scaling）：也称为多队列接收，它基于硬件的多个接收队列，来分配网络接收进程，这样可以让多个 CPU 来处理接收到的网络包。

-  VXLAN 卸载：也就是让网卡来完成 VXLAN 的组包功能。
      最后，对于网络接口本身，也有很多方法，可以优化网络的吞吐量。
      比如，你可以开启网络接口的多队列功能。这样，每个队列就可以用不同的中断号，调度到不同 CPU 上执行，从而提升网络的吞吐量。
      再如，你可以增大网络接口的缓冲区大小，以及队列长度等，提升网络传输的吞吐量（注意，这可能导致延迟增大）。

  你还可以使用 Traffic Control 工具，为不同网络流量配置 QoS

# 精彩评论

1. 最大连接数是不是受限于 65535 个端口

服务器端的理论最大连接数，可以达到 2 的 48 次方（IP 为 32 位，端口号为 16 位），远大于 65535,综合来看，客户端最大支持 65535 个连接，而服务器端可支持的连接数是海量的。

2. 中断不均，连接跟踪打满

嗯 这是最常见的两个问题， 网络报文传需要在用户态和内核态来回切换，导致性能下降。业界使用零拷贝或intel的dpdk来提高性能。

3. 半连接状态不只这一个参数net.ipv4.tcp_max_syn_backlog 控制，实际是 这个，还有系统somaxcon，以及应用程序的backlog，三个一起控制的。具体可以看下相关的源码

4. 实际上，根据 IP 地址反查域名、根据端口号反查协议名称，是很多网络工具默认的行为，而这往往会导致性能工具的工作缓慢。所以，通常，网络性能工具都会提供一个选项（比如 -n 或者 -nn），来禁止名称解析。

# Reference

[关于TCP 半连接队列和全连接队列](http://jm.taobao.org/2017/05/25/525-1/)

[TCP SOCKET中backlog参数的用途是什么？](https://www.cnxct.com/something-about-phpfpm-s-backlog/)
