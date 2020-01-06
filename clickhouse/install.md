# 检查环境

1.  `grep sse4_2 /proc/cpuinfo |`
   `grep -q sse4_2 /proc/cpuinfo && echo "SSE 4.2 supported" || echo "SSE 4.2 not supported"`
>    use a CPU with x86_64 architecture and support for SSE 4.2 instructions. To run ClickHouse with processors that do not support SSE 4.2 or have AArch64 or PowerPC64LE architecture, you should build ClickHouse from sources. // 不支持SSE4.2的请从源码编译

2. `swapoff -a`
>  Disable the swap file for production environments.
同时去关报警

3. Huge Pages
`echo 'never' | sudo tee /sys/kernel/mm/transparent_hugepage/enabled`
>Use `perf top` to watch the time spent in the kernel for memory management. Permanent huge pages also do not need to be allocated.


# install
[官方安装]( https://clickhouse.yandex/docs/en/getting_started/install/ )

# 配置文件

按官方要求修改配置文件

rpm装的 请修改/var/lib/clickhouse 为/apply 注意用户权限为clickhouse