1. 引入一个长度为0的chan来 阻止main退出, 还见过空的select{} 来阻塞

2.  初始化 klog 日志框架

   ```go
   import "k8s.io/klog"
   klog.InitFlags(nil)
   defer klog.Flush()
   ```

3. 解析变量

4. 设定CPU

5. 建立sourceFactory 

   解析--source="kubernetes:https://kubernetes.default"

   和所有k8s api建连
   
6. 建立sink

   现在的做法是: dingtalk, kafka, etc各自一个启动器

   正在进行中的是webhook

7. 建立manage来管理 source 和sink

8. chan 堵塞

总结: 每个source, sink 都是一个工厂模式, 

最后一个manage管理一对source, sink



---

编程技巧

1. 建立空 struct来 确定method 来 实现interface



# net/url

```go
type Error struct {
	Op  string
	URL string
	Err error
}

func (e *Error) Unwrap() error { return e.Err } // 未封装的
func (e *Error) Error() string { return fmt.Sprintf("%s %q: %s", e.Op, e.URL, e.Err) }
```

- 建立一个Error的结构体来处理错误

  ```go
  type Error struct {
  	Op  string
  	URL string
  	Err error
  }
  
  func (e *Error) Unwrap() error { return e.Err }
  func (e *Error) Error() string { return fmt.Sprintf("%s %q: %s", e.Op, e.URL, e.Err) }
  
  func (e *Error) Timeout() bool {
  	t, ok := e.Err.(interface {  // 判断是否实现了接口
  		Timeout() bool
  	})
  	return ok && t.Timeout()
  }
  
  func (e *Error) Temporary() bool {
  	t, ok := e.Err.(interface {
  		Temporary() bool
  	})
  	return ok && t.Temporary()
  }
  ```

  

- 工厂模式