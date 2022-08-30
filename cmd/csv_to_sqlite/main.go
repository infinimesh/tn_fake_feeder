package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/infinimesh/tn_fake_feeder/pkg/db"
)

func main() {
	fmt.Println("Checking track.csv")
	track_file, err := os.Open("track.csv")
	if err != nil {
		panic(err)
	}

	r := csv.NewReader(track_file)
	header, _ := r.Read()
	fmt.Println("Header", header)

	for {
		line, err := r.Read()
		if err != nil {
			fmt.Println("EOF or different error", err.Error())
			return
		}
		lat, _ := strconv.ParseFloat(line[0], 64)
		lng, _ := strconv.ParseFloat(line[1], 64)
		db.DB.Create(&db.Point{
			Lat: lat, Lng: lng,
		})
	}
}
