package util

import (
	"log"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
// Viper uses the mapstructure package under the hood for unmarshal of values.
type Config struct {
	Environment       string `mapstructure:"ENVIRONMENT"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	MigrationURL      string `mapstructure:"MIGRATION_URL"`
	LogLevel          string `mapstructure:"LOG_LEVEL"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	// Tell Viper the location of the config file.
	viper.AddConfigPath(path)

	// Tell Viper what file name (e.g app) to look for when loading config file.
	viper.SetConfigName("app")

	// Tell Viper what type of configuration to load, e.g env, json, yaml etc.
	viper.SetConfigType("env")

	// Tell Viper to automatically override values that it read from config
	// file with the values of the corresponding environment variables.
	viper.AutomaticEnv()

	// Set some decent defaults
	viper.SetDefault("ENVIRONMENT", "production")
	viper.SetDefault("DB_SOURCE", "postgresql://postgres:postgres@localhost:5432/bankdb?sslmode=disable")
	viper.SetDefault("MIGRATION_URL", "file://migrations")
	viper.SetDefault("LOG_LEVEL", "INFO")
	viper.SetDefault("HTTP_SERVER_ADDRESS", "0.0.0.0:8080")

	// Tell Viper to read config from file.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("load config: optional config file not provided")
		} else {
			log.Fatalln("load config:", err)
			// Config file was found but another error was produced
		}
	}

	err = viper.Unmarshal(&config)

	return
}
