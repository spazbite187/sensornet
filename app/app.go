// Package app contains the app specific data structure.
package app

import (
	"log"

	"github.com/asdine/storm"
	"github.com/spazbite187/sensornet"
)

// Data contains all the global configuration data including the loggers.
type Data struct {
	DIR, Assets, Version  string
	DB                    *storm.DB
	Log, ErrLog, DebugLog *log.Logger         // global logger, error logger and debug logger.
	CachedSensors         []*sensornet.Sensor // global cache
}
