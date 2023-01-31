package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
	"os"
)

var host = os.Getenv("DATABASE_HOST")
var username = os.Getenv("DATABASE_USERNAME")
var password = os.Getenv("DATABASE_PASSWORD")
var dbName = os.Getenv("DATABASE_DBNAME")
var timezone = os.Getenv("DATABASE_TIMEZONE")

func New() (*gorm.DB, error) {
	dsn := "host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=%s"
	dsn = fmt.Sprintf(dsn, host, username, password, dbName, timezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "mancala-db",
		RefreshInterval: 15,
		StartServer:     true,
		HTTPServerPort:  8080,
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to connect gorm to prometheus: %v", err)
	}

	return db, nil
}
