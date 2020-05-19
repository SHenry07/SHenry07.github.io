Mux是`multiplexor` 的缩写，就是多路传输的意思（请求传过来，根据某种判断，分流到后端多个不同的地方）。`ServeMux` 可以注册多了 URL 和 handler 的对应关系，并自动把请求转发到对应的 handler 进行处理。

`func(ResponseWriter, *Request))`
// ResponseWriter是个interface, 只是为了规范数据的返回
// Request 包含的是用户请求的所有数据, METHOD BODY等等, 所以要传指针

```go
// 什么是Handler 看下文的时候一定要明确
type Handler interface {
   ServeHTTP(ResponseWriter, *Request)
}
```


# 第一种 自己实现handler 更灵活些
```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

只要实现了 `ServeHTTP` 方法的对象都可以作为 Handler, 然后注册到**对应的路由上**即可

```go
package main

import (
    "io"
    "net/http"
)

type helloHandler struct{}

func (h *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, world!"))
}

func main() {
    http.Handle("/", &helloHandler{})
    http.ListenAndServe(":12345", nil)
}
```
# 第二种 handleFunc 封装好的handler函数

上面的代码没有什么问题，但是有一个不便：每次写 Handler 的时候，都要定义一个类型，然后编写对应的 `ServeHTTP` 方法，这个步骤对于所有 Handler 都是一样的。重复的工作总是可以抽象出来，`net/http` 也正这么做了，它提供了 `http.HandleFunc` 方法，允许直接把特定类型的函数作为 handler。

```go
func helloHandler(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "hello, world!\n")
}

func main() {
    http.HandleFunc("/", helloHandler)
    http.ListenAndServe(":12345", nil)
}
```

其实`HandleFunc`只是一个适配器

```go
// HandleFunc registers the handler function for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}
// HandleFunc registers the handler function for the given pattern.
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, HandlerFunc(handler))
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers.  If f is a function
// with the appropriate 恰当的 signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```

自动给 `f` 函数添加了 `HandlerFunc` 这个壳，最终调用的还是 `ServerHTTP`，只不过会直接使用 `f(w, r)`。这样封装的好处是：使用者可以专注于业务逻辑的编写，省去了很多重复的代码处理逻辑。如果只是简单的 Handler，会直接使用函数；如果是需要传递更多信息或者有复杂的操作，会使用上部分的方法。

如果需要我们自己写的话，是这样的：

```go
package main

func helloHandler(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "hello, world!\n")
}

func main() {
    // 通过 HandlerFunc 把函数转换成 Handler 接口的实现对象
    hh := http.HandlerFunc(helloHandler)
    http.Handle("/", hh)
    http.ListenAndServe(":12345", nil)
}
```
汇总对比版
```go
package main

func main() {
// 1. handlefunc处理

	http.HandleFunc("/hello/", NameFunc)
// HanderFunc 封装函数
//   type HandlerFunc func(ResponseWriter, *Request)
	http.Handle("/shouthello/", http.HandlerFunc(NameFunc2))//  http.HandlerFunc(NameFunc2)只是转换了类型
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
	fmt.Println("listen on 9000")
}


func NameFunc(w http.ResponseWriter, req *http.Request) {
	remPartOfURL := req.URL.Path[len("/hello/"):] //get everything after the /hello/ part of the URL
	fmt.Fprintf(w, "Hello %s!\n", remPartOfURL)
	fmt.Println("test01 /hello/Name")
	fmt.Fprintf(w,"<h1>%s</h1> ", req.URL.Path[7:])
}

func NameFunc2(w http.ResponseWriter, req *http.Request) {
	remPartOfURL := req.URL.Path[len("/shouthello/"):] //get everything after the /shouthello/ part of the URL
	fmt.Fprintf(w, "Hello %s!", strings.ToUpper(remPartOfURL))
}

```



# 第三种 自己实现Mux, 不使用默认的

mux主要用来控制路由router

```go
type myHandler struct{}
func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
   io.WriterString(w,"URL: "+ r.URL.String())
}
func main(){
	// mux 多路复用器
	mux := http.NewServeMux()
	mux.Handle("/", &myHandler{})
	
	log.Fatal(http.ListenAndServe(":8080",mux))
    // ListenAndServe(addr string, handler Handler)
}
```

```go
// ListenAndServe listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// The handler is typically nil, in which case the DefaultServeMux is used.
//
// ListenAndServe always returns a non-nil error.
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
```

这段代码和之前的代码有两点区别：

1. 通过 `NewServeMux` 生成了 `ServerMux` 结构，URL 和 handler 是通过它注册的
2. `http.ListenAndServe` 方法第二个参数变成了上面的 `mux` 变量

还记得我们之前说过，`http.ListenAndServe` 第二个参数应该是 Handler 类型的变量吗？这里为什么能传过来 `ServeMux`？嗯，估计你也猜到啦：`ServeMux` 也是是 `Handler` 接口的实现，也就是说它实现了 `ServeHTTP` 方法，我们来看一下：

```go
type ServeMux struct {
        // contains filtered or unexported fields
}

func NewServeMux() *ServeMux
func (mux *ServeMux) Handle(pattern string, handler Handler) // 
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request))
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string)
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request)// mux实现了Handler
```

哈！果然，这里的方法我们大都很熟悉，除了 `Handler()` 返回某个请求的 Handler。`Handle` 和 `HandleFunc` 这两个方法 `net/http` 也提供了，后面我们会说明它们之间的关系。而 `ServeHTTP` 就是 `ServeMux` 的核心处理逻辑：**根据传递过来的 Request，匹配之前注册的 URL 和处理函数，找到最匹配的项，进行处理。**可以说 `ServeMux` 是个特殊的 Handler，它负责路由和调用其他后端 Handler 的处理方法。

```go
type Ownhandler struct {}
func (h handler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}
func main() {
	var h Ownhandler
	http.ListenAndServe("127.0.0.1:9000", h)
	
}
```

# 第四种  更深一层的serve

```go
package main 
import (
	"io"
    "log"
    "net/http"
    "time"
)
// 自己的mux实现路由转发
var mux map[string]func(http.ResponseWriter, *http.Request)
func main() {
    server := http.Server{
        Addr: ":8080",
        Handler: &myHandler{},
        ReadTimeout: time.Second * 5,
    }
    mux = make(map[string]func(http.ResponseWriter, *http.Request))
    mux["/hello"] = sayHello
    mux["/bye"] = sayBye
    err := server.ListenAndServe()
    if err != nil {
        log.Fatal(err)
    }
}

type myHandler struct{}
func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if h, ok := mux[r.URL.String()]; ok {
        h(w, r)
        return 
    }
    io.WriterString(w,"URL: "+ r.URL.String())
}

func sayHello(w http.RsponseWriter, r *http.Request) {
    io.WriteString(w, "Hello world, this is version 4")
}

func sayBye(w http.RsponseWriter, r *http.Request) {
    
}
```

[更多阅读](https://cizixs.com/2016/08/17/golang-http-server-side/)

# 完整案例

```go
package main

import (
	"bytes"
	"expvar"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var helloRequests = expvar.NewInt("hello-requests")

// flags
var webroot = flag.String("root","/home/user","web root directory")

// simple flag server
var booleanflag = flag.Bool("boolean",true, "another flag for testing")

// 简单的服务器计数器, 发布它将设置值
type Counter struct {
	n int
}
type Chan chan int
func main() {
	flag.Parse()
	http.Handle("/", http.HandlerFunc(Logger))
	http.Handle("/go/hello", http.HandlerFunc(HelloServer))

	ctr := new(Counter)
	expvar.Publish("counter", ctr)
	http.Handle("/counter", ctr)

	http.Handle("/go/", http.StripPrefix("/go/", http.FileServer(http.Dir(*webroot))))

	http.Handle("/flags", http.HandlerFunc(FlagServer))
	http.Handle("/args", http.HandlerFunc(ArgServer))


	http.Handle("/chan", chanCreate())
	http.Handle("/date", http.HandlerFunc(DateServer))

	log.Fatal(http.ListenAndServe(":12345", nil))
}


func Logger(w http.ResponseWriter, req *http.Request) {
	log.Print(req.URL.String())
	w.WriteHeader(404)
	w.Write([]byte("oops"))
}


func HelloServer(w http.ResponseWriter, req *http.Request) {
	helloRequests.Add(1)
	io.WriteString(w, "hello, world\n")
}

func (ctr *Counter)String() string {
	return fmt.Sprintf("%d", ctr.n)
}
// ctr 实现了 ServeHTTP 方法，就实现了 Handler 接口，可以看到示例中，就不需要再通过 HandlerFunc 了，因为它自己就已经是一个 Handler 了
func (ctr *Counter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ctr.n++
	case "POST":
		buf := new(bytes.Buffer)
		io.Copy(buf, r.Body)
		body := buf.String()
		if n, err := strconv.Atoi(body); err != nil {
			fmt.Fprintf(w, "bad POST: %v\nbody: [%v]\n", err,body)
		}else{
			ctr.n = n
			fmt.Fprint(w, "counter reser\n")
		}
	}
	fmt.Fprintf(w, "counter = %d\n", ctr.n)
}


func FlagServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charest=utf-8")
	fmt.Fprint(w, "Flags:\n")
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() != f.DefValue {
			fmt.Fprintf(w, "%s = %s [default = %s]\n", f.Name, f.Value.String(), f.DefValue)
		}else{
			fmt.Fprintf(w, "%s = %s\n", f.Name, f.Value.String())
		}
	})
}


func ArgServer(w http.ResponseWriter, req *http.Request) {
	for _, s := range os.Args {
		fmt.Fprint(w,s," ")
	}
}

func chanCreate() Chan {
	c := make(Chan)
	go func (c Chan) {
		for x := 0; ; x++ {
			c <- x

		}
	}(c)
	return  c
}

func  (ch Chan)ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	io.WriteString(w, fmt.Sprintf("channel send #%d\n", <-ch))
}


// 执行一个程序, 输出重定向

func DateServer(rw http.ResponseWriter, req *http.Request){
	rw.Header().Set("Content-Type", "text/plain; charest=utf-8")
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintf(rw, "pipe: %s\n", err)
		return
	}

	p, err := os.StartProcess("/bin/date", []string{"date"}, &os.ProcAttr{Files: []*os.File{nil,w,w}})
	defer r.Close()
	w.Close()

	if err != nil {
		fmt.Fprintf(rw, "fork/exec: %s\n", err)
		return
	}

	defer p.Release()
	io.Copy(rw,r)
	wait, err := p.Wait()
	if err != nil {
		fmt.Fprintf(rw, "wait: %s\n", err)
		return
	}

	if !wait.Exited() {
		fmt.Fprintf(rw, "date: %v\n", wait)
		return
	}
}

```



