| 转义字符escape char | 作用                                         |
| ------------------- | -------------------------------------------- |
| \n                  | 换行                                         |
| \t                  | 制表符 用于对齐 排版                         |
| \r                  | 回车,从当前行的最前面开始输出,覆盖掉以前内容 |

`fmt.Println("天龙八部雪山飞狐\r张飞")`


| 缩写 | 类型含义     |
| ---- | ------------ |
| T    | 类型         |
| v    | 值           |
| s    | 字符串       |
| d    | 整数         |
| f    | 浮点数       |
| t    | 布尔值         |
| p   | 指针  |
| c   | 字符类型character     |
| b    | 二进制整数   |
| o    | 八进制整数   |
| x    | 十六进制整数 |
|q|a single-quoted character literal safely escaped with Go syntax.该值对应的单引号括起来的go语法字符字面值，必要时会采用安全的转义表示|
| #加缩写 | 详细输出   |
| +| 可能也是详细输出/总是为数字打印一个符号 always print a sign for numeric values; guarantee ASCII-only output for %q (%+q)|
|- | 左对齐 pad with spaces on the right rather than the left (left-justify the field)|
|0-7任意缩写|宽度0-7|
|7|宽度7|
|%|输出%|

```go
// Decimal 保留2位数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
```

[官方fmt包解释](https://pkg.go.dev/fmt?tab=doc) 