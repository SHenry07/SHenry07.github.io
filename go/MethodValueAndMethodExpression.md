最早是在unknown的视频中听到,后查阅资料, 感觉上像是go后面对其作了统一



[github地址](https://github.com/unknwon/go-fundamental-programming/blob/master/lectures/lecture11.md)

```go
// 一个简单的例子
type TZ int

func main() {
    var a TZ
    // Method Value
    a.Print()
    
    // Method Expression 方法解析式  // 一种类似函数的方式,第一个参数传递的就是reciever
    (*TZ).Print(&a)
}

func (a *TZ)Print() {
    fmt.Println("TZ")
}
//本质上这就是一种语法糖，方法调用如下：
//instance.method(args) -> (type).func(instance, args)
```

[延申](https://blog.csdn.net/hittata/article/details/24346949)

[延申2](https://blog.csdn.net/q1007729991/article/details/80472622)