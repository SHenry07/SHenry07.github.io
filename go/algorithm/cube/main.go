package main

/*
题目描述
形如 a^3= b^3+ c^3+ d^3的等式被称为完美立方。例如 12^3= 6^3+ 8^3+ 10^3。编写一个程序，对任给的正整数 (N≤100) ，
寻找所有的四元组 (a, b, c, d)，使得 a^3= b^3+ c^3+ d^3， 其中 a,b,c,d大于 1, 小于等于N，且 b<=c<=d。

输入
一个正整数 N (N≤100)

输出
每行输出一个完美立方。格式为： Cube = a, Triple (b,c,d )其 中 a,b,c,d a,b,c,d 所在位置分别用实际求出四元组值代入。
要求： 请按照 a的值，从小到大依次输出。当两个完美立方 等式中 a的值相同，则 b值小的优先输出、仍相同 则c值小的优先输出、
再相同则 d值小的先输出。
*/
import "fmt"

func main() {
	var input int
	fmt.Print("Please type into a number(<=100): ")
	fmt.Scanln(&input)
	if input > 100 {
		fmt.Println("invalid input, valid range is from 1 to 100.")
		return
	}
	// for a := 2; a <= input; a++ { //外层a  其中 a,b,c,d大于 1,
	// 	a3 := a * a * a
	// 	for b := a - 1; b > 1; b-- { // b 一定小于a  因为3个数的立方相加 才等于a的立方
	// 		b3 := b * b * b
	// 		for c := b; c > 1; c-- { // c 且 b<=c<=d。
	// 			c3 := c * c * c
	// 			for d := c; d > 1; d-- {
	// 				if a3 == d*d*d+b3+c3 {
	// 					fmt.Printf("cube = %d,tringer = %d,%d.%d \n", a, b, c, d)
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	// 因为输出的要求所以下面更合理
	for a := 2; a <= input; a++ { //外层a  其中 a,b,c,d大于 1,
		a3 := a * a * a
		for b := 2; b < a; b++ { // b 一定小于a  因为3个数的立方相加 才等于a的立方
			b3 := b * b * b
			for c := b; c < a; c++ { // c 且 b<=c<=d。
				c3 := c * c * c
				for d := c; d < a; d++ {
					if a3 == d*d*d+b3+c3 {
						fmt.Printf("cube = %d,tringer = %d,%d.%d \n", a, b, c, d)
					}
				}
			}
		}
	}
}
