package storage

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/spazbite187/sensornet"
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

// UpdateSensorLocation ...
func UpdateSensorLocation(sensor *sensornet.Sensor, db *storm.DB) error {
	// update sensor
	err := db.UpdateField(&sensornet.Sensor{ID: sensor.ID}, "Location", sensor.Location)
	if err != nil {
		return err
	}

	return nil
}

func updateSensorNumReadings(sensor *sensornet.Sensor, db *storm.DB) error {
	err := db.UpdateField(&sensornet.Sensor{ID: sensor.ID}, "NumReadings", sensor.NumReadings)
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

// GetSensor ...
func GetSensor(ID string, db *storm.DB) (sensornet.Sensor, error) {
	var sensor sensornet.Sensor

	err := db.One("ID", ID, &sensor)
	if err != nil {
		return sensor, err
	}

	return sensor, nil
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

// GetAllSensorData ...
func GetAllSensorData(ID string, db *storm.DB) ([]sensornet.SensorData, error) {
	var sensorData []sensornet.SensorData
	err := db.Find("ID", ID, &sensorData, storm.Reverse())
	if err != nil {
		return sensorData, err
	}

	return sensorData, nil
}

// GetLimitedSensorData ...
func GetLimitedSensorData(ID string, db *storm.DB) ([]sensornet.SensorData, error) {
	var sensorData []sensornet.SensorData
	err := db.Find("ID", ID, &sensorData, storm.Limit(5000), storm.Reverse())
	if err != nil {
		return sensorData, err
	}

	return sensorData, nil
}

// GetSensors ...
func GetSensors(db *storm.DB) ([]sensornet.Sensor, error) {
	var sensors []sensornet.Sensor

	err := db.All(&sensors)
	if err != nil {
		return sensors, err
	}

	return sensors, nil
}

// GetTempFAvg ...
func GetTempFAvg(ID string, db *storm.DB) (float64, error) {
	var total float64
	var count float64

	data, err := GetAllSensorData(ID, db)
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

// GetSignalAvg ...
func GetSignalAvg(ID string, db *storm.DB) (float64, error) {
	var total float64
	var count float64

	data, err := GetAllSensorData(ID, db)
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

// GetHighTemp ...
// TODO: Fix issues with going from postive numbers to negative numbers
func GetHighTemp(ID string, db *storm.DB) (sensornet.SensorData, error) {
	sensorData, err := getHighest("TempF", ID, db)
	if err != nil {
		return sensorData, err
	}

	return sensorData, nil
}

// GetLowTemp ...
func GetLowTemp(ID string, db *storm.DB) (sensornet.SensorData, error) {
	sensorData, err := getLowest("TempF", ID, db)
	if err != nil {
		return sensorData, err
	}

	return sensorData, nil
}

// GetHighSignal ...
func GetHighSignal(ID string, db *storm.DB) (sensornet.SensorData, error) {
	sensorData, err := getHighest("Signal", ID, db)
	if err != nil {
		return sensorData, err
	}

	return sensorData, nil
}

// GetLowSignal ...
func GetLowSignal(ID string, db *storm.DB) (sensornet.SensorData, error) {
	sensorData, err := getLowest("Signal", ID, db)
	if err != nil {
		return sensorData, err
	}

	return sensorData, nil
}

// GetNumReadings ...
func GetNumReadings(ID string, db *storm.DB) (int, error) {
	var sensor sensornet.Sensor
	err := db.One("ID", ID, &sensor)
	if err != nil {
		return 0, err
	}

	return sensor.NumReadings, nil
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

// UpdateNumReadings ...
func UpdateNumReadings(db *storm.DB) error {
	sensors, err := GetSensors(db)
	if err != nil {
		return err
	}

	for _, v := range sensors {
		updateNumReads(v.ID, db)
	}

	return nil
}

// updateNumReads ...
func updateNumReads(ID string, db *storm.DB) error {
	numReadings, err := getNumReads(ID, db)
	if err != nil {
		return err
	}

	err = updateSensorNumReadings(&sensornet.Sensor{ID: ID, NumReadings: numReadings}, db)
	if err != nil {
		return err
	}

	return nil
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
	sensors, err := GetSensors(db)
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
