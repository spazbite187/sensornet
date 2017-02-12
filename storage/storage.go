package storage

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/spazbite187/sensornet"
	"github.com/spazbite187/sensornet/app"
	"github.com/spazbite187/sensornet/graphs"
)

// GetDatabase ...
func GetDatabase(dbFile string) (*storm.DB, error) {
	db, err := storm.Open(dbFile)
	if err != nil {
		return &storm.DB{}, err
	}

	return db, nil
}

// StoreSensor ...
func StoreSensor(sensor *sensornet.Sensor, db *storm.DB) error {
	// store sensor
	err := db.Save(sensor)
	if err != nil {
		return err
	}

	return nil
}

// UpdateSensors ...
func UpdateSensors(appData *app.Data) error {

	// get sensor models from db
	sensors, err := getSensors(appData.DB)
	if err != nil {
		return err
	}

	// update sensor models
	for i, v := range sensors {
		sensor, err := GetSensorData(v.ID, appData.DB)
		if err != nil {
			return err
		}
		allSensorData, err := getAllSensorData(sensor.ID, appData.DB)
		if err != nil {
			return err
		}

		if sensors[i].Location == "" || sensors[i].Location == "NEW" {
			sensors[i].Location = sensor.Location
		}
		sensors[i].Uptime = sensor.Uptime
		sensors[i].IP = sensor.IP
		sensors[i].SSID = sensor.SSID
		sensors[i].Signal = sensor.Signal
		sensors[i].NumReadings, err = getNumReads(v.ID, appData.DB)
		if err != nil {
			return err
		}
		sensors[i].TempC = sensor.TempC
		sensors[i].TempF = sensor.TempF
		sensors[i].HighTemp, err = getHighTemp(v.ID, appData.DB)
		if err != nil {
			return err
		}
		sensors[i].HighTemp.LastUpdate, err = sensornet.ToLocalTime(sensors[i].HighTemp.LastUpdate)
		if err != nil {
			return err
		}
		sensors[i].LowTemp, err = getLowTemp(v.ID, appData.DB)
		if err != nil {
			return err
		}
		sensors[i].LowTemp.LastUpdate, err = sensornet.ToLocalTime(sensors[i].LowTemp.LastUpdate)
		if err != nil {
			return err
		}
		sensors[i].AvgTemp, err = getTempFAvg(v.ID, appData.DB)
		if err != nil {
			return err
		}
		sensors[i].TempGraph, err = graphs.GetTempGraph(allSensorData)
		if err != nil {
			return err
		}

		sensors[i].Signal = sensor.Signal
		sensors[i].HighSignal, err = getHighSignal(v.ID, appData.DB)
		if err != nil {
			return err
		}
		sensors[i].HighSignal.LastUpdate, err = sensornet.ToLocalTime(sensors[i].HighSignal.LastUpdate)
		if err != nil {
			return err
		}
		sensors[i].LowSignal, err = getLowSignal(v.ID, appData.DB)
		if err != nil {
			return err
		}
		sensors[i].LowSignal.LastUpdate, err = sensornet.ToLocalTime(sensors[i].LowSignal.LastUpdate)
		if err != nil {
			return err
		}
		sensors[i].AvgSignal, err = getSignalAvg(v.ID, appData.DB)
		if err != nil {
			return err
		}
		sensors[i].SignalGraph, err = graphs.GetSignalGraph(allSensorData)
		if err != nil {
			return err
		}
		sensors[i].LastUpdate, err = sensornet.ToLocalTime(sensor.LastUpdate)
		if err != nil {
			return err
		}
	}

	// update sensors in global memory cache
	appData.CachedSensors = sensors

	return nil
}

// UpdateSensorLocation ...
func UpdateSensorLocation(sensor *sensornet.Sensor, db *storm.DB) error {
	// update sensor
	err := db.UpdateField(&sensornet.Sensor{ID: sensor.ID}, "Location", sensor.Location)
	if err != nil {
		return err
	}

	return nil
}

// StoreSensorData ...
func StoreSensorData(data *sensornet.SensorData, db *storm.DB) error {
	// store data
	err := db.Save(data)
	if err != nil {
		return err
	}

	return nil
}

// GetSensorData ...
func GetSensorData(ID string, db *storm.DB) (sensornet.SensorData, error) {
	var sensorData []sensornet.SensorData
	err := db.Find("ID", ID, &sensorData, storm.Limit(1), storm.Reverse())
	if err != nil {
		return sensornet.SensorData{}, err
	}

	return sensorData[0], nil
}

// getAllSensorData ...
func getAllSensorData(ID string, db *storm.DB) ([]*sensornet.SensorData, error) {
	var sensorData []*sensornet.SensorData
	err := db.Find("ID", ID, &sensorData, storm.Reverse())
	if err != nil {
		return sensorData, err
	}

	return sensorData, nil
}

// getSensors ...
func getSensors(db *storm.DB) ([]*sensornet.Sensor, error) {
	var sensors []*sensornet.Sensor

	err := db.All(&sensors)
	if err != nil {
		return sensors, err
	}

	return sensors, nil
}

// getTempFAvg ...
func getTempFAvg(ID string, db *storm.DB) (float64, error) {
	var total float64
	var count float64

	data, err := getAllSensorData(ID, db)
	if err != nil {
		return 0, nil
	}
	for _, v := range data {
		total = total + v.TempF
		count = count + 1
	}
	avg := sensornet.ToFixed(total/count, 2)

	return avg, nil
}

// getSignalAvg ...
func getSignalAvg(ID string, db *storm.DB) (float64, error) {
	var total float64
	var count float64

	data, err := getAllSensorData(ID, db)
	if err != nil {
		return 0, nil
	}
	for _, v := range data {
		total = total + float64(v.Signal)
		count = count + 1
	}
	avg := sensornet.ToFixed(total/count, 2)

	return avg, nil
}

// getHighTemp ...
func getHighTemp(ID string, db *storm.DB) (sensornet.SensorData, error) {
	sensorData, err := getHighest("TempF", ID, db)
	if err != nil {
		return sensorData, err
	}

	return sensorData, nil
}

// getLowTemp ...
func getLowTemp(ID string, db *storm.DB) (sensornet.SensorData, error) {
	sensorData, err := getLowest("TempF", ID, db)
	if err != nil {
		return sensorData, err
	}

	return sensorData, nil
}

// getHighSignal ...
func getHighSignal(ID string, db *storm.DB) (sensornet.SensorData, error) {
	sensorData, err := getHighest("Signal", ID, db)
	if err != nil {
		return sensorData, err
	}

	return sensorData, nil
}

// getLowSignal ...
func getLowSignal(ID string, db *storm.DB) (sensornet.SensorData, error) {
	sensorData, err := getLowest("Signal", ID, db)
	if err != nil {
		return sensorData, err
	}

	return sensorData, nil
}

// getNumReads ...
func getNumReads(ID string, db *storm.DB) (int, error) {
	var sensorData []sensornet.SensorData
	err := db.Find("ID", ID, &sensorData)
	if err != nil {
		return 0, err
	}

	return len(sensorData), nil
}

func getLowest(sort, ID string, db *storm.DB) (sensornet.SensorData, error) {
	var sensorData sensornet.SensorData
	err := db.Select(q.Eq("ID", ID)).OrderBy(sort).First(&sensorData)
	if err != nil {
		return sensorData, err
	}

	return sensorData, nil
}

func getHighest(sort, ID string, db *storm.DB) (sensornet.SensorData, error) {
	var sensorData sensornet.SensorData
	err := db.Select(q.Eq("ID", ID)).OrderBy(sort).Reverse().First(&sensorData)
	if err != nil {
		return sensorData, err
	}

	return sensorData, nil
}

// CleanDB ...
func CleanDB(max int, db *storm.DB) error {
	sensors, err := getSensors(db)
	if err != nil {
		return err
	}
	// get number of readings per sensor
	for _, sensor := range sensors {
		if sensor.NumReadings > max {
			diff := sensor.NumReadings - max
			// remove sensor readings
			var sensorsData []sensornet.SensorData
			err := db.Find("ID", sensor.ID, &sensorsData, storm.Limit(diff))
			if err != nil {
				return err
			}

			for _, data := range sensorsData {
				err := db.DeleteStruct(&data)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
