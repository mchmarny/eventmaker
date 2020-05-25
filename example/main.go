package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/mchmarny/eventmaker/pkg/mock"
	"github.com/mchmarny/eventmaker/pkg/publish/stdout"
)

func main() {

	metrics := []event.MetricTemplate{
		{
			Label:     "temperature",
			Unit:      "celsius",
			Frequency: time.Duration(1 * time.Second),
			Template: event.ValueTemplate{
				Type: "float",
				Min:  39.1,
				Max:  73.5,
			},
		},
		{
			Label:     "humidity",
			Unit:      "percent",
			Frequency: time.Duration(1 * time.Minute),
			Template: event.ValueTemplate{
				Type: "int",
				Min:  0,
				Max:  100,
			},
		},
	}

	ctx := context.Background()
	target, err := stdout.NewEventSender(ctx)
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	mocker, err := mock.New("demo-device-1", metrics, target)
	if err != nil {
		panic(err)
	}

	cancel, errCh := mocker.Start(ctx)

	for {
		select {
		case e := <-errCh:
			cancel()
			panic(e)
		case <-sigChan:
			cancel()
			return
		}
	}

}
