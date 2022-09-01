package common

import (
	"fmt"
	"sync"
	"time"

	"github.com/infinimesh/tn_fake_feeder/pkg/db"
	wr "github.com/mroth/weightedrand"
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
	Gps    []float64 `json:"gps"`
	Sent   string    `json:"sent"`
	Status string    `json:"status"`
}

func (t *Truck) Stop() {
	t.stop = true
}

func (t *Truck) Start(wg *sync.WaitGroup) {
	t.stop = false
	t.wg = wg

	status := t.Status()
	hold := 0
	for {

		if status.SkipMove {
			t.Point -= 1
		}

		np, n := t.Move(t.Point)
		fmt.Printf("Moving Truck %s to Point %d(%.4f, %.4f)\n", t.Uuid, t.Point, np.Lat, np.Lng)
		t.Point = n

		if status.Report {
			t.Report(t.Uuid, TruckReport{
				Gps:    []float64{np.Lat, np.Lng},
				Sent:   time.Now().Format(TN_TIME_FORMAT),
				Status: status.Key,
			})
		} else {
			fmt.Printf("Truck(%s) is %s, won't report\n", t.Uuid, status.Key)
		}

		time.Sleep(t.Speed)
		if t.stop {
			break
		}
		if status.HoldStatus > hold {
			fmt.Printf("Truck(%s) is in status %s, holding %d(%d left)\n", t.Uuid, status.Key, status.HoldStatus, status.HoldStatus-hold)
			hold++
		} else {
			status = t.Status()
			hold = 0
		}
	}

	wg.Done()
}

type StatusWithProb struct {
	Key        string
	SkipMove   bool
	Report     bool
	HoldStatus int
}

var status_chooser, _ = wr.NewChooser(
	wr.Choice[StatusWithProb]{Item: StatusWithProb{"online", false, true, 0}, Weight: 80},
	wr.Choice[StatusWithProb]{Item: StatusWithProb{"inactive", true, true, 20}, Weight: 10},
	wr.Choice[StatusWithProb]{Item: StatusWithProb{"offline", false, true, 10}, Weight: 7},
	wr.Choice[StatusWithProb]{Item: StatusWithProb{"dead_offline", true, false, 10}, Weight: 3},
)

func (t *Truck) Status() StatusWithProb {
	s := status_chooser.Pick()
	return s
}
