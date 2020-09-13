# float 转int 如何正确保留大小

```go
package main

import (
    "fmt"
)

func main() {
    nf := []float64{-1.9999, -2.0001, -2.0, 0, 1.9999, 2.0001, 2.0}

    //round
    fmt.Printf("Round : ")
    for _, f := range nf {
        fmt.Printf("%d ", round(f))
    }
    fmt.Printf("\n")

    //rounddown ie. math.floor
    fmt.Printf("RoundD: ")
    for _, f := range nf {
        fmt.Printf("%d ", roundD(f))
    }
    fmt.Printf("\n")

    //roundup ie. math.ceil
    fmt.Printf("RoundU: ")
    for _, f := range nf {
        fmt.Printf("%d ", roundU(f)) 
    }
    fmt.Printf("\n")

}

func roundU(val float64) int {
    if val > 0 { return int(val+1.0) }
    return int(val)
}

func roundD(val float64) int {
    if val < 0 { return int(val-1.0) }
    return int(val)
}

func round(val float64) int {
    if val < 0 { return int(val-0.5) }
    return int(val+0.5)
}
/*
Round : -2 -2 -2 0 2 2 2 
RoundD: -2 -3 -3 0 1 2 2 
RoundU: -1 -2 -2 0 2 3 3 
*/

func round(a float64) int {
    var a1 int;
    var raz float64;
    a1 = int(a);
    raz = a - float64(a1);
    if raz >= 0.5{
        return a1+1;
    } else if raz <= -0.5 {
        return a1-1;
    } else {
        return a1;
    }
}

```
