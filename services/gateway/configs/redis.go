package configs

import (
	"log"

	"github.com/spf13/viper"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	Username string
}

func LoadRedisConfig() *RedisConfig {
	cfg := &RedisConfig{
		Host:     viper.GetString("REDIS_HOST"),
		Port:     viper.GetString("REDIS_PORT"),
		Password: viper.GetString("REDIS_PASSWORD"),
		Username: viper.GetString("REDIS_USERNAME"),
	}
	if cfg.Host == "" || cfg.Port == "" {
		log.Printf("[WARN] Some Redis config fields are empty. Please check your environment variables or .env file.")
	}
	return cfg
}
