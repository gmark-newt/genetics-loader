package configuration

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Env         string `json:"env"`
	Environment string `json:"environment"`
	Home        string `json:"home"`

	Debug       bool                   `json:"debug_mode"`
	AppDatabase map[string]interface{} `json:"database_app"`
}

var config *Configuration = nil

func init() {
	config_location := os.Getenv("INTAKE_CONFIG")
	if config_location == "" {
		config_location = ".newtopia.json"
	}

	log.Print("Loading configuration from: ", config_location)

	file, err := os.Open(config_location)
	if err != nil {
		log.Fatalf("Error opening config file: %v.\nError: %v", config_location, err)
	}
	decoder := json.NewDecoder(file)

	config = &Configuration{}
	err = decoder.Decode(config)
	if err != nil {
		log.Fatal("Error here ? loading config: ", err)
	}

	log.Printf("Configuration loaded.\n\tEnv: %v\n\tEnvironment: %v", config.Env, config.Environment)
}

func GetConfiguration() *Configuration {
	return config
}

func IsDebug() bool {
	return config.Debug
}
