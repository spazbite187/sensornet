package sensornet

import (
	"math"
	"time"
)

const (
	// Version defines the current version number (major.minor.patch).
	Version = "1.1.0-DEV"
	// DBFilename defines the static name for the sensornet db file.
	DBFilename = "sensornet.db"
	// MaxNumReadings to keep in DB. Storing a reading every 60 seconds is 1440 readings
	// a day, or 10,080 for a 7 days.
	MaxNumReadings = 10080
	// BckgrndInterval defines the number of seconds to run background services.
	BckgrndInterval = 60 // seconds
)

// CalcTimeSince takes a timestamp as a string and returns the amount of
// time since the timestamp as a string.
func CalcTimeSince(input string) (string, error) {
	updateTime, err := parseANSICTimeStr(input)
	if err != nil {
		return "", err
	}

	return time.Since(updateTime).String(), nil
}

// ToLocalTime ...
func ToLocalTime(input string) (string, error) {
	timeObj, err := parseANSICTimeStr(input)
	if err != nil {
		return "", err
	}

	return timeObj.Local().String(), nil
}

func parseANSICTimeStr(timeStr string) (time.Time, error) {
	parsedTime, err := time.Parse(time.ANSIC, timeStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

// TempCtoF ...
func TempCtoF(c float64) float64 {
	return c*1.8 + 32
}

// ToFixed ...
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
