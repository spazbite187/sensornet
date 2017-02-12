package sensornet

import "testing"

//
// test structs
type tempTestPair struct {
	celsius    float64
	fahrenheit float64
}

//
// test global vars
var tempTests = []tempTestPair{
	{-100, -148},
	{-32, -25.6},
	{-20, -4},
	{-10, 14},
	{0, 32},
	{10, 50},
	{20, 68},
	{32, 89.6},
	{100, 212},
}

//
// test functions
func TestTempCtoF(t *testing.T) {
	for _, pair := range tempTests {
		v := TempCtoF(pair.celsius)
		if v != pair.fahrenheit {
			t.Error("For", pair.celsius,
				"expected", pair.fahrenheit,
				"got", v)
		}
	}
}
