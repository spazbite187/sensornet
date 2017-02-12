package handlers

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// GetTempGraph ...
func (d *Data) GetTempGraph(c *gin.Context) {
	htmlResp := make(map[string]interface{}) // html response map
	defer d.Data.DebugLog.Printf("html resp data: %v", htmlResp)

	var image []byte
	for _, sensor := range d.Data.CachedSensors {
		if sensor.ID == c.Param("sensorid") {
			image = sensor.TempGraph
		}
	}

	c.Header("Content-Type", "image/svg+xml")
	c.Data(http.StatusOK, "", image)
}

// GetSignalGraph ...
func (d *Data) GetSignalGraph(c *gin.Context) {
	htmlResp := make(map[string]interface{}) // html response map
	defer d.Data.DebugLog.Printf("html resp data: %v", htmlResp)

	var image []byte
	for _, sensor := range d.Data.CachedSensors {
		if sensor.ID == c.Param("sensorid") {
			image = sensor.SignalGraph
		}
	}

	c.Header("Content-Type", "image/svg+xml")
	c.Data(http.StatusOK, "", image)
}
