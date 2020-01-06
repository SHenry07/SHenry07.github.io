3. 【初级】通过指针变量 p 访问其成员变量 name，下面语法正确的是（）
   A. p.name
   B. (*p).name
   C. (&p).name
   D. p->name

参考答案：AB

4.   【初级】关于接口和类的说法，下面说法正确的是（）
A. 一个类只需要实现了接口要求的所有函数，我们就说这个类实现了该接口
B. 实现类的时候，只需要关心自己应该提供哪些方法，不用再纠结接口需要拆得多细才合理
C. 类实现接口时，需要导入接口所在的包
D. 接口由使用方按自身需求来定义，使用方无需关心是否有其他模块定义过类似的接口

参考答案：ABD


7. 【中级】关于init函数，下面说法正确的是（）
   A. 一个包中，可以包含多个init函数
   B. 程序编译时，先执行导入包的init函数，再执行本包内的init函数
   C. main包中，不能有init函数
   D. init函数可以被其他函数调用

参考答案：AB

【中级】对于函数定义：
```go
func add(args ...int) int {

	sum :=0

	for _,arg := range args {

	    sum += arg

	 }

	 return sum
}
```
下面对add函数调用正确的是（）
A. add(1, 2)
B. add(1, 3, 7)
C. add([]int{1, 2})
D. add([]int{1, 3, 7}...)

参考答案：ABD

\19. 【初级】关于局部变量的初始化，下面正确的使用方式是（）
A. var i int = 10
B. var i = 10
C. i := 10
D. i = 10

参考答案：ABC

\36. 【中级】下面赋值正确的是（）
A. var x = nil
B. var x interface{} = nil
C. var x string = nil
D. var x error = nil

38. 【中级】从切片中删除一个元素，下面的算法实现正确的是（）
A.
```go
 func (s *Slice)Remove(value interface{})error {

 for i, v := range *s {

    if isEqual(value, v) {

        if i== len(*s) - 1 {
          *s = (*s)[:i]
        }else {
          *s = append((*s)[:i],(*s)[i + 2:]...)
        }
        return nil
  }
 }
 return ERR_ELEM_NT_EXIST
}
```
B.
```go
func (s*Slice)Remove(value interface{}) error {

    for i, v:= range *s {
     if isEqual(value, v) {    
        *s =append((*s)[:i],(*s)[i + 1:])    
        return nil    
    }
}
return ERR_ELEM_NT_EXIST
}
```
C.
```go
func (s*Slice)Remove(value interface{}) error {

for i, v:= range *s {
    if isEqual(value, v) {    
        delete(*s, v)    
        return nil    
    }
}
return ERR_ELEM_NT_EXIST
}
```
D.

```go
func (s*Slice)Remove(value interface{}) error {

for i, v:= range *s {
	if isEqual(value, v) {
    	*s =append((*s)[:i],(*s)[i + 1:]...)
    	return nil
	}
}
return ERR_ELEM_NT_EXIST
}
```


