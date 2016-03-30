package main

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type CheckRecord struct {
	Id        int       `json:"id"`
	RunId     int64     `json:"run_id"`
	Service   string    `json:"service"`
	Check     string    `json:"check"`
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Duration  float64   `json:"duration"`
	Timestamp time.Time `json:"timestamp"`
}

var DB *gorm.DB

func dbSetup() error {
	db, err := gorm.Open("sqlite3", DbPath)
	if err != nil {
		return err
	}

	model := &CheckRecord{}

	db.AutoMigrate(model)
	db.Model(model).AddIndex("idx_service", "service")
	db.Model(model).AddIndex("idx_service", "service", "timestamp")

	DB = db
	return nil
}
