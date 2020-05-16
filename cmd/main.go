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
	hw "github.com/mchmarny/eventmaker/pkg/provider/hardware"
)

var (
	ps = []provider.Provider{
		hw.NewCPUMetricProvider(),
		hw.NewLoadMetricProvider(),
		hw.NewRAMMetricProvider(),
	}

	freq = time.Duration(1 * time.Second)
)

func main() {
	// load providers

	// client
	c, err := newClient()
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}

	// setup signals
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// send
	for _, p := range ps {
		log.Printf("metric: %+v", p.Describe())
		go send(ctx, wg, c, p)
	}

	// wait
	<-sigChan
	cancel()
	wg.Wait()
}

func send(ctx context.Context, wg *sync.WaitGroup, c *iotdevice.Client, p provider.Provider) {
	err := p.Provide(ctx, wg, "client-1", freq, func(e *event.SimpleEvent) {
		data, _ := json.Marshal(e)
		log.Printf("%s", string(data))
		opts := []iotdevice.SendOption{
			iotdevice.WithSendMessageID(e.ID),
		}
		if err := c.SendEvent(ctx, data, opts...); err != nil {
			log.Printf("error on publish: %v ", err)
			return
		}
	})

	if err != nil {
		log.Printf("error on provide: %v ", err)
		return
	}
}

func newClient() (*iotdevice.Client, error) {
	client, err := iotdevice.NewFromConnectionString(
		mqtt.New(), os.Getenv("DEV1_CONN"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing connection from env IOTHUB_DEVICE_CONNECTION_STRING")
	}

	err = client.Connect(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "error connecting")
	}

	return client, nil
}
