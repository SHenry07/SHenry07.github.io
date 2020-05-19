## 我的容器被终止了

您的容器可能因为资源枯竭而被终止了。要查看容器是否因为遇到资源限制而被杀死，请在相关的 Pod 上调用 `kubectl describe pod`：

`kubectl get pod -o go-template='{{range.status.containerStatuses}}{{"Container Name: "}}{{.name}}{{"\r\nLastState: "}}{{.lastState}}{{end}}'  容器名字 -n prod`

# 容器处于invalid
创建的容器有个创建的时长的选项age，测试环境发现有的容器时间是invalid，通过代码查看原来是时间不对
因为每个kubelet是定时上报的，上报是基于本地当前时间的。

`wildfly-rc-1l9qv   0/1       ContainerCreating   0          <invalid>`
解决办法，在每个机器安装ntp服务。保持每个节点时间同步

# 容器处于Evicted

处理办法: `kubectl get pods --all-namespaces -ojson | jq -r '.items[] | select(.status.reason!=null) | select(.status.reason | contains("Evicted")) | .metadata.name + " " + .metadata.namespace' | xargs -n2 -l bash -c 'kubectl delete pods $0 --namespace=$1'`



```bash
#!/bin/bash
kubectl get pods --all-namespaces -o go-template='{{range .items}} \
{{if eq .status.phase "Failed"}} {{if eq .status.reason "Evicted"}} {{.metadata.name}}{{" "}} {{.metadata.namespace}}{{"\n"}} \
{{end}} \
{{end}} \
{{end}}' | while read epod namespace; do kubectl -n $namespace delete pod $epod; done
```



结局办法: 

```txt
I suppose, this issue can be closed, because the evicted pods deletion can be controlled through settings in kube-controller-manager.

For those k8s users who hit the kube-apiserver or etcd performance issues due to too many evicted pods, i would recommend updating the kube-controller-manager config to set --terminated-pod-gc-threshold 100 or similar small value. The default GC threshold is 12500, which is too high for most etcd installations. Reading 12500 pod records from etcd takes seconds to complete.

Also ask yourself why are there so many evicted pods? Maybe your kube-scheduler keeps scheduling pods on a node which already reports DiskPressure or MemoryPressure? This could be the case if the kube-scheduler is configured with a custom --policy-config-file which has no CheckNodeMemoryPressure or CheckNodeDiskPressure in the list of policy predicates.
```

```sh
$ kube-controller-manager --help 2>&1|grep terminated
      --terminated-pod-gc-threshold int32                                 Number of terminated pods that can exist before the terminated pod garbage collector starts deleting terminated pods. If <= 0, the terminated pod garbage collector is disabled. (default 12500)
```

`kubectl get pods -n test | grep -v Running | awk '{print $1}' | xargs kubectl delete pod --grace-period=0 --force -n test`

`kubectl get pods --no-headers -o custom-columns=":metadata.name" `

## 重写label

`kubectl label --overwrite pods $(kubectl get pod -n test -l app=banksteel-assistant-spi -o=custom-columns=':.metadata.name') app=banksteel-assistant-spi-bak -n test`



- 如果下载速度比较慢，可以使用aliyun的镜像：

  ```
  java -jar arthas-boot.jar --repo-mirror aliyun --use-http
  ```

- 如果从github下载有问题，可以使用gitee镜像

  ```
  curl -O https://arthas.gitee.io/arthas-boot.jar
  ```