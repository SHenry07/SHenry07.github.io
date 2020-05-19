KMP算法:

模式字符串才是问题所在 ，看起来很像

前缀[i]

后缀[j]

NEXT数组: 当模式匹配串T失配的时候,NEXT数组对应的元素指导应该用T串的那个元素惊醒下一轮的匹配

找出最长公共前缀后缀的长度，然后右移得到next数组

```txt
// 伪代码
i := 1
j := 0
next[1] = 0 
for i < T[0] {
	if 0 == j || T[i] == T[j] {
 			 i ++ 
  		j ++
 		 next[i] = j
	}else {
	  // 回溯
	  j = next[j]
	}
}
// 前缀是固定 后缀是不固定的
```

```go
func getNext(T string) (next []int) {
	pLen := len(T)
	next = make([]int, pLen)
	i := 1
	j := 0
	// next[0] = len(T)
	next[0] = -1
	for i < pLen-1 { // 这里是因为go 语言的字符串index0保存的是真实的数据不和c一样保存的是长度
		if 0 == j || T[i] == T[j] {
			i++
			j++
			// 优化算法
			if T[i] != T[j] { // T 到3就结束了
				next[i] = j
			} else {
				next[i] = next[j]
			}
			/*
			 0   1  2  3
			 -1  0  1  2
			*/
		} else {
			j = next[j]
		}
	}
	return
}
```



