 在 Go http包的Server中，每一个请求在都有一个对应的 goroutine 去处理。请求处理函数通常会启动额外的 goroutine 用来访问后端服务，比如数据库和RPC服务。用来处理一个请求的 goroutine 通常需要访问一些与请求特定的数据，比如终端用户的身份认证信息、验证相关的token、请求的截止时间。 当一个请求被取消或超时时，所有用来处理该请求的 goroutine 都应该迅速退出，然后系统才能释放这些 goroutine 占用的资源。 

## Context 两个变量,4个函数

1. ### 全局变量方式

2. ### 通道方式

3. ### 官方版的方案 import "context"

   ```go
   	ctx, cancel := context.WithCancel(context.Background())
   	cancel() // 通知子goroutine结束
   	case <-ctx.Done(): // 等待上级通知
   			break LOOP
   ```

   

   ```go
   var wg sync.WaitGroup
   
   func worker(ctx context.Context) {
   	go worker2(ctx)
   LOOP:
   	for {
   		fmt.Println("worker")
   		time.Sleep(time.Second)
   		select {
   		case <-ctx.Done(): // 等待上级通知
   			break LOOP
   		default:
   		}
   	}
   	wg.Done()
   }
   
   func worker2(ctx context.Context) {
   LOOP:
   	for {
   		fmt.Println("worker2")
   		time.Sleep(time.Second)
   		select {
   		case <-ctx.Done(): // 等待上级通知
   			break LOOP
   		default:
   		}
   	}
   }
   func main() {
   	ctx, cancel := context.WithCancel(context.Background())
   	wg.Add(1)
   	go worker(ctx)
   	time.Sleep(time.Second * 3)
   	cancel() // 通知子goroutine结束
   	wg.Wait()
   	fmt.Println("over")
   }
   ```

## Background()和TODO()

   Go内置两个函数：`Background()`和`TODO()`，这两个函数分别返回一个实现了`Context`接口的`background`和`todo`。我们代码中最开始都是以这两个内置的上下文对象作为最顶层的`partent context`，衍生出更多的子上下文对象。

   `Background()`主要用于main函数、初始化以及测试代码中，作为Context这个树结构的最顶层的Context，也就是根Context。

   `TODO()`，它目前还不知道具体的使用场景，如果我们不知道该使用什么Context的时候，可以使用这个。

   `background`和`todo`本质上都是`emptyCtx`结构体类型，是一个不可取消，没有设置截止时间，没有携带任何值的Context。

## With系列函数

此外，`context`包中还定义了四个With系列函数。

### WithCancel

`WithCancel`的函数签名如下：

```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
```

`WithCancel`返回带有新Done通道的父节点的副本。当调用返回的cancel函数或当关闭父上下文的Done通道时，将关闭返回上下文的Done通道，无论先发生什么情况。

取消此上下文将释放与其关联的资源，因此代码应该在此上下文中运行的操作完成后立即调用cancel。

### WithDeadline

`WithDeadline`的函数签名如下：

```go
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
```

返回父上下文的副本，并将deadline调整为不迟于d。如果父上下文的deadline已经早于d，则WithDeadline(parent, d)在语义上等同于父上下文。当截止日过期时，当调用返回的cancel函数时，或者当父上下文的Done通道关闭时，返回上下文的Done通道将被关闭，以最先发生的情况为准。

取消此上下文将释放与其关联的资源，因此代码应该在此上下文中运行的操作完成后立即调用cancel。

### WithTimeout

`WithTimeout`的函数签名如下：

```go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```

`WithTimeout`等于是`WithDeadline(parent, time.Now().Add(timeout))`。

code应该在上下文中运行的操作完成后立即调用cancel,  释放与其相关的资源.  通常用于数据库或者网络连接的超时控制。具体示例如下：

```go
func slowOperationWithTimeout(ctx context.Context) (Result, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
    defer cancel()  // releases resources if slowOperation completes before timeout elapses(流逝) 释放资源 如果slowOperation在时间耗尽前完成
	return slowOperation(ctx)
}
```

下面这个例子 通过context的timeout来通知 在时间耗尽后, 被调用的函数要抛弃它的工作.

```go
// Pass a context with a timeout to tell a blocking function that it
// should abandon its work after the timeout elapses.
ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
defer cancel()
select {
case <-time.After(1 * time.Second):// 等待1秒后
	fmt.Println("overslept")
case <-ctx.Done():
	fmt.Println(ctx.Err()) // prints "context deadline exceeded"
}
/*
context deadline exceeded  // 50 调大 就会输出 overslept
*/
```

### WithValue

`WithValue`函数能够将请求作用域的数据与 Context 对象建立关系。声明如下：

```go
func WithValue(parent Context, key, val interface{}) Context
```

`WithValue`返回父节点的副本，其中与key关联的值为val。

仅对API和进程间传递请求域的数据使用上下文值，而不是使用它来传递可选参数给函数。

所提供的key必须是可比较的，并且不应该是`string`类型或任何其他内置类型，以避免使用上下文在包之间发生冲突。`WithValue`的用户应该为键定义自己的类型。为了避免在赋值给interface{}时进行分配，context keys应该有 具体的类型`struct{}`。或者，导出的context.key变量的静态类型需是**pointer或interface**。

下面这个列子展示了 一个值如何被传递 到 context 和 当它存在时,如何去获取它,

```go
type favContextKey string

f := func(ctx context.Context, k favContextKey) {
	if v := ctx.Value(k); v != nil {
		fmt.Println("found value:", v)
		return
	}
	fmt.Println("key not found:", k)
}

k := favContextKey("language")
ctx := context.WithValue(context.Background(), k, "Go")

f(ctx, k)
f(ctx, favContextKey("color"))
/*
found value: Go
key not found: color
*/
```

## 使用Context的注意事项

- 推荐以参数的方式显示传递Context
- 以Context作为参数的函数方法，应该把Context作为第一个参数。
- 给一个函数方法传递Context的时候，不要传递nil，如果不知道传递什么，就使用context.TODO()
- Context的Value相关方法应该传递请求域的必要数据，不应该用于传递可选参数
- Context是线程安全的，可以放心的在多个goroutine中传递