// used packages - 1)fmt, 2)math, 3)testing, 4)net/http

package main

import (
	"fmt"
	// "testing"
	"math"
	"net/http"
)

type dataForBMI struct {
	weight float64
	height float64
}

func main() {
	fmt.Printf("Starting server at port 8080\n")
	LoadServer()
}

func LoadServer() {
	http.HandleFunc("/", bmi)
	http.ListenAndServe(":8080", nil)
}

func bmi(w http.ResponseWriter, req *http.Request) {
	data := dataForBMI{weight: 20, height: 0.8}
	myBmi, myResult := GetBMI(data.weight, data.height)
	fmt.Fprintf(w, "Your body mass index is %v. %v", myBmi, myResult)

}

func GetBMI(weight float64, height float64) (float64, string) {
	bmiRes := weight / math.Pow(height, 2)
	res := ""
	if bmiRes < 18.5 {
		res = "Underweight( Eat more cookies."
	} else if bmiRes < 24.9 {
		res = "Everything is okay:)"
	} else {
		res = "Overweight( Eat less cookies."
	}
	return bmiRes, res
}
