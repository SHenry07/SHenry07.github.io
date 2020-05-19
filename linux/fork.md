fork函数通过系统调用创建一个与原来进程几乎完全相同的进程，一个进程调用fork函数后，系统先给新的进程分配资源，例如存储数据和代码的空间。

- 在父进程中，fork返回新创建子进程的进程ID
- 在子进程中，fork返回0
- 如果出现错误，fork返回一个负值

> 其实就相当于链表，进程形成了链表，父进程的fpid(p 意味point)指向子进程的进程id, 因为子进程没有子进程，所以其fpid为0.



[https://de4dcr0w.github.io/%E5%88%86%E6%9E%90fork%E7%B3%BB%E7%BB%9F%E8%B0%83%E7%94%A8.html](https://de4dcr0w.github.io/分析fork系统调用.html)

https://blog.csdn.net/jason314/article/details/5640969