package database

import (
	"log"
	"os"

	"github.com/mayank12gt/ProgressTracker/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDB() {

	db, err := gorm.Open(sqlite.Open("progress_tracker.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed")
		os.Exit(2)
	}
	db = db.Exec("PRAGMA foreign_keys = ON")

	log.Println("Connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Migrations")

	db.AutoMigrate(&model.List{}, &model.Item{}, &model.User{})

	Database = DbInstance{Db: db}
}
