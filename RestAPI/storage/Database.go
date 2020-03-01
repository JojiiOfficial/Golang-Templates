package storage

import (
	"Golang-Templates/RestAPI/models"
	"strconv"

	log "github.com/sirupsen/logrus"

	dbhelper "github.com/JojiiOfficial/GoDBHelper"
)

//ConnectDB connects to MySQL
func ConnectDB(config *models.Config, isDebug, nocolor bool) (*dbhelper.DBhelper, error) {
	log.Debug("Connecting to DB")
	db, err := dbhelper.NewDBHelper(dbhelper.Mysql).Open(
		config.Server.Database.Username,
		config.Server.Database.Pass,
		config.Server.Database.Host,
		strconv.Itoa(config.Server.Database.DatabasePort),
		config.Server.Database.Database,
		"parseTime=True",
	)

	if err != nil {
		return nil, err
	}

	log.Info("Connected successfully")

	//Only debugMode if logLevel is debug
	db.Options.Debug = isDebug
	db.Options.UseColors = !nocolor

	return db, updateDB(db)
}

func updateDB(db *dbhelper.DBhelper) error {
	db.AddQueryChain(getInitSQL())
	return db.RunUpdate()
}

func getInitSQL() dbhelper.QueryChain {
	return dbhelper.QueryChain{
		Name:    "initChain",
		Order:   0,
		Queries: dbhelper.CreateInitVersionSQL(),
	}
}
