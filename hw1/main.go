package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := 1.0
	cnt := 0
	for i := 0; i <= int(x); i++ {
		fmt.Println("number of iterations:", cnt, z)
		if z == z - (z * z - x) / (2 * z) {
			break
		}
		z -= (z * z - x) / (2 * z)
		cnt++
	}
	return z
}

func main() {
	fmt.Println(Sqrt(901))
	fmt.Println(math.Sqrt(901))
}
