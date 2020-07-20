```
# 建立临时文件
$ mktemp  
/tmp/tmp.QXXgqELfD7
```

# docker

```
--link 共享网络空间
--cap-add=NET_ADMIN # 限定上下文
--dns 8.8.8.8
--rm
–platform arm
 manifest inspect --verbose 镜像名
```

# LVM

https://www.cnblogs.com/blog-lhong/p/11712069.html
扩展完LV后， 要手动扩容

-  运行扩容命令，对/dev/vg_template/lv_root逻辑卷进行扩容
` lvextend -l +100%FREE /dev/mapper/vg_template-lv_root`

-  查看分区格式，如果分区格式是ext,用resize2fs扩容；如果分区格式是xfs，用xfs_growfs扩容。命令如下`df -Th`

- 然后对容量重新扩容。
  `resize2fs /dev/mapper/vg_template-lv_root`

  xfs格式的分区用xfs_growfs 命令对容量重新扩容。
  `xfs_growfs /dev/mapper/centos-root`

