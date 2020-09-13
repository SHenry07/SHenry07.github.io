# 获取参数/选项

### 手工处理

- $0 : 在用sh 或者 ./执行脚本时，指的是脚本名，用source或.执行时，永运是bash，这也反应了sh 或者 ./执行脚本的原理和source的方式是不同的.
- $1 : -v,第一个参数.
- $2 : -f
- $3 : -out
- $4 : /test.log
- 依次类推 $5 $6 …
- $# : 参数的个数，不包括命令本身，上例中\$#为5.
- $@ : 参数本身的列表，也不包括命令本身，如上例为 -v -f -out /test.log –prefix=/home
- $* : 参数本身的列表，也不包括命令本身，但”$*” 和”\$@”(加引号)并不同，”$*“将所有的参数解释成一个字符串，而”$@”是一个参数数组。如下例所示：

```
#!/bin/bash

for arg in "$*"
do
    echo $arg
done

for arg in "$@"
do
    echo $arg
done
-v -f -out /test.log --prefix=/home
-v
-f
-out
/test.log
--prefix=/home
```

**也就是说手工处理方式高度依赖命令行中参数的位置, 只适合简单的参数较少的命令, 手工处理方式能满足大多数的简单需求，配合shift使用也能构造出强大的功能**

### getopts

先来看看参数传递的典型用法:

- ./test.sh -a -b -c ： 短选项，各选项不需参数
- ./test.sh -abc ： 短选项，和上一种方法的效果一样，只是将所有的选项写在一起。
- ./test.sh -a args -b -c ：短选项，其中-a需要参数，而-b -c不需参数。
- ./test.sh –a-long=args –b-long ：长选项

#### getopts 用法

#### 变量

- **OPTIND**: getopts 在解析传入 Shell 脚本的参数时（也就是 $@），并不会执行 shift 操作，而是通过变量 OPTIND 来记住接下来要解析的参数的位置。

- **OPTARG**: getopts 在解析到选项的参数时，就会将参数保存在 OPTARG 变量当中；如果 getopts 遇到不合法的选项，择把选项本身保存在 OPTARG 当中。

  ```
  getopts OPTSTRING VARNAME [ARGS...]
  ```

- **OPTSTRING** 记录合法的选项列表（以及参数情况)

- **VARNAME** 则传入一个 Shell 变量的名字，用于保存 getopts 解析到的选项的名字（而不是参数值，参数值保存在 OPTARG 里）

- **ATGS…** 是可选的，默认是 $@，即传入 Shell 脚本的全部参数

通常来说，我们会将 getopts 放在 while 循环的条件判断式中。getopts 在顺利解析到参数的时候，会返回 TRUE；否则返回 FALSE，用以结束循环.

```
while getopts ...; do
    ...
done
```



getopts 在两种情况下会停止解析并返回 FALSE：

- getopts 读入不以 - 开始的字符串；比如: sh test.sh flag
- getopts 读入连续的两个 - (i.e. –)

#### OPTSTRING

通过 **OPTSTRING** getopts 知道哪些参数是合法的，哪些参数又是需要接受参数的。
OPTSTRING 的格式很简单，就是一个简单的字符串。字符串里，每一个字母（大小写均可，但区分大小写）都是一个选项的名字。

**值得一提的是冒号 (:)**
在 OPTSTRING 中，冒号有两种含义：

- 首位的 : 表示「不打印错误信息」；

- 紧邻字母（选项名字）的 : 表示该选项接收一个参数。

  例如:

  ```
  getopts aBcD VARNAME // 该脚本接受四个标签-a, -B, -c, -D, 均不接受参数
  getopts :aB:Cd VARNAME // 该脚本接受两个标签-a, -B, 两个短选项-C, -d
  ```

下面是实例:

```
#!/bin/bash
echo "$@"
while getopts ":a:bc:" opt; do #不打印错误信息, -a -c需要参数 -b 不需要传参  
  case $opt in
    a)
      echo "-a arg:$OPTARG index:$OPTIND" #$OPTIND指的下一个选项的index
      ;;
    b)
      echo "-b arg:$OPTARG index:$OPTIND"
      ;;
    c) 
      echo "-c arg:$OPTARG index:$OPTIND"
      ;;
    :)
      echo "Option -$OPTARG requires an argument." 
      exit 1
      ;;
    ?) #当有不认识的选项的时候arg为?
      echo "Invalid option: -$OPTARG index:$OPTIND"
      ;;
    
  esac
done
```



```
$ ./test.sh -a ssss -b ssss -c
  >>
  -a ssss -b ssss -c
  -a arg:ssss index:3
  -b arg: index:4 #-b并不接受参数, 解析到ssss时直接停止解析

$ ./test.sh  -c xxx -b -a ssssss
  >>
  -c xxx -b -a ssssss
  -c arg:xxx index:3
  -b arg: index:4
  -a arg:ssssss index:6

$ ./test.sh  -c -b -a ssssss  // -c 后面没有参数 -b会解析成-c参数
  >>
  -c -b -a ssssss
  -c arg:-b
  -a arg:ssssss
  
$ ./test.sh -a
  >> 
  a
  Option -a requires an argument.
```

# Reference

[shell - 参数解析三种方式(手工, getopts, getopt)](https://bummingboy.top/2017/12/19/shell%20-%20%E5%8F%82%E6%95%B0%E8%A7%A3%E6%9E%90%E4%B8%89%E7%A7%8D%E6%96%B9%E5%BC%8F(%E6%89%8B%E5%B7%A5,%20getopts,%20getopt)/)