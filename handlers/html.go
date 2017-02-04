package handlers

import (
	"net/http"

	"github.com/asdine/storm"
	"github.com/spazbite187/sensornet"
	"github.com/spazbite187/sensornet/storage"
	"gopkg.in/gin-gonic/gin.v1"
)

// GetIndex returns the index page of the app.
func (d *Data) GetIndex(c *gin.Context) {
	htmlResp := make(map[string]interface{})
	defer d.Data.DbugLog.Printf("html resp data: %v", htmlResp)

	c.HTML(http.StatusOK, "index.tmpl", htmlResp)
}

// GetSensors returns the index page of the app.
func (d *Data) GetSensors(c *gin.Context) {
	htmlResp := make(map[string]interface{})
	defer d.Data.DbugLog.Printf("html resp data: %v", htmlResp)

	sensors, err := storage.GetSensors(d.Data.DB)
	if err != nil {
		d.Data.ErrLog.Printf("getting sensors: %s", err.Error())
		c.HTML(http.StatusInternalServerError, "err.tmpl", htmlResp)
		return
	}

	var sensorsData []sensornet.SensorData
	for _, v := range sensors {
		data, err := storage.GetSensorData(v.ID, d.Data.DB)
		if err != nil {
			if err != storm.ErrNotFound {
				d.Data.ErrLog.Printf("getting sensor data: %s", err.Error())
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		} else {
			data.TimeSince, err = sensornet.CalcTimeSince(data.LastUpdate)
			if err != nil {
				d.Data.ErrLog.Printf("parsing time: %s", err.Error())
				c.HTML(http.StatusInternalServerError, "err.tmpl", htmlResp)
				return
			}
			sensorsData = append(sensorsData, data)
		}
	}

	htmlResp["sensors"] = sensorsData
	c.HTML(http.StatusOK, "sensors.tmpl", htmlResp)
}

// GetSensor ...
func (d *Data) GetSensor(c *gin.Context) {
	htmlResp := make(map[string]interface{})
	defer d.Data.DbugLog.Printf("html resp data: %v", htmlResp)

	ID := c.Param("sensorid")
	sensorData, err := storage.GetSensorData(ID, d.Data.DB)
	if err != nil {
		d.Data.ErrLog.Printf("getting data: %s", err.Error())
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

	sensorData.LastUpdate, err = sensornet.ToLocalTime(sensorData.LastUpdate)
	if err != nil {
		d.Data.ErrLog.Printf("parsing time: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	htmlResp["sensor"] = sensorData

	htmlResp["readings"], err = storage.GetNumReadings(ID, d.Data.DB)
	if err != nil {
		d.Data.ErrLog.Printf("getting num of readings: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	highTemp, err := storage.GetHighTemp(ID, d.Data.DB)
	if err != nil {
		d.Data.ErrLog.Printf("getting high temp: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	highTemp.LastUpdate, err = sensornet.ToLocalTime(highTemp.LastUpdate)
	if err != nil {
		d.Data.ErrLog.Printf("%s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	htmlResp["highTemp"] = highTemp

	lowTemp, err := storage.GetLowTemp(ID, d.Data.DB)
	if err != nil {
		d.Data.ErrLog.Printf("getting low temp: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	lowTemp.LastUpdate, err = sensornet.ToLocalTime(lowTemp.LastUpdate)
	if err != nil {
		d.Data.ErrLog.Printf("%s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	htmlResp["lowTemp"] = lowTemp

	htmlResp["tempFavg"], err = storage.GetTempFAvg(ID, d.Data.DB)
	if err != nil {
		d.Data.ErrLog.Printf("getting temp avg: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	lowSig, err := storage.GetLowSignal(ID, d.Data.DB)
	if err != nil {
		d.Data.ErrLog.Printf("getting low sig: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	lowSig.LastUpdate, err = sensornet.ToLocalTime(lowSig.LastUpdate)
	if err != nil {
		d.Data.ErrLog.Printf("%s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	htmlResp["lowSig"] = lowSig

	highSig, err := storage.GetHighSignal(ID, d.Data.DB)
	if err != nil {
		d.Data.ErrLog.Printf("getting high sig: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	highSig.LastUpdate, err = sensornet.ToLocalTime(highSig.LastUpdate)
	if err != nil {
		d.Data.ErrLog.Printf("%s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	htmlResp["highSig"] = highSig

	htmlResp["signalAvg"], err = storage.GetSignalAvg(ID, d.Data.DB)
	if err != nil {
		d.Data.ErrLog.Printf("getting signal avg: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "sensor.tmpl", htmlResp)
}
