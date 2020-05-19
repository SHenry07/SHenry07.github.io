# filebeat直接抽取日志
```txt
'{"@timestamp":"$time_iso8601",'
                      '"@source":"$server_addr",'
                      '"@nginx_fields":{'
                      '"remote_addr":"$remote_addr",'
                      '"remote_user":"$remote_user",'
                      '"body_bytes_sent":"$body_bytes_sent",'
                      '"request_time":"$request_time",'
                      '"status":"$status",'
                      '"host":"$host",'
                      '"uri":"$uri",'
                      '"server":"$server_name",'
                      '"port":"$server_port",'
                      '"protocol":"$server_protocol",'
                      '"request_uri":"$request_uri",'
                      '"request_body":"$request_body",'
                      '"request_method":"$request_method",'
                      '"http_referrer":"$http_referer",'
                      '"body_bytes_sent":"$body_bytes_sent",'
                      '"http_x_forwarded_for":"$http_x_forwarded_for",'
                      '"http_user_agent":"$http_user_agent",'
                      '"upstream_response_time":"$upstream_response_time",'
                      '"upstream_addr":"$upstream_addr"}},';
```

所有的可以引用的参数: http://nginx.org/en/docs/varindex.html