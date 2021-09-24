package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := 1.0
	cnt := 0
	for {
		cnt++
		fmt.Println("number of iterations:", cnt, z) 
		if math.Abs(z - z - (z * z - x) / (2 * z)) <= 1e-6 {
			break
		}
		z -= (z * z - x) / (2 * z)
	}
	return z
}

func main() {
	fmt.Println(Sqrt(555))
	fmt.Println(math.Sqrt(555))
}
