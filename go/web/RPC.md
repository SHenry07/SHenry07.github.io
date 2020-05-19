# RPC（Remote Procedure Call Protocol）——远程过程调用协议

RPC就是想实现函数调用模式的网络化。客户端就像调用本地函数一样，然后客户端把这些参数打包之后通过网络传递到服务端，服务端解包到处理过程中执行，然后执行的结果反馈给客户端。

## RPC工作原理

![img](https://github.com/astaxie/build-web-application-with-golang/raw/master/zh/images/8.4.rpc.png?raw=true)

图8.8 RPC工作流程图

运行时,一次客户机对服务器的RPC调用,其内部操作大致有如下十步：

- 1.调用客户端句柄；执行传送参数
- 2.调用本地系统内核发送网络消息
- 3.消息传送到远程主机
- 4.服务器句柄得到消息并取得参数
- 5.执行远程过程
- 6.执行的过程将结果返回服务器句柄
- 7.服务器句柄返回结果，调用远程系统内核
- 8.消息传回本地主机
- 9.客户句柄由内核接收消息
- 10.客户接收句柄返回的数据

## Go RPC

Go标准包中已经提供了对RPC的支持，而且支持三个级别的RPC：TCP、HTTP、JSONRPC。但Go的RPC包是独一无二的RPC，它和传统的RPC系统不同，它只支持Go开发的服务器与客户端之间的交互，因为在内部，它们采用了Gob来编码。

Go RPC的函数只有符合下面的条件才能被远程访问，不然会被忽略，详细的要求如下：

- 函数必须是导出的(首字母大写)
- 必须有两个导出类型的参数，
- 第一个参数是接收的参数，第二个参数是返回给客户端的参数，第二个参数必须是指针类型的
- 函数还要有一个返回值error

举个例子，正确的RPC函数格式如下：

```
func (t *T) MethodName(argType T1, replyType *T2) error
```

T、T1和T2类型必须能被`encoding/gob`包编解码。

- type.go 封装client server公用的变量, 建议放到一个包里

```go
package main

import "errors"

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Math int
func (m *Math)Multiply(args *Args, reply *int) error{
	*reply = args.A * args.B
	return nil
}

func (m *Math)Divide(args *Args, quo *Quotient)error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

```

- server.go

  ```go
  package main
  
  import (
  	"log"
  	"net"
  	"net/rpc"
  	"net/rpc/jsonrpc"
  	"os"
  	"time"
  )
  
  func main() {
  	arith := new(Math)
  	rpc.Register(arith)
  	// http
  	//rpc.HandleHTTP()
  	//log.Fatal(http.ListenAndServe(":12345", nil))
  
  
  	//tcp
  	//tcpAddr, e  := net.ResolveTCPAddr("tcp", ":12345")
  	//if e != nil {
  	//	log.Println("Resolve Address error:", e)
  	//	os.Exit(2)
  	//}
  	//listen, e := net.ListenTCP("tcp",tcpAddr)
  	//if e != nil {
  	//	log.Println("starting RPC-server -listen error:", e)
  	//	os.Exit(2)
  	//}
  	//go func() {for {
  	//	conn, e := listen.Accept()
  	//	if e != nil {
  	//		log.Println("establish error:", e, "will wait seconds")
  	//		time.Sleep(time.Millisecond * 5)
  	//		continue
  	//	}
  	//	rpc.ServeConn(conn)
  	//
  	//}}()
  	//select {
  	//}
  
  	// json
  	tcpAddr, e  := net.ResolveTCPAddr("tcp", ":12345")
  	if e != nil {
  		log.Println("Resolve Address error:", e)
  		os.Exit(2)
  	}
  	listen, e := net.ListenTCP("tcp",tcpAddr)
  	if e != nil {
  		log.Println("starting RPC-server -listen error:", e)
  		os.Exit(2)
  	}
  
  	for {
  		conn, e := listen.Accept()
  		if e != nil {
  			log.Println("establish error:", e, "will wait seconds")
  			time.Sleep(time.Millisecond * 5)
  			continue
  		}
  		jsonrpc.ServeConn(conn)
  
  	}
  }
  
  
  ```

  

- client.go

  ```go
  package main
  
  import (
  	"log"
  	"net/rpc/jsonrpc"
  )
  
  func main() {
  	// http
  	//client,e := rpc.DialHTTP("tcp","localhost:12345")
  	// tcp
  	//client,e := rpc.Dial("tcp","localhost:12345")
  	// json
  	client,e := jsonrpc.Dial("tcp","localhost:12345")
  	if e != nil {
  		log.Fatalf("starting RPC-client error %s", e.Error())
  	}
  	args := &Args{10,8}
  	//clientArgs := new(Args)
  	//clientArgs.B = 5
  	//clientArgs.A = 6
  	var reply int
  	// 同步
  	//client.Call("Args.Multiply", args, &reply )
  	// 异步 如果done==nil 会返回一个新的chan
  	call1 := client.Go("Math.Multiply", args, &reply, nil)
  	<- call1.Done // 这里不知道为什么必须要接受了 Error 才有值
  	if call1.Error != nil {
  		log.Println(call1.Error.Error())
  	}
  	//log.Println(<-call1.Done)
  	log.Printf("Args: %d * %d = %d\n", args.A, args.B, reply)
  
  	quo := &Quotient{}
  	client.Call("Math.Divide",args, quo)
  	log.Printf("Args: %d / %d = %d, Remainer%d\n", args.A, args.B, quo.Quo, quo.Rem)
  }
  ```

  