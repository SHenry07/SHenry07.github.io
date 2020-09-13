# 深入理解排序
sort 包中有一个 sort.Interface 接口，该接口有三个方法 Len() 、 Less(i,j) 和 Swap(i,j) 。 通用排序函数 sort.Sort 可以排序任何实现了 sort.Inferface 接口的对象(变量)。对于 [] int 、[] float64 和 [] string 除了使用特殊指定的函数外，还可以使用改装过的类型 IntSclice 、 Float64Slice 和 StringSlice ， 然后直接调用它们对应的 Sort() 方法；因为这三种类型也实现了 sort.Interface 接口， 所以可以通过 sort.Reverse 来转换这三种类型的 Interface.Less 方法来实现逆向排序， 这就是前面最后一个排序的使用。

```go
package main

import (
	"fmt"
	"sort"
)

// Reverse 自定义的类型
type Reverse struct {
	sort.Interface // 这样，Reverse可以接纳任何实现了sort.Interface的对象
}

// Less Reverse只是将其中的 Inferface.Less 的顺序对调了一下
func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

func main() {
	ints := []int{5, 2, 7, 3, 1, 4}

	sort.Ints(ints) // 特殊排序函数，升序
	fmt.Println("after sort by Ints:\t", ints)

	doubles := []float64{2.3, 3.2, 6.7, 10.9, 5.4, 1.8}

	sort.Float64s(doubles)
	fmt.Println("after sort by Float64s:\t", doubles) // [1.8 2.3 3.2 5.4 6.7 10.9]

	strings := []string{"hello", "good", "students", "morning", "people", "world"}
	sort.Strings(strings)
	fmt.Println("after sort by Strings:\t", strings) // [good hello mornig people students world]

	ipos := sort.SearchInts(ints, 6) // int 搜索 属于一个不在数组里的数,返回插入该数,该数所在的索引
	fmt.Printf("pos of 6 is %d th\n", ipos)

	dpos := sort.SearchFloat64s(doubles, 21.8) // float64 搜索
	fmt.Printf("pos of 21.8 is %d th\n", dpos)

	fmt.Printf("doubles is asc ? %v\n", sort.Float64sAreSorted(doubles))

	doubles = []float64{3.5, 4.2, 8.9, 100.98, 20.14, 79.32}
	// sort.Sort(sort.Float64Slice(doubles))    // float64 排序方法 2
	// fmt.Println("after sort by Sort:\t", doubles)    // [3.5 4.2 8.9 20.14 79.32 100.98]
	(sort.Float64Slice(doubles)).Sort()           // float64 排序方法 3  ()可加可不加 类似类型断言
	fmt.Println("after sort by Sort:\t", doubles) // [3.5 4.2 8.9 20.14 79.32 100.98]

	sort.Sort(Reverse{sort.Float64Slice(doubles)})         // float64 逆序排序
	fmt.Println("after sort by Reversed Sort:\t", doubles) // [100.98 79.32 20.14 8.9 4.2 3.5]
}
// 二维数组
type IntList [][]int

func (s IntList)Less(i,j int) bool {
	return s[i][0] <= s[j][0]
}
func (s IntList)Swap(i,j int) {
	s[i], s[j] = s[j], s[i]
}
func (s IntList) Len() int  { return len(s) }


```

我前面一直认为， sort.Interface.Swap 方法不仅是多余的， 简直就是多余 ！ 难道么人赶脚到么 ？ 因为 a b 的交换就是将 a, b = b, a !

但是我后来思量了一下， Less 和 Swap 的参数不是待比较和交换的两个元素， 而是两个序数 i 和 j ；再结合上面结构体排序方法 2 的例子， Sort 的未必就是数组或是分片， 可以是一个普通的对象中的一部分，那么，结合着想一想， Sort 岂不是可以只针对某些位置上的进行排序，比如奇数位上的 ？比如 10 个元素的数组，长度是 5， 对于 Less 和 Swap 的 i 和 j ，me 分别针对 2*i 和 2*j 操作，这样岂不就是只对偶数位上进行排序了？！

实际上确实如此 ！ 下面的程度就是对于一个 int 数组， 偶数位上递增排序， 奇数位上递减排序(数组的第一位序数是 0 ，所以是偶数位)。
```go
 package main  
import (
	"fmt"
    "sort"
)
  
type IntList struct{
  data [] int
  cmp func(p, q int) bool
}
  
func (ints IntList) Len() int{
	return (len(ints.data) + 1) / 2
}
 
func (ints IntList) Less (i, j int) bool {
   return ints.cmp(ints.data[2*i], ints.data[2*j])
}
  
func (ints IntList) Swap (i, j int){
   ints.data[2*i], ints.data[2*j] = ints.data[2*j], ints.data[2*i]
}
  
func SortInts(ints [] int, cmp func(p, q int) bool){
   sort.Sort(IntList{ints, cmp});
}
 
func main() {
 intList := [] int {2, 4, 3, 7, 8, 9, 5, 1, 0, 6}
  
 SortInts(intList, func (p, q int) bool { return p < q;});   // 偶数位置上递增
  SortInts(intList[1:], func (p, q int) bool { return q < p;});   // 奇数位上递减
  fmt.Println(intList)   // [0 9 2 7 3 6 5 4 8 1] 
}
```



# map排序的方法

- 最好的方法是创建一个struct
  要对golang map按照value进行排序，思路是直接不用map，用struct存放key和value，实现sort接口，就可以调用sort.Sort进行排序了。

`func SliceStable(slice interface{}, less func(i, j int) bool)`

```go
people := []struct {
	Name string
	Age  int
}{
	{"Alice", 25},
	{"Elizabeth", 75},
	{"Alice", 75},
	{"Bob", 75},
	{"Alice", 75},
	{"Bob", 25},
	{"Colin", 25},
	{"Elizabeth", 25},
}

// Sort by name, preserving original order
sort.SliceStable(people, func(i, j int) bool { return people[i].Name < people[j].Name })
fmt.Println("By name:", people)

// Sort by age preserving name order
sort.SliceStable(people, func(i, j int) bool { return people[i].Age < people[j].Age })
fmt.Println("By age,name:", people)
```
```txt
Output:

By name: [{Alice 25} {Alice 75} {Alice 75} {Bob 75} {Bob 25} {Colin 25} {Elizabeth 75} {Elizabeth 25}]
By age,name: [{Alice 25} {Bob 25} {Colin 25} {Elizabeth 25} {Alice 75} {Alice 75} {Bob 75} {Elizabeth 75}]
```



- golang的map不保证有序性，所以按key排序需要取出key，对key排序，再遍历输出value

  [go sort](https://www.jianshu.com/p/01adb0e2a69f)