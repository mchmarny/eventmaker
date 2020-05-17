package main

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/amenzhinsky/iothub/iotdevice"
	"github.com/amenzhinsky/iothub/iotdevice/transport/mqtt"
	"github.com/pkg/errors"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/mchmarny/eventmaker/pkg/provider"
)

func newClient() (*iotdevice.Client, error) {
	client, err := iotdevice.NewFromConnectionString(mqtt.New(), clientConnStr)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing connection from env CONN_STR")
	}

	err = client.Connect(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "error connecting")
	}

	return client, nil
}

type job struct {
	clientID  string
	ctx       context.Context
	wg        *sync.WaitGroup
	client    *iotdevice.Client
	provider  provider.Provider
	frequency time.Duration
}

func send(j *job) {
	logger.Printf("sending metric: %+v", j.provider.Describe())
	err := j.provider.Provide(j.ctx, j.wg, j.clientID, j.frequency, func(e *event.SimpleEvent) {
		data, _ := json.Marshal(e)
		logger.Printf("%s", string(data))
		opts := []iotdevice.SendOption{
			iotdevice.WithSendMessageID(e.ID),
			iotdevice.WithSendQoS(1),
			iotdevice.WithSendCorrelationID(e.SrcID),
			iotdevice.WithSendProperty("uom", e.Unit),
			iotdevice.WithSendProperty("src", j.clientID),
		}
		if err := j.client.SendEvent(j.ctx, data, opts...); err != nil {
			if !errors.Is(err, context.Canceled) {
				logger.Printf("error on publish: '%+v' with %v", e, opts)
			}
		}
	})

	if err != nil {
		logger.Fatalf("error initializing provide: %v", err)
	}
}
