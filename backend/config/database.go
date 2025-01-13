package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=yourpassword dbname=carbon_calculator port=5432 sslmode=disable"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
} 