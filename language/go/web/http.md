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

# 调用流程

![img](./3.3.illustrator.png)



我们来梳理一下整个的代码执行过程。

- 首先调用Http.HandleFunc

  按顺序做了几件事：

  1 调用了DefaultServeMux的HandleFunc

  2 调用了DefaultServeMux的Handle

  3 往DefaultServeMux的map[string]muxEntry中增加对应的handler和路由规则

- 其次调用http.ListenAndServe(":9090", nil)

  按顺序做了几件事情：

  1 实例化Server

  2 调用Server的ListenAndServe()

  3 调用net.Listen("tcp", addr)监听端口

  4 启动一个for循环，在循环体中Accept请求

  5 对每个请求实例化一个Conn，并且开启一个goroutine为这个请求进行服务go c.serve()

  6 读取每个请求的内容w, err := c.readRequest()

  7 判断handler是否为空，如果没有设置handler，handler就设置为DefaultServeMux

  8 调用handler的ServeHttp

  9 在前两个例子中，下面就进入到DefaultServeMux.ServeHttp

  10 根据request选择handler，并且进入到这个handler的ServeHTTP

  ```
    mux.handler(r).ServeHTTP(w, r)
  ```

  11 选择handler：

  A 判断是否有路由能满足这个request（循环遍历ServeMux的muxEntry）

  B 如果有路由满足，调用这个路由handler的ServeHTTP

  C 如果没有路由满足，调用NotFoundHandler的ServeHTTP



# 第一种 自己实现handler 更灵活些

golang 的标准库 `net/http` 提供了 http 编程有关的接口，封装了内部TCP连接和报文解析的复杂琐碎的细节，使用者只需要和 `http.request` 和 `http.ResponseWriter` 两个对象交互就行。也就是说，我们只要写一个 handler，请求会通过参数传递进来，而它要做的就是根据请求的数据做处理，把结果写到 Response 中。

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

只要**实现**了 `ServeHTTP` 方法的对象都可以作为 Handler, 然后注册到**对应的路由上**即可

```go
package main

import (
    "io"
    "net/http"
)
// 定义一个类型/结构体
type helloHandler struct{}

// type编写ServeHTTP method 以实现 Handler interface
func (h *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, world!"))
}

func main() {
    http.Handle("/", &helloHandler{}) // 注册到对应的路由上
    http.ListenAndServe(":12345", nil)
}
```

## Handle 和 HandleFunc的不同

一个接受的是Handler，一个接受的是Func。

优缺点: Handler 更灵活，可以附带额外的method

HandleFunc 更便捷，不需要创建额外的结构体/类型，不需要实现ServeHTTP方法

```go
// Handle registers the handler for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched. 
func Handle(pattern string, handler Handler) { 
    DefaultServeMux.Handle(pattern, handler) 
}

// HandleFunc registers the handler function for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}
```



# 第二种  将普通函数封装成的handler函数 handleFunc

上面的代码没有什么问题，但是有一个不便：每次写 Handler 的时候，都要定义一个类型，然后编写对应的 `ServeHTTP` 方法，以实现handle接口，这个步骤对于所有 Handler 都是一样的。重复的工作总是可以抽象出来，`net/http` 也正这么做了，它提供了 `http.HandleFunc` 方法，允许直接把特定类型的函数作为 handler。

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
	mux.Handle(pattern, HandlerFunc(handler))//强制转换
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

将函数 `f` 强制转化为 `HandlerFunc`类型 ，这个类型默认就实现了`ServeHTTP` interface，这样f就拥有了ServeHTTP方法。`DefaultServeMux`路由器里面就创建并存储了相应的路由规则。这样封装的好处是：使用者可以专注于业务逻辑的编写，省去了很多重复的代码处理逻辑。如果只是简单的 Handler，会直接使用函数；如果是需要传递更多信息或者有复杂的操作，会使用上部分的方法。

如果需要我们自己写的话，是这样的：

```go
package main

func helloHandler(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "hello, world!\n")
}

func main() {
    // 即我们调用了HandlerFunc(helloHandler),强制类型转换helloHandler成为HandlerFunc类型 由上可知type  HandlerFunc默认实现了ServeHTTP 接口 这样f就拥有了ServeHTTP方法。
    // 即强制转换类型后，可以通过Handle方法绑定到DefaultServeMux上

    hh := http.HandlerFunc(helloHandler)
    http.Handle("/", hh)
    http.ListenAndServe(":12345", nil)
}
```
汇总对比版
```go
package main

func main() {
// 1. Handlefunc处理

	http.HandleFunc("/hello/", NameFunc)
// HanderFunc 封装函数  注意多了一个r
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

`Mux` 是 `multiplexor` 的缩写

mux主要用来控制路由router


```go
type myHandler struct{}
func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
   io.WriterString(w,"URL: "+ r.URL.String())
}
func echoHandler(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, r.URL.Path)
}
func main(){
	// mux 多路复用器
	mux := http.NewServeMux()
	mux.Handle("/", &myHandler{})
    mux.HandleFunc("/hello",echoHandler)
    
	log.Fatal(http.ListenAndServe(":8080",mux))
    // ListenAndServe(addr string, handler Handler)
}
```

这段代码和之前的代码有两点区别：

1. 通过 `NewServeMux` 生成了 `ServerMux` 结构，URL 和 handler 是通过它注册的
2. `http.ListenAndServe` 方法第二个参数变成了上面的 `mux` 变量

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
func (h Ownhandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}
func main() {
	var h Ownhandler
	http.ListenAndServe("127.0.0.1:9000", h)
	
}
```

## 关于`ServeMux` ，有几点要说明：

- URL 分为两种，末尾是 `/`：表示一个子树，后面可以跟其他子路径； 末尾不是 `/`，表示一个叶子，固定的路径
- 以`/` 结尾的 URL 可以匹配它的任何子路径，比如 `/images` 会匹配 `/images/cute-cat.jpg`
- 它采用最长匹配原则，如果有多个匹配，一定采用匹配路径最长的那个进行处理
- 如果没有找到任何匹配项，会返回 404 错误
- `ServeMux` 也会识别和处理 `.` 和 `..`，正确转换成对应的 URL 地址

你可能会有疑问？我们之间为什么没有使用 `ServeMux` 就能实现路径功能？那是因为 `net/http` 在后台默认创建使用了 `DefaultServeMux`

# 第四种  更深一层  serve/ServeMux的完全自定义


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

# 关于net库的更多信息

嗯，上面基本覆盖了编写 HTTP 服务端需要的所有内容。这部分就分析一下，它们的源码实现，加深理解，以后遇到疑惑也能通过源码来定位和解决。

### Server

首先来看 `http.ListenAndServe()`:

```go
func ListenAndServe(addr string, handler Handler) error {
    server := &Server{Addr: addr, Handler: handler}
    return server.ListenAndServe()
}
```

这个函数其实也是一层封装，创建了 `Server` 结构，并调用它的 `ListenAndServe` 方法，那我们就跟进去看看：

```go
// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.
type Server struct {
    Addr           string        // TCP address to listen on, ":http" if empty
    Handler        Handler       // handler to invoke, http.DefaultServeMux if nil
    ......
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.  If
// srv.Addr is blank, ":http" is used.
func (srv *Server) ListenAndServe() error {
    addr := srv.Addr
    if addr == "" {
        addr = ":http"
    }
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}
```

`Server` 保存了运行 HTTP 服务需要的参数，调用 `net.Listen` 监听在对应的 tcp 端口，`tcpKeepAliveListener` 设置了 TCP 的 `KeepAlive` 功能，最后调用 `srv.Serve()`方法开始真正的循环逻辑。我们再跟进去看看 `Serve` 方法：

```go
// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each.  The service goroutines read requests and
// then call srv.Handler to reply to them.
func (srv *Server) Serve(l net.Listener) error {
    defer l.Close()
    var tempDelay time.Duration // how long to sleep on accept failure
    // 循环逻辑，接受请求并处理
    for {
         // 有新的连接
        rw, e := l.Accept()
        if e != nil {
            if ne, ok := e.(net.Error); ok && ne.Temporary() {
                if tempDelay == 0 {
                    tempDelay = 5 * time.Millisecond
                } else {
                    tempDelay *= 2
                }
                if max := 1 * time.Second; tempDelay > max {
                    tempDelay = max
                }
                srv.logf("http: Accept error: %v; retrying in %v", e, tempDelay)
                time.Sleep(tempDelay)
                continue
            }
            return e
        }
        tempDelay = 0
         // 创建 Conn 连接
        c, err := srv.newConn(rw)
        if err != nil {
            continue
        }
        c.setState(c.rwc, StateNew) // before Serve can return
         // 启动新的 goroutine 进行处理
        go c.serve()
    }
}
```

最上面的注释也说明了这个方法的主要功能：

- 接受 `Listener l` 传递过来的请求
- 为每个请求创建 goroutine 进行后台处理
- goroutine 会读取请求，调用 `srv.Handler`

```go
func (c *conn) serve() {
    origConn := c.rwc // copy it before it's set nil on Close or Hijack

      ...

    for {
        w, err := c.readRequest()
        if c.lr.N != c.server.initialLimitedReaderSize() {
            // If we read any bytes off the wire, we're active.
            c.setState(c.rwc, StateActive)
        }

         ...

        // HTTP cannot have multiple simultaneous active requests.[*]
        // Until the server replies to this request, it can't read another,
        // so we might as well run the handler in this goroutine.
        // [*] Not strictly true: HTTP pipelining.  We could let them all process
        // in parallel even if their responses need to be serialized.
        serverHandler{c.server}.ServeHTTP(w, w.req)

        w.finishRequest()
        if w.closeAfterReply {
            if w.requestBodyLimitHit {
                c.closeWriteAndWait()
            }
            break
        }
        c.setState(c.rwc, StateIdle)
    }
}
```

看到上面这段代码 `serverHandler{c.server}.ServeHTTP(w, w.req)`这一句了吗？它会调用最早传递给 `Server` 的 Handler 函数：

```go
func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
    handler := sh.srv.Handler
    if handler == nil {
        handler = DefaultServeMux
    }
    if req.RequestURI == "*" && req.Method == "OPTIONS" {
        handler = globalOptionsHandler{}
    }
    handler.ServeHTTP(rw, req)
}
```

哇！这里看到 `DefaultServeMux` 了吗？如果没有 handler 为空，就会使用它。`handler.ServeHTTP(rw, req)`，Handler 接口都要实现 `ServeHTTP` 这个方法，因为这里就要被调用啦。

也就是说，无论如何，最终都会用到 `ServeMux`，也就是负责 URL 路由的家伙。前面也已经说过，它的 `ServeHTTP` 方法就是根据请求的路径，把它转交给注册的 handler 进行处理。这次，我们就在源码层面一探究竟。

### ServeMux

我们已经知道，`ServeMux` 会以某种方式保存 URL 和 Handlers 的对应关系，下面我们就从代码层面来解开这个秘密：

```go
type ServeMux struct {
    mu    sync.RWMutex //锁，由于请求涉及到并发处理，因此这里需要一个锁机制
    m     map[string]muxEntry  // 存放路由信息的字典！\(^o^)/ 由规则，一个string对应一个mux实体，这里的string就是注册的路由表达式
    es    []muxEntry // slice of entries sorted from longest to shortest.
    hosts bool // whether any patterns contain hostnames // 是否在任意的规则中带有host信息
}

// 下面看一下muxEntry

type muxEntry struct {
	explicit bool   // 是否精确匹配
	h        Handler // 这个路由表达式对应哪个handler
	pattern  string  //匹配字符串
}

// 接着看一下Handler的定义
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)  // 路由实现器
}
```

没错，数据结构也比较直观，和我们想象的差不多，路由信息保存在字典中，接下来就看看几个重要的操作：路由信息是怎么注册的？`ServeHTTP` 方法到底是怎么做的？路由查找过程是怎样的？

```go
// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (mux *ServeMux) Handle(pattern string, handler Handler) {
    mux.mu.Lock()
    defer mux.mu.Unlock()

    // 边界情况处理
    if pattern == "" {
        panic("http: invalid pattern " + pattern)
    }
    if handler == nil {
        panic("http: nil handler")
    }
    if mux.m[pattern].explicit {
        panic("http: multiple registrations for " + pattern)
    }

    // 创建 `muxEntry` 并添加到路由字典中
    mux.m[pattern] = muxEntry{explicit: true, h: handler, pattern: pattern}

    if pattern[0] != '/' {
        mux.hosts = true
    }

    // 这是一个很有用的小技巧，如果注册了 `/tree/`， `serveMux` 会自动添加非精准匹配的 `/tree` 并重定向到 `/tree/`。当然这个 `/tree` 会被精准匹配的路由信息覆盖。
    // Helpful behavior:
    // If pattern is /tree/, insert an implicit permanent redirect for /tree.
    // It can be overridden by an explicit registration.
    n := len(pattern)
    if n > 0 && pattern[n-1] == '/' && !mux.m[pattern[0:n-1]].explicit {
        // If pattern contains a host name, strip it and use remaining
        // path for redirect.
        path := pattern
        if pattern[0] != '/' {
            // In pattern, at least the last character is a '/', so
            // strings.Index can't be -1.
            path = pattern[strings.Index(pattern, "/"):]
        }
        mux.m[pattern[0:n-1]] = muxEntry{h: RedirectHandler(path, StatusMovedPermanently), pattern: pattern}
    }
}
```

路由注册没有什么特殊的地方，很简单，也符合我们的预期，注意最后一段代码对类似 `/tree` URL 重定向的处理。

```go
// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
    if r.RequestURI == "*" {
        if r.ProtoAtLeast(1, 1) {
            w.Header().Set("Connection", "close")
        }
        w.WriteHeader(StatusBadRequest)
        return
    }
    h, _ := mux.Handler(r)
    h.ServeHTTP(w, r)
}
```

如上所示路由器接收到请求之后，如果是*那么关闭链接，不然调用mux.Handler(r)返回对应设置路由的处理Handler，然后调用它的 `ServeHTTP` 方法。

那么mux.Handler(r)怎么处理的呢？ 

```go
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string) {
	if r.Method != "CONNECT" {
		if p := cleanPath(r.URL.Path); p != r.URL.Path {
			_, pattern = mux.handler(r.Host, p)
			return RedirectHandler(p, StatusMovedPermanently), pattern
		}
	}	
	return mux.handler(r.Host, r.URL.Path)
}

func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	// Host-specific pattern takes precedence over generic ones
	if mux.hosts {
		h, pattern = mux.match(host + path)
	}
	if h == nil {
		h, pattern = mux.match(path)
	}
	if h == nil {
		h, pattern = NotFoundHandler(), ""
	}
	return
}
```

它最终会调用 `mux.match()` 方法，我们来看一下它的实现：

```go
// Does path match pattern?
func pathMatch(pattern, path string) bool {
    if len(pattern) == 0 {
        // should not happen
        return false
    }
    n := len(pattern)
    if pattern[n-1] != '/' {
        return pattern == path
    }
    // 匹配的逻辑很简单，path 前面的字符和 pattern 一样就是匹配
    return len(path) >= n && path[0:n] == pattern
}

// Find a handler on a handler map given a path string
// Most-specific (longest) pattern wins
func (mux *ServeMux) match(path string) (h Handler, pattern string) {
    var n = 0
    for k, v := range mux.m {
        if !pathMatch(k, path) {
            continue
        }
         // 最长匹配的逻辑在这里
        if h == nil || len(k) > n {
            n = len(k)
            h = v.h
            pattern = v.pattern
        }
    }
    return
}
```

`match` 会根据用户请求的URL和路由器里面存储的map去匹配路由，找到所有匹配该路径最长的那个。当匹配到之后返回存储的handler，调用这个handler的ServeHTTP接口就可以执行到相应的函数了。

这就是DefaultServeMux

### Request

最后一部分，要讲讲 Handler 函数接受的两个参数：`http.Request` 和 `http.ResponseWriter`。

Request 就是封装好的客户端请求，包括 URL，method，header 等等所有信息，以及一些方便使用的方法：

```go
// A Request represents an HTTP request received by a server
// or to be sent by a client.
//
// The field semantics differ slightly between client and server
// usage. In addition to the notes on the fields below, see the
// documentation for Request.Write and RoundTripper.
type Request struct {
    // Method specifies the HTTP method (GET, POST, PUT, etc.).
    // For client requests an empty string means GET.
    Method string

    // URL specifies either the URI being requested (for server
    // requests) or the URL to access (for client requests).
    //
    // For server requests the URL is parsed from the URI
    // supplied on the Request-Line as stored in RequestURI.  For
    // most requests, fields other than Path and RawQuery will be
    // empty. (See RFC 2616, Section 5.1.2)
    //
    // For client requests, the URL's Host specifies the server to
    // connect to, while the Request's Host field optionally
    // specifies the Host header value to send in the HTTP
    // request.
    URL *url.URL

    // The protocol version for incoming requests.
    // Client requests always use HTTP/1.1.
    Proto      string // "HTTP/1.0"
    ProtoMajor int    // 1
    ProtoMinor int    // 0

    // A header maps request lines to their values.
    // If the header says
    //
    //    accept-encoding: gzip, deflate
    //    Accept-Language: en-us
    //    Connection: keep-alive
    //
    // then
    //
    //    Header = map[string][]string{
    //        "Accept-Encoding": {"gzip, deflate"},
    //        "Accept-Language": {"en-us"},
    //        "Connection": {"keep-alive"},
    //    }
    //
    // HTTP defines that header names are case-insensitive.
    // The request parser implements this by canonicalizing the
    // name, making the first character and any characters
    // following a hyphen uppercase and the rest lowercase.
    //
    // For client requests certain headers are automatically
    // added and may override values in Header.
    //
    // See the documentation for the Request.Write method.
    Header Header

    // Body is the request's body.
    //
    // For client requests a nil body means the request has no
    // body, such as a GET request. The HTTP Client's Transport
    // is responsible for calling the Close method.
    //
    // For server requests the Request Body is always non-nil
    // but will return EOF immediately when no body is present.
    // The Server will close the request body. The ServeHTTP
    // Handler does not need to.
    Body io.ReadCloser

    // ContentLength records the length of the associated content.
    // The value -1 indicates that the length is unknown.
    // Values >= 0 indicate that the given number of bytes may
    // be read from Body.
    // For client requests, a value of 0 means unknown if Body is not nil.
    ContentLength int64

    // TransferEncoding lists the transfer encodings from outermost to
    // innermost. An empty list denotes the "identity" encoding.
    // TransferEncoding can usually be ignored; chunked encoding is
    // automatically added and removed as necessary when sending and
    // receiving requests.
    TransferEncoding []string

    // Close indicates whether to close the connection after
    // replying to this request (for servers) or after sending
    // the request (for clients).
    Close bool

    // For server requests Host specifies the host on which the
    // URL is sought. Per RFC 2616, this is either the value of
    // the "Host" header or the host name given in the URL itself.
    // It may be of the form "host:port".
    //
    // For client requests Host optionally overrides the Host
    // header to send. If empty, the Request.Write method uses
    // the value of URL.Host.
    Host string

    // Form contains the parsed form data, including both the URL
    // field's query parameters and the POST or PUT form data.
    // This field is only available after ParseForm is called.
    // The HTTP client ignores Form and uses Body instead.
    Form url.Values

    // PostForm contains the parsed form data from POST or PUT
    // body parameters.
    // This field is only available after ParseForm is called.
    // The HTTP client ignores PostForm and uses Body instead.
    PostForm url.Values

    // MultipartForm is the parsed multipart form, including file uploads.
    // This field is only available after ParseMultipartForm is called.
    // The HTTP client ignores MultipartForm and uses Body instead.
    MultipartForm *multipart.Form

    ...

    // RemoteAddr allows HTTP servers and other software to record
    // the network address that sent the request, usually for
    // logging. This field is not filled in by ReadRequest and
    // has no defined format. The HTTP server in this package
    // sets RemoteAddr to an "IP:port" address before invoking a
    // handler.
    // This field is ignored by the HTTP client.
    RemoteAddr string
    ...
}
```

Handler 需要知道关于请求的任何信息，都要从这个对象中获取，一般不会直接修改这个对象（除非你非常清楚自己在做什么）！

### ResponseWriter

ResponseWriter 是一个接口，定义了三个方法：

- `Header()`：返回一个 Header 对象，可以通过它的 `Set()` 方法设置头部，注意最终返回的头部信息可能和你写进去的不完全相同，因为后续处理还可能修改头部的值（比如设置 `Content-Length`、`Content-type` 等操作）
- `Write()`： 写 response 的主体部分，比如 `html` 或者 `json` 的内容就是放到这里的
- `WriteHeader()`：设置 status code，如果没有调用这个函数，默认设置为 `http.StatusOK`， 就是 `200` 状态码

```go
// A ResponseWriter interface is used by an HTTP handler to
// construct an HTTP response.
type ResponseWriter interface {
    // Header returns the header map that will be sent by WriteHeader.
    // Changing the header after a call to WriteHeader (or Write) has
    // no effect.
    Header() Header

    // Write writes the data to the connection as part of an HTTP reply.
    // If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
    // before writing the data.  If the Header does not contain a
    // Content-Type line, Write adds a Content-Type set to the result of passing
    // the initial 512 bytes of written data to DetectContentType.
    Write([]byte) (int, error)

    // WriteHeader sends an HTTP response header with status code.
    // If WriteHeader is not called explicitly, the first call to Write
    // will trigger an implicit WriteHeader(http.StatusOK).
    // Thus explicit calls to WriteHeader are mainly used to
    // send error codes.
    WriteHeader(int)
}
```

实际上传递给 Handler 的对象是:

```go
// A response represents the server side of an HTTP response.
type response struct {
    conn          *conn
    req           *Request // request for this response
    wroteHeader   bool     // reply header has been (logically) written
    wroteContinue bool     // 100 Continue response was written

    w  *bufio.Writer // buffers output in chunks to chunkWriter
    cw chunkWriter
    sw *switchWriter // of the bufio.Writer, for return to putBufioWriter

    // handlerHeader is the Header that Handlers get access to,
    // which may be retained and mutated even after WriteHeader.
    // handlerHeader is copied into cw.header at WriteHeader
    // time, and privately mutated thereafter.
    handlerHeader Header
    ...
    status        int   // status code passed to WriteHeader
    ...
}
```

它当然实现了上面提到的三个方法，具体代码就不放到这里了，感兴趣的可以自己去看。



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



# Reference

- [柴大的教学](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/03.4.md)

- 无闻的视频