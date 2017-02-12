package sensornet_test

import (
	"testing"

	"time"

	"github.com/spazbite187/sensornet"
)

//
// test structs
type calcTimeSinceTestPair struct {
	input  string
	output string
}

type toLocalTimeTestPair struct {
	input  string
	output string
}

type tempTestPair struct {
	input  float64
	output float64
}

type toFixedTestPair struct {
	input  float64
	output float64
}

//
// test global vars
var calcTimeSinceTests = []calcTimeSinceTestPair{
	{"Sun Feb 12 18:48:54 2017", "42s"},
	{"Not a vaid time", "Me too"},
}

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

var toFixedTests = []toFixedTestPair{
	{12.8919291, 12.9},
	{12.8919291, 12.89},
	{12.8919291, 12.892},
	{12.8919291, 12.8919},
	{12.8919291, 12.89193},
	{12.8919291, 12.891929},
	{12.8919291, 12.8919291},
}

var toLocalTimeTests = []toLocalTimeTestPair{
	{"Sun Feb 12 06:34:57 2017", "2017-02-11 22:34:57 -0800 PST"},
	{"Sun Feb 12 06:35:42 2017", "2017-02-11 22:35:42 -0800 PST"},
	{"Sun Feb 12 99:35:42 2017", ""},
}

//
// test functions
func TestCalcTimeSince(t *testing.T) {
	for _, pair := range calcTimeSinceTests {
		v, _ := sensornet.CalcTimeSince(pair.input)
		// convert back to time.Time struct
		vTime, _ := time.ParseDuration(v)
		outputTime, _ := time.ParseDuration(pair.output)
		if vTime < outputTime {
			t.Error("For", pair.input,
				"expected", pair.output,
				"got", v)
		}
	}
}

func TestToLocalTime(t *testing.T) {
	for _, pair := range toLocalTimeTests {
		v, _ := sensornet.ToLocalTime(pair.input)
		if v != pair.output {
			t.Error("For", pair.input,
				"expected", pair.output,
				"got", v)
		}
	}
}

func TestTempCtoF(t *testing.T) {
	for _, pair := range tempTests {
		v := sensornet.TempCtoF(pair.input)
		if v != pair.output {
			t.Error("For", pair.input,
				"expected", pair.output,
				"got", v)
		}
	}
}

func TestToFixed(t *testing.T) {
	for i, pair := range toFixedTests {
		v := sensornet.ToFixed(pair.input, i+1)
		if v != pair.output {
			t.Error("For", pair.input,
				"expected", pair.output,
				"got", v)
		}
	}
}
