# http的两种处理方式

1. 

```go
http.HandleFUnc("/", HelloServer)

func HelloServer(w http.ResponseWriter, req *http.Request) {
    fmt.Println("Inside HelloServer handler")
    fmt.Fprintf(w, "Hello,"+ req.URL.Path[1:])
    // Fprintf 实现了io.writer方法
    // io.WriteString(w, "hello, world!\n")
    // Path[1:] 去除了'/'
}
```

```go
http.HandleFunc("/", Hfunc) 中的 HFunc 是一个处理函数，如下：

func HFunc(w http.ResponseWriter, req *http.Request) {
    ...
}
```

2. 第二种处理方法

```go
http.Handle("/", http.HandlerFunc(HFunc))
```

这里的 `HandlerFunc` 只是一个类型名称，它定义如下：

```go
type HandlerFunc func(ResponseWriter, *Request)
```

```
/ The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
它是一个可以把普通的函数当做 HTTP 处理器的适配器。如果 f 函数声明的合适，HandlerFunc(f) 就是一个执行了 f 函数的处理器对象
```

http.Handle 的第二个参数也可以是 T 的一个 obj 对象：http.Handle("/", obj) 给 T 提供了 ServeHTTP 方法，实现了 http 的 Handler 接口：

```go
func (obj *Typ) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    ...
}
```