## 中文

有时候我们想通过表单元素获取一个用户的中文名字，但是又为了保证获取的是正确的中文，我们需要进行验证，而不是用户随便的一些输入。对于中文我们目前有两种方式来验证，可以使用 `unicode` 包提供的 `func Is(rangeTab *RangeTable, r rune) bool` 来验证，也可以使用正则方式来验证，这里使用最简单的正则方式，如下代码所示

```
if m, _ := regexp.MatchString("^\\p{Han}+$", r.Form.Get("realname")); !m {
	return false
}
```

[更多研效方式](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/04.2.md)