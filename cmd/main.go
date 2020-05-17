package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/amenzhinsky/iothub/iotdevice"
	mqtt "github.com/amenzhinsky/iothub/iotdevice/transport/mqtt"
	"github.com/pkg/errors"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/mchmarny/eventmaker/pkg/provider"
	"github.com/mchmarny/eventmaker/pkg/provider/hardware"
	"github.com/mchmarny/gcputil/env"
)

var (
	logger = log.New(os.Stdout, "", 0)

	// Version will be overritten during build
	Version = "v0.0.1-default"

	ps = []provider.Provider{
		hardware.NewCPUMetricProvider(),
		hardware.NewLoadMetricProvider(),
		hardware.NewRAMMetricProvider(),
	}

	freq     = time.Duration(1 * time.Second)
	clientID = env.MustGetEnvVar("CLIENT_ID", "client-1")
)

func main() {
	logger.Printf("version: %s", Version)

	// client
	c, err := newClient()
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}
	defer c.Close()

	// signals
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// send
	for _, p := range ps {
		go func() {
			wg.Add(1)
			go send(ctx, wg, c, p)
		}()
	}

	// wait
	<-sigChan
	cancel()
	wg.Wait()
}

func send(ctx context.Context, wg *sync.WaitGroup, c *iotdevice.Client, p provider.Provider) {
	logger.Printf("metric: %+v", p.Describe())
	err := p.Provide(ctx, wg, clientID, freq, func(e *event.SimpleEvent) {
		data, _ := json.Marshal(e)
		logger.Printf("%s", string(data))
		opts := []iotdevice.SendOption{
			iotdevice.WithSendMessageID(e.ID),
			iotdevice.WithSendQoS(1),
			iotdevice.WithSendCorrelationID(e.SrcID),
			iotdevice.WithSendProperty("uom", e.Unit),
			iotdevice.WithSendProperty("src", clientID),
		}
		if err := c.SendEvent(ctx, data, opts...); err != nil {
			logger.Printf("error on publish: '%+v' with %v", e, opts)
		}
	})

	if err != nil {
		logger.Fatalf("error initializing provide: %v", err)
	}
}

func newClient() (*iotdevice.Client, error) {
	client, err := iotdevice.NewFromConnectionString(
		mqtt.New(), os.Getenv("CONN_STR"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing connection from env CONN_STR")
	}

	err = client.Connect(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "error connecting")
	}

	return client, nil
}
