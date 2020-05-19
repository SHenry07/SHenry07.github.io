IPVS provides more options for balancing traffic to backend Pods; these are:

- `rr`: round-robin 轮询
- `lc`: least connection (smallest number of open connections) 最少连接
- `dh`: destination hashing 目的hash
- `sh`: source hashing 资源hash 权重
- `sed`: shortest expected delay 最低延迟
- `nq`: never queue 





**static-rr-->不支持动态hash,没有后端数量限制,基本不用**

------

**leastconn-->类似于lvs中的wlc**

   不过这里只考虑活动连接数,即选择活动连接数少的,另外,最好在长连接会话中使用,如sql,ldap

------

**source-->基于hash表的算法,类似于nginx中的iphash**

   键:原IP地址的hash/值:挑选过的server,应用于动态服务器,保持会话

  hash-type map-based   静态hash    对于ip的hash取余

  hash-type consistent   动态hash   使用hash环,原理查看( [memcached多个memcached的解决方法](http://www.cnblogs.com/aaa103439/p/wiz:open_document?guid=360b182b-10bc-4f80-adae-a178fa0ee582&kbguid=) )

------

**uri-->基于uri生成hash表的算法,主要用于后端是缓存服务器**

  基于uri来进行选择,比如客户访问了 [http://test.com/a.jpg](http://test.com/a.jpg,那么在此uri缓存失效之前,任何一个客户访问这个路径,一定访问的就是a.jpg) , 那么在此uri缓存失效之前,任何一个客户访问这个路径,一定访问的就是a.jpg

  len    基于多少个字符的uri

  depth   基于多少个目录层次的uri

  例子: http://test.com/a/b/c/d/e/a.jpg

​     len 3   hash的uri是/a/

​     depth 3   hash的uri是/a/b/c/

------

**url_params-->根据url的参数来调度,用于将同一个用户的信息,都发送到同一个后端server**

  参数指那一部分,先看下url格式,其中绿色加粗的就是

  <scheme>://<user>:<passwd>@<host>:<port>/<path>; **** ?<query>#<frag>

  <query>表示php程序请求的查询信息

  <frag>表示当前页的片段页,即跳转到当前页的某个部分

  例子: http://test.com/hammers:sale=false/index.html:graphics=ture

  这里的意思就是如果访问的是hammers,则传递参数sale(并赋值false),如果访问的是index.html,则传递参数graphics(并赋值ture)

------

**hdr(name)-->header基于首部的信息来构建hash表 [HTTP之报文|首部](http://www.cnblogs.com/aaa103439/p/wiz:open_document?guid=12d5d429-ce56-4406-98f0-d8f7e51a8959&kbguid=)**

  hdr(Host)   基于用户请求的主机名进行调度 



四、负载均衡算法
1、轮询法

“风水轮流转 今年到我家”，将请求按顺序轮流地分配到后端服务器上，它均衡地对待后端的每一台服务器，而不关心服务器实际的连接数和当前的系统负载。

这样看似我们解决了资源调配的问题，但是资源真正合理利用了吗，愿望很美丽现实却很骨感，实际应用中同一个服务会部署到不同的硬件环境，性能各不相同。若直接使用简单轮询调度算法，给每个服务实例相同的负载，那么，必然会出现有的资源紧张有的浪费的情况。因此为了避免这种情况发生，业界就提出了加权轮询算法的解决方案。

2、加权轮询法

“能力越大，责任越大”，不同的后端服务器可能机器的配置和当前系统的负载并不相同，因此它们的抗压能力也不相同。给配置高、负载低的机器配置更高的权重，让其处理更多的请；而配置低、负载高的机器，给其分配较低的权重，降低其系统负载，加权轮询能很好地处理这一问题，并将请求顺序且按照权重分配到后端

3、随机法

“点兵点将,大兵大将,小兵小将…”，通过系统的随机算法，根据后端服务器的列表大小值来随机选取其中的一台服务器进行访问。看似随机，其实当请求数增大到一定程度的时候，随机优势不再明显，而最终的分配效果可能更加贴近于轮询方式。

4、加权随机法

与加权轮询法一样，加权随机法也根据后端机器的性能配置，给系统的负载分配不同的权重值，调度者会依据权重值进行随机抽取相应的服务并进行请求，而不是同普通随机法是按照点兵点将的方式盲目提供。

5、哈希值寻址法

源地址哈希是利用客户端的IP，通过哈希函数计算得出一个hash数，用该数值对服务器个数进行取模，取模后的余数即为访问服务器的序号。采用源地址哈希法进行负载均衡，同一IP地址的客户端，当后端服务器列表不变时，它每次都会映射到同一台后端服务器进行访问，这个算法也适用于一些特殊场景业务的需要例如session的存取等。

6、最小连接数法

“谁有空谁来”，最小连接数算法想对来说已经从死板的轮询和加权中脱离出来，它具有一定的智能性，还是源于我们的后端服务器的配置各不相同，对于请求的处理有快有慢，对于请求的堆积情况有多有少，调度服务需要获取众多服务器当前的实际链接情况，动态地选取其中当前积压连接数最少的一台服务器来处理当前的请求，尽可能地提高后端服务的利用效率，将负责合理地分流到每一台服务器
————————————————
版权声明：本文为CSDN博主「老猫的TOM」的原创文章，遵循 CC 4.0 BY-SA 版权协议，转载请附上原文出处链接及本声明。
原文链接：https://blog.csdn.net/qq355667166/article/details/86626564