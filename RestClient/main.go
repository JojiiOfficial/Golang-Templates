package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"Golang-Templates/RestClient/constants"
	"Golang-Templates/RestClient/models"
	"Golang-Templates/RestClient/server"

	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	appName = "clapp"
	version = "1.0.0"

	//EnVarPrefix prefix for env vars
	EnVarPrefix = "GOLOG"

	//Datapath the default path for data files
	Datapath = "./data"
	//DefaultConfig the default config file
	DefaultConfig = "config.yaml"
)

var (
	//DefaultConfigPath default config path
	DefaultConfigPath = path.Join(Datapath, DefaultConfig)
)

//App commands
var (
	app          = kingpin.New(appName, "A Logging server")
	appLogLevel  = app.Flag("log-level", "Enable debug mode").HintOptions(constants.LogLevels...).Envar(getEnVar(EnVarLogLevel)).Short('l').Default(constants.LogLevels[2]).String()
	appNoColor   = app.Flag("no-color", "Disable colors").Envar(getEnVar(EnVarNoColor)).Bool()
	appYes       = app.Flag("yes", "Skips confirmations").Short('y').Envar(getEnVar(EnVarYes)).Bool()
	appVerbosity = app.Flag("verbose", "Set the verbosity level").Short('v').Counter()
	appCfgFile   = app.
			Flag("config", "the configuration file for the app").
			Envar(getEnVar(EnVarConfigFile)).
			Short('c').String()

	//Commands

	//Server start
	pingCmd = app.Command("ping", "Test the connection to the server")
)

var (
	config  *models.Config
	isDebug = false
)

func main() {
	app.HelpFlag.Short('h')
	app.Version(version)

	//parsing the args
	parsed := kingpin.MustParse(app.Parse(os.Args[1:]))

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: false,
		TimestampFormat:  time.Stamp,
		FullTimestamp:    true,
		ForceColors:      !*appNoColor,
		DisableColors:    *appNoColor,
	})

	log.Debugf("LogLevel: %s\n", *appLogLevel)

	//set app logLevel
	switch *appLogLevel {
	case constants.LogLevels[0]:
		//Debug
		log.SetLevel(log.DebugLevel)
		isDebug = true
	case constants.LogLevels[1]:
		//Info
		log.SetLevel(log.InfoLevel)
	case constants.LogLevels[2]:
		//Warning
		log.SetLevel(log.WarnLevel)
	case constants.LogLevels[3]:
		//Error
		log.SetLevel(log.ErrorLevel)
	default:
		fmt.Println("LogLevel not found!")
		os.Exit(1)
		return
	}

	//print verbosity level if greater than 1
	if *appVerbosity > 1 {
		log.Debugf("Verbosity set to %d\n", *appVerbosity)
	}

	//Init config
	var err error
	config, err = models.InitConfig(DefaultConfigPath, *appCfgFile)
	if err != nil {
		log.Error(err)
		return
	}

	if config == nil {
		log.Info("New config cerated")
		return
	}

	//Run specified command
	switch parsed {
	case pingCmd.FullCommand():
		pingServer(config)
	}
}

//Env vars
const (
	//EnVarPrefix prefix of all used env vars
	EnVarLogLevel   = "LOG_LEVEL"
	EnVarNoColor    = "NO_COLOR"
	EnVarYes        = "SKIP_CONFIRM"
	EnVarConfigFile = "CONFIG"
)

//Return the variable using the server prefix
func getEnVar(name string) string {
	return fmt.Sprintf("%s_%s", EnVarPrefix, name)
}

func pingServer(config *models.Config) {
	var response server.StringResponse
	res, err := server.
		NewRequest(server.EPPing, server.PingRequest{Payload: "ping"}, config).
		Do(&response)

	if err != nil {
		log.Error(err.Error())
		return
	}

	if res.Status == server.ResponseSuccess {
		fmt.Println("Ping success:", response.String)
	} else {
		log.Errorf("Error (%d) %s\n", res.HTTPCode, res.Message)
	}
}
