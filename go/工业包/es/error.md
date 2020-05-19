# filebeat7.6.1修改索引名字后elasticsearch中没有生成新索引

filebeat7.6.1修改索引名字后，比如下面这样，shop-api xxx

```
output.elasticsearch:
  hosts: ["http://192.168.0.10:9200"]
  index: "shop-api-%{[agent.version]}-%{+yyyy.MM.dd}"
setup.template.name: "shop-api"
setup.template.pattern: "shop-api-*"
```

启动后，发现在es中并没有生成新的索引。
filebeat有打印下面的日志信息

```
2020-05-07T02:40:53.010Z    INFO    [index-management]  idxmgmt/std.go:258  Auto ILM enable success.
2020-05-07T02:40:53.012Z    INFO    [index-management.ilm]  ilm/std.go:139  do not generate ilm policy: exists=true, overwrite=false
2020-05-07T02:40:53.012Z    INFO    [index-management]  idxmgmt/std.go:271  ILM policy successfully loaded.
2020-05-07T02:40:53.012Z    INFO    [index-management]  idxmgmt/std.go:410  Set setup.template.name to '{filebeat-7.6.1 {now/d}-000001}' as ILM is enabled.
```

提示开启了ILM策略

翻官方文档（https://www.elastic.co/guide/en/beats/filebeat/current/elasticsearch-output.html）后发现：
index配置部分中提示 The index setting is ignored when index lifecycle management is enabled
意思就是index设置的参数在索引生命周期管理（ilm）开启后会忽略。

查看ilm文档 https://www.elastic.co/guide/en/beats/filebeat/current/ilm.html 提示：
Starting with version 7.0, Filebeat uses index lifecycle management by default when it connects to a cluster that supports lifecycle management
从7.0版本开始，当elasticsearch支持生命周期管理时，filebeat默认使用索引生命周期管理，这样就导致自己修改的日志文件名无效了。

关闭ilm功能即可(setup.ilm.enabled: false)。

```yaml
filebeat.inputs:
- type: log
  paths:
   - /mnt/logs/*.log
  fields:
   java: true
  fields_under_root: true
  multiline.pattern: '^[0-9]{2}:[0-9]{2}:[0-9]{2}.* \[http-nio'
  multiline.negate: true
  multiline.match: after
 
setup.ilm.enabled: false
output.elasticsearch:
  hosts: ["http://192.168.0.10:9200"]
  index: "shop-api-%{[agent.version]}-%{+yyyy.MM.dd}"
setup.template.name: "shop-api"
setup.template.pattern: "shop-api-*"
```