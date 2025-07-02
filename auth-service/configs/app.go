package configs

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port string
	Env  string
}

// LoadAppConfig loads application configuration using viper.
func LoadAppConfig() *AppConfig {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Read config file (optional)
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("[INFO] No .env file found or error reading config: %v", err)
	}

	cfg := &AppConfig{
		Port: viper.GetString("APP_PORT"),
		Env:  viper.GetString("APP_ENV"),
	}
	return cfg
}
