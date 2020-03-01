package main

import (
	"Golang-Templates/RestAPI/services"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	dbhelper "github.com/JojiiOfficial/GoDBHelper"
)

//Services
var (
	apiService *services.APIService //Handle endpoints
)

func startAPI() {
	log.Info("Starting version " + version)

	//Setting up database
	db.SetErrHook(func(err error, query, prefix string) {
		logMessage := prefix + query + ": " + err.Error()

		//Warn only on production
		if isDebug {
			log.Error(logMessage)
		} else {
			log.Warn(logMessage)
		}
	}, dbhelper.ErrHookOptions{
		Prefix:         "Query: ",
		ReturnNilOnErr: false,
	})

	//Create the APIService and start it
	apiService = services.NewAPIService(db, config)
	apiService.Start()

	//Startup done
	log.Info("Startup completed")

	//Start loop to tick the services
	go (func() {
		for {
			time.Sleep(time.Hour)

		}
	})()

	awaitExit(db, apiService)
}

//Shutdown server gracefully
func awaitExit(db *dbhelper.DBhelper, httpServer *services.APIService) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)

	// await os signal
	<-signalChan

	// Create a deadline for the await
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	log.Info("Shutting down server")

	if httpServer.HTTPServer != nil {
		httpServer.HTTPServer.Shutdown(ctx)
		log.Info("HTTP server shutdown complete")
	}

	if httpServer.HTTPTLSServer != nil {
		httpServer.HTTPTLSServer.Shutdown(ctx)
		log.Info("HTTPs server shutdown complete")
	}

	if db != nil && db.DB != nil {
		db.DB.Close()
		log.Info("Database shutdown complete")
	}

	log.Info("Shutting down complete")
	os.Exit(0)
}
