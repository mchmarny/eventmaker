package main

import (
	"context"
	"encoding/json"

	"github.com/amenzhinsky/iothub/iotdevice"
	"github.com/amenzhinsky/iothub/iotdevice/transport/mqtt"
	"github.com/pkg/errors"

	"github.com/mchmarny/eventmaker/pkg/event"
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
	client   *iotdevice.Client
	provider event.Provider
	request  *event.InvokerRequest
}

func send(j *job) {
	err := j.provider.Provide(j.request, func(e *event.Reading) {
		data, _ := json.Marshal(e)
		logger.Printf("%s", string(data))
		opts := []iotdevice.SendOption{
			iotdevice.WithSendMessageID(e.ID),
			iotdevice.WithSendQoS(1),
			iotdevice.WithSendCorrelationID(e.SrcID),
			iotdevice.WithSendProperty("uom", e.Unit),
			iotdevice.WithSendProperty("src", j.request.Source),
		}
		if err := j.client.SendEvent(j.request.Context, data, opts...); err != nil {
			if !errors.Is(err, context.Canceled) {
				logger.Printf("error on publish: '%+v' with %v", e, opts)
			}
		}
	})

	if err != nil {
		logger.Fatalf("error initializing provide: %v", err)
	}
}
