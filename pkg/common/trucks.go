package common

import (
	"sync"
	"time"
)

type Truck struct {
	Uuid  string
	Point int64
	Speed time.Duration

	stop bool
	wg   *sync.WaitGroup
}

func (t *Truck) Stop() {
	t.stop = true
}

func (t *Truck) Start(wg *sync.WaitGroup) {
	t.stop = false
	t.wg = wg

	wg.Done()
}
