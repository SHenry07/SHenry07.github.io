- [关键字](#关键字)
- [工作空间](#工作空间)
- [环境变量](#GoPATH环境变量)
- [基本语法](#基本语法)
- [指针](#指针)
- [经典应用](#经典应用)
- [Error](#Error)

# 关键字

able 1.2. Go 中的关键字
|   |             |           |         |        |
| ---- | ---- | ---- | ---- | ---- |
| import | package | func | const | var |
|defer|fallthrough <br />在switch中即使匹配成功还可以继续|interface|chan|go|
|else|map|for|range|return|
|if| continue    |break|default|select|
|switch|case| type      |struct|goto|


Table 1.3. Go 中的预定义函数
|  |  |  |  |
|--|--|--|--|
|close|delete|len|cap|
|new|make|append|copy|
|recover|print|println|complex|
|real|imag|||

# 工作空间

Go代码必须放在**工作空间**内。它其实就是一个目录，其中包含三个子目录：

- src 目录包含Go的源文件，它们被组织成包（每个目录都对应一个包），
- pkg 目录包含包对象，
- bin 目录包含可执行命令。

go 工具用于构建源码包，并将其生成的二进制文件安装到 pkg 和 bin 目录中。

src 子目录通常包会含多种版本控制的代码仓库（例如Git或Mercurial）， 以此来跟踪
一个或多个源码包的开发。

以下例子展现了实践中工作空间的概念：

```shell
bin/
	streak                         # 可执行命令
	todo                           # 可执行命令
pkg/
	linux_amd64/
		code.google.com/p/goauth2/
			oauth.a                # 包对象
		github.com/nf/todo/
			task.a                 # 包对象
src/
	code.google.com/p/goauth2/
		.hg/                       # mercurial 代码库元数据
		oauth/
			oauth.go               # 包源码
			oauth_test.go          # 测试源码
	github.com/nf/
		streak/
		.git/                      # git 代码库元数据
			oauth.go               # 命令源码
			streak.go              # 命令源码
		todo/
		.git/                      # git 代码库元数据
			task/
				task.go            # 包源码
			todo.go                # 命令源码
```

此工作空间包含三个代码库（goauth2、streak 和 todo），两个命令（streak 和 todo） 以及两个库
（oauth 和 task）。

##  GOPATH环境变量

GOPATH 环境变量指定了你的工作空间位置。它或许是你在开发Go代码时， 唯一需要设置的环境变量。

首先创建一个工作空间目录，并设置相应的 GOPATH。你的工作空间可以放在任何地方， 在此文档中我们使用 
$HOME/work。注意，它绝对不能和你的Go安装目录相同。 （另一种常见的设置是 GOPATH=$HOME。）

```
$ mkdir $HOME/work
$ export GOPATH=$HOME/work
```
作为约定，请将此工作空间的 bin 子目录添加到你的 PATH 中：

`$ export PATH=$PATH:$GOPATH/bin`

**包路径**

标准库中的包有给定的短路径，比如 "fmt" 和 "net/http"。 
对于你自己的包，你必须选择一个基本路径，来保证它不会与将来添加到标准库， 或其它扩展库中的包相冲突。
如果你将你的代码放到了某处的源码库，那就应当使用该源码库的根目录作为你的基本路径。
例如，若你在 GitHub 上有账户 github.com/user 那么它就应该是你的基本路径。

比如孙恒自己的格式
`mkdir -p $GOPATH/src/github.com/HengHughSun`

之后在自己的工作空间中建立自己的相应的包目录  类似danjgo的APP

[包导入的过程](package.md)

# 基本语法

### **print家族**

print将它的参数显示在命令窗口，并将输出光标定位在所显示的最后一个字符之后。

println 将它的参数显示在命令窗口，并在结尾加上换行符，将输出光标定位在下一行的开始。

printf是格式化输出的形式。("%T", d)

### 标签与 goto

```go
package main

import "fmt"

func main() {

LABEL1:
    for i := 0; i <= 5; i++ {
        for j := 0; j <= 5; j++ {
            if j == 4 {
                continue LABEL1
            }
            fmt.Printf("i is: %d, and j is: %d\n", i, j)
        }
    }

}
```

本例中，continue 语句指向 `LABEL1`，当执行到该语句的时候，就会跳转到 `LABEL1 `标签的位置。

您可以看到当 j==4 和 j==5 的时候，没有任何输出：标签的作用对象为外部循环，因此 i 会直接变成下一个循环的值，而此时 j 的值就被重设为 0，即它的初始值。如果将 continue 改为 break，则不会只退出内层循环，而是直接退出外层循环了。另外，还可以使用 goto 语句和标签配合使用来模拟循环。

Go 有 goto 语句——明智的使用它。用 goto 跳转到一定是当前函数内定义的标签。例如假设这样一个循环:
```go
func myfunc() {
	i := 0
Here:  # 这行的第一个词,以冒号结束作为标签
	println(i)
	i++
	goto Here  #跳转
}
```
（**标签名是大小写敏感的**，为了提升可读性，一般建议使用全部大写字母）

<u>如果您必须使用 goto，应当只使用正序的标签（标签位于 goto 语句之后），但注意标签和 goto 语句之间不能出现定义新变量的语句，否则会导致编译失败。</u>

### **循环**
Go 只有一种循环结构：for 循环。

基本的 for 循环由三部分组成，它们用分号隔开：

初始化语句：在第一次迭代前执行
条件表达式：在每次迭代前求值
后置语句：在每次迭代的结尾执行

`for i := 0; i < 10; i++ { }`

注意：和 C、Java、JavaScript 之类的语言不同，Go 的 for 语句后面的三个构成部分外没有小括号，
大括号 { } 则是必须的。且{}必须和for if 在同一行
**初始化语句和后置语句是可选的**

for 是 Go 中的 “while”

`for ; sum < 1000;` 与 `for sum < 1000 `

前者C中的for 后者 c中的while 只有‘;’的不同

此时你可以去掉分号，因为 C 的 while 在 Go 中叫做 for。

if 与 for 相同 无需() 而{}则是必须的

```go
if {

}else if{

}else{

}
```

### break

循环嵌套循环时,可以在 break 后指定标签。用标签决定 哪个 循环被终止:

```go
J: for j := 0; j < 5; j++ {
	for i := 0; i < 10; i++ {
		if i > 5 {
			break J   #← 现在终止的是 j 循环,而不是 i 的那个
		}
		println(i)
	}
}
```



`switch case default` 可以将长的if-then-else写的更清楚

```go
switch condition; var{   // confition var 都可省 可以定义在全局
case ' ', '?', '&', '=', '#', '+':  //  用,表示 ”or”
case "":
default:
}
```

### **延时调用函数(Deferred Function Calls)**

defer 语句会将函数延迟到外层函数返回后执行,延迟的函数是按照**后进先出 (LIFO) 的顺序**执行 如果有多个defer调用 **是先进后出的顺序, 
类似于入栈出栈一样:**

### **基本类型type**

Go 的基本类型有

```go
bool

string

//带符号:  首位0 表示- 1表示+ 
int  int8  int16  int32  int64
//无符号:
uint uint8 uint16 uint32 uint64 uintptr

//int是带符号的，表示范围是：-2147483648到2147483648，即-2^31到2^31次方。

//uint则是不带符号的，表示范围是：2^32即0到4294967295。

byte // uint8 的别名

rune // int32 的别名
    // 表示一个 Unicode 码点
uintptr	// 无符号整型，用于存放一个指针

float32 float64

complex64 complex128
```

**声明变量**
一般形式是使用var关键字 `var identifier type` 也可以一次声明多个变量

1. **这种因式分解关键字的写法一般用于声明全局变量**
   
```go
   var (
    vname1 v_type1
    vname2 v_type2
   )
```

   声明变量同时赋值`var  s1 string = “who”`

2. 类型推断(根据值推断变量类型)

   `var identifier = “test”`

3. **声明类型**
   type 类型名 真实的type   // 自定义类型

   > `type yourInt  = int` 类型别名 比如 rune byte 编译完还是 int32 和 uint8
   
 4. **短变量声明**

在函数中，简洁赋值语句 := 可在类型明确的地方代替 var 声明。

注意说明:
- vname1, vname2, vname3 := v1, v2, v3 // 出现在 := 左侧的变量不应该是已经被声明过的，否则会导致编译错(即同一作用域中变量不能二次声明)

- **函数外的每个语句都必须以关键字开始（var, func 等等），因此 := 结构不能在函数外使用。** `:=`只用来声明临时变量, 初始化全局变量需使用`var`关键字 

- 局部变量声明不使用是无法编译的,全局变量可以通过

- **定义后赋值和变量初始化的不同**

  ```go
  var n int//
  n = 9 // 先定义，然后进行赋值操作， 在函数体外不被允许； 而`var n = 9`是变量初始化操作，在函数体外是允许的。变量初始化和变量赋值是两个不同的概念
  ```

#### 常量

常量是一个简单值的标识符，在程序运行时，不会被修改的量。

常量的声明与变量类似，只不过是使用 const 关键字。

```go
const identifier [type] = value
```

你可以省略类型说明符 [type]，因为编译器可以根据变量的值来推断其类型。

- 显式类型定义： `const b string = "abc"`
- 隐式类型定义： `const b = "abc"`

多个相同类型的声明可以简写为：

```go
const c_name1, c_name2 = value1, value2
```

常量可以是字符、字符串、布尔值或数值。

常量**不能用 :=** 语法声明。

**数值常量**

数值常量是高精度的值。

一个未指定类型的常量由上下文来决定其类型。

再尝试一下输出 needInt(Big) 吧。

（int 类型最大可以存储一个 64 位的整数，有时会更小。）

（int 可以存放最大64位的整数，根据平台不同有时会更少。）

**赋值可以用八进制、十六进制或科学计数法: 077 , 0xFF , 1e3 或者 6.022e23 这些都
是合法的。**

### 枚举
iota，特殊常量，可以认为是一个可以被编译器修改的常量。

iota 在 **const关键字出现时将被重置为 0(const 内部的第一行之前)**，const 中每**新增一行常量声明**将使 iota 计数一次(iota 可理解为 const 语句块中的行索引)。 可以使用 iota 生成枚举值
`const (
a = iota //0
b = iota //1
)`

省略 Go 重复的 = iota

`const (
a = iota
b             # ← Implicitly b
)`
如果需要,可以明确指定常量的类型：
`const (
a = 0        ← Is an int now
b string = "0"
)`

### 字符串

Go语言中字符串是用双引号""包裹的

字符使用单引号''包裹的

```go
// 字符串
s := "Hello"
// 单独的字母、汉字、符号表示一个字符
c1 := 'h'
c2 := '1'
c3 := '中'
// 1字节(Byte） = 8位(bit)
// 一个字符'A' = 1个字节
// 1个utf8编码的汉字 = 一般占3个字节
```

#### 多行字符串

Go语言中要定义一个多行字符串时，就必须使用`字符：

```go
s1 := `第一行
第二行
第三行
`
fmt.Println(s1)
```

反引号间换行   将被作为字符串中的换行，所有的转义字符均无效，文本将会原样输出。

#### 字符串的常用操作

|                方法                 |             介绍              |
| :---------------------------------: | :---------------------------: |
|              len(str)               |            求长度             |
|           +或fmt.Sprintf            |          拼接字符串           |
|            strings.Split            |             分割              |
|          strings.contains           |         判断是否包含          |
| strings.HasPrefix/strings.HasSuffix |         前缀/后缀判断         |
| strings.Index()/strings.LastIndex() | 子串出现的位置/最后出现的位置 |
| strings.Join(a[]string, sep string) |  join操作(连接操作类似拼接)   |

## byte和rune类型

组成每个字符串的元素叫做“字符”，可以通过遍历或者单个获取字符串元素获得字符。 字符用单引号（’）包裹起来，如：

```go
var a := '中'
var b := 'x'
```

Go 语言的字符有以下两种：

1. `uint8`类型，或者叫 byte 型，代表了`ASCII码`的一个字符。
2. `rune`类型，代表一个 `UTF-8字符`。

当需要处理中文、日文或者其他复合字符时，则需要用到`rune`类型。`rune`类型实际是一个`int32`。

Go 使用了特殊的 rune 类型来处理 Unicode，让基于 Unicode 的文本处理更为方便，也可以使用 byte 型进行默认字符串处理，性能和扩展性都有照顾。

```go
// 遍历字符串
func traversalString() {
	s := "hello沙河"
	for i := 0; i < len(s); i++ { //byte
		fmt.Printf("%v(%c) ", s[i], s[i])
	}
	fmt.Println()
	for _, r := range s { //rune
		fmt.Printf("%v(%c) ", r, r)
	}
	fmt.Println()
}
```

输出：

```bash
104(h) 101(e) 108(l) 108(l) 111(o) 230(æ) 178(²) 153() 230(æ) 178(²) 179(³) 
104(h) 101(e) 108(l) 108(l) 111(o) 27801(沙) 27827(河) 
```

因为UTF8编码下一个中文汉字由3~4个字节组成，所以我们不能简单的按照字节去遍历一个包含中文的字符串，否则就会出现上面输出中第一行的结果。

字符串底层是一个byte数组，所以可以和`[]byte`类型相互转换。字符串是不能修改的 字符串是由byte字节组成，所以字符串的长度是byte字节的长度。 rune类型用来表示utf8字符，一个rune字符由一个或多个byte组成。

## 运算符

### 算术运算符

a ++

b -- 是一个单独的语句,不能像python一样放在=的右边进行赋值

### 关系运算符

### 位目运算符

&   与  00 0  01  0  11 1  （两位均为1才为1）

|  或     00 0  01 1   11 1    （两位有一个为1就为1）  

^ 异或00 0  01  1   11 0  （两位不一样则为1）

<< 二进制数右移X位   相当于乘以2的X倍
\>>二进制数左移X位   相当于除以2的X倍

### 逻辑运算符

&& and   || or   ! NOT

### 赋值运算符

------

### 类型别名得到的新类型并非和原类型完全相同，新类型不会拥有原类型所附带的方法

# 内存操作

## 指针

指针保存了值的内存地址

type *T 是一个指向T类型值的指针,其零值为nil.

`var p *int` 现在p是一个指向整数值的指针

让指针指向某些内容,可以使用取址操作符 ( & )

& 操作符会生成一个指向其操作数的指针 **作用:生成指针指向操作数相当于拿到地址值(一个路牌)**

```go
i := 42
p = &i 此时p是一个内存地址形如0xc0000010o78
变量p此时是一个*int的类型
```

\* 操作符表示指针指向的 **数据值  可以拿到数据**

```go
fmt.Prrintln(*p)  //通过指针p 读取 i
*p = 21           //通过指针p 设置 i
```

这就是常说的"间接引用"or"重定向"  与c不同,Go **没有指针运算**

如果这样写: `*p++` ,它表示 `(*p)++` :首先获取指针指向的值,然后对这个值**加一**

## 内存分配

go有两个内存分配原语 new 和 make

new 分配;make 初始化
可以简单总结为:

- new(T) 返回 *T 指向一个零值 T

- make(T) 返回初始化后的 T
  当然 make 仅适用于 slice,map 和 channel。
  
  另外一种总结方式：
  
  1. make 和new都是用来申请内存的
  2. new很少用，一般用来给基本数据类型申请内存,string、int返回的是对应类型的指针(\*string,\*int)
  3. make 是用来给slice，map和channel申请内存的,make函数返回的是对应的这三个类型本身


- var p1 Person 分配了 Person - 值 给 p1 。 p1 的类型是 Person 。

 - p2 := new(Person) 分配了内存并且将 指针 赋值给 p2 。 p2 的类型是*Person 。

   ```go
   // 下面两个内存分配的区别是什么?
   func Set(t *T) {
   x = t
   }
   // 和
   func Set(t T) {
   x= &t
   }
   // 上面x得到的是值  下面x得到的是指针
   ```

# 进阶结构

## **Array 数组**

[n]T 表示拥有n个T类型的值的数组  var a [10]int    a为有10个整数的数组

> **数组的长度是其 '类型'的一部分**，因此数组不能改变大小
>
> 数组是值类型,赋值和传参会复制整个数组。因此改变副本的值，不会改变本身的值。

初始化方式:

```go
 var a [5]int
    fmt.Println("emp:", a)
    a[4] = 100
    fmt.Println("set:", a)
var XX = [3]bool{true,true,true}
var XX = [...]int{0,23,2315,666,7777}
var XX = [5]int{0 : 1, 4 : 2}   //--> [1,0,0,0,2]
b := [5]int{1, 2, 3, 4, 5}
```

多维数组
```go
var XX [3][2]int
XX = {
   <[2]int可省>{1,2}，
   {3,4},
   {5,6},
}
```
注意： 多维数组只有第一层(最外围数组)可以使用`...`来让编译器推导数组长度。


## **切片slices**

[]T 表示一个元素类型为T的切片。

切片的本质就是对底层数组的封装,包含了**底层数组的指针，长度，容量**

切片通过两个下标来界定，即一个上届和一个下届，二者以:为分隔符a[low : hight]

### 左闭右开 

类似一个半开区间[0,n) 包含第一个元素，不包含最后一个元素a[1:4] 表示1~3

**切片就像数组的引用 只是描述底层数组中的一段** (注意:引用)

但是更改切片的元素 **会修改**其底层数的组中的元素，与它共享的数组的切片都会引用新的修改

### 切片语法slice literals

This is an array literals `[3]bool{true,true,false}`

This is a slice literals `[]bool{true,true,false}`

**注意[n]T和[]T是数组和切片声明方式的不同** 

切片看起来像创建一个和上面一模一样的数组，然后引用了它,所以切片的语法类似与no length的数组语法

>对于数组var s [10]int  a[0:10] 等价于 a[:10] a[0:] a[:]

切片提供了计算容量的函数 cap() 可以测量切片最长可以达到多少：它等于切片的长度 + 数组除切片之外的长度

1. 切片拥有length 和 capacity 容量

**切片的长度就是the number of elements it contains**

切片的容量是 the number of elements in the underlying底层,counting from  the first elements in the slice
**从其第一个元素开始数，到其底层数组元素末尾的个数**

切片s的长度和容量可通过表达式len(s) 和 cap(s) 来获取

可以通过重新切片来扩展一个切片，给它提供足够的容量，但不能超过底层元素个数

2. 切片的零值是 nil。

- 切片之间不能比较
- nil 切片的长度和容量为 0 且没有底层数组。
- 因为`s3 := make(int[],0)  //len(s3)=0;cap(s3)=0;s3!=nil`所以要判断一个切片是否是空的，要是用`len(s) == 0`来判断，不应该使用`s == nil`来判断。

#### 切片的扩容策略

1. 如果申请的容量大于原来的2倍,那就直接扩容至新申请的容量
2. 如果数组原容量小于1024，那就直接2倍
3. 如果数组原容量大于1024，那就以原容量的1/4循环扩容
4. 具体存储的值类型不同，扩容策略也有部分不同

### **用make函数创建切片creating s slice with make**

内建函数make来创建动态数组

make函数会分配一个**元素为零值**的**数组**并**返回**一个**引用了它的切片** 特别注意 是0
`a := make([]int, 5)` // len(a)=5  make返回的切片 [0,0,0,0,0]

要指定它的容量，需向make传入第三个参数:
`b := make([]int, 0, 5) `// len(b)=0                                                      

`b = b[:cap(b)]` // len(b)=5 cap(b)=5

`b = b[1:]` // len(b)=4  cap(b)=4

`b := make([]int, 2, 5)` // b len=2 cap=5 [0 0]

切片的切片: 切片可包含任何类型，包括其他切片详见slice-of-slice.go

### 向切片追加元素 appending to a slice

`func append(s []T, vs .....T) []T`

`s1 = append(s1, "广州")`

```go
ss := []string{"武汉","西安","苏州"}

s1 = append(s1, ss...)  // ...表示拆开
```

Go 提供了内建的append 函数，append 第一个参数s是一个元素类型为T的切片 T 可以是int string  其余的T类型会追加到切片的末尾

append 的结果是一个包含原切片所有元素加上新添加元素的切片

当s的底层数组太小，不足以容纳所有给定的值时,它就会分配一个更大的数组。返回的切片会指向这个新分配的数组

### copy 复制切片

`func copy (dst []Type, src []Type) int`

详见slice.md

### 切片与数组的关系

1. 切片不保存具体的值，是一个引用

2. 切片一定对应一个底层数组,(make 也是先申请数组空间再返回切片)

3. 底层数组占有一块连续的内存

### 传递不定长值的例子
```go
package main
import "fmt"

func MaxSubseqSum1 (seq ...int64) (M int64) {  thisnum , maxnum := 0,0}

func main() {   
   var slice = []int64{3,-1,5,-2,3,4,-5,2,4,-1}   
   MaxNum := MaxSubseqSum1(slice...) 
   fmt.Println(MaxNum)
}
```

## **映射Maps**

map 将keys 映射到 values `map[KeyType]ValueType`

其中

- `KeyType`表示键的类型

- `ValueType`表示键对应的值的类型

map类型的变量初始值为nil,需要使用make()函数来分配内存，语法为:`make(map[KeyType]ValueType),[cap])`
其中cap表示map的容量,该参数虽然不是必须的,但是我们应该在初始化map的时候就为其制定一个合适的容量

一个map 的零值为nil 一个nil map 没有keys 也不能添加keys 

make 函数返回给定类型的map 初始化和ready for use

map 的语法与struct类似 但是必须有键名keys

若顶级类型只是一个类型名，你可以在文法的元素中省略它。

**如果不存在这个key，返回这个key键的类型的零值**

### 删除

`delete(map, key)`

#### **for 循环的 range 形式可遍历切片或者映射**

当使用for 循环遍历切片时,每次迭代都会返回两个值, 第一个值为当前元素的下标,第二个值为该下标所对应元素的一个副本 即vaule

可以将下标或值赋予 _ 来忽略它。

```go
for i, _ := range pow  // 只要下标key  写成 for i := range pow 也可
for  _, value := range pwo // 只要值
```

## **结构体structs**约等于函数

结构体就是一组字段(field)

```go
type person {
    name stting
    age  int 
    gender bool
    city string
} // type:main.person
```

字段可以用.(dot)来访问

```go
// 匿名结构体 多用于临时场景
var person struct {
    name stting
    age  int 
    gender bool
    city string
} // type:struct
```

**结构体指针Pointers to structs**

(*p).X 语言允许隐式间接引用，直接写 p.X 就可以。

### 创建指针类型结构体

我们还可以通过使用`new`关键字对结构体进行实例化，得到的是结构体的地址。 格式如下：

```go
var p2 = new(person)
fmt.Printf("%T\n", p2)     //*main.person
fmt.Printf("p2=%#v\n", p2) //p2=&main.person{name:"", city:"", age:0}
```

从打印的结果中我们可以看出`p2`是一个结构体指针。

需要注意的是在Go语言中支持对结构体指针直接使用`.`来访问结构体的成员。

```go
var p2 = new(person)
p2.name = "小王子"
p2.age = 28
p2.city = "上海"
fmt.Printf("p2=%#v\n", p2) //p2=&main.person{name:"小王子", city:"上海", age:28}
```

### 取结构体的地址实例化

使用`&`对结构体进行取地址操作相当于对该结构体类型进行了一次`new`实例化操作。

```go
p3 := &person{}
fmt.Printf("%T\n", p3)     //*main.person
fmt.Printf("p3=%#v\n", p3) //p3=&main.person{name:"", city:"", age:0}
p3.name = "七米"
p3.age = 30
p3.city = "成都"
fmt.Printf("p3=%#v\n", p3) //p3=&main.person{name:"七米", city:"成都", age:30}
```

`p3.name = "七米"`其实在底层是`(*p3).name = "七米"`，这是Go语言帮我们实现的语法糖。

 [额外的扩展-与json](struct.md)

`json.Marshal(c1)``json.Unmarshal([]byte(str), &c1)`

注意：1. 结构体中的变量名字 首字母应该为 大写 2.反序列化时应该传指针

【进阶知识点】关于Go语言中的内存对齐推荐阅读:[在 Go 中恰到好处的内存对齐](https://segmentfault.com/a/1190000017527311?utm_campaign=studygolang.com&utm_medium=studygolang.com&utm_source=studygolang.com) 


### 结构体初始化

```go
type person struct {
	name  string
	age   int
}
var p1 person //结构体实例化
p1.name = "周林"
p1.age  = 9000

p2 := person{"邢浩杰",18} // 结构体实例化

p4 := newPerson("邱骥"，36) // 调用构造函数生成person类型变量
```



# Func函数,一等公民

type mytype int  ← type声明新的类型
func (p mytype) func name(q int) (r,s int) { return 0,0 }
0           1    	               2                3           4              5

- 0 关键字 func 用于定义一个函数;
- 1 函数可以绑定到特定的类型上。这叫做 接收者 。有接收者的函数被称作 method。
- 2 funcname 是你函数的名字
- 3 int 类型的变量 q 作为输入参数。参数用 pass-by-value 方式传递,意味着它们会被复制;
- 4 变量 r 和 s 是这个函数的 命名返回值。// 返回参数
- 5 这是函数体。注意 return 是一个语句,所以包裹参数的括号是可选的

## 匿名函数

```go
// 第一种 把匿名函数赋值给一个变量
sayhello := func() {
         fmt.Println(“匿名函数”）
}
sayhello()
// 第二种
func(x int) {
/* ... */
}(5) //为输入参数 x 赋值5   ()定义并执行了匿名函数
最后的一对括号表示对该匿名函数的调用
```

## 匿名函数和defer的用法

```go
// 在这个 (匿名) 函数中,可以访问任何命名返回参数:
//在 defer 中访问返回值
func f() (ret int) {  //← ret 被初始化为零
	defer func() {
	ret++  //← 将 ret 加一
	}()
	return 0   //← 返回的是 1 而 不是 0 !
}
```

## 变参/可变长参数

接受不定数量的参数的函数叫做变参函数。为了使其接受变参需要进行如下定义:
`func myfunc(arg ...int) {}`

`arg ... int` 告诉 Go 这个函数接受不定数量的参数。注意,这些参数的类型全部是
`int` ,可以为空也就是一个也不传入 。在函数体中,变量 `arg` 是一个 int 类型的 slice:

`func prtthem(numbers ... int) {} // numbers 现在是整数类型的 slice`

### 如何传入变参

```go
slice := []int{1,2,3,45,6}
prtthem(slice...)
```
### 疑问

s ...string
s []string 
有何不同

尝试解答: ...string可以为空  []string必须要传入切片

## 函数值

函数也是值. 它们可以像其他值一样传递。

**函数值可以用作函数的参数或者返回值**

```go
func compute(fn func(float64, float64) float64) float64 {
//{需要传入的变量:name                  type}   {函数返回的类型} 
	return fn(3, 4)
}
// 个人理解 不一定对
fn func(float64, float64)     相当与  fn := func(float64,float64)/ var fn  func   || 需要传入的变量

(fn func(float64, float64)  float64)    相当与  (func(float64, float64) float64)

fmt.Println(compute(math.Pow)) ==  math.Pow为fn 

return  fn(3, 4)  相当于  return  math.Pow(3,4)
```

## **闭包function closures**闭包 = 函数 + 外部变量的引用

Go 函数可以是一个闭包。闭包是一个函数值，它引用了**其函数体之外的变量**。
该函数可以访问并赋予其引用的变量的值，换句话说，该函数被这些变量“绑定”在一起。

闭包封闭的不是内部的状态，封闭的是外部的状态，当外部的scope失效后，还有一个留在内部状态里。
闭包可使得某个函数捕捉到一些外部状态，例如：函数被创建时的状态。另一种表示方式为：一个闭包继承了函数所声明时的作用域。这种状态（作用域内的变量）都被共享到闭包的环境中，因此这些变量可以在闭包中被操作，直到被销毁，

例如，函数 adder 返回一个闭包。每个闭包都被绑定在其各自的 sum 变量上。

```go
package main

import "fmt"

func adder() func(int) int {
	sum := 0
	return func(x int) int {  //闭包部分
		sum += x
		return sum
	}   // 闭包部分
}

func main() {
// pos,neg 在这里被赋值为匿名函数
	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
		fmt.Println(
			pos(i),
			neg(-2*i),
		)
	}
}
```
闭包函数保存并积累其中的变量的值，不管外部函数退出与否，它都能够继续操作外部函数中的局部变量。

函数作为值也会被用在其他地方,例如 map。这里将整数转换为函数:
Listing 2.12. 使用 map 的函数作为值
```go
var xs = map[int]func() int{
	1: func() int { return 10 },
	2: func() int { return 20 },
	3: func() int { return 30 },← 必须有逗号
	/* ... */
}
```

# Methods 方法

GO 没有类 A method is a function with a special ***receiver***
方法就是一类带特殊的 **接受者/参数**的函数，作用于特定类型的函数

> The receiver 在自己的参数列表内 位于func关键字和the method name方法名之间

可以在任意类型上定义方法(cd除了非本地类型,包括内建类型: int 类型不能有方法)。然而可以新建一个拥有方法的整数类型。
**约定成俗: 多用类型名首字母小写表示**

1. 可以为结构体类型定义方法
    `func (v Vertex) Abs() float64 {}`

   **Abs method 方法abs** has 一个名为v类型，type 为Vertex的接受者

方法abs 与regular function `func Abs(v Vertex) float64 {}` 效果是一样的

2. 你也可以为非结构体类型声明 method方法`func (f MyFloat) Abs() float64 {}`

接收者的类型定义和方法声明必须在同一包内；更不能为内建类型声明方法(such as int)

3. 为指针接收者声明方法

这意味着对于某类型 T，接收者的类型可以用 *T 的文法。（此外，T 不能是像 *int 这样的指针。）

例如，这里为 *Vertex 定义了 Scale 方法。

指针接收者的方法可以修改接收者指向的值（就像 Scale 在这做的）。由于方法经常需要修改它的接收者，指针接收者比值接收者更常用。

试着移除第 16 行 Scale 函数声明中的 *，观察此程序的行为如何变化。

若使用值接收者，那么 Scale 方法会对原始 Vertex 值的副本进行操作。（对于函数的其它参数也是如此。）Scale 方法必须用指针接受者来更改 main 函数中声明的 Vertex 的值。

## 指针与重定向

- 带指针参数的函数必须接受一个指针
- 以指针为接收者的方法被调用时，接收者既能为值又能为指针(推荐用这种方式) Go 会将语句 v.Scale(5) 解释为 (&v).Scale(5)。

- 接受一个值作为参数的函数必须接受一个指定类型的值
- 以值为接收者的方法被调用时，接收者既能为值又能为指  方法调用 p.Abs() 会被解释为 (*p).Abs()

使用指针接收者的原因有二：

首先，方法能够修改其接收者指向的值。

其次，这样可以避免在每次调用方法时复制该值。若值的类型为大型结构体时，这样做会更加高效。

# 接口 Interfaces

An ***interface type*** is defined as a set of method signatures
接口**类型** 是由**一组**方法签名定义的集合

隐式实现

接口也是值，它们也可以像其他值一样传递。可以用作函数的参数或者返回值

在内部，接口值可以看做包含值和具体类型的元组：

`(value, type)`

接口值保存了一个具体底层类型的具体值。

接口值调用方法时会执行其底层类型的同名方法。

底层值为 nil 的接口值
即便接口内的具体值为 nil，方法仍然会被 nil 接收者调用。

在一些语言中，这会触发一个空指针异常，但在 Go 中通常会写一些方法来优雅地处理它（如本例中的 M 方法）。

注意: 保存了 nil 具体值的接口其自身并不为 nil。

**制定0个方法的接口值为“空接口”** `interface{}`

**空接口可保存任何类型的值。（因为每个类型都至少实现了零个方法。）**

**空接口被用来处理未知类型的值。例如，fmt.Print 可接受类型为 interface{} 的任意数量的参数。**

[额外的解释](interface.md)

## 类型断言 Type assertions

A type assertions 可以访问  接口值的底层具体值

`t := i.(T)`
该语句断言接口值 i 保存了具体类型 T，并将其底层类型为 T 的值赋予变量 t

若 i 并未保存 T 类型的值，该语句就会触发一个panic恐慌。

为了 判断 一个接口值是否保存了一个特定的类型，类型断言可返回两个值：其底层值以及一个报告断言是否成功的布尔值。

`t, ok := i.(T)`
若 i 保存了一个 T，那么 t 将会是其底层值，而 ok 为 true。

否则，ok 将为 false 而 **t 将为 T 类型的零值**，程序并不会产生panic恐慌。

>请注意这种语法和读取一个映射时的相同之处。

### 类型选择

一种按顺序从几个类型断言中选择分支的结构.与一般的switch语句相似,不过case为类型(int float64)

```go
switch v := i.(type) {
case T:
    // v 的类型为 T
case S:
    // v 的类型为 S
default:
    // 没有匹配，v 与 i 的类型相同
}
```

类型选择中的声明与类型断言i.(T)的语法相同只不过把 具体类型T换成了keyword `type`
## 常见包中的接口

1. fmt 包中定义的 Stringer 是最普遍的接口之一。

```go
type Stringer interface {
    String() string
}
```

2.  go-error

Go程序使用error值来表示错误状态 same as fmt.Stringer.
`error`类型是一个内建接口:

```go
type error interface {
     Error() string
}
```

经典应用3

3. io 包中的 `io.Reader`

   Go 标准库包含了该接口的[许多实现](https://go-zh.org/search?q=Read#Global)，包括文件、网络连接、压缩和加密等等。

`io.Reader` 接口有一个 `Read` 方法：

```
func (T) Read(b []byte) (n int, err error)
```

`Read` 用数据填充给定的字节切片并返回填充的字节数和错误值。在遇到数据流的结尾时，它会返回一个 `io.EOF` 错误。

示例代码创建了一个 [`strings.Reader`](https://go-zh.org/pkg/strings/#Reader) 并以每次 8 字节的速度读取它的输出。

4. 图像

[`image`](https://go-zh.org/pkg/image/#Image) 包定义了 `Image` 接口：

```
package image

type Image interface {
    ColorModel() color.Model
    Bounds() Rectangle
    At(x, y int) color.Color
}
```

**注意:** `Bounds` 方法的返回值 `Rectangle` 实际上是一个 [`image.Rectangle`](https://go-zh.org/pkg/image/#Rectangle)，它在 `image` 包中声明。

（请参阅[文档](https://go-zh.org/pkg/image/#Image)了解全部信息。）

`color.Color` 和 `color.Model` 类型也是接口，但是通常因为直接使用预定义的实现 `image.RGBA` 和 `image.RGBAModel` 而被忽视了。这些接口和类型由 [`image/color`](https://go-zh.org/pkg/image/color/) 包定义。

## 反射 reflect

待补充

# panic 和 recover 恐慌和恢复

## Panic
是一个内建函数,可以中断原有的控制流程,进入一个令人恐慌的流程中。当函
数 F 调用 panic ,函数 F 的执行被中断,并且 F 中的延迟函数会正常执行,然后
F 返回到调用它的地方。在调用的地方, F 的行为就像调用了 panic 。这一过程
继续向上,直到当前的 goroutine 返回,这时程序崩溃。
恐慌可以直接调用 panic 产生。也可以由运行时错误产生,例如访问越界的数
组。

## Recover

是一个内建的函数,可以让进入令人恐慌的流程中的 goroutine 恢复过来。 recover
仅 在延迟函数中有效。在正常的执行过程中,调用 recover 会返回 nil 并且没
有其他任何效果。如果当前的 goroutine 陷入恐慌,调用 recover 可以捕获到
panic 的输入值,并且恢复正常的执行。

```go
//这个函数检查作为其参数的函数在执行时是否会产生 panic c :
func throwsPanic(f func()) (b bool) { 0
	defer func() { 1
		if x := recover(); x != nil {
			b = true
		}
	}()
	f() 2
	return 3
}
```
- 0 定义一个新函数 throwsPanic
  接受一个函数作为参数 (参看“函数作为值”)。函数 f 产生 panic,就返回 true,否则返回 false;
- 1 定义了一个利用 recover 的 defer 函数。
  如果当前的 goroutine 产生了 panic,这个 defer 函数能够发现。当 recover() 返回非 nil 值,设置 b 为 true;
- 2 调用作为参数接收的函数;
- 3 返回 b 的值。由于 b 是命名返回值(第 28 页),因此无须指定 b 。


# Go程和channel

## goroutine

是由Go运行时管理的轻量级线程.工作在用户态

```
go f(x, y, z)
```

会启动一个新的 Go 程并执行

```
f(x, y, z)
```

`f`, `x`, `y` 和 `z` 的求值发生在当前的 Go 程中，而 **`f` 的执行发生在新的 Go 程中**。

> goroutine 可能是按照从下往上的执行顺序来的，不确定待考察

### M:N
将m个goroutine调度到n个操作系统线程上,n默认是操作系统的逻辑核心数

### 可增长的栈

OS线程（操作系统线程）一般都有固定的栈内存（通常为2MB）,一个`goroutine`的栈在其生命周期开始时只有很小的栈（典型情况下2KB），`goroutine`的栈不是固定的，他可以按需增大和缩小，`goroutine`的栈大小限制可以达到1GB，虽然极少会用到这个大。所以在Go语言中一次创建十万左右的`goroutine`也是可以的。

### 调度模型

**GPM是Go语言运行时（runtime）层面的实现，是go语言自己实现的一套调度系统。区别于操作系统调度OS线程。**

- G很好理解，就是个goroutine的，里面除了存放本goroutine信息外 还有与所在P的绑定等信息。
- P管理着一组goroutine队列，P里面会存储当前goroutine运行的上下文环境（函数指针，堆栈地址及地址边界），P会对自己管理的goroutine队列做一些调度（比如把占用CPU时间较长的goroutine暂停、运行后续的goroutine等等）当自己的队列消费完了就去全局队列里取，如果全局队列里也消费完了会去其他P的队列里抢任务。
- M（machine）是Go运行时（runtime）对操作系统内核线程的虚拟， M与内核线程一般是一一映射的关系， 一个goroutine最终是要放到M上执行的；
  P与M一般也是一一对应的。他们关系是： P管理着一组G挂载在M上运行。当一个G长久阻塞在一个M上时，runtime会新建一个M，阻塞G所在的P会把其他的G 挂载在新建的M上。当旧的G阻塞完成或者认为其已经死掉时 回收旧的M。

P的个数是通过`runtime.GOMAXPROCS`设定（最大256），Go 1.5版本之后默认为物理线程数。 在并发量大的时候会增加一些P和M，但不会太多，切换太频繁的话得不偿失。

## Channel信道

信道是带有类型的管道，你可以通过它用信道操作符 `<-` 来发送或者接收值。
顾名思义开辟出来的道路，有接收方和发送方，不关心接收方和发送方在哪里，是否在一个函数中

```go
ch <- v    // 将 v 发送至信道 ch。
v := <-ch  // 从 ch 接收值并赋予 v。
<- ch      // 等待, 直到从channel上接收到一个值，值被忽略
ci <- 1    // 发送 整数 1 到 channel c
close(ch)  // 关闭  关闭后的通道还是可以取值,取完之后返回的是类型零值
```

（“箭头”就是数据流的方向。）
Go语言的并发模型是`CSP（Communicating Sequential Processes）`，提倡**通过通信共享内存**而不是**通过共享内存而实现通信**。

### Channel类型

Channel类型的定义格式如下：

```
ChannelType = ( "chan" | "chan" "<-" | "<-" "chan" ) ElementType .
```

它包括三种类型的定义。可选的`<-`代表channel的方向。如果没有指定方向，那么Channel就是双向的，既可以接收数据，也可以发送数据。

```go
chan T          // 可以接收和发送类型为 T 的数据
// 单项通道
chan<- float64  // 只可以用来发送 float64 类型的数据
<-chan int      // 只可以用来接收 int 类型的数据
```

`<-`总是优先和最左边的类型结合。(The <- operator associates with the leftmost chan possible)

```go
chan<- chan int    // 等价 chan<- (chan int)
chan<- <-chan int  // 等价 chan<- (<-chan int)
<-chan <-chan int  // 等价 <-chan (<-chan int)
chan (<-chan int)
```

![channel异常总结](E:/banksteel/simple-ci-description.wiki/images/channel01.png)

和映射与切片一样，信道在使用前必须创建：

```go
ch := make(chan int)
```

### 带缓冲的信道

信道可以是 *带缓冲的*。将缓冲长度作为第二个参数提供给 `make` 来初始化一个带缓冲的信道：

```
ch := make(chan int, 100)
```

仅当信道的缓冲区填满后，向其发送数据时才会阻塞。当缓冲区为空(即被填满)时，接受方会阻塞。

缓冲区填满后继续发送会抛出`fatal error: all goroutines are asleep - deadlock!`错误

uber的代码建议:要不使用无缓冲通道,要不使用缓冲为1的channel 

### close和 如何判断一个channel已经关闭(range)  

发送者可通过 `close` 关闭一个信道来表示没有需要发送的值了。

1. 接收者可以通过为接收表达式分配第二个参数来测试信道是否被关闭：若没有值可以接收且信道已被关闭，那么在执行完

```
v, ok := <-ch
```

之后 `ok` 会被设置为 `false`。

2. 循环 `for i := range c` 会不断从信道接收值，直到它被关闭。

***注意：* 只有发送者才能关闭信道，而接收者不能。向一个已经关闭的信道发送数据会引发程序恐慌（panic）。**

**关闭一个已经关闭的`channel`也会引发panic**

*还要注意：* 信道与文件不同，通常情况下无需关闭它们。只有在必须告诉接收者不再有需要发送的值时才有必要关闭，例如终止一个 `range` 循环。

### select 语句

`select` 语句使一个 Go 程可以等待多个通信操作。

`select` 会阻塞到某个分支可以继续执行为止，这时就会执行该分支。当**多个分支都准备好**时会随机选择一个执行。

 `select`语句选择**一组可能的send操作和receive操作**去处理。它类似`switch`,但是只是用来处理通讯(communication)操作。
-  它的`case`可以是send语句，也可以是receive语句，亦或者`default`。

- `receive`语句可以将值赋值给一个或者两个变量。它必须是一个receive操作。

- 最多允许有一个`default case`,它可以放在case列表的任何位置，尽管我们大部分会将它放在最后。

channel的 receive支持 *multi-valued assignment*，如

```
v, ok := <-ch
```

它可以用来检查Channel是否已经被关闭了。

#### 默认选择

如果没有case需要处理，则会选择`default`去处理，如果`default case`存在的情况下。如果没有`default case`，则`select`语句会阻塞，直到某个case需要处理。

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}
```

## sync.Mutex互斥锁

我们已经看到信道非常适合在各个 Go 程间进行通信。

但是如果我们并不需要通信呢？比如说，若我们只是想保证每次只有一个 Go 程能够访问一个共享的变量，从而避免冲突？

这里涉及的概念叫做 *互斥（mutual*exclusion）* ，我们通常使用 *互斥锁（Mutex）* 这一数据结构来提供这种机制。

Go 标准库中提供了 [`sync.Mutex`](https://go-zh.org/pkg/sync/#Mutex) 互斥锁类型及其两个方法：

- `Lock`
- `Unlock`

我们可以通过在代码前调用 `Lock` 方法，在代码后调用 `Unlock` 方法来保证一段代码的互斥执行。参见 `Inc` 方法。

我们也可以用 `defer` 语句来保证互斥锁一定会被解锁。参见 `Value` 方法。

### 值类型

`sync.Mutex`是一个结构体,是值类型，给函数传参属的时候要传指针

## 读写互斥锁

#### 应用场景

适用于读多写少场景

#### 特点

1. read goroutine拿到RLock，后续的goroutine能读不能写
2. write goroutine拿到Lock, 后续的goroutine不管是read是wirte都要等待

## sync.WaitGroup

用来等goroutine执行完再继续

详见package/sync.md

## sync.Map

详见package/sync.md

## sync.Atomic 原子操作

详见package/sync.md

# Testing测试

## 单元测试
测试函数需要以`Test目标函数`开头

## 基准测试

## 示例函数

## Setup和Teardown

setup准备环境

teardown拆解环境

# pprof调试工具





# go 的命令行操作

1. 只要在代码集中执行 go vet，就可以发现所有的无标签的语法

2. `go test`

   ```go
   go test
   go test -v
   go test -run=RegExp  // 例如-run="More"  
   go test -v -run=Split/simple  // 子测试t.Run()
   go test -cover
   go test -cover -coverprofile=c.out
   // 扩展
   go tool cover -html=c.out
   // 基准测试
   go test -bench=Split -benchmem
   ```

   

# 经典应用

### 1. 空白标识符_在函数返回值时的使用

```go
package main

import "fmt"

func main() {
  _,numb,strs := numbers() //只获取函数返回值的后两个
  fmt.Println(numb,strs)
}

//一个可以返回多个值的函数
func numbers()(int,int,string){
  a , b , c := 1 , 2 , "str"
  return a,b,c
}
```

输出结果: 2 str

### 2. 函数作为参数传递，实现回调。

```go
package main
import "fmt"

// 声明一个函数类型
type cb func(int) int

func main() {
    testCallBack(1, callBack)
    testCallBack(2, func(x int) int {
        fmt.Printf("我是回调，x：%d\n", x)
        return x
    })
}

func testCallBack(x int, f cb) {
    f(x)
}

func callBack(x int) int {
    fmt.Printf("我是回调，x：%d\n", x)
    return x
}
```

输出结果: 

我是回调，1

我是回调，2

### 3. 练习ERROR

   从[之前的练习](https://tour.go-zh.org/flowcontrol/8)中复制 `Sqrt` 函数，修改它使其返回 `error` 值。

   `Sqrt` 接受到一个负数时，应当返回一个非 nil 的错误值。复数同样也不被支持。

   创建一个新的类型

   ```
   type ErrNegativeSqrt float64
   ```

   并为其实现

   ```
   func (e ErrNegativeSqrt) Error() string
   ```

   方法使其拥有 `error` 值，通过 `ErrNegativeSqrt(-2).Error()`调用该方法应返回 `"cannot Sqrt negative number: -2"`。

   ```go
   package main
   import ("fmt")
   
   type ErrNegativeSqrt float64
   
   func (e ErrNegativeSqrt) Error() string {
       return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))#这里为什么必须要转换
   }
   func Sqrt(x float64) (float64 error) {
       if x < 0 {
           return 0, ErrNegativeSqrt(x)
       }
       z := x
       for i := 0 ; i < 10; i++ {
           z -= (z*z -x) / (2*z) 
       }
       return z , nil
   }
   func main() {
       fmt.Printf(Sqrt(2))
       fmt.Printf(Sqrt(-2))
   }
   ```

   在 `Error` 方法内调用 `fmt.Sprint(e)` 会让程序陷入死循环。可以通过先转换 `e` 来避免这个问题：`fmt.Sprint(float64(e))`。这是为什么呢？

   解: ErrNegativeSqrt 也是type 类似与int和string 不转换的话等于调用e本身的函数

### 4. 从未排序的切片中移除元素的有效方法

      ```go
      package main 
      func main() {
        scores := []int{1, 2, 3, 4, 5}
        scores = removeAtIndex(scores, 2)
        fmt.Println(scores) // [1 2 5 4]
      }
      
      // 不会保持顺序
      func removeAtIndex(source []int, index int) []int {
        lastIndex := len(source) - 1
        // 交换最后一个值和想去移除的值
        source[index], source[lastIndex] = source[lastIndex], source[index]
        return source[:lastIndex]
      }
      ```
### 5 interface的妙用,判断传入的类型
```go
package main

import "fmt"

func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Twice %v is %v\n", v, v*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", v, len(v))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}

func main() {
	do(21)
	do("hello")
	do(true)
}
```

# method调用和function调用的区别

可以对新定义的类型创建函数以便操作,可以通过两种途径:
1. 创建一个函数接受这个类型的参数。
    `func doSomething(n1 *NameAge, n2 int) { /* */ }`
    (你可能已经猜到了)这是 函数调用 。

2. 创建一个工作在这个类型上的函数(参阅在 2.1 中定义的 接收方 ):
    `func (n1 *NameAge) doSomething(n2 int) { /* */ }`
    这是 方法调用 ,可以类似这样使用:

  ```go
  var n *NameAge
  n.doSomething(2)
  ```
  使用函数还是方法完全是由程序员个人决定,但是若需要满足接口(参看下一章)就
必须使用方法。如果没有这样的需求,那就完全由习惯来决定是使用函数还是方法。

# 复合声明

构造函数与复合声明
有时零值不能满足需求,必须要有一个用于初始化的构造函数,例如这个来自 os 包的例子。

```go
func NewFile(fd int, name string) *File {
	if fd < 0 {
		return nil
	}
	f := new(File)
	f.fd = fd
	f.name = name
	f.dirinfo = nil
	f.nepipe = 0
	return f
}
```
有许多冗长的内容。可以使用复合声明使其更加简洁,每次只用一个表达式创建一个新的实例。
```go
func NewFile(fd int, name string) *File {
	if fd < 0 {
		return nil
	}
	f := File{fd, name, nil, 0}  // Create a new File
	return &f  //  返回 f 的地址
	}
```
返回本地变量的地址没有问题;在函数返回后,相关的存储区域仍然存在。
事实上,从复合声明获取分配的实例的地址更好,因此可以最终将两行缩短到一行

> 从复合声明中获取地址,意味着告诉编译器在堆中分配空间,而不是栈中

`return &File{fd, name, nil, 0}`
The items (called of a composite +literal are laid out in order and must all be 所有的项目(称作 字段)都必须按顺序全部写上。然而,通过对元素用字段: 值成对的标识,初始化内容可以按任意顺序出现,并且可以省略初始化为零值的字段。因此可以这样:` return &File{fd: fd, name: name}`

# fallthrough

```go
package main

import "fmt"

func main() {
	x := 2
	switch x {
	case 1:
		fmt.Print("test 1")
	case 2:
		fmt.Print("test 2")
		fallthrough //会有穿透性
	case 3:
		fmt.Print("test 2+3")
	default:
		fmt.Print("test 3")
	}
}
```



# Error

1. main redeclared in this block previous declaration at filename.go
   解决方案：在分别建立两个文件夹hello和sandbox，把文件放进去，再次BR两个文件，就没问题了。

   

   原因是：同一个目录下面不能有个多 package mainConcurrency

2. panic: reflect.Value.SetString using value obtained using unexported field

   原因是: pubilc的struct中没有public的var

   解决方案:

   ```go
   type Person struct {
    	name string
    	age int 
   }
   ----->
   type Person struct {
   		Name string
   		age int
   }
   // 之类只有Person.Name可以被外部函数引用、反射
   ```


- 新起一行输入fmt.，然后ctrl+x, ctrl+o，Vim 会弹出补齐提示下拉框，不过并非实时跟随的那种补齐，这个补齐是由gocode提供的。
- 输入一行代码：time.Sleep(time.Second)，执行:GoImports，Vim会自动导入time包。
  -
- 将光标移到Sleep函数上，执行:GoDef或命令模式下敲入gd，Vim会打开$GOROOT/src/time/sleep.go中 的Sleep函数的定义。执行:b 1返回到hellogolang.go。
  -
- 执行:GoLint，运行golint在当前Go源文件上。
  -
- 执行:GoDoc，打开当前光标对应符号的Go文档。
  -
- 执行:GoVet，在当前目录下运行go vet在当前Go源文件上。
  -
- 执行:GoRun，编译运行当前main package。
  -
- 执行:GoBuild，编译当前包，这取决于你的源文件，GoBuild不产生结果文件。
  -
- 执行:GoInstall，安装当前包。
  -
- 执行:GoTest，测试你当前路径下地_test.go文件。
  -
- 执行:GoCoverage，创建一个测试覆盖结果文件，并打开浏览器展示当前包的情况。
  -
- 执行:GoErrCheck，检查当前包种可能的未捕获的errors。
  -
- 执行:GoFiles，显示当前包对应的源文件列表。
  -
- 执行:GoDeps，显示当前包的依赖包列表。
  -
- 执行:GoImplements，显示当前类型实现的interface列表。
- 执行:GoRename [to]，将当前光标下的符号替换为[to]
