```txt
active_anon + inactive_anon = anonymous memory + file cache for tmpfs + swap cache = rss + file cache for tmpfs 
active_file + inactive_file = cache - size of tmpfs
working set = usage - total_inactive(k8s根据workingset 来判断是否驱逐pod)
```

https://www.orchome.com/6745