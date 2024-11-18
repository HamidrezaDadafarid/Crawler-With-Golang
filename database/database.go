package database

import (
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var lock sync.Mutex
var instance *DatabaseManager

type DatabaseManager struct {
	Db *gorm.DB
}

func GetInstnace() *DatabaseManager {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			dsn := "host=localhost user=postgres password=1234 dbname=crawlerdb port=5432 sslmode=disable TimeZone=Asia/Tehran"

			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
			if err != nil {
				log.Fatal("Failed to connect to database. \n", err)
			}

			log.Println("connected")

			instance = &DatabaseManager{
				Db: db,
			}
		}
	}

	return instance
}
