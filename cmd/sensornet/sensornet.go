// Command sensornet ...
package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spazbite187/sensornet/app"
	"github.com/spazbite187/sensornet/handlers"
	"github.com/spazbite187/sensornet/storage"
	"gopkg.in/gin-gonic/gin.v1"
)

const (
	version        = "1.0.0"
	dbFilename     = "sensornet.db"
	maxNumReadings = 10080 // 7 days worth if sensors measure every minute
	bckgrdInterval = 30    // seconds
)

var (
	appData *app.Data
	router  *gin.Engine
	hndlers handlers.Data
)

func init() {
	appData = &app.Data{
		Version: version,
		Log:     log.New(os.Stdout, "sensornet ", log.LstdFlags),
		ErrLog:  log.New(os.Stderr, "sensornet - error ", log.LstdFlags|log.Lshortfile),
		DbugLog: log.New(ioutil.Discard, "", 0),
	}
}

func main() {
	// setup aliases for logging
	gLog := appData.Log
	eLog := appData.ErrLog

	gLog.Println("version", version)

	mode := strings.ToLower(os.Getenv("MODE"))
	if mode == "debug" {
		appData.DbugLog = log.New(os.Stdout, "sensornet - debug ", log.LstdFlags|log.Lshortfile)
		appData.DbugLog.Printf("----- DEBUG MODE -----")
	} else {
		gin.SetMode(gin.ReleaseMode) // set gin to release mode
	}
	dLog := appData.DbugLog // alias for debug logger

	port := strings.ToLower(os.Getenv("PORT"))
	if port == "" {
		eLog.Fatalln("environment vairable PORT not defined")
	}

	appDir := strings.ToLower(os.Getenv("APPDIR"))
	if appDir == "" {
		appDir = "./"
	}

	// set remaining config
	appData.DIR = appDir
	appData.Assets = appData.DIR + "assets/"

	// hndlers is used for passing global data to the handlers
	hndlers.Data = appData

	db, err := storage.GetDatabase(dbFilename)
	if err != nil {
		eLog.Fatal(err)
	}
	defer db.Close()
	appData.DB = db

	// setup web app with a router using some built in gin middleware
	router = gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// add custom middlewares here

	// load html templates
	router.LoadHTMLGlob(appData.Assets + "/templates/*")

	// Initialize the routes
	initRoutes()

	// display run data
	gLog.Printf("version %s starting service on port %s", appData.Version, port)
	dLog.Printf("assets dir %s", appData.Assets)
	dLog.Printf("app dir %s", appData.DIR)

	// start background service for updating sensor readings
	go func() {
		// loop
		for {
			time.Sleep(bckgrdInterval / 2 * time.Second) // udpate half the bckgrdInterval
			dLog.Println("updating sensor number of readings...")
			err := storage.UpdateNumReadings(db)
			if err != nil {
				eLog.Printf("updating sensors - %s", err.Error())
			}
			dLog.Println("done updating sensor number of readings.")
		}
	}()

	// start background service for cleaning up db
	go func() {
		// loop
		for {
			time.Sleep(bckgrdInterval * time.Second) // interval
			dLog.Printf("removing sensor data over max  num (%d)...", maxNumReadings)
			err = storage.CleanDB(maxNumReadings, db)
			if err != nil {
				eLog.Printf("cleaning up db - %s", err.Error())
			}
			dLog.Println("done cleaning up db.")

		}
	}()

	// start service
	if err := router.Run(); err != nil {
		eLog.Fatal(err)
	}
}
