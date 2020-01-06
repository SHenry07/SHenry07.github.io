package main

import (
	"fmt"
	"sync"
	//"time"
)

const N = 10

var wg = &sync.WaitGroup{}
var wg1 sync.WaitGroup

func main() {

	fmt.Printf("%T,%T\n", wg, wg1)

}
