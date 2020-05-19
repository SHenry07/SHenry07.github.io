# kind

- object
- list
- simple

# apiVersion

# Metadata

- namespace

- name

- uid

- labels

- annotations

  用户自定义的'注解', key value都是string的map,被kubernetes的内部进程或者某些外部工具使用,用于存储和获取关于该对象的特定元数据

- resourceVersion

- creatTimestamp

- deleteTimestamp

- selfLink

  通过API访问资源滋生的URL

# spec

# Status

Status用于纪录对象,在系统中的当前状态信息, 他也是集合类元素类型,status在一个自动处理的进程中被持久化,可以在流转的过程中生成,如果观察到一个资源丢失了它的status,则该丢失的状态可能被重新构造,以pod为例,Pod的status信息主要包括conditions, containerStatuses, hostIP, phase, podIP, startTime等, 其中比较重要的两个状态属性如下:

1. phase: 描述对象所处的生命周期阶段, pharse的典型值是Pending, Running, Active(正在运行中)或Terminated(已终结)

2. condition: 表示条件, 由条件类型和状态值组成, 目前仅有一种条件类型: Ready

   对应的值: True, False, Unkonwn, 