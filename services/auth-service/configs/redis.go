package configs

import (
	"log"

	"github.com/spf13/viper"
)

// RedisConfig holds Redis configuration.
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	Username string 
}

// LoadRedisConfig loads Redis configuration using viper.
func LoadRedisConfig() *RedisConfig {
	cfg := &RedisConfig{
		Host:     viper.GetString("REDIS_HOST"),
		Port:     viper.GetString("REDIS_PORT"),
		Password: viper.GetString("REDIS_PASSWORD"),
		Username: viper.GetString("REDIS_USERNAME"),
	}
	// Optionally log warning if any field is empty
	if cfg.Host == "" || cfg.Port == "" {
		log.Printf("[WARN] Some Redis config fields are empty. Please check your environment variables or .env file.")
	}
	return cfg
}
