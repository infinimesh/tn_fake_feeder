package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Point struct {
	Lat string
	Lng string

	gorm.Model
}

func init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("track.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&Point{})
}
