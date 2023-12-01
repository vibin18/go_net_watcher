package database

import (
	"go_net_watcher/internal/netwatcher"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type SQL struct {
	Db *gorm.DB
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

	db.AutoMigrate(&netwatcher.NetDevices{})

	Database = SQL{
		Db: db,
	}
}
