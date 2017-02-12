package handlers

import (
	"net/http"

	"github.com/spazbite187/sensornet"
	"gopkg.in/gin-gonic/gin.v1"
)

// GetIndex returns the index page of the app.
func (d *Data) GetIndex(c *gin.Context) {
	htmlResp := make(map[string]interface{})
	defer d.Data.DebugLog.Printf("html resp data: %v", htmlResp)

	c.HTML(http.StatusOK, "index.tmpl", htmlResp)
}

// GetSensors returns the index page of the app.
func (d *Data) GetSensors(c *gin.Context) {
	htmlResp := make(map[string]interface{})
	defer d.Data.DebugLog.Printf("html resp data: %v", htmlResp)

	htmlResp["sensors"] = d.Data.CachedSensors
	c.HTML(http.StatusOK, "sensors.tmpl", htmlResp)
}

// GetSensor ...
func (d *Data) GetSensor(c *gin.Context) {
	htmlResp := make(map[string]interface{})
	defer d.Data.DebugLog.Printf("html resp data: %v", htmlResp)

	sensor := &sensornet.Sensor{}
	for _, v := range d.Data.CachedSensors {
		if v.ID == c.Param("sensorid") {
			sensor = v
		}
	}

	htmlResp["sensor"] = sensor
	c.HTML(http.StatusOK, "sensor.tmpl", htmlResp)
}
