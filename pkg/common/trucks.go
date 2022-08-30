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

	Move func(p int64) (db.Point, int64)

	stop bool
	wg   *sync.WaitGroup
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

		t.Point = n
		time.Sleep(t.Speed)
		if t.stop {
			break
		}
	}

	wg.Done()
}
