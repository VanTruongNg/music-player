package configs

import (
	"log"

	"github.com/spf13/viper"
)

// DBConfig holds PostgreSQL database configuration.
type DBConfig struct {
	Host     string // Database host
	Port     string // Database port
	User     string // Database user
	Password string // Database password
	DBName   string // Database name
}

// LoadDBConfig loads PostgreSQL configuration using viper.
func LoadDBConfig() *DBConfig {
	cfg := &DBConfig{
		Host:     viper.GetString("POSTGRES_HOST"),
		Port:     viper.GetString("POSTGRES_PORT"),
		User:     viper.GetString("POSTGRES_USER"),
		Password: viper.GetString("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("POSTGRES_DB"),
	}
	// Optionally log warning if any field is empty
	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.Password == "" || cfg.DBName == "" {
		log.Printf("[WARN] Some database config fields are empty. Please check your environment variables or .env file.")
	}
	return cfg
}
