# time.Timer

Timer是一个定时器, 代表未来的一个单一事件,你可以告诉timer你要等待多长时间

```go
type Timer struct {
   C <-chan Time
   r runtimeTimer
}
```

提供一个channel, 在定时时间到达之前, 没有数据写入Timer.C会一直阻塞, 直到定时时间到, 系统会自动向timer.C 这个channel中写入当前时间,阻塞即被解除

# time.Ticker 周期性定时

```go
type Ticker struct {
   C <-chan Time // 'ticks' was delivered 滴答
   r runtimeTimer
}
```

