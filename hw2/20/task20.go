package main

import (
	"fmt"
	"math"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func Sqrt(x float64) (float64, error) {
	z := 1.0
	cnt := 0
	if x < 0 {
		return -1, ErrNegativeSqrt(x) 
	} else if x == 0 {
		return 0, nil
	}
	for {
		cnt++
		fmt.Println("number of iterations:", cnt, z) 
		if math.Abs(z - z - (z * z - x) / (2 * z)) <= 1e-6 || cnt > 50{
			break
		}
		z -= (z * z - x) / (2 * z)
	}
	return z, nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
	fmt.Println(Sqrt(0))
}