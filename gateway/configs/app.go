package configs

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port     string
	Env      string
	GRPCPort string
}

func LoadAppConfig() *AppConfig {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Print("Error reading .env file")
	}

	config := &AppConfig{
		Port:     viper.GetString("APP_PORT"),
		Env:      viper.GetString("APP_ENV"),
		GRPCPort: viper.GetString("GRPC_PORT"),
	}

	return config
}
