package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Point struct {
	Lat float64
	Lng float64

	gorm.Model
}

func (Point) Count() int64 {
	var count int64
	DB.Model(&Point{}).Count(&count)
	return count
}

func init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("track.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&Point{})
}
