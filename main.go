package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"strconv"
	"sync"
	"time"

	pb "github.com/infinimesh/proto/node"
	devpb "github.com/infinimesh/proto/node/devices"
	"github.com/infinimesh/proto/node/namespaces"
	"github.com/infinimesh/proto/shadow"
	"github.com/infinimesh/tn_fake_feeder/pkg/common"
	"github.com/infinimesh/tn_fake_feeder/pkg/db"
	faker "github.com/jaswdr/faker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Host     string `yaml:"infinimesh"`
	Insecure bool   `yaml:"insecure"`
	Token    string `yaml:"token"`
}

var (
	fk faker.Faker
)

var country_codes = []string{
	"D", "NL", "PL", "FIN", "BY",
}

func init() {
	rand.Seed(time.Now().UnixNano())

	fk = faker.New()
}

func main() {
	n_trucks := 1
	if len(os.Args) < 2 {
		panic("Not enough arguments, namespace and optionally number of simulated trucks must be passed(defaults to 1)")
	}
	namespace := os.Args[1]
	if len(os.Args) > 2 {
		i, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Passed amount of trucks is not an integer, exiting")
			return
		}
		n_trucks = i
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	config_path := path.Join(home, ".default.infinimesh.yaml")
	fmt.Printf("Config used from path: %s\n", config_path)

	config_bytes, err := os.ReadFile(config_path)
	if err != nil {
		fmt.Println("Error reading infinimesh CLI config, make sure you have it installed and logged in(`inf login`)")
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(config_bytes, &config)
	if err != nil {
		fmt.Println("Error parsing infinimesh CLI config")
		panic(err)
	}

	fmt.Printf("Connecting to infinimesh at %s...\n", config.Host)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+config.Token)
	var opts []grpc.DialOption
	if config.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}
	conn, err := grpc.DialContext(ctx, config.Host, opts...)
	if err != nil {
		fmt.Println("Error Dialing infinimesh")
		panic(err)
	}
	defer conn.Close()

	nsc := pb.NewNamespacesServiceClient(conn)
	ns, err := nsc.Get(ctx, &namespaces.Namespace{
		Uuid: namespace,
	})
	if err != nil {
		panic(fmt.Errorf("error retrieving Namespace: %v", err))
	}
	if ns.Access.Level < 3 {
		panic("Not enough Access Rights to Create Devices in choosen Namespace")
	}

	devc := pb.NewDevicesServiceClient(conn)
	shad := pb.NewShadowServiceClient(conn)

	fmt.Println("gRPC Connection Established")

	rows := db.Point.Count(db.Point{})
	fmt.Printf("Amount of waypoints found: %d\n", rows)

	var retrieve_func = func(p int64) (r db.Point, _ int64) {
		db.DB.First(&r, p)
		if p >= rows {
			return r, 1
		}
		return r, p + 1
	}

	var postCtx context.Context
	var report_func = func(uuid string, report common.TruckReport) {
		jsonb, _ := json.Marshal(report)
		data := &structpb.Struct{}
		data.UnmarshalJSON(jsonb)

		_, err := shad.Patch(postCtx, &shadow.Shadow{
			Device: uuid,
			Reported: &shadow.State{
				Data: data,
			},
		})
		if err != nil {
			fmt.Printf("[WARN] Couldn't patch Reported state for device %s: %v", uuid, err)
		}
	}

	pool := make([]*common.Truck, n_trucks)
	uuids := make([]string, n_trucks)

	for i := 0; i < n_trucks; i++ {

		country := fk.RandomStringElement(country_codes)

		tags := []string{
			"tn:simulated",
			fmt.Sprintf("tn:number_plate_truck:%s_%s", country, fk.Car().Plate()),
		}

		if rand.Intn(10)%2 == 0 {
			tags = append(tags, fmt.Sprintf("tn:number_plate_trailer:%s_%s", country, fk.Car().Plate()))
		}

		res, err := devc.Create(ctx, &devpb.CreateRequest{
			Device: &devpb.Device{
				Title:   fmt.Sprintf("sim-truck-%d", i),
				Enabled: true,
				Tags:    tags,
				Certificate: &devpb.Certificate{
					PemData: common.SAMPLE_CERT,
				},
			},
			Namespace: namespace,
		})
		if err != nil {
			panic(err)
		}
		defer func(dev *devpb.Device) {
			fmt.Printf("Deleting device %s(%s)\n", dev.Title, dev.Uuid)
			_, err := devc.Delete(ctx, dev)
			if err != nil {
				fmt.Printf("[WARN] Error deleting device: %v", err)
			}
		}(res.GetDevice())

		truck := &common.Truck{
			Uuid:  res.GetDevice().GetUuid(),
			Point: rand.Int63n(rows),
			Speed: time.Duration(rand.Intn(5)+1) * time.Second,

			Move:   retrieve_func,
			Report: report_func,
		}

		pool[i] = truck
		uuids[i] = truck.Uuid
	}

	res, err := devc.MakeDevicesToken(ctx, &pb.DevicesTokenRequest{
		Devices: uuids,
		Post:    true,
	})
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return
	}

	postCtx = metadata.AppendToOutgoingContext(context.Background(), "Authorization", "Bearer "+res.GetToken())

	wg := sync.WaitGroup{}
	for _, truck := range pool {
		wg.Add(1)
		go truck.Start(&wg)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	fmt.Println("\nInterrupt caught, gracefully stopping trucks")

	for _, truck := range pool {
		truck.Stop()
	}

	wg.Wait()

	fmt.Println("Done")
}
