package storage

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	return db
}

func InitDB(dsn string) (*gorm.DB, error) {
	var err error
	once.Do(func() {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// Configure the Connection Pool
			sqlDB, _ := db.DB()
			sqlDB.SetMaxIdleConns(10)   // Keep 10 connections idle
			sqlDB.SetMaxOpenConns(100)  // Max 100 concurrent connections
			sqlDB.SetConnMaxLifetime(0) // Reuse connections forever
		}
	})
	return db, err
}

type PostgresRepository struct {
	db *gorm.DB
}


