![20170315101923](20170315101923.jpg)



`kubernetes List-Watch`用于实现数据同步的代码逻辑: 客户端在初始化的时，先调用`Kubernetes List API`获得某种`resource`的全部`Object`，缓存在`内存`中; 然后启动对应`resource`对象的Watch协程去维护这份缓存, 在接收到Watch时间后,再根据事件的类型(比如新增,修改,删除)对内存中的全量资源对象列表做出相对应的同步修改; 最后，客户端就不再调用`Kubernetes`的任何 API。

- 为了让kubernetes中的其他组件在不访问底层etcd数据库的情况下,也能及时获取资源对象的变化,API Server模仿etcd的Watch API接口提供了自己的Watch接口

https://zhuanlan.zhihu.com/p/59660536

https://yq.aliyun.com/articles/679797

