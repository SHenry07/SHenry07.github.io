

https://github.com/rootdeep/rootdeep.github.io

# Free 命令结果变化

最近，在用IOzone 对磁盘进行压力测试的时候发现，IOzone 一开启（216个线程同时写36个8T的磁盘，每个磁盘大概写6T的数据），使用free 命令查看系统内存时，64G的内存，free的内存为0 ， buff/cache 很高，并且系统的kswapd0 对CPU的占用率很高，高达90%+ 。随后查看了kswapd 的作用，大概是：kswapd是个内核进程，虚拟内存管理中负责换页的，当OS的可用内存小于阀值时，kswapd会将部分进程的页从物理内存交换到swap上，这个换页操作是十分消耗主机CPU资源的。

buffer/cache 的使用值这么高，又是文件的相关的读写，于是想起了文件系统中的Page Cache 概念。那么free命令显示的buff/cache 跟文件系统的这个page cache 有什么关系呢？

首先看一下free 命令显示的统计维度在内核版本中的变化：

在内核3.14 之前，使用free 命令结果是这样的：

```
free -m
             total       used       free     shared    buffers     cached
Mem:          2016       1973         42          0        163       1497
-/+ buffers/cache:        312       1703
Swap:         4094          0       4094
```

可用内存 = Mem行中的free + buffers + cached = 42+ 163+1497 = 1702 = (-/+ buffers/cache) 行的free

已使用内存 = Mem行中的used – buffers – cached = 1973 - 163 - 1497 = 313 = (-/+ buffers/cache) 行的used

总内存 = Mem行中的free +used

在内核3.14 中（2014年4月发布），结果是这样的：

```
free -h
               total       used        free      shared     buff/cache      available
Mem:           993M        253M        334M        39M        405M             556M
Swap:          2.0G        82M         1.9G
```

可用内存 即 available 的值

总内存 total = Mem 行中的 used + free + buff /cache

即：主要的变化是buff/cache被合并为一列，**并且增加了available这一列**

为何会有这样的变化？内核提交的相关commit 的解释是这样的：

相关commit: https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/commit/?id=34e431b0ae398fc54ea69ff85ec700722c9da773

Many load balancing and workload placing programs check /proc/meminfo to estimate how much free memory is available. They generally do this by adding up “free” and “cached”, which was fine ten years ago, but is pretty much guaranteed to be wrong today. It is wrong because Cached includes memory that is not freeable as page cache, for example shared memory segments, tmpfs, and ramfs, and it does not include reclaimable slab memory, which can take up a large fraction of system memory on mostly idle systems with lots of files.Currently, the amount of memory that is available for a new workload,without pushing the system into swap, can be estimated from MemFree, Active(file), Inactive(file), and SReclaimable, as well as the”low” watermarks from /proc/zoneinfo. However, this may change in the future, and user space really should not be expected to know kernel internals to come up with an estimate for the amount of free memory. It is more convenient to provide such an estimate in /proc/meminfo. If things change in the future, we only have to change it in one place.

大概意思就是说 原来计算可用内存的方式 free +buffers+cached 的方法在现在数十年前是可行的，但是现在是错误的方式。因为在现代的内存系统中，Cache 中的一些空间是不能被算作 freeable 的，比如：tmpfs,ramfs，memory segments 等，但是可释放的内存，如slab memory ，又没算到freeable 的内存中。所以，现在需要一种新的计算可用内存的方式。对于可用内存的定义是：在不适用交互分区下，可给应用程序分配的剩余内存大小。这些剩余内存包括了MemFree, Active(file), Inactive(file), and SReclaimable。

对于 3.14 及以后的内核， available才是机器的”可用内存” , 而不是像过去那样简单的把 free, cached 和 buffers 加起来。并且，可以想象，available 一定是小于 free+buffers+cached 的值了。

对于这样的改进，有必要对buff 和cache 的来龙去脉做个了解。

有的文章说cache 是负责读内容的缓存，buffer是负责写内容的缓存，这样的解释说服不了自己对文件系统cache 的理解。翻阅了很多文章，发现一片好文，感觉解释的很合理。摘其核心内容，以作备忘及副本，根据自己实验环境，多处略有改动。

------

参考文章来自：https://mp.weixin.qq.com/s?biz=MzA3NjYxOTA0MQ==&mid=2653965212&idx=2&sn=8b63bee228bdc5cb7441e0679c03808c&scene=0#rd

# 经学习后以下文章 有大量错误与猜想

# 什么是buffer/cache？

buffer和cache是两个在计算机技术中被用滥的名词，放在不通语境下会有不同的意义。在Linux的内存管理中，这里的buffer指Linux内存的：block buffer/Buffer cache。这里的cache指Linux内存中的：Page cache。翻译成中文可以叫做缓冲区缓存和页面缓存。

理论上，一个文件读首先到Block Buffer, 然后到Page Cache。有了文件系统才有了Page Cache. 在老的Linux上这两个Cache是分开的。 那这样对于文件数据，会被Cache两次。这种方案虽然简单，但低效。后期Linux把这两个Cache统一了。对于文件，Page Cache指向Block Buffer，对于非文件则是Block Buffer。这样就如文件实验的结果，文件操作，只影响Page Cache，Raw操作，则只影响Buffer. 比如一台VM虚拟机，则会越过File System，只接操作 Disk, 常说的Direct IO.

**什么是page cache**

Page cache主要用来作为文件系统上的文件数据的缓存来用，尤其是针对当进程对文件有read／write操作的时候。

**什么是buffer cache/block buffer**

Buffer cache则主要是设计用来在系统对块设备进行读写的时候，对块进行数据缓存的系统来使用。这意味着某些对块的操作会使用buffer cache进行缓存，比如我们在格式化文件系统的时候。一般情况下两个缓存系统是一起配合使用的，比如当我们对一个文件进行写操作的时候，page cache的内容会被改变，而buffer cache则可以用来将page标记为不同的缓冲区，并记录是哪一个缓冲区被修改了。这样，内核在后续执行脏数据的回写（writeback）时，就不用将整个page写回，而只需要写回修改的部分即可。

# 如何回收cache？

人工触发缓存清除的操作：

`echo 1 > /proc/sys/vm/drop_caches`表示清除page cache。

`echo 2 > /proc/sys/vm/drop_caches`:表示清除回收slab分配器中的对象（包括目录项缓存和inode缓存）。slab分配器是内核中管理内存的一种机制，其中很多缓存数据实现都是用的page cache。

`echo 3 > /proc/sys/vm/drop_caches`:表示清除page cache和slab分配器中的缓存对象。

```shell
          To free pagecache, use:

              echo 1 > /proc/sys/vm/drop_caches

          To free dentries and inodes, use:

              echo 2 > /proc/sys/vm/drop_caches

          To free pagecache, dentries and inodes, use:

              echo 3 > /proc/sys/vm/drop_caches
```
Because writing to this file is a nondestructive operation and dirty objects are not freeable, the  user should run sync(1) first.

It appears you are working with memory caching of directory structures. 当我们正在使用内存缓存目录结构时，

An ***inode*** in your context is a data structure that represents a file. A ***dentries*** is a data structure that represents a directory. 在这个上下文中的**inode是表示文件的数据结构，而dentries是表示目录的数据结构**。

These structures could be used to build a memory cache that represents the file structure on a disk. To get a directly listing, the OS could go to the dentries–if the directory is there–list its contents (a series of inodes). If not there, go to the disk and read it into memory so that it can be used again. 这些结构可用于构建表示磁盘上的文件结构的内存高速缓存。为了直接获得列表，操作系统可以去dentries那里（如果目录在那里的话）列出其内容（一系列inode）。如果没有，则会去磁盘上将其读入内存，以便它可以再次使用。

The ***page cache*** could contain any memory mappings to blocks on disk. That could conceivably be buffered I/O, memory mapped files, paged areas of executables–anything that the OS could hold in memory from a file. 页面缓存（page cache）可以包含磁盘块的任何内存映射。这可以是缓冲I/O，内存映射文件，可执行文件的分页区域——操作系统可以从文件保存在内存中的任何内容。

# 哪些cache 不能被回收？

在前边说到的改变可用内存的计算方式是因为不是所有cache都能被回收变成free 状态的内存，并列举了tmpfs, ramfs, memory segments 三种不可被操作系统当作缓存释放的内存空间。下边，跟随实验，一起见证一下，当前实验环境的操作系统版本是：Centos 7.4_1708

**tmpfs**

linux提供一种“临时”文件系统叫做tmpfs，它可以将内存的一部分空间拿来当做文件系统使用，使内存空间可以当做目录文件来用。现在绝大多数Linux系统都有一个叫做/dev/shm的tmpfs目录，就是这样一种存在。当然，我们也可以手工创建一个自己的tmpfs，方法如下：

![mount_tmpfs](https://rootdeep.github.io/Linux-Cache-Analysis/mount_tmpfs.png)

上图创建了一个大小为10G 的tmpfs。我们可以在/tmp/tmpfs中创建一个10G以内的文件，文件实际占用的空间是内存空间。因为实验环境中用到的内存总大小只有1G，所以实验中，我们用dd 命令创建了一个300M 大小的文件。

![mount_tmpfs](https://rootdeep.github.io/Linux-Cache-Analysis/write_tmpfs.png)

从上图中，可以看到 buff/cache 部分确实是增长了300M，并且，shared 的部分也增加了300M, 并且available 的值减少了300M。但该文件具体是用的是buff方式存储 ,还是cache方式存储， 用atop 工具进一步查看内存使用情况,命令：`atop -l`

![mount_tmpfs](https://rootdeep.github.io/Linux-Cache-Analysis/atop.png)

上图中可以看到buff部分的大小是0, cache 占用较多（cache 不一致时因为刚才用该环境做了其他操作，如安装atop 工具），说明刚才使用dd 命令生成的文件确实放在了内存里并且内核使用的是cache作为存储。那么这部分cache能够被手动释放吗？

![drop_cache](https://rootdeep.github.io/Linux-Cache-Analysis/drop_cache.png)

可以看到，buff/cache 部分并没有释放多少，shared 的大小也没有变化。**那么tmpfs占用的cache空间什么时候会被释放呢？是在其文件被删除的时候.**如果不删除文件，无论内存耗尽到什么程度，内核都不会自动帮你把tmpfs中的文件删除来释放cache空间。

![rm_tmpfs](https://rootdeep.github.io/Linux-Cache-Analysis/rm_testfile.png)

可以看到，删除testfile后，shared的部分减少了300M， buff/cache 部分接近300M。但是发现一个奇怪的问题，为什么available的值比free 还小一点呢？

以上是我们分析的第一种cache不能被回收的情况。还有其他情况，比如：

**共享内存**

共享内存是系统提供给我们的一种常用的进程间通信（IPC）方式，但是这种通信方式不能在shell中申请和使用，所以我们需要一个简单的测试程序，代码如下：

```
cat shm.c 

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/ipc.h>
#include <sys/shm.h>
#include <string.h>

#define MEMSIZE 300*1024*1024  #申请300M的内存空间

int
main()
{
    int shmid;
    char *ptr;
    pid_t pid;
    struct shmid_ds buf;
    int ret;

    shmid = shmget(IPC_PRIVATE, MEMSIZE, 0600);
    if (shmid<0) {
        perror("shmget()");
        exit(1);
    }

    ret = shmctl(shmid, IPC_STAT, &buf);
    if (ret < 0) {
        perror("shmctl()");
        exit(1);
    }

    printf("shmid: %d\n", shmid);
    printf("shmsize: %d\n", buf.shm_segsz);

    buf.shm_segsz *= 2;

    ret = shmctl(shmid, IPC_SET, &buf);
    if (ret < 0) {
        perror("shmctl()");
        exit(1);
    }

    ret = shmctl(shmid, IPC_SET, &buf);
    if (ret < 0) {
        perror("shmctl()");
        exit(1);
    }

    printf("shmid: %d\n", shmid);
    printf("shmsize: %d\n", buf.shm_segsz);


    pid = fork();
    if (pid<0) {
        perror("fork()");
        exit(1);
    }
    if (pid==0) {
        ptr = shmat(shmid, NULL, 0);
        if (ptr==(void*)-1) {
            perror("shmat()");
            exit(1);
        }
        bzero(ptr, MEMSIZE);
        strcpy(ptr, "Hello!");
        exit(0);
    } else {
        wait(NULL);
        ptr = shmat(shmid, NULL, 0);
        if (ptr==(void*)-1) {
            perror("shmat()");
            exit(1);
        }
        puts(ptr);
        exit(0);
    }
}
```

程序功能很简单，就是申请一段300M的共享内存，然后打开一个子进程对这段共享内存做一个初始化操作，父进程等子进程初始化完之后输出一下共享内存的内容，然后退出。但是退出之前并没有删除这段共享内存。我们来看看这个程序执行前后的内存使用：

![run_shm](https://rootdeep.github.io/Linux-Cache-Analysis/run_shm.png)

同样，可以看到程序执行后，shared ,buff/cache 部分增加了300M。那么这段cache能被回收么？继续测试：

![free_shm](https://rootdeep.github.io/Linux-Cache-Analysis/free_shm.png)

可以观察到，人工释放的效果不明显，buff/cache 几乎没变。所以这段共享内存即使没人使用，仍然会长期存放在cache中，直到其被删除。删除方法有两种，一种是程序中使用shmdt()，另一种是使用ipcrm命令。我们来删除试试：

![ipcrm_shm](https://rootdeep.github.io/Linux-Cache-Analysis/ipcrm_shm.png)

删除共享内存后，cache被正常释放了。这个行为与tmpfs的逻辑类似。内核底层在实现共享内存（shm）、消息队列（msg）和信号量数组（sem）这些POSIX:XSI的IPC机制的内存存储时，使用的都是tmpfs。这也是为什么共享内存的操作逻辑与tmpfs类似的原因。当然，一般情况下是shm占用的内存更多，所以我们在此重点强调共享内存的使用。说到共享内存，Linux还给我们提供了另外一种共享内存的方法，就是：

mmap**

mmap()是一个非常重要的系统调用，这仅从mmap本身的功能描述上是看不出来的。从字面上看，mmap就是将一个文件映射进进程的虚拟内存地址，之后就可以通过操作内存的方式对文件的内容进行操作。

实际上mmap 这个调用的用途是很广泛的,比如：

1. 当malloc申请内存时，小段内存内核使用sbrk处理，而大段内存就会使用mmap。
2. 当系统调用exec族函数执行时，因为其本质上是将一个可执行文件加载到内存执行，所以内核很自然的就可以使用mmap方式进行处理。

在此，我们在此仅仅考虑一种情况，就是使用mmap进行共享内存的申请时，会不会跟shmget()一样也使用cache？同样，我们也需要一个简单的测试程序：

```
cat mmap.c 

#include <stdlib.h>
#include <stdio.h>
#include <strings.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <fcntl.h>
#include <unistd.h>

#define MEMSIZE 300*1024*1024
#define MPFILE "./mmapfile"

int main()
{
    void *ptr;
    int fd;

    fd = open(MPFILE, O_RDWR);
    if (fd < 0) {
        perror("open()");
        exit(1);
    }

    ptr = mmap(NULL, MEMSIZE, PROT_READ|PROT_WRITE, MAP_SHARED|MAP_ANON, fd, 0);
    if (ptr == NULL) {
        perror("malloc()");
        exit(1);
    }

    printf("%p\n", ptr);
    bzero(ptr, MEMSIZE);

    sleep(100);

    munmap(ptr, MEMSIZE);
    close(fd);

    exit(1);
}
```

这次我们干脆不用什么父子进程的方式了，就一个进程，申请一段30M 的mmap共享内存，然后初始化这段空间之后等待100秒，再解除影射，所以我们需要在它sleep这100秒内检查我们的系统内存使用，看看它用的是什么空间？当然在这之前要先创建一个300G的文件./mmapfile。结果如下：

![dd_mmap](https://rootdeep.github.io/Linux-Cache-Analysis/dd_mmap.png)

然后执行测试程序

![run_test](https://rootdeep.github.io/Linux-Cache-Analysis/run_test.png)

我们可以看到，在程序执行期间，cached一直为543M(542M)，比之前涨了300M，并且此时这段cache仍然无法被回收。然后我们等待100秒之后程序结束。

![test_result](https://rootdeep.github.io/Linux-Cache-Analysis/test_result.png)

程序运行结束之时，使用unmap，释放cached占用的空间被释放。这样我们可以看到，使用mmap申请标志状态为MAP_SHARED的内存，内核也是使用的cache进行存储的。在进程对相关内存没有释放之前，这段cache也是不能被正常释放的。实际上，mmap的MAP_SHARED方式申请的内存，在内核中也是由tmpfs实现的。由此我们也可以推测，由于共享库的只读部分在内存中都是以mmap的MAP_SHARED方式进行管理，实际上它们也都是要占用cache且无法被释放的。

我们通过三个测试例子，发现Linux系统内存中的cache并不是在所有情况下都能被释放当做空闲空间用的。并且也也明确了，即使可以释放cache，也并不是对系统来说没有成本的。总结一下要点，我们应该记得这样几点：

- 当cache作为文件缓存被释放的时候会引发IO变高，这是cache加快文件访问速度所要付出的成本。
- tmpfs中存储的文件会占用cache空间，除非文件删除否则这个cache不会被自动释放。
- 使用shmget方式申请的共享内存会占用cache空间，除非共享内存被ipcrm或者shmdt，否则相关的cache空间都不会被自动释放。
- 使用mmap方法申请的MAP_SHARED标志的内存会占用cache空间，除非进程将这段内存munmap，否则相关的cache空间都不会被自动释放。
- 实际上shmget、mmap的共享内存，在内核层都是通过tmpfs实现的，tmpfs实现的存储用的都是cache。

当理解了这些的时候，希望大家对free命令的理解可以达到一个新的层次。我们应该明白，内存的使用并不是简单的概念，cache也并不是真的可以当成空闲空间用的。如果我们要真正深刻理解你的系统上的内存到底使用的是否合理，是需要理解清楚很多更细节知识，并且对相关业务的实现做更细节判断的。

------

# buffer/cache 的广义理解

上边说道在内存管理系统中对buffer 和cache的定义，功能。在更宽泛的意义上，个人理解，可以从两个角度区别buffer和cache：

1. 作为缓冲（buffer)，其内容有头有尾，类似于一个很深的队列，能容量一定的数据量。其作用就是解决生产者消费者两边处理速度的差异，缓冲内的内容将会在一定时间内迅速消耗掉的而不是长期呆在里面，而缓存(cache) 如果没有内存释放的需要，就会一直停留。
2. 缓存(cache) 缓存的一个重要的作用就是提升数据读写的效率（命中率）。把新数据写入cache，而不是buffer，可以让数据留在缓存，以待近期的读操作。
3. 从提高读写的命中率来说，缓存不仅仅有page cache ,还有 CPU的各级缓存，memcache, redis ，数据库等。
4. 那么，网卡用到的一段内存算cache，还是buffer呢？基于以上理解，我觉得是cache，这样就不用死记啦！