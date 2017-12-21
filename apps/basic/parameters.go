package main

import "fmt"

func Map(f func(int) int, xs []int) []int {
	ys := make([]int, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

func callMap() {
	fmt.Println("result :", Map(func(x int) int { return x * x }, []int{1, 2, 3}))
}

func main() {
	callMap()
}
