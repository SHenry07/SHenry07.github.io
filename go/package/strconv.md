## 字符串 和 数字的转换

```go
n, err := strconv.Atoi("123")
if err != nil {
    fmt.Println("转换错误",err)
}else {
    fmt.Println("转换的结果是",n)
}
// 可以用来验证输入
```





## 10进制转化为2, 8, 16进制

```go
Func FormatInt(i int64, base int) string
```

返回i的base进制的字符串表示,base必须在2到36之间,结果中会使用小写字母'a'到'z'表示大于10的数字

`Func FormatUint(i int64, base int) string`无符号版本