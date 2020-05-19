

首 先通过kubectl label命令给目标Node打上一些标签：

`kubectl label nodes <node-name> <label-key>=<label-value>`
这里，我们为k8s-node1节点打上一个zone=north标签，表明它是一个“北方”的一个节点

`kubectl label nodes k8s-node1 zone=north`



# NodeSeletor

除了用户可以自行给Node添加标签，Kubernetes也会给Node预定义一些标签，包括：

- kubernetes.io/hostname
- beta.kubernetes.io/os（从1.14版本开始更新为稳定版，到1.18版本删除）
- beta.kubernetes.io/arch（从1.14版本开始更新为稳定版，到1.18版本删除）
- kubernetes.io/os（从1.14版本开始启用）
- kubernetes.io/arch（从1.14版本开始启用）

# NodeAffinity亲和性调度 // NodeSeletor的未来

​	NodeAffinity意为Node亲和性的调度策略，是用于替换NodeSelector的全新调度策略。目前有两种节点亲和性表达。

- RequiredDuringSchedulingIgnoredDuringExecution：必须满足指定的规则才可以调度Pod到Node上（功能与nodeSelector很像，但是使用的是不同的语法），相当于硬限制。

- PreferredDuringSchedulingIgnoredDuringExecution：强调优先满足指定规则，调度器会尝试调度Pod到Node上，但并不强求，相当于软限制。多个优先级规则还可以设置权重（weight）值，以定义执行的先后顺序。



​	IgnoredDuringExecution的意思是：如果一个Pod所在的节点在Pod运行期间标签发生了变更，不再符合该Pod的节点亲和性需求，则系统将忽略Node上Label的变化，该Pod能继续在该节点运行。

下面的例子设置了NodeAffinity调度的如下规则。

- requiredDuringSchedulingIgnoredDuringExecution要求只运行在amd64的节点上（beta.kubernetes.io/arch In amd64)
- preferredDuringSchedulingIgnoredDuringExecution的要求是尽量运行在磁盘类型为ssd（disk-type In ssd）的节点上。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: with-node-affinity
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:// 
        nodeSelectorTerms:
        - matchExpressions:
          - key: beta.kubernetes.io/arch
            operator: In
            values:
            - amd64
      preferredDuringSchedulingIgnoredDuringExecution:// 
      - weight: 1
        preference:
          matchExpressions:
          - key: disk-type
            operator: In
            values:
            - ssd
  containers:
  - name: with-node-affinity
    image: tomcat
```

   从上面的配置中可以看到In操作符，NodeAffinity语法支持的操作符包括`In、NotIn、Exists、DoesNotExist、Gt、Lt`。虽然没有节点排斥功能，但是用NotIn和DoesNotExist就可以实现排斥的功能了。

NodeAffinity规则设置的注意事项如下：

- 如果同时定义了nodeSelector和nodeAffinity，那么必须两个条件都得到满足，Pod才能最终运行在指定的Node上。
- 如果nodeAffinity指定了多个nodeSelectorTerms，那么其中一个能够匹配成功即可。
- 如果在nodeSelectorTerms中有多个matchExpressions，则一个节点必须满足所有matchExpressions才能运行该Pod。

# PodAffinity：Pod亲和与互斥调度策略
​	Pod间的亲和与互斥从Kubernetes 1.4版本开始引入。这一功能让用户从另一个角度来限制Pod所能运行的节点：根据在节点上正在运行的Pod的标签而不是节点的标签进行判断和调度，要求对节点和Pod两个条件进行匹配。这种规则可以描述为：如果在具有标签X的Node上运行了一个或者多个符合条件Y的Pod，那么Pod应该（如果是互斥的情况，那么就变成拒绝）运行在这个Node上。
​	这里X指的是一个集群中的节点、机架、区域等概念，通过Kubernetes内置节点标签中的key来进行声明。这个key的名字为topologyKey，意为表达节点所属的topology范围。

- kubernetes.io/hostname

- failure-domain.beta.kubernetes.io/zone

- failure-domain.beta.kubernetes.io/region
	



​	与节点不同的是，Pod是属于某个命名空间的，所以条件Y表达的是一个或者全部命名空间中的一个Label Selecor。
​	和节点亲和相同，Pod亲和与互斥的条件设置也是`requiredDuringSchedulingIgnoredDuringExecution`和`preferredDuringSchedulingIgnoredDuringExecution`。Pod的亲和性被定义于PodSpec的affinity字段下的podAffinity子字段中。Pod间的互斥性则被定义于同一层次的podAntiAffinity子字段中。Pod间的互斥性则被定义与同一层次的PodAntiAffinity字段中。



​	下面通过实例来说明Pod间的亲和性和互斥性策略设置。

1.参照目标Pod

​	首先，创建一个名为pod-flag的Pod，带有标签security=S1和app=nginx，后面的例子将使用pod-flag作为Pod亲和互斥的目标Pod：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod-flag
  labels:
    security: "S1"
    app: "nginx"
spec:
  containers:

  - name: nginx
    image: nginx
```

2.Pod 的亲和性调度
     下面创建第2个Pod来说明Pod的亲和性调度，这里定义的亲和标签是security=S1，对应上面的Pod“pod-flag”，topologyKey的值被设置为“kubernetes.io/hostname”

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod-affinity
spec:
  affinity://
    podAffinity://
      requiredDuringSchedulingIgnoredDuringExecution://
      - labelSelector:
          matchExpressions:
          - key: security
            operator: In
            values:
            - S1
        topologyKey: kubernets.io/hostname  // 
  containers:
  - name: with-pod-affinity
    image: tomcat
```
创建Pod之后，使用kubectl get pods -o wide命令可以看到，这两个Pod在同一个Node上运行。
有兴趣的读者还可以测试一下，在创建这个Pod之前，删掉这个节点的kubernetes.io/hostname标签，重复上面的创建步骤，将会发现Pod一直处于Pending状态，这是因为找不到满足条件的Node了
3.Pod的互斥性调度

创建第3个Pod，我们希望它不与目标Pod运行在同一个Node上：
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: anti-affinity
spec:
  affinity:
    podAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchExpressions:
          - key: security
            operator: In
            values:
            - S1
        topologyKey: failure-domain.beta.kubernetes.io/zone
        podAntiAffinity:
            requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
          matchExpressions:
          - key: app
            operator: In
            values:
            - nginx
        topologyKey: kubernetes.io/hostname
    containers:
    - name: with-pod-affinity
      image: tomcat
```
​	这里要求这个新Pod与security=S1的Pod为同一个zone，但是不与app=nginx的Pod为同一个Node。创建Pod之后，同样用`kubectl get pods -o wide`来查看，会看到新的Pod被调度到了同一Zone内的不同Node上。（ps:由于我们security=S1只有一个node节点，并且上面运行着niginx，所以不能找到第二个security=S1的node，显然我们的条件不够，导致Pod处于Pending状态。）



​	与节点亲和性类似，Pod亲和性的操作符也包括`In、NotIn、Exists、DpesNotExist、Gt、Lt`
​     原则上，topologyKey可以使用任何合法的标签Key赋值，但是出于性能和安全方面的考虑，对topologyKey有如下限制。

- 在Pod亲和性和RequiredDuringScheduling的Pod互斥性的定义中，不允许使用空的topologyKey。

- 如果Admission controller包含了LimitPodHardAntiAffinityTopology，那么针对Required DuringScheduling的Pod互斥性定义就被限制为kubernetes.io/hostname，要使用自定义的topologyKey，就要改写或禁用该控制器。

- 在PreferredDuringScheduling类型的Pod互斥性定义中，空的topologyKey会被解释为kubernetes.io/hostname、failure-domain.beta.kubernetes.io/zone及failuredomain.beta.kubernetes.io/region的组合。

- 如果不是上述情况，就可以采用任意合法的topologyKey了。

  

PodAffinity规则设置的注意事项如下。

- 除了设置Label Selector和topologyKey，用户还可以指定Namespace列表来进行限制，同样，使用Label Selector对Namespace进行选择。Namespace的定义和Label Selector及topologyKey同级。省略Namespace的设置，表示使用定义了affinity/anti-affinity的Pod所在的Namespace。如果Namespace被设置为空值（""），则表示所有Namespace。

- 在所有关联requiredDuringSchedulingIgnoredDuringExecution的matchExpressions全都满足之后，系统才能将Pod调度到某个Node上。

关于Pod亲和性和互斥性调度的更多信息可以参考其设计文档，网址为https://github.com/kubernetes/kubernetes/blob/master/docs/design/podaffinity.md。




