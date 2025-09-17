package db

import (
	"auth-service/configs"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewGormDB creates a new GORM DB connection using DBConfig.
func NewGormDB(cfg *configs.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
