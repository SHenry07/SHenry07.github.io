runtime.Gosched() // 出让当前goroutine所占用的cpu时间片, 当再次获得cpu时,从出让位置继续恢复运行

runtime.Goexit // 退出当前goroutine

runtime.GOMAXPROCS() // 设置可以用来计算的CPU核数的最大值,并返回之前的值

