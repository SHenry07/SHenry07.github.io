golang的字符称为rune，等价于C中的char，可直接与整数转换

```go
    var c rune='a' 
    var i int =98
    i1:=int(c)
    fmt.Println("'a' convert to",i1)
    c1:=rune(i)
    fmt.Println("98 convert to",string(c1))

    //string to rune
    for _, char := range []rune("世界你好") {
        fmt.Println(string(char))
    }
```

rune实际是整型，必需先将其转换为string才能打印出来，否则打印出来的是一个整数

```go
c:='a'
fmt.Println(c)
fmt.Println(string(c))
fmt.Println(string(97))
```

输出

```txt
97
a
a
```