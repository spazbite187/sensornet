// Package app contains the app specific data structure.
package app

import (
	"log"

	"github.com/asdine/storm"
)

// Data contains all the global configuration data including the loggers.
type Data struct {
	DIR, Assets, Version string
	DB                   *storm.DB
	Log, ErrLog, DbugLog *log.Logger // global logger, error logger and debug logger.
}
