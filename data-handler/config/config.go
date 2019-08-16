package config

import (
	"flag"
	"github.com/spf13/viper"
	"log"
)

type AppConfig struct {
	InputFile        string
	DataBaseUrl      string
	DataBaseUserName string
	DataBasePassword string
}

var appConfig *AppConfig

func init() {
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "config.yml", "absolute path to the configuration file")
	flag.Parse()
	viper.SetConfigFile(configFilePath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("error reading config file")
	}
	appConfig = CreateAppConfig()
}

func GetConfig() *AppConfig {
	return appConfig
}

func CreateAppConfig() *AppConfig {
	config := new(AppConfig)
	config.InputFile = viper.GetString("inputFile")
	config.DataBaseUrl = viper.GetString("database.url")
	config.DataBaseUserName = viper.GetString("database.username")
	config.DataBasePassword = viper.GetString("database.password")
	return config
}
