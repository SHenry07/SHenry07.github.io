## 单向链表

```go
package main

import (
	"math/rand"
	"fmt"
	"time"
)
// 头插法
type node struct {
	index int
	value int
	next *node
}

// Newnode _
func Newnode(i,v int, n *node) *node {
	return &node{
		index:  i ,
		value: v,
		next: n,
	}
}
func headlist() *node {
rand.Seed(time.Now().UnixNano())
	// 生成一个20个元素的链表
	head := &node{
		index: 0,
		value: rand.Int(),
		next: nil,
	}
	for i :=1 ; i < 20 ; i ++ {
		head = Newnode(i,rand.Int(),head)
	}
	return head
}

// 尾插法
func newnode(i,v int) node {
	return node{
		index:  i ,
		value: v,
		next: nil,
	}

}

func taillist() *node {
	rand.Seed(time.Now().UnixNano())
	// 生成一个20个元素的链表
	var current *node
	head := &node{
		value: rand.Int(),
		next: nil,
	}
	current = head

	for i :=1 ; i < 20 ; i ++ {
		wudi := newnode(i,rand.Int())
		current.next = &wudi
		current = &wudi
	}
	return head
}


func main() {
	tmp := headlist()
	for {
		fmt.Println(tmp.index,tmp.value)	
		if tmp.next == nil {
			break
		}
		tmp = tmp.next
	}

	test := taillist()
		fmt.Println("tail")
	for {
		fmt.Println(test.index,test.value)
		if test.next == nil {
			break
		}
		test = test.next
	}
}
```

