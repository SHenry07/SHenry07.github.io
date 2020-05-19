

# 翻转字符串

```go
func reverseWords(s string) string {
	// 分割字符串
    seg := strings.Fields(s)
	
    // 反转字符串数组
	var reserseSeg []string
    for i := len(seg) - 1; i >= 0; i-- {
		reserseSeg = append(reserseSeg, seg[i])
	}
    
    // 将字符串数组里的元素拼接成一个字符串并返回
	return  strings.Join(reserseSeg, " ")
}

func main() {
	fmt.Println(reverseWords("the sky is blue"))
	fmt.Println(reverseWords("   hello world!"))
}
```

# builder

## 误用字符串

[ Go 入门指南](https://learnku.com/docs/the-way-to-go) /

当需要对一个字符串进行频繁的操作时，谨记在 go 语言中字符串是不可变的（类似 java 和 c#）。使用诸如 `a += b` 形式连接字符串效率低下，尤其在一个循环内部使用这种形式。这会导致大量的内存开销和拷贝。**应该使用一个字符数组代替字符串，将字符串内容写入一个缓存中。** 例如以下的代码示例：

```go
var b bytes.Buffer
...
for condition {
    b.WriteString(str) // 将字符串str写入缓存buffer
}
    return b.String()
```

注意：由于编译优化和依赖于使用缓存操作的字符串大小，**当循环次数大于 15 时，效率才会更佳。**

## 拼接字符串的三种方式

```go
func compressString(S string) string {
    b := []byte(S)
    pLen :=  len(b)
    if 0 == pLen {
        return ""
    }
    var slice string
    ch := b[0]
    ans := 1
    for i := 1 ;i < pLen; i++{
       if ch == b[i] {
          ans += 1
       }else{
           // 可以用builder优化
          slice = slice + string(ch)
          slice = slice + strconv.Itoa(ans) 
          ans = 1
          ch = b[i]
       }
    }
    slice = slice + string(b[pLen-1])
    slice = slice + strconv.Itoa(ans)
    if len(slice) >= len(b) {
        return S
    }else{
        return string(slice)
    }
}
```
```go
func compressString(S string) string {
    b := []byte(S)
    pLen :=  len(b)
    if 0 == pLen {
        return ""
    }   
    var slice []byte 
    ch := b[0]
    ans := 1
    for  i := 1 ;i < pLen; i++{
       if ch == b[i] {
          ans += 1
       }else{
          slice = append(slice,ch)
          slice = append(slice,strconv.Itoa(ans)...)
          ans = 1
          ch = b[i]

       }
    }
    slice = append(slice,ch)
    slice = append(slice,strconv.Itoa(ans)...)
    if len(slice) >= len(b) {
        return S
    }else{
        return string(slice)
    }
}

```
```go
func compressString(S string) string {
    if ""  == S {
        return S 
    }
    var sb strings.Builder
    ch := S[0]
    ans := 1
    for  i := 1 ;i < len(S); i++{
       if ch == S[i] {
          ans += 1
       }else{
          sb.WriteByte(ch)
          sb.WriteString(strconv.Itoa(ans))
          ans = 1
          ch = S[i]
       }
    }
    sb.WriteByte(ch)
    sb.WriteString(strconv.Itoa(ans))
    if sb.Len() >= len(S) {
        return S
    }else{
        return sb.String()
    }
}
```

# Contains

查找子串是否在指定的字符串中

`strings.Contains("seafood", "foo")` // true

# Count
统计一个字符串中有几个指定的子串
`strings.Count("cheese","e")` // 3

# EqualFold

不区分大小写的字符串比较(==是区分字母大小写)

`fmt.Println(strings.EqualFold("abc", "Abc"))` // true

# Index

返回子串在字符串第一次出现的index值, 如果没有返回-1

`strings.Index("NLT_abc", "abc")`// 4

# LastIndex

子串sep在字符串s中最后一次出现的位置，不存在则返回-1。

`strings.LastIndex("go goland", "go")`

# Replace

```
func Replace(s, old, new string, n int) string
```

返回将s中前n个不重叠old子串都替换为new的新字符串，如果n<0会替换所有old子串。

# Split

```
func Split(s, sep string) []string
```

用去掉s中出现的sep的方式进行分割，会分割到结尾，并返回生成的所有片段组成的切片（每一个sep都会进行一次切割，即使两个sep相邻，也会进行两次切割）。如果sep为空字符，Split会将s切分成每一个unicode码值一个字符串。

# 转换大小写

```
func ToLower(s string) string
```

返回将所有字母都转为对应的小写版本的拷贝。

```
func ToUpper(s string) string
```

返回将所有字母都转为对应的大写版本的拷贝。

# 去掉字符串左右两边的指定字符

```
func Trim(s string, cutset string) string
```

返回将s前后端所有cutset包含的utf-8码值都去掉的字符串。

> 更多操作如 去掉空格, 去掉!, 去掉func 去看标准库

# 字符串是否以指定的前缀开头或 后缀结尾

```
func HasSuffix(s, suffix string) bool
```

判断s是否有后缀字符串suffix。

```
func HasPrefix(s, prefix string) bool
```

判断s是否有前缀字符串prefix。