IPVS provides more options for balancing traffic to backend Pods; these are:

- `rr`: round-robin 轮询
- `lc`: least connection (smallest number of open connections) SessionAffinity 基于客户端IP地址进行会话保持的模式
- `dh`: destination hashing
- `sh`: source hashing
- `sed`: shortest expected delay
- `nq`: never queue

# 无Label的service

一个不带标签选择器的service,即无法选择后端的pod,系统不会自动创建Endpoint,

需要手动创建一个和该service**同名**的endpoint,用于指向实际的后端访问地址,

比如: 集群外的mysql服务器, 另外一个集群或namespace内的服务

# Headless Service

不为Service设置**ClusterIP(入口IP地址),**仅通过Label Selector将后端的Pod列表返回给调用的客户端

但是要建立一个seed provider,来查询service的值

1. 通过环境变量
2. DNS // 使用rest-api去API-server中查找pod列表

# Node Port

