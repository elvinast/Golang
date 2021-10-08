package main

import "testing"

func TestingResult(myBMI float64, correctBMI float64, tst *testing.T) {
	if myBMI == correctBMI {
		tst.Log("Correct")
	} else {
		tst.Error("Not correct")
	}
}

func Test(tst *testing.T) {
	TestingResult(GetBMI(20, 0.8), 31.249999999999993, tst)
}