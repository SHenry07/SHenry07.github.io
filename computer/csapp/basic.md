# Textbooks
- Computer Systems: A Programmer's Perspective, Third Edition (cs:app3e) Pearson,2016
- http://csapp.cs.cmu.edu/
- The C Programming Language, Second Edition, Prentice Hall, 1988

I think a  really good strategy for studying and preparing for this course would be to **read each chapter three times**

# Aim

- 了解编译系统

  - **优化程序性能**

    比如：

    1. 一个switch语句是否总是比一系列的if-else语句高效的多？
    2. 一个函数调用的开销有多大？
    3. while循环比for循环更有效吗?
    4. 指针引用比数组索引更有效吗？
    5. 为什么将循环求和的结果放到一个本地变量中，会比将其放到一个通过引用索引传递过来的参数中，运行起来快很多呢？
    6. 为什么我们只是简单的重新排列一下算术表达式的括号就能让函数运行的更快？

  - 理解链接时出现的错误

    1. 连接器报告说它无法解析一个引用，这是什么意思
    2. 静态变量和全局变量的区别是什么？
    3. 如果你在不同的C文件中定义了名字相同的两个全局变量会发生什么？
    4. **静态库和动态 库的区别是什么？**
    5. 我们在命令行上排列库的顺序有什么影响？
    6. 为什么有些链接错误直到运行时才会出现

  - 避免安全漏洞。 常见的漏洞攻击: 缓冲区溢出攻击


# 第一周(chapter I)

**计算需要先有菜谱，然后有执行步骤，最后有数据结构与算法，最后才是实现**

信息就是位+上下文

源程序就是一个由值0和1组成的位(又成bit比特)，8位被组织成一组，称为字节byte

> 一定要理解bit和byte的关系，最高有效**byte**字节包含bit[x<sub>w-1</sub>,x<sub>w-2</sub>x<sub>w-3</sub> ….. x<sub>w-8</sub>], 最低有效字节包含位[x<sub>7</sub>, x<sub>6</sub>, …., x<sub>0</sub>]

上下文context: 操作系统保持追踪进程运行所需的所有状态信息. 包括许多信息,比如PC和寄存器文件的当前值,以及主存的内容.

在任何一个时刻,单处理器系统都只能执行一个进程的代码,当操作系统决定要把控制权从当前进程转移到某个新进程时,就会进行**上下文切换context switch** 即保存当前进程的上下文,恢复新进程的上下文,然后将控制权传递到新进程.

示例场景中有两个并发的进程,shell进程和hello进程,最开始只有shell进程在运行,即等待命令行上的输入. 当我们让它运行hello程序时,shell通过调用一个专门的函数,即**系统调用**,来执行我们的请求,系统调用会将控制权传递给操作系统, 操作系统保存shell进程的上下文,创建一个新的hell进程及其上下文,然后将控制权传给新的hello进程, hello进程终止后,操作系统恢复shell进程的上下文, 并将控制权传回给它.<sup>P12</sup>

从一个进程到另一个进程的转换是**由操作系统内核(kernel)管理的** 内核是操作系统代码常驻内存的部分. 当应用程序需要操作系统的某些操作时, 比如读写文件, 它就执行一条特殊的system call 系统调用指令,将控制权传递给内核, 然后内核执行被请求的操作并返回应用程序, 注意,内核不是一个独立的进程,相反,他是系统管理全部进程所用代码和数据结构的集合

## 进制转换

字母'a' ~ 'z' 的ASCII码为 0x61 ~ 0x7Az

## 整数和浮点数

除了最高有效位1， 整数的所有未都嵌在浮点数中。

例子:

12345的整型为0x00003039, 浮点数0x4640E400

3510593为0x00359141, 3510593.0为0x4A564504

​              0011 0101 1001 0001 0100 0001
​             001 101 0110 0100 0101 0000 01
0100 1001 0101 0110 0100 0101 0000 0100

## 布尔代数 Boolean algebra

|signal | logical operations  | |
|--|--|--|
|~ tilde | Complement 补集 取反| A&B = 1 when both A=1 and B=1 |
|& ampersand | Intersection 与 相同为1/0 不同为0|  A|B = 1 when either A = 1 or B = 1 |
|\| | Union OR 或 相同为0/1 不同为1| ~A =1 when A = 0 |
|^ | symmetric difference  EXCLUSIVE-OR 异或 相同为0， 不同为1| A^B = 1 when either A = 1 Or B = 1, but not A = B |


十进制计算的时候: ~y = -y - 1

symmetric difference  对称差:两个集合的对称差是只属于其中一个集合，而不属于另一个集合的元素组成的集合 

**当p=1 and q = 0 Or p =1 and q=1 是 P^Q =1, 即 什么数和1异或都等于1**

布尔代数和整数运算有很多相似之处，具有分配律，集合律等

整数运算的一个属性是每个值x都有一个加法逆元(additive inverse)-x,  使得x+(-x)=0 

布尔代数这里的"加法"运算是^ , 也就是**对于任何值a来说，这里a^a=0**

应用:[除自己之外数组的乘积](https://leetcode-cn.com/problems/product-of-array-except-self/)

## 位级运算和逻辑运算，移位运算

首先要分清它们之间的优先级(按照从左至右结合性规则依次递减，!,移位>><<, &,^,|,&&,||)，如果不清楚，请**加上括号**

逻辑运算:

! 经常读作bang, 表示非

> 认为所有非0的参数都表示**True** 返回1/0x01, 参数0表示**flase** 返回0/0x00, <sup>P37-39</sub> 

而Golang中是不支持integer 与boolean 互相转换的

> Early Termination: 如果对第一个参数求值就能确定表达式的结果, 逻辑运算符就不会对第二个参数求值

应用:`a&&5/a`将不会造成被零除, p&&*p++也不会导致间接引用空指针

shift Operations:
> 从左至右结合
> 逻辑右移logical shift 和算术右移 Arithmetic shift 的区别: 算术右移是在左端补k个最高有效位的值,即是0就补0,1就补1

# 第二周

## C语言中的有符号数与无符号数

**TMAX<sub>w</sub> = 2<sup>w-1</sup> - 1** 

**TMIN<sub>w</sub> = -2<sup>w-1</sup>** 

> 在本节中，我们描述用位来编码整数的两种不同的方式: 一种只能表示非负数(无符号数)，一种能够表示负数、零和正数.<sup>P41</sup> 
>
> C 语言标准中没有指定**有符号数**要采用某种表示，但是几乎所有的机器都**使用补码 two's Complement**

有符号数是用补码来表示的

> ```c
> short sx = -12345;
> unsigned uy = sx;
> ```
>
> X当把short 4byte 16bit 转换成 unsigned 8bytes 32bit时， 我们要先改变大小，之后再完成从有符号到无符号的转换，也就是说(unsigned) sx 等价于 (unsigned)(int)sx。

改变大小，就要扩展这个数字sx的位表示，从16位cf f7到32位ff ff cf f7，这个规则是C语言标准要求的

> x=9 和y=12的位表示分别为[1001]和[1100], 它们的和是21,5位的表示为[10101], 丢弃最高位是[0101] 也就是十进制的5,这和值21 mod 16 = 5是一致的

截断一个数字可能会改变它的值–>溢出的一种形式(a form of overflow)

> 通过观察发现x+y>=x, 因此如果s没有溢出,我们能够肯定s>=x. 另一方面,如果s确实溢出了,我们就有s=x+y-2<sup>w</sup>. 假设y<2<sup>w</sup> < 0, 因此s=x+(y-2<sup>w</sup>) < x



### 补码的非

![image-20200924173422629](./image-20200924173422629.png)

**当一个无符号数y等于TMIN时， -y也等于TMIN**， 因为转换成了TMATX+1，又溢出成了TMIN

#### 补码非的位级表示

- 在C语言中，**对于任意整数值x,计算表达式`-x` 和`~x+1`得到的结果完全一样**

  [1000] -8 ~x [0111] 7  incr(-x) -8

  ```c
  int x;
  int y;
  unsigned ux = x;
  unsigned uy = y;
  ```
  判读真假: `x*~y + uy*ux = == -x` 

- 第二种方法: 建立在将位向量分为两部分的基础上， 假设k是最右边的1的位置，(只要x != 0 就能找到这样的k)。我们对位k左边的所有位取反

### 如何判断溢出

无符号数加法溢出: 在范围0<=x, y <= UMax<sub>w</sub> 中的x和y, 当且仅当s<x(等价地s<y)时,发生了溢出

```c
// unsigned scomplement
int uadd_ok(unsigned x, unsigned y) {
    unsigned sum = x + y;
    printf("%u %u %u\n", sum, x, y);
    return sum >= x;
}
```

有符号数加法溢出:  对满足Tmin<sub>w</sub>	<= x, y <= Tmax<sub>w</sub> 的x 和 y ， 令 s = x + $ \mathop{{}}\nolimits^{{t}}\nolimits_{{w}}y$. 当且仅当 x > 0, y > 0, 但 s <=0 时，计算s发生了正溢出，当且仅当x<0, y<0， 但 s >=0 时 s发生了负溢出

```c
// two's complement
int tadd_ok(int x, int y) {
    int sum = x + y;
    // 正溢出
    int pos_over = x >=0 && y >= 0 && sum < 0;
    // 负溢出 true 返回 1 false 返回 0
    int neg_over = x < 0 && y < 0 && sum >= 0 ;
    //      !1 && !1
    return !pos_over && !neg_over ;
}
```

有符号减法溢出: 

乘法溢出:

```c
int tmult_ok(int x, int y) {
   int p = x*y;
   /* Either x is zero, or dividing p by x gives y */
    return !x || p/x == y;
}
```

乘法溢出: 

```c
// c/c++ 其他高级语言基本不允许 不同type的value 比大小
int tmult_ok(int x, int y) {
    int64_t pll = (int64_t) x * y;

    return pll == (int) pll;
}
```
> c 中 利用 不同类型进行比较的例子
>
> ```c
> #define int size_t;
> int ele_cnt;
> 
> uint64_t required_size = ele_cnt * (uint64_t) ele_size;
> size_t request_size = (size_t) required_size;
> if (required_size != request_size )
>     /* Overflow must have occurred. Abort operation */
>     return NULL;
> void *result = malloc(request_size);
> if (result == NULL) 
>     /* malloc failed */
>     return NULL;
> ```
>
> 

### 运算操作

#### 加法
有符号数在补码表示下, 无论加法是否溢出 (x+y)-y总是会求值得到x. 这是因为**补码加**会形成一个阿贝尔群

> 在大多数机器上，整数乘法指令相当慢，需要10个或者更多的时钟周期，而其他整数运算(例如 加法、减法、危机运算、和移位)只需要1个时钟周期， 因此 **编译器**试着用移位和加法运算的组合来代替乘以常数因子的乘法


#### 乘法

先说结论: 有符号数和无符号数**截断后的乘积的bit表示是相同的**
```
        1110 -2 14  
        1101 -3 13
             6   182 
             0x6   0xb6
   1011 0110 
    // w =4 时 bit是一样的
```

### 关于整数运算的最后思考

计算机执行的"整数"运算 实际上是一种模运算形式

## 浮点数

> **注意,二进制下 形如0.111...1<sub>2</sub>, 表示刚好小于1的数**
>
> 二进制小数的表示的一个简单办法是: **讲一个数转化为形如$\frac{x}{2^k}$的小数**. 假分数的整数部分按2的幂从小到大表示,分数部分按2的负幂从大到小表示. 即`8 4 2 1 . 1/2 1/4 1/8`

### IEEE

$V=(-1)^s*M*x^E$

| 缩写 | 作用                          |
| ---- | ----------------------------- |
| s    | sign bit 符号位               |
| e    | exponent 阶码                 |
| f    | fraction/frac 分数            |
| M    | implicit leading 1, 隐式的1+f |
| Bias | 偏移量 $2^{{k=1}}-1$          |

float: s=1 k=8 f=23  Bias=127

double: s=1 k=11 f=52 Bias=1023

##### Normalized Form

Let's illustrate with an example, suppose that the 32-bit pattern is `1 1000 0001 011 0000 0000 0000 0000 0000`, with:

- `S = 1`
- `E = 1000 0001`
- `F = 011 0000 0000 0000 0000 0000`

In the *normalized form*, the actual fraction is normalized with an implicit leading 1 in the form of `1.F`. In this example, the actual fraction is `1.011 0000 0000 0000 0000 0000 = 1 + 1×2^-2 + 1×2^-3 = 1.375D`.

The sign bit represents the sign of the number, with `S=0` for positive and `S=1` for negative number. In this example with `S=1`, this is a negative number, i.e., `-1.375D`.

In normalized form, the actual exponent is `E-127` (so-called excess"-127" or bias"-127"). This is because we need to represent both positive and negative exponent. With an 8-bit E, ranging from 0 to 255, the excess-127 scheme could provide actual exponent of -127 to 128. In this example, `E-127=129-127=2D`.

Hence, the number represented is `-1.375×2^2=-5.5D`.

## 第二章 Representing and Manipulating Information 总结

unsigned int 用 bit 直接表述

signed int 计算机用two's complement 实现 

溢出这种问题要注意

shift Operations: 在计算字节大小， 传输数据的大小， 硬件设备的大小时相当有用

kib >> 10 >> 10

Kib    Mib   Gib

浮点数要考虑的情况较多， 这里只记录常规形式 理解在计算机里是如何表达的即可