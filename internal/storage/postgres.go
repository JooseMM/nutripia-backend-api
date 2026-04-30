package storage

import (
	"sync"

	bodyMeasurementTypes "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	return db
}

func InitDB(dsn string) (*gorm.DB, error) {
	var errBuffer error
	once.Do(func() {
		db, errBuffer = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if errBuffer == nil {
			// Configure the Connection Pool
			sqlDB, _ := db.DB()
			sqlDB.SetMaxIdleConns(10)   // Keep 10 connections idle
			sqlDB.SetMaxOpenConns(100)  // Max 100 concurrent connections
			sqlDB.SetConnMaxLifetime(0) // Reuse connections forever


			if err := db.AutoMigrate(&bodyMeasurementTypes.BodyMeasurement{}); err != nil {
				errBuffer = err
				return
			}
		}
	})
	return db, errBuffer
}

type PostgresRepository struct {
	db *gorm.DB
}
