package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/infinimesh/tn_fake_feeder/pkg/common"
	"github.com/infinimesh/tn_fake_feeder/pkg/db"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v2"

	pb "github.com/infinimesh/proto/node"
	devpb "github.com/infinimesh/proto/node/devices"
	"github.com/infinimesh/proto/node/namespaces"
)

type Config struct {
	Host     string `yaml:"infinimesh"`
	Insecure bool   `yaml:"insecure"`
	Token    string `yaml:"token"`
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
		panic(fmt.Errorf("Error retrieving Namespace: %v", err))
	}
	if ns.Access.Level < 3 {
		panic("Not enough Access Rights to Create Devices in choosen Namespace")
	}

	devc := pb.NewDevicesServiceClient(conn)
	// shad := pb.NewShadowServiceClient(conn)

	fmt.Println("gRPC Connection Established")

	rows := db.Point.Count(db.Point{})
	fmt.Printf("Amount of waypoints found: %d\n", rows)

	wg := sync.WaitGroup{}

	var pool []*common.Truck
	for i := 0; i < n_trucks; i++ {
		res, err := devc.Create(ctx, &devpb.CreateRequest{
			Device: &devpb.Device{
				Title:   fmt.Sprintf("sim-truck-%d", i),
				Enabled: true,
				Tags:    []string{"tn:simulated"},
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

		truck := &common.Truck{}
		truck.Uuid = res.GetDevice().GetUuid()
		truck.Point = rand.Int63n(rows)
		truck.Speed = time.Duration(rand.Intn(10)) * time.Second

		wg.Add(1)
		pool = append(pool, truck)
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
