package main

func initRoutes() {
	// routes - static assets
	router.Static("/assets", appData.Assets)

	// html front-end routes
	router.GET("/", hndlers.GetIndex)
	router.GET("/sensors", hndlers.GetSensors)
	router.GET("/sensor/:sensorid", hndlers.GetSensor)

	// json api routes
	// Create
	router.POST("/api/v1/sensor", hndlers.CreateSensorJSON)

	// Read
	router.GET("/api/v1/sensors", hndlers.GetSensorsJSON)
	router.GET("/api/v1/sensor/:sensorid", hndlers.GetSensorDataJSON)
	router.GET("/api/v1/sensor/:sensorid/readings", hndlers.GetNumReadingsJSON)
	router.GET("/api/v1/sensor/:sensorid/temp/graph", hndlers.GetTempGraph)
	router.GET("/api/v1/sensor/:sensorid/temp/high", hndlers.GetHighTempJSON)
	router.GET("/api/v1/sensor/:sensorid/temp/low", hndlers.GetLowTempJSON)
	router.GET("/api/v1/sensor/:sensorid/signal/high", hndlers.GetHighSigJSON)
	router.GET("/api/v1/sensor/:sensorid/signal/low", hndlers.GetLowSigJSON)
	router.GET("/api/v1/sensor/:sensorid/signal/graph", hndlers.GetSignalGraph)

	// Update
	router.PUT("/api/v1/sensor/:sensorid", hndlers.UpdateSensorDataJSON)
	router.PUT("/api/v1/sensor/:sensorid/location", hndlers.UpdateSensorLocationJSON)

	// Delete
}
