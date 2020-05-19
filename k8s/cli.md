1. 删除命令行,一直等待
    `kubectl delete pods <pod> --grace-period=0 --force`
    `kubectl delete pod --force --grace-period=0 --wait=false some_pod`
    [github介绍](https://github.com/kubernetes/kubernetes/issues/66478)

2. ```shell
   kubectl apply -f https://k8s.io/examples/application/deployment.yaml --record
   ```

3. 获取端口

     ```export NODE_PORT=$(kubectl get services/kubernetes-bootcamp -o go-template='{{(index .spec.ports 0).nodePort}}') echo NODE_PORT=$NODE_PORT```