package common

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/color"

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

type TruckReport struct {
	Sent string `json:"sent"`

	Gps   []float64 `json:"gps"`
	Sats  int       `json:"sats"`
	Speed int       `json:"speed"`

	Status string `json:"status"`

	Cell []string `json:"cell"`

	Sensors map[string]bool `json:"sensors"`
}

func (t *Truck) Stop() {
	t.stop = true
}

func init() {
	rand.New(
		rand.NewSource(time.Now().UnixNano()),
	)

	go func() {
		for now := range time.NewTicker(time.Hour * 24).C {
			rand.New(
				rand.NewSource(now.UnixNano()),
			)
		}
	}()
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

			sensors := make(map[string]bool)

			if left := sensor_chooser.Pick(); left != nil {
				sensors["left"] = *left
			}
			if right := sensor_chooser.Pick(); right != nil {
				sensors["right"] = *right
			}

			t.Report(t.Uuid, TruckReport{
				Gps:    []float64{np.Lat, np.Lng},
				Sent:   time.Now().Format(time.RFC3339),
				Status: status.Key,

				Sats:  sats_chooser.Pick(),
				Cell:  t.Cell(),
				Speed: int(120 * 1 / t.Speed.Seconds()),

				Sensors: sensors,
			})
		} else {
			fmt.Printf("Truck(%s) is %s, won't report\n", t.Uuid, status.Output(status.Key))
		}

		time.Sleep(t.Speed)
		if t.stop {
			break
		}
		if status.HoldStatus > hold {
			fmt.Printf("Truck(%s) is in status %s, holding %d(%d left)\n", t.Uuid, status.Output(status.Key), status.HoldStatus, status.HoldStatus-hold)
			hold++
		} else {
			status = t.Status()
			hold = 0
		}

		if acceleration := rand.Intn(2) - 1; acceleration != 0 {
			t.Speed = t.Speed + time.Duration(acceleration)*time.Second
			if t.Speed == 0 || t.Speed > time.Second*5 {
				t.Speed = time.Second * 3
			}
		}
	}

	wg.Done()
}

type StatusWithProb struct {
	Key        string
	Output     func(a ...interface{}) string
	SkipMove   bool
	Report     bool
	HoldStatus int
}

var status_chooser, _ = wr.NewChooser(
	wr.Choice[StatusWithProb]{Item: StatusWithProb{"online", color.New(color.FgHiGreen).SprintFunc(), false, true, 1}, Weight: 80},
	wr.Choice[StatusWithProb]{Item: StatusWithProb{"inactive", color.New(color.FgYellow).SprintFunc(), true, true, 20}, Weight: 10},
	wr.Choice[StatusWithProb]{Item: StatusWithProb{"offline", color.New(color.FgRed).SprintFunc(), false, true, 10}, Weight: 7},
	wr.Choice[StatusWithProb]{Item: StatusWithProb{"dead_offline", color.New(color.FgRed, color.BgWhite).SprintFunc(), true, false, 10}, Weight: 3},
)

func (t *Truck) Status() StatusWithProb {
	s := status_chooser.Pick()
	return s
}

var cell_chooser, _ = wr.NewChooser(
	wr.Choice[string]{Item: "Telekom", Weight: 944},
	wr.Choice[string]{Item: "vodafone.de", Weight: 913},
	wr.Choice[string]{Item: "O2", Weight: 874},
)

var cell_protocol_chooser, _ = wr.NewChooser(
	wr.Choice[string]{Item: "offline", Weight: 5},
	wr.Choice[string]{Item: "5G", Weight: 30},
	wr.Choice[string]{Item: "4G", Weight: 40},
	wr.Choice[string]{Item: "3G", Weight: 20},
	wr.Choice[string]{Item: "2G", Weight: 5},
)

var (
	_true  = true
	_false = false
)

var sensor_chooser, _ = wr.NewChooser(
	wr.Choice[*bool]{Item: nil, Weight: 5},
	wr.Choice[*bool]{Item: &_false, Weight: 10},
	wr.Choice[*bool]{Item: &_true, Weight: 85},
)

func (t *Truck) Cell() []string {
	proto := cell_protocol_chooser.Pick()
	if proto == "offline" {
		return []string{}
	}

	return []string{
		proto,
		cell_chooser.Pick(),
		strconv.Itoa(rand.Intn(1200) - 600),
	}
}

var sats_chooser, _ = wr.NewChooser(
	wr.Choice[int]{Item: 0, Weight: 10},
	wr.Choice[int]{Item: 2, Weight: 20},
	wr.Choice[int]{Item: 3, Weight: 30},
	wr.Choice[int]{Item: 4, Weight: 30},
	wr.Choice[int]{Item: 8, Weight: 20},
)
