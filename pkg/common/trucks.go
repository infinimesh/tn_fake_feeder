package common

import (
	"fmt"
	"sync"
	"time"

	"github.com/infinimesh/tn_fake_feeder/pkg/db"
)

type Truck struct {
	Uuid  string
	Point int64
	Speed time.Duration

	Move   func(p int64) (db.Point, int64)
	Report func(uuid string, report TruckReport)

	stop bool
	wg   *sync.WaitGroup
}

var TN_TIME_FORMAT = "2006-02-01T15:04:05.999Z"

type TruckReport struct {
	Gps  []float64 `json:"gps"`
	Sent string    `json:"sent"`
}

func (t *Truck) Stop() {
	t.stop = true
}

func (t *Truck) Start(wg *sync.WaitGroup) {
	t.stop = false
	t.wg = wg

	for {
		np, n := t.Move(t.Point)
		fmt.Printf("Moving Truck %s to Point %d(%.4f, %.4f)\n", t.Uuid, t.Point, np.Lat, np.Lng)
		t.Report(t.Uuid, TruckReport{
				Gps:    []float64{np.Lat, np.Lng},

		t.Point = n
		time.Sleep(t.Speed)
		if t.stop {
			break
		}
	}

	wg.Done()
}
