# CRUD

# Proxy

主要用于管理目的, 

- 逐一排查Service的pod副本,
- 检查哪些pod的服务存在异常

```
/api/v1/nodes/{name}/proxy/pods
/api/v1/nodes/{name}/proxy/stats # 列出指定结点内物理资源的统计信息
/api/v1/nodes/{name}/proxy/spec #列出指定结点的概要信息
```

## --enable-dubuggins-handlers=true

`curl http://lcoalhost:8080/api/v1/namespaces/default/pods/myweb-g9pmm/prxoy` # 访问pod

`curl http://lcoalhost:8080/api/v1/namespaces/default/pods/myweb-g9pmm/prxoy/{path:*}` # 访问pod服务的URL路径

就能直接访问特定的pod

## Service也有自己的proxy

