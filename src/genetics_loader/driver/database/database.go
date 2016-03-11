package database

import (
	"fmt"
	"genetics_loader/configuration"
	"github.com/dabfleming/gorm"
	_ "github.com/lib/pq"
	"log"
)

// TODO Rename this connection "Intake" for clarity?
var App *gorm.DB

// Connection URL in the format used by github.com/mattes/migrate/migrate
var ConnectionURL string

func init() {
	var appDBConnection gorm.DB
	var config *configuration.Configuration

	config = configuration.GetConfiguration()

	ConnectionURL = getConnectionURL(config.AppDatabase)

	appDBConnection = connect(config.AppDatabase)
	App = &appDBConnection
}

func getConnectionURL(config map[string]interface{}) string {
	return fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s&password=%s", config["user"].(string), config["host"].(string), config["port"].(string), config["dbname"].(string), config["sslmode"].(string), config["password"].(string))
}

func connect(config map[string]interface{}) gorm.DB {
	dbConnection, err := gorm.Open("postgres", "user="+config["user"].(string)+
		" host="+config["host"].(string)+
		" port="+config["port"].(string)+
		" password="+config["password"].(string)+
		" dbname="+config["dbname"].(string)+
		" sslmode="+config["sslmode"].(string))
	if err != nil {
		log.Fatalf("Could not connect to database: ", err)
	}

	if config["log_to_console"] != nil {
		dbConnection.LogMode(config["log_to_console"].(bool))
	}

	return dbConnection
}
