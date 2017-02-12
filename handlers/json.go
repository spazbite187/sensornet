package handlers

import (
	"net/http"
	"time"

	"github.com/asdine/storm"
	"github.com/spazbite187/sensornet"
	"github.com/spazbite187/sensornet/storage"
	"gopkg.in/gin-gonic/gin.v1"
)

// CreateSensorJSON ...
func (d *Data) CreateSensorJSON(c *gin.Context) {
	var sensor sensornet.Sensor

	err := c.BindJSON(&sensor)
	if err != nil {
		d.Data.ErrLog.Printf("bad req: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = storage.StoreSensor(&sensor, d.Data.DB)
	if err != nil {
		d.Data.ErrLog.Printf("storing sensor - %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, "ok")
}

// UpdateSensorDataJSON ...
func (d *Data) UpdateSensorDataJSON(c *gin.Context) {
	var sensorData sensornet.SensorData

	err := c.BindJSON(&sensorData)
	if err != nil {
		d.Data.ErrLog.Printf("bad req: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	sensorData.ID = c.Param("sensorid")

	tempUptime, err := time.ParseDuration(sensorData.Uptime)
	if err != nil {
		d.Data.ErrLog.Printf("%s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	sensorData.Uptime = tempUptime.String()
	sensorData.TempF = sensornet.TempCtoF(sensorData.TempC)
	sensorData.TempC = sensornet.ToFixed(sensorData.TempC, 2)
	sensorData.TempF = sensornet.ToFixed(sensorData.TempF, 2)

	now := time.Now().UTC()
	sensorData.LastUpdate = now.Format(time.ANSIC)

	retrievedSensor := &sensornet.Sensor{}
	for _, v := range d.Data.CachedSensors {
		if v.ID == sensorData.ID {
			retrievedSensor = v
		}
	}

	if retrievedSensor.Location == "" {
		retrievedSensor.Location = "NEW"
	}

	sensorData.Location = retrievedSensor.Location
	d.Data.DebugLog.Printf("sensor id: %s, uptime: %s",
		sensorData.ID,
		sensorData.Uptime,
	)
	d.Data.DebugLog.Printf("ssid: %s, signal: %d, ip: %s", sensorData.SSID, sensorData.Signal, sensorData.IP)
	d.Data.DebugLog.Printf("temp(c): %.2f, temp(f): %.2f", sensorData.TempC, sensorData.TempF)
	d.Data.DebugLog.Printf("last update: %s", sensorData.LastUpdate)

	err = storage.StoreSensorData(&sensorData, d.Data.DB)
	if err != nil {
		d.Data.ErrLog.Printf("%s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, "ok")
}

// UpdateSensorLocationJSON ...
func (d *Data) UpdateSensorLocationJSON(c *gin.Context) {
	var sensor sensornet.Sensor

	err := c.BindJSON(&sensor)
	if err != nil {
		d.Data.ErrLog.Printf("bad req: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	sensor.ID = c.Param("sensorid")

	if sensor.Location == "" {
		d.Data.ErrLog.Println("bad req")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	d.Data.DebugLog.Printf("sensor id: %s, location: %s", sensor.ID, sensor.Location)

	err = storage.UpdateSensorLocation(&sensor, d.Data.DB)
	if err != nil {
		d.Data.ErrLog.Printf("%s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, "ok")
}

// GetSensorsJSON ...
func (d *Data) GetSensorsJSON(c *gin.Context) {
	var IDs []string

	for _, v := range d.Data.CachedSensors {
		IDs = append(IDs, v.ID)
	}
	c.JSON(http.StatusOK, IDs)
}

// GetSensorDataJSON ...
func (d *Data) GetSensorDataJSON(c *gin.Context) {
	sensorData, err := storage.GetSensorData(c.Param("sensorid"), d.Data.DB)
	if err != nil {
		if err == storm.ErrNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		d.Data.ErrLog.Printf("%s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	timeSince, err := sensornet.CalcTimeSince(sensorData.LastUpdate)
	if err != nil {
		d.Data.ErrLog.Printf("fixing time: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	sensorData.TimeSince = timeSince

	c.JSON(http.StatusOK, sensorData)
}

// GetHighTempJSON ...
func (d *Data) GetHighTempJSON(c *gin.Context) {

	retrievedSensor := &sensornet.Sensor{}
	for _, v := range d.Data.CachedSensors {
		if v.ID == c.Param("sensorid") {
			retrievedSensor = v
		}
	}

	c.JSON(http.StatusOK, retrievedSensor.HighTemp)
}

// GetLowTempJSON ...
func (d *Data) GetLowTempJSON(c *gin.Context) {

	retrievedSensor := &sensornet.Sensor{}
	for _, v := range d.Data.CachedSensors {
		if v.ID == c.Param("sensorid") {
			retrievedSensor = v
		}
	}

	c.JSON(http.StatusOK, retrievedSensor.LowTemp)
}

// GetHighSigJSON ...
func (d *Data) GetHighSigJSON(c *gin.Context) {

	retrievedSensor := &sensornet.Sensor{}
	for _, v := range d.Data.CachedSensors {
		if v.ID == c.Param("sensorid") {
			retrievedSensor = v
		}
	}

	c.JSON(http.StatusOK, retrievedSensor.HighSignal)
}

// GetLowSigJSON ...
func (d *Data) GetLowSigJSON(c *gin.Context) {

	retrievedSensor := &sensornet.Sensor{}
	for _, v := range d.Data.CachedSensors {
		if v.ID == c.Param("sensorid") {
			retrievedSensor = v
		}
	}

	c.JSON(http.StatusOK, retrievedSensor.LowSignal)
}

// GetNumReadingsJSON ...
func (d *Data) GetNumReadingsJSON(c *gin.Context) {

	retrievedSensor := &sensornet.Sensor{}
	for _, v := range d.Data.CachedSensors {
		if v.ID == c.Param("sensorid") {
			retrievedSensor = v
		}
	}

	c.JSON(http.StatusOK, retrievedSensor.NumReadings)
}
