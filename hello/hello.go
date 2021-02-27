package main

import (
	"fmt"
)

func swap(a, b int) (int, int) {
	a, b = b, a
	return a, b
}

func main() {
	var foo func(a, b int) (int, int)

	foo = swap
	var a, b = foo(1, 2)
	fmt.Printf("%d %d\n", a, b)
	// 印出: 3

}
