// used packages - 1)fmt, 2)math, 3)testing, 4)net/http

package main

import (
	"fmt"
	// "testing"
	"net/http"
	"math"
)

type dataForBMI struct {
	weight float64
	height float64
}

func main () {
	fmt.Printf("Starting server at port 8080\n")
	LoadServer()
}

func LoadServer () {
	http.HandleFunc("/", bmi)
	http.ListenAndServe(":8080", nil)
}

func bmi(w http.ResponseWriter, req *http.Request) {
	data := dataForBMI{weight: 20, height: 0.8}
	fmt.Fprintf(w, "%v", GetBMI(data.weight, data.height))
}

func GetBMI(weight float64, height float64) float64 {
	return weight / math.Pow(height, 2)
}