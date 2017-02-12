// Command sensornet ...
package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spazbite187/sensornet"
	"github.com/spazbite187/sensornet/app"
	"github.com/spazbite187/sensornet/handlers"
	"github.com/spazbite187/sensornet/storage"
	"gopkg.in/gin-gonic/gin.v1"
)

var (
	appData *app.Data
	router  *gin.Engine
	hndlers handlers.Data
)

func init() {
	appData = &app.Data{
		Version:       sensornet.Version,
		Log:           log.New(os.Stdout, "[SensorNet] ", log.LstdFlags),
		ErrLog:        log.New(os.Stderr, "[SensorNet-error] ", log.LstdFlags|log.Llongfile),
		DebugLog:      log.New(ioutil.Discard, "", 0),
		CachedSensors: []*sensornet.Sensor{},
	}
}

func main() {
	// setup aliases for logging
	gLog := appData.Log
	eLog := appData.ErrLog

	mode := strings.ToLower(os.Getenv("MODE"))
	if mode == "debug" {
		appData.DebugLog = log.New(os.Stdout, "[SensorNet-debug] ", log.LstdFlags|log.Lshortfile)
		gLog.Printf("----- DEBUG MODE -----")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	dLog := appData.DebugLog // alias for debug logger

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

	db, err := storage.GetDatabase(sensornet.DBFilename)
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
	dLog.Printf("assets dir %s", appData.Assets)
	dLog.Printf("app dir %s", appData.DIR)

	// start background services
	bckgrndServices(appData)

	// start service
	initRoutes()
	gLog.Printf("running web service on port %s", port)
	if err := router.Run(); err != nil {
		eLog.Fatal(err)
	}
}

func bckgrndServices(appData *app.Data) {
	// start background service for updating sensor data
	go func() {
		// loop
		for {
			appData.DebugLog.Println("updating the sensor models...")
			err := storage.UpdateSensors(appData)
			if err != nil {
				appData.ErrLog.Printf("updating sensors - %s", err.Error())
			}
			appData.DebugLog.Println("completed updating the sensor models")
			time.Sleep(sensornet.BckgrndInterval * time.Second)
		}
	}()

	// start background service for cleaning up db
	go func() {
		// loop
		for {
			time.Sleep(sensornet.BckgrndInterval * time.Second)
			appData.DebugLog.Printf("removing sensor data over the max number (%d)...", sensornet.MaxNumReadings)
			err := storage.CleanDB(sensornet.MaxNumReadings, appData.DB)
			if err != nil {
				appData.ErrLog.Printf("cleaning up db - %s", err.Error())
			}
			appData.DebugLog.Printf("completed removing sensor data over the max number (%d)", sensornet.MaxNumReadings)
		}
	}()

}
