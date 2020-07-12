# 前置知识

物理内存也称为为主存, 为了更加有效地管理存储器并且少出错，现代系统提供了一种对主存的抽象概念，叫做**虚拟存储器（VM）** 大多数计算机用的主存都是动态随机访问内存(DRAM), 只有内核才可以直接访问物理内存, 

Linux 内核给每个进程 都提供了一个独立的虚拟地址空间 ,并且这个地址空间是连续的

# 名词

| word       | explain                                                      |
| ---------- | ------------------------------------------------------------ |
| MMU        | 内存管理单元(Memory Management Unit 把虚拟地址转换为物理地址的硬件设备 |
| buffers    | Memory used by kernel buffers (Buffers in /proc/meminfo) <br />Relatively temporary storage for raw disk blocks that shouldn't get tremendously large (20MB or so).对原始**磁盘**块的临时存储，也就是用来缓存磁盘的数据，通常不会特别大（20MB 左右）。这样，内核就可以把分散的写集中起来，统一优化磁盘的写入，比如可以把多次小的写合并成单次大的写等等 |
| cache      | Memory used by the page cache and slabs (Cached and SReclaimable in /proc/meminfo)      |
|Cacheed |In-memory cache for files read from the disk (the page cache).  Doesn't include SwapCached. **从磁盘读取文件**的页缓存，也就是用来缓存从文件读取的数据。  |
|SReclaimable|Part of Slab, that might be reclaimed, such as caches.  slab可以回收的部分 |
|SUnreclaim | Part of Slab, that cannot be reclaimed on memory pressure. slab不可以回收的部分 |
| rss/res    | Resident Set Size: number of pages the process has in real memory.  This is just the pages which count toward  text,  data,  or  stack  space.<br/>                        This does not include pages which have not been demand-loaded in, or which are swapped out. |
| kmem_cache | 高速缓存                                                     |
|            |                                                              |
|            |                                                              |
|            |                                                              |



### 虚拟地址

虚拟地址空间的内部 又被分为**内核空间和用户空间**两部分, 不同字长(也就是单个CPU指令可以处理数据的最大长度)的处理器. 地址空间的范围也不同. 比如最常见的32位和64位系统

<img src="https://static001.geekbang.org/resource/image/ed/7b/ed8824c7a2e4020e2fdd2a104c70ab7b.png" alt="img" style="zoom: 50%;" />

物理内存只有进程真正去访问虚拟地址，**发生缺页中断时，才分配实际的物理页面**，建立物理内存和虚拟内存的映射关系。

应用程序操作的是虚拟内存；而CPU处理器直接操作的是物理内存。

当应用程序访问虚拟地址，必须将虚拟地址转化为物理地址，处理器才能解析地址访问请求。

主存中的每个字节有两个地址：一个选自虚拟地址空间的虚拟地址，一个选自物理地址空间的物理地址。

### 内存映射

将**虚拟内存地址**映射到**物理内存地址**. 为了完成内存映射, 内核为每个进程都维护了一张页表, 记录虚拟地址与物理地址的映射关系

<img src="https://static001.geekbang.org/resource/image/fc/b6/fcfbe2f8eb7c6090d82bf93ecdc1f0b6.png" alt="img" style="zoom:50%;" />

页表实际上存储在CPU的MMU(Memory Management Unit 内存管理单元)中, 这样,正常情况下,处理器就可以直接通过硬件,找出要访问的内存

而当进程访问的虚拟地址在页表中查不到时，系统会产生一个**缺页异常**(Page Fault)，进入内核空间分配物理内存、更新进程页表，最后再返回用户空间，恢复进程的运行。

 TLB（Translation Lookaside Buffer，转译后备缓冲器）会影响 CPU 的内存访问性能，在这里其实就可以得到解释。TLB 其实就是 MMU 中页表的高速缓存。由于进程的虚拟地址空间是独立的，而 TLB 的访问速度又比 MMU 快得多，所以，通过减少进程的上下文切换，减少 TLB 的刷新次数，就可以提高 TLB 缓存的使用率，进而提高 CPU 的内存访问性能。不过要注意，MMU 并不以字节为单位来管理内存，而是规定了一个内存映射的最小单位，也就是页，通常是 4 KB 大小。这样，每一次内存映射，都需要关联 4 KB 或者 4KB 整数倍的内存空间。



页的大小只有 4 KB ，导致的另一个问题就是，整个页表会变得非常大。比方说，仅 32 位系统就需要 100 多万个页表项（4GB/4KB），才可以实现整个地址空间的映射。为了解决页表项过多的问题，Linux 提供了两种机制，也就是**多级页表和大页（HugePage）**。

- 多级页表就是把内存分成区块来管理，将原来的映射关系改成区块索引和区块内的偏移。由于虚拟内存空间通常只用了很少一部分，那么，多级页表就只保存这些使用中的区块，这样就可以大大地减少页表的项数。

  

  Linux 用的正是四级页表来管理内存页，如下图所示，虚拟地址被分为 5 个部分，前 4 个表项用于选择页，而最后一个索引表示页内偏移。

  <img src="https://static001.geekbang.org/resource/image/b5/25/b5c9179ac64eb5c7ca26448065728325.png" alt="img" style="zoom:50%;" />
  
- 大页，顾名思义，就是比普通页更大的内存块，常见的大小有 2MB 和 1GB。大页通常用在使用大量内存的进程上，比如 Oracle、DPDK 等。通过这些机制，在页表的映射下，进程就可以通过虚拟地址来访问物理内存了。那么具体到一个 Linux 进程中，这些内存又是怎么使用的呢

## 虚拟内存空间分布

<img src="https://static001.geekbang.org/resource/image/71/5d/71a754523386cc75f4456a5eabc93c5d.png" alt="32位虚拟内存空间分布" style="zoom:50%;" />

用户空间内存，从低到高分别是五种不同的内存段。

1. 只读段，包括代码和常量等。(const) 由于是只读的，不会再去分配新的内存，不会产生内存泄漏。

2. 数据段，包括全局变量和静态变量。(var) 这些变量在定义时就已经确定了大小，所以也不会产生内存泄漏

3. 堆，包括动态分配的内存，从低地址开始向上增长。

4. 文件映射段(/内存映射段)，包括动态库、共享内存等，从高地址开始向下增长。其中共享内存由程序动态分配和管理。所以，如果程序在分配后忘了回收，就会导致跟堆内存类似的泄漏问题。

5. 栈，包括局部变量和函数调用的上下文等。栈的大小是固定的，一般是 8 MB, 栈内存由系统自动分配和管理。一旦程序运行超出了这个局部变量的作用域，栈内存就会被系统自动回收，所以不会产生内存泄漏的问题。

在这五个内存段中，**堆和文件映射段的内存是动态分配的**。比如说，使用 C 标准库的 malloc() 或者 mmap() ，就可以分别在堆和文件映射段动态分配内存。
[Hack the Virtual Memory: malloc, the heap & the program break](https://blog.holbertonschool.com/hack-the-virtual-memory-malloc-the-heap-the-program-break/)

### 内存分配与回收

malloc() 是 C 标准库提供的内存分配函数，对应到系统调用上，有两种实现方式，即 brk() 和 mmap()。

- 对小块内存（小于 128K），C 标准库使用 brk() 来分配，也就是**通过移动堆顶的位置来分配内存**。这些内存释放后并不会立刻归还系统，而是被缓存起来，这样就可以重复使用。

- 而大块内存（大于 128K），则直接使用内存映射 mmap() 来分配，也就是在**文件映射段找一块空闲内存分配出去。**

这两种方式，自然各有优缺点。

**brk() 方式的缓存，可以减少缺页异常的发生，提高内存访问效率。不过，由于这些内存没有归还系统，在内存工作繁忙时，频繁的内存分配和释放会造成内存碎片。**

**而 mmap() 方式分配的内存，会在释放时直接归还系统，所以每次 mmap 都会发生缺页异常。在内存工作繁忙时，频繁的内存分配会导致大量的缺页异常，使内核的管理负担增大。这也是 malloc 只对大块内存使用 mmap  的原因。**

##### 缺页异常

**当这两种调用发生后，其实并没有真正分配内存。这些内存，都只在首次访问时才分配，也就是通过缺页异常进入内核中，再由内核来分配内存。**

缺页异常又分为下面两种场景。

- 可以直接从物理内存中分配时，被称为次缺页异常。
- 需要磁盘 I/O 介入（比如 Swap）时，被称为主缺页异常。

整体来说，Linux 使用伙伴系统(buddy)来管理内存分配。这些内存在 MMU 中以页为单位进行管理，并且会通过**相邻页的合并，减少内存碎片化**（比如 brk 方式造成的内存碎片）。

#### 比页更小的对象

###### 比如不到 1K 的时候，该怎么分配内存呢？

实际系统运行中，确实有大量比页还小的对象，如果为它们也分配单独的页，那就太浪费内存了。

在用户空间，malloc 通过 brk() 分配的内存，在释放时并不立即归还系统，而是**缓存起来重复利用**。

在内核空间，Linux 则通过 slab 分配器来管理小内存。你可以把 slab 看成构建在伙伴系统上的一个缓存，主要作用就是**分配并释放内核中的小对象**。

#### 回收

对内存来说，如果只分配而不释放，就会造成内存泄漏，甚至会耗尽系统内存。所以，在应用程序用完内存后，还需要调用 free() 或 unmap() ，来释放这些不用的内存。

系统也不会任由某个进程用完所有内存。在发现内存紧张时，系统就会通过一系列机制来回收内存，比如下面这三种方式：

- 回收缓存，比如使用 LRU（Least Recently Used）算法，回收最近使用最少的内存页面；
- 回收不常访问的内存，把不常用的内存通过交换分区直接写到磁盘中；
- 杀死进程，内存紧张时系统还会通过 OOM（Out of Memory），直接杀掉占用大量内存的进程。

前两种方式，**缓存回收和 Swap 换出**，实际上都是基于 LRU 算法，也就是优先回收不常访问的内存。

##### LRU

LRU 回收算法，实际上维护着 active 和 inactive 两个双向链表，其中：active 记录活跃的内存页；inactive 记录非活跃的内存页。越接近链表尾部，就表示内存页越不常访问。这样，在回收内存时，系统就可以根据活跃程度，优先回收不活跃的内存。

活跃和非活跃的内存页，按照类型的不同，又分别分为**文件页和匿名页**，对应着缓存回收和 Swap 回收。当然，你可以从 /proc/meminfo 中，查询它们的大小，

```
# grep表示只保留包含active的指标（忽略大小写）
# sort表示按照字母顺序排序
$ cat /proc/meminfo | grep -i active | sort
Active(anon):     167976 kB
Active(file):     971488 kB
Active:          1139464 kB
Inactive(anon):      720 kB
Inactive(file):  2109536 kB
Inactive:        2110256 kB
```

**LRU回收的是缓存，Swap换出的是不可回收的内存，比如进程的堆内存**

##### swap

第二种方式回收不常访问的内存时，会用到交换分区（以下简称 Swap）

Swap就是把一块磁盘空间或者一个本地文件（以下讲解以磁盘为例），当成内存来使用。它包括换出和换入两个过程。

- 换出，就是把进程暂时不用的内存数据存储到磁盘中，并释放这些数据占用的内存。
- 换入，则是在进程再次访问这些内存的时候，把它们从磁盘读到内存中来。

通常只在内存不足时，才会发生 Swap 交换。并且由于磁盘读写的速度远比内存慢，Swap 会导致严重的内存性能问题。

有新的大块内存分配请求，但是剩余内存不足。这个时候系统就需要回收一部分内存（比如前面提到的缓存），进而尽可能地满足新内存请求。这个过程通常被称为**直接内存回收**。除了直接内存回收，还有一个专门的内核线程用来定期回收内存，也就是 kswapd0。为了衡量内存的使用情况，kswapd0 定义了三个内存阈值（watermark，也称为水位），分别是**页最小阈值（pages_min）、页低阈值（pages_low）和页高阈值（pages_high）**。剩余内存，则使用 pages_free 表示。

![img](https://static001.geekbang.org/resource/image/c1/20/c1054f1e71037795c6f290e670b29120.png)

kswapd0 定期扫描内存的使用情况，并根据剩余内存落在这三个阈值的空间位置，进行内存的回收操作。

- 剩余内存小于页最小阈值，说明进程可用内存都耗尽了，只有内核才可以分配内存。

- 剩余内存落在页最小阈值和页低阈值中间，说明内存压力比较大，剩余内存不多了。这时 kswapd0 会执行内存回收，直到剩余内存大于高阈值为止。

- 剩余内存落在页低阈值和页高阈值中间，说明内存有一定压力，但还可以满足新内存请求。

- 剩余内存大于页高阈值，说明剩余内存比较多，没有内存压力。我们可以看到，一旦剩余内存小于页低阈值，就会触发内存的回收。

这个页低阈值，其实可以通过内核选项 /proc/sys/vm/min_free_kbytes 来间接设置。min_free_kbytes 设置了页最小阈值，而其他两个阈值，都是根据页最小阈值计算生成的，计算方法如下 ：

```
pages_low = pages_min*5/4
pages_high = pages_min*3/2
```

##### OOM

当系统发现内存不足以分配新的内存请求时，就会尝试直接内存回收。这种情况下，如果回收完文件页和匿名页后，内存够用了，当然皆大欢喜，把回收回来的内存分配给进程就可以了。但如果内存还是不足，OOM（Out of Memory） 就要登场了。

OOM其实是内核的一种保护机制。它监控进程的内存使用情况，并且使用 `oom_score` 为每个进程的内存使用情况进行评分：一个进程消耗的内存越大，`oom_score` 就越大；一个进程运行占用的 CPU 越多，`oom_score` 就越小。这样，进程的 `oom_score` 越大，代表消耗的内存越多，也就越容易被 OOM 杀死，从而可以更好保护系统。



当然，为了实际工作的需要，管理员可以通过 `/proc` 文件系统，手动设置进程的 `oom_adj` ，从而调整进程的 `oom_score``。oom_adj` 的范围是 [-17, 15]，数值越大，表示进程越容易被 OOM 杀死；数值越小，表示进程越不容易被 OOM 杀死，其中 -17 表示禁止 OOM。比如用下面的命令，你就可以把 sshd 进程的 `oom_adj `调小为 -16，这样， sshd 进程就不容易被 OOM 杀死。

`echo -16 > /proc/$(pidof sshd)/oom_adj`




# 总结

https://my.oschina.net/fileoptions/blog/968320

<img src="../image/memory.png" alt="img" style="zoom: 33%;" />

# 内存结构

## page

页(page)是内核的内存管理基本单位。（linux/mm_types.h）

```
struct page {
       page_flags_t flags;  页标志符
       atomic_t _count;    页引用计数
       atomic_t _mapcount;     页映射计数
       unsigned long private;    私有数据指针
       struct address_space *mapping;    该页所在地址空间描述结构指针，用于内容为文件的页帧
       pgoff_t index;               该页描述结构在地址空间radix树page_tree中的对象索引号即页号
       struct list_head lru;        最近最久未使用struct slab结构指针链表头变量
       void *virtual;               页虚拟地址
};
```

- flags：页标志包含是不是脏的，是否被锁定等等，每一位单独表示一种状态，可同时表示出32种不同状态，定义在<linux/page-flags.h>
- _count：计数值为-1表示未被使用。
- virtual：页在虚拟内存中的地址，对于不能永久映射到内核空间的内存(比如高端内存)，该值为NULL；需要事必须动态映射这些内存。

尽管处理器的最小可寻址单位通常为字或字节，但内存管理单元(MMU，把虚拟地址转换为物理地址的硬件设备)通常以页为单位处理。内核用struct page结构体表示每个物理页，struct page结构体占40个字节，假定系统物理页大小为4KB，对于4GB物理内存，1M个页面，故所有的页面page结构体共占有内存大小为40MB，相对系统4G，这个代价并不高。

##### 脏页

被应用程序修改过，并且暂时还没写入磁盘的数据

脏页，一般可以通过两种方式写入磁盘。

- 可以在应用程序中，通过系统调用 fsync  ，把脏页同步到磁盘中；
- 也可以交给系统，由内核线程 pdflush 负责这些脏页的刷新。

##### 匿名页Anonymous Page

应用程序动态分配的堆内存

- **Anonymous Paging**: Paging is the movement of pages: small units of memory (eg, 4 Kbytes). The term anonymous refers to it being working memory, and having no named location in a file system. Linux calls this type of paging "swapping" (which means something else on other OSes).

##### 共享内存

通过 tmpfs 实现的，所以它的大小也就是 tmpfs 使用的内存大小。tmpfs 其实也是一种特殊的缓存。

##### cache 

包括两部分，一部分是磁盘读取文件的页缓存，用来缓存从磁盘读取的数据，可以加快以后再次访问的速度。另一部分，则是 Slab 分配器中的可回收内存。

##### buffer

缓冲区是对原始磁盘块的临时存储，用来缓存将要写入磁盘的数据。这样，内核就可以把分散的写集中起来，统一优化磁盘写入。

###### buffer and cache 利用内存，充当起慢速磁盘与快速 CPU 之间的桥梁，可以加速 I/O 的访问速度。

磁盘是一个块设备，可以划分为不同的分区；在分区之上再创建文件系统，挂载到某个目录，之后才可以在这个目录中读写文件

buffer是用于不带文件系统的直接操作**块设备**(一般是磁盘)的缓存(裸IO)

> 裸IO 详见IO篇的图

cache是带文件系统的某个目录下的文件。

目的都是积累一定的量统一读写，减少IO次数。

在读写普通文件时，会经过文件系统，由文件系统负责与磁盘交互；而读写磁盘或者分区时，就会跳过文件系统，也就是所谓的“裸I/O“。这两种读写方式所使用的缓存是不同的，也就是文中所讲的 Cache 和 Buffer 区别。

第一个问题，Buffer 的文档没有提到这是磁盘读数据还是写数据的缓存，而在很多网络搜索的结果中都会提到 Buffer 只是对将要写入磁盘数据的缓存。那反过来说，它会不会也缓存从磁盘中读取的数据呢？

答: 会, 变化很小

第二个问题，文档中提到，Cache 是对从文件读取数据的缓存，那么它是不是也会缓存写文件的数据呢？

答:会, 变化很小

写文件时会用到 Cache 缓存数据，而写磁盘则会用到 Buffer 来缓存数据

Buffer 既可以用作“将要写入磁盘数据的缓存”，也可以用作“从磁盘读取数据的缓存”。Cache 既可以用作“从文件读取数据的页缓存”，也可以用作“写文件的页缓存”。

Buffer 是对磁盘数据的缓存，而 Cache 是文件数据的缓存，它们既会用在读请求中，也会用在写请求中。

**直接IO/`direct IO`是跳过Buffer，裸IO是跳过文件系统（还是有buffer的）**

## PSS RSS VSS

```txt
/proc/[pid]/smaps (since Linux 2.6.14)
              This file shows memory consumption for each of the process's mappings.  For each of mappings there is a series of lines such as the following:

                  00400000-0048a000 r-xp 00000000 fd:03 960637       /bin/bash
                  Size:                552 kB
                  Rss:                 460 kB
                  Pss:                 100 kB
                  Shared_Clean:        452 kB
                  Shared_Dirty:          0 kB
                  Private_Clean:         8 kB
                  Private_Dirty:         0 kB
                  Referenced:          460 kB
                  Anonymous:             0 kB
                  AnonHugePages:         0 kB
                  Swap:                  0 kB
                  KernelPageSize:        4 kB
                  MMUPageSize:           4 kB
                  KernelPageSize:        4 kB
                  MMUPageSize:           4 kB
                  Locked:                0 kB
                  VmFlags: rd ex mr mw me dw
```

 The following lines show the size of the mapping, the amount of the  mapping  that  is  currently resident  in  RAM  ("Rss"), 

the process's proportional share of this mapping ("Pss"), the number of clean and dirty shared pages  in  the mapping, and the number of clean and dirty private pages in the mapping. 私有内存+共享内存按比例属于自己计算的那一部分

 "Referenced" indicates the amount of memory currently  marked as  referenced  or accessed. 
 "Anonymous" shows the amount of memory that does not belong to any file. 
"Swap" shows how  much  would-be-anonymous memory is also used, but out on swap.

- **VSS**- Virtual Set Size 虚拟耗用内存（包含共享库占用的内存）
- **RSS**- Resident Set Size 实际使用物理内存（包含共享库占用的内存）
- **PSS**- Proportional Set Size 实际使用的物理内存（比例分配共享库占用的内存）
- **USS**- Unique Set Size 进程独自占用的物理内存（不包含共享库占用的内存）

**一般来说内存占用大小有如下规律：VSS >= RSS >= PSS >= USS**其中VSS>RSS>PSS>USS
 对于RSS PSS USS的统计，可以通过代码简单的
 `system/extras/libpagemap/pm_map.c`

```php
usage.rss += (count >= 1) ? map->proc->ker->pagesize : (0);
usage.pss += (count >= 1) ? (map->proc->ker->pagesize / count) : (0);
usage.uss += (count == 1) ? (map->proc->ker->pagesize) : (0);
```

通过上面的算法可以看出，USS部分是进程独占的部分，PSS是USS+共享/count， RSS是USS+共享部分
 如果我们要对进程退出后内存回收的统计的话，应该使用USS部分

**VSS** (reported as VSZ from ps) is the total accessible address space of a process.This size also includes memory that may not be resident in RAM like mallocs that have been allocated but not written to. VSS is of very little use for determing real memory usage of a process.

**RSS** is the total memory actually held in RAM for a process.RSS can be misleading, because it reports the total all of the shared libraries that the process uses, even though a shared library is only loaded into memory once regardless of how many processes use it. RSS is not an accurate representation of the memory usage for a single process.

**PSS** differs from RSS in that it reports the proportional size of its shared libraries, i.e. if three processes all use a shared library that has 30 pages, that library will only contribute 10 pages to the PSS that is reported for each of the three processes. PSS is a very useful number because when the PSS for all processes in the system are summed together, that is a good representation for the total memory usage in the system. When a process is killed, the shared libraries that contributed to its PSS will be proportionally distributed to the PSS totals for the remaining processes still using that library. In this way PSS can be slightly misleading, because when a process is killed, PSS does not accurately represent the memory returned to the overall system.

**USS** is the total private memory for a process, i.e. that memory that is completely unique to that process.USS is an extremely useful number because it indicates the true incremental cost of running a particular process. When a process is killed, the USS is the total memory that is actually returned to the system. USS is the best number to watch when initially suspicious of memory leaksin a process.

## slab层

slab分配器的作用：

- 对于频繁地分配和释放的数据结构，会缓存它；
- 频繁分配和回收比如导致内存碎片，为了避免，空闲链表的缓存会连续的存放，已释放的数据结构又会放回空闲链表，不会导致碎片；
- 让部分缓存专属单个处理器，分配和释放操作可以不加SMP锁；

slab层把不同的对象划分为高速缓存组，每个高速缓存组都存放不同类型的对象，每个对象类型对应一个高速缓存。kmalloc接口监理在slab层只是，使用一组通用高速缓存。

每个高速缓存都是用kmem_cache结构来表示

- kmem_cache_crreate：创建高速缓存
- kmem_cache_destroy: 撤销高速缓存
- kmem_cache_alloc: 从高速缓存中返回一个指向对象的指针
- kmem_cache_free：释放一个对象

实例分析： 内核初始化期间，/kernel/fork.c的fork_init()中会创建一个名叫task_struct的高速缓存； 每当进程调用fork（）时，会通过dup_task_struct()创建一个新的进程描述符，并调用do_fork(),完成从高速缓存中获取对象。

## NUMA 与 Swap

### Swap 升高 可是系统剩余内存还多着呢

这正是处理器的 NUMA （Non-Uniform Memory Access）架构导致的。

在 NUMA 架构下，多个处理器被划分到不同 Node 上，且每个 Node 都拥有自己的本地内存空间。而同一个 Node 内部的内存空间，实际上又可以进一步分为不同的内存域（Zone），比如直接内存访问区（DMA）、普通内存区（NORMAL）、伪内存区（MOVABLE）等，如下图所示：

![img](https://static001.geekbang.org/resource/image/be/d9/be6cabdecc2ec98893f67ebd5b9aead9.png)

先不用特别关注这些内存域的具体含义，我们只要会查看阈值的配置，以及缓存、匿名页的实际使用情况就够了。既然 NUMA 架构下的每个 Node 都有自己的本地内存空间，那么，在分析内存的使用时，我们也应该针对每个 Node 单独分析。你可以通过 numactl 命令，来查看处理器在 Node 的分布情况，以及每个 Node 的内存使用情况。比如，下面就是一个 numactl 输出的示例：

```shell
$ numactl --hardware
available: 1 nodes (0)
node 0 cpus: 0 1
node 0 size: 7977 MB
node 0 free: 4416 MB
...
```

这个界面显示，我的系统中只有一个 Node，也就是 Node 0 ，而且编号为 0 和 1 的两个 CPU， 都位于 Node 0 上。另外，Node 0 的内存大小为 7977 MB，剩余内存为 4416 MB。了解了 NUNA 的架构和 NUMA 内存的查看方法后，你可能就要问了这跟 Swap 有什么关系呢？实际上，前面提到的三个内存阈值（页最小阈值、页低阈值和页高阈值），都可以通过内存域在 proc 文件系统中的接口 /proc/zoneinfo 来查看。比如，下面就是一个 /proc/zoneinfo 文件的内容示例：

```shell
$ cat /proc/zoneinfo
...
Node 0, zone   Normal
 pages free     227894
       min      14896
       low      18620
       high     22344
...
     nr_free_pages 227894
     nr_zone_inactive_anon 11082
     nr_zone_active_anon 14024
     nr_zone_inactive_file 539024
     nr_zone_active_file 923986
...
```

- pages 处的 min、low、high，就是上面提到的三个内存阈值，

-  free 是剩余内存页数，它跟后面的 nr_free_pages 相同。

- nr_zone_active_anon 和 nr_zone_inactive_anon，分别是活跃和非活跃的匿名页数。

- nr_zone_active_file 和 nr_zone_inactive_file，分别是活跃和非活跃的文件页数。

从这个输出结果可以发现，剩余内存远大于页高阈值，所以此时的 kswapd0 不会回收内存。当然，某个 Node 内存不足时，系统可以从其他 Node 寻找空闲内存，也可以从本地内存中回收内存。具体选哪种模式，你可以通过 /proc/sys/vm/zone_reclaim_mode 来调整。它支持以下几个选项：

- 默认的 0 ，也就是刚刚提到的模式，表示既可以从其他 Node 寻找空闲内存，也可以从本地回收内存
- 1、2、4 都表示只回收本地内存，2 表示可以回写脏数据回收内存，4 表示可以用 Swap 方式回收内存



# 缓存命中率

直接通过缓存获取数据的请求次数，占所有数据请求次数的百分比。

# 如何查看内存

## top

- total 是总内存大小；

- used 是已使用内存的大小，包含了共享内存；

- free 是未使用内存的大小；

- shared 是共享内存的大小；

- buff/cache 是缓存和缓冲区的大小；

- available 是新进程可用内存的大小。

  > available 不仅包含未使用内存，还包括了可回收的缓存

- VIRT(virtual address space) 是进程虚拟内存的大小，只要是进程申请过的内存，即便还没有真正分配物理内存，也会计算在内。

- RES (Resident Memory Size)是常驻内存的大小，也就是进程实际使用的物理内存大小，但不包括 Swap 和共享内存。

  >  等于ps中的rss (Resident Set Size)

- SHR 是共享内存的大小，比如与其他进程共同使用的共享内存、加载的动态链接库以及程序的代码段等。

- %MEM (即RES的占比)是进程使用物理内存占系统总内存的百分比。

除了要认识这些基本信息，在查看 top 输出时，你还要注意两点。

第一，虚拟内存通常并不会全部分配物理内存。从上面的输出，你可以发现每个进程的虚拟内存都比常驻内存大得多。

第二，共享内存 SHR 并不一定是共享的，比方说，程序的代码段、非共享的动态链接库，也都算在 SHR 里。当然，SHR 也包括了进程间真正共享的内存。所以在计算多个进程的内存使用时，不要把所有进程的 SHR 直接相加得出结果。

## 清理缓存

```shell
# 清理page cache、目录项、Inodes等各种缓存
$ echo 3 > /proc/sys/vm/drop_caches
```

