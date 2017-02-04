package handlers

import (
	"net/http"

	"github.com/spazbite187/sensornet/graphs"
	"github.com/spazbite187/sensornet/storage"
	"gopkg.in/gin-gonic/gin.v1"
)

// GetTempGraph ...
func (d *Data) GetTempGraph(c *gin.Context) {
	htmlResp := make(map[string]interface{}) // html response map
	defer d.Data.DbugLog.Printf("html resp data: %v", htmlResp)

	sensors, err := storage.GetAllSensorData(c.Param("sensorid"), d.Data.DB)
	if err != nil {
		d.Data.ErrLog.Printf("getting data: %s", err.Error())
		c.HTML(http.StatusInternalServerError, "err.tmpl", htmlResp)
		return
	}

	image, err := graphs.GetTempGraph(sensors)
	if err != nil {
		d.Data.ErrLog.Printf("getting graph: %s", err.Error())
		c.HTML(http.StatusInternalServerError, "err.tmpl", htmlResp)
		return
	}

	c.Header("Content-Type", "image/svg+xml")
	c.Data(http.StatusOK, "", image)
}

// GetSignalGraph ...
func (d *Data) GetSignalGraph(c *gin.Context) {
	htmlResp := make(map[string]interface{}) // html response map
	defer d.Data.DbugLog.Printf("html resp data: %v", htmlResp)

	sensors, err := storage.GetAllSensorData(c.Param("sensorid"), d.Data.DB)
	if err != nil {
		d.Data.ErrLog.Printf("getting data: %s", err.Error())
		c.HTML(http.StatusInternalServerError, "err.tmpl", htmlResp)
		return
	}

	image, err := graphs.GetSignalGraph(sensors)
	if err != nil {
		d.Data.ErrLog.Printf("getting graph: %s", err.Error())
		c.HTML(http.StatusInternalServerError, "err.tmpl", htmlResp)
		return
	}

	c.Header("Content-Type", "image/svg+xml")
	c.Data(http.StatusOK, "", image)
}
