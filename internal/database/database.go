package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type SQL struct {
	Db *gorm.DB
}

type Device struct {
	MAC  string `json:"mac"`
	IP   string `json:"ip"`
	Name string `json:"name"`
	ID   uint   `json:"id" gorm:"primaryKey"`
}

var Database SQL

func ConnectDB() {
	db, err := gorm.Open(
		sqlite.Open("net_watcher.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err.Error())
	}

	log.Printf("Conected to db sucessfully!")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Printf("Running DB migrations..")

	db.AutoMigrate(&Device{})

	Database = SQL{
		Db: db,
	}
}
