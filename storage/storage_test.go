package storage_test

import (
	"os"
	"testing"

	"github.com/asdine/storm"

	"github.com/spazbite187/sensornet"
	"github.com/spazbite187/sensornet/storage"
)

const testDBFilename = "sensornetTest.db"

var (
	testDB *storm.DB
)

func init() {
	testDB, _ = storage.GetDatabase(testDBFilename)
}

//
// test structs
type storeSensorTestPair struct {
	Sensor      sensornet.Sensor
	ErrorString string
}

type storeSensorDataTestPair struct {
	SensorData  sensornet.SensorData
	ErrorString string
}

type getLatestSensorDataTestPair struct {
	input  string
	output sensornet.SensorData
}

type updateSensorLocationTestPair struct {
	input  sensornet.Sensor
	output string
}

type cleanDBTestPair struct {
	input  int
	output int
}

//
// test global vars
var (
	testSensor01 = sensornet.Sensor{
		ID:          "12345",
		Location:    "location01",
		Uptime:      "177h8m42s",
		IP:          "192.168.1.100",
		SSID:        "SensorNet",
		Signal:      -39,
		NumReadings: 2,
		TempC:       23.31,
		TempF:       73.96,
		HighTemp:    testSensor01Data01,
		LowTemp:     testSensor01Data01,
		AvgTemp:     73.96,
		TempGraph:   []byte{},
		HighSignal:  testSensor01Data01,
		LowSignal:   testSensor01Data01,
		AvgSignal:   -39,
		SignalGraph: []byte{},
		LastUpdate:  "Sun Feb 12 06:35:42 2017",
	}

	testSensor02 = sensornet.Sensor{
		ID:       "54321",
		Location: "location02",
	}

	testSensor03 = sensornet.Sensor{
		ID:       "11111",
		Location: "location03",
	}

	testSensor01Data01 = sensornet.SensorData{
		Pk:         1,
		ID:         "12345",
		Location:   "location01",
		Uptime:     "177h8m42s",
		IP:         "192.168.1.100",
		SSID:       "SensorNet",
		Signal:     -39,
		TempC:      23.31,
		TempF:      73.96,
		LastUpdate: "Sun Feb 12 06:35:42 2017",
		TimeSince:  "42s",
	}

	testSensor01Data02 = sensornet.SensorData{
		Pk:         2,
		ID:         "12345",
		Location:   "location01",
		Uptime:     "180h8m42s",
		IP:         "192.168.1.100",
		SSID:       "SensorNet",
		Signal:     -40,
		TempC:      23.31,
		TempF:      73.96,
		LastUpdate: "Sun Feb 12 06:35:52 2017",
		TimeSince:  "52s",
	}

	testSensor01Data03 = sensornet.SensorData{
		Pk:         3,
		ID:         "12345",
		Location:   "location01",
		Uptime:     "180h8m42s",
		IP:         "192.168.1.100",
		SSID:       "SensorNet",
		Signal:     -38,
		TempC:      23.31,
		TempF:      73.96,
		LastUpdate: "Sun Feb 12 06:35:52 2017",
		TimeSince:  "52s",
	}

	testSensor02Data01 = sensornet.SensorData{
		Pk:         4,
		ID:         "54321",
		Location:   "location02",
		Uptime:     "180h8m42s",
		IP:         "192.168.1.100",
		SSID:       "SensorNet",
		Signal:     -39,
		TempC:      23.31,
		TempF:      73.96,
		LastUpdate: "Sun Feb 12 06:35:52 2017",
		TimeSince:  "52s",
	}

	testSensor02Data02 = sensornet.SensorData{
		Pk:         5,
		ID:         "54321",
		Location:   "location02",
		Uptime:     "181h8m42s",
		IP:         "192.168.1.100",
		SSID:       "SensorNet",
		Signal:     -40,
		TempC:      23.31,
		TempF:      73.96,
		LastUpdate: "Sun Feb 12 06:35:52 2017",
		TimeSince:  "54s",
	}

	testSensor02Data03 = sensornet.SensorData{
		Pk:         6,
		ID:         "54321",
		Location:   "location02",
		Uptime:     "181h8m42s",
		IP:         "192.168.1.100",
		SSID:       "SensorNet",
		Signal:     -40,
		TempC:      23.31,
		TempF:      73.96,
		LastUpdate: "Sun Feb 12 06:35:52 2017",
		TimeSince:  "54s",
	}

	testSensor03Data01 = sensornet.SensorData{
		Pk:         7,
		ID:         "11111",
		Location:   "location03",
		Uptime:     "180h8m42s",
		IP:         "192.168.1.100",
		SSID:       "SensorNet",
		Signal:     -37,
		TempC:      23.31,
		TempF:      73.96,
		LastUpdate: "Sun Feb 12 06:35:52 2017",
		TimeSince:  "52s",
	}

	testSensor03Data02 = sensornet.SensorData{
		Pk:         8,
		ID:         "11111",
		Location:   "location03",
		Uptime:     "180h8m42s",
		IP:         "192.168.1.100",
		SSID:       "SensorNet",
		Signal:     -37,
		TempC:      23.31,
		TempF:      73.96,
		LastUpdate: "Sun Feb 12 06:35:52 2017",
		TimeSince:  "52s",
	}

	testSensor03Data03 = sensornet.SensorData{
		Pk:         9,
		ID:         "11111",
		Location:   "location03",
		Uptime:     "180h8m42s",
		IP:         "192.168.1.100",
		SSID:       "SensorNet",
		Signal:     -37,
		TempC:      23.31,
		TempF:      73.96,
		LastUpdate: "Sun Feb 12 06:35:52 2017",
		TimeSince:  "52s",
	}

	testSensors = []sensornet.Sensor{testSensor01}

	storeSensorTests = []storeSensorTestPair{
		{testSensor01, ""},
		{testSensor02, ""},
		{testSensor03, ""},
	}

	storeSensorDataTests = []storeSensorDataTestPair{
		{testSensor01Data01, ""},
		{testSensor01Data02, ""},
		{testSensor01Data03, ""},
		{testSensor02Data01, ""},
		{testSensor02Data02, ""},
		{testSensor02Data03, ""},
		{testSensor03Data01, ""},
		{testSensor03Data02, ""},
		{testSensor03Data03, ""},
	}

	getLatestSensorDataTests = []getLatestSensorDataTestPair{
		{"1", sensornet.SensorData{}},
		{"12345", testSensor01Data03},
		{"54321", testSensor02Data03},
		{"11111", testSensor03Data03},
	}

	updateSensorLocationTests = []updateSensorLocationTestPair{
		{testSensor01, ""},
		{testSensor02, ""},
		{testSensor03, ""},
	}

	cleanDBTests = []cleanDBTestPair{
		{3, 3},
		{2, 2},
		{1, 1},
		{0, 0},
	}
)

//
// test functions
func TestStoreSensor(t *testing.T) {
	for _, pair := range storeSensorTests {
		v := storage.StoreSensor(&pair.Sensor, testDB)
		if v != nil {
			t.Error("For", pair.Sensor,
				"expected", pair.ErrorString,
				"got", v)
		}
	}
}

func TestStoreSensorData(t *testing.T) {
	for _, pair := range storeSensorDataTests {
		v := storage.StoreSensorData(&pair.SensorData, testDB)
		if v != nil {
			t.Error("For", pair.SensorData,
				"expected", pair.ErrorString,
				"got", v)
		}
	}
}

func TestGetLatestSensorData(t *testing.T) {
	for _, pair := range getLatestSensorDataTests {
		v, _ := storage.GetLatestSensorData(pair.input, testDB)
		if v != pair.output {
			t.Error("For", pair,
				"expected", pair.output,
				"got", v)
		}
	}
}

func TestUpdateSensorLocation(t *testing.T) {
	for _, pair := range updateSensorLocationTests {
		v := storage.UpdateSensorLocation(&pair.input, testDB)
		if v != nil {
			t.Error("For", pair,
				"expected", pair.output,
				"got", v)
		}
	}
}

func TestCleanDB(t *testing.T) {
	for _, pair := range cleanDBTests {

		storage.CleanDB(pair.input, testDB)
		var sensorDataAfter []sensornet.SensorData
		testDB.Find("ID", "12345", &sensorDataAfter)
		dbNumAfter := len(sensorDataAfter)

		if dbNumAfter != pair.input {
			t.Error("For", pair,
				"expected", pair.output,
				"got", dbNumAfter)
		}
	}
}

// utility functions
func TestUtilityFuncDeleteDB(t *testing.T) {
	var err = os.Remove(testDBFilename)
	if err != nil {
		t.Errorf("Error deleting %s", testDBFilename)
	}
}
