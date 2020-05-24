package iothub

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/amenzhinsky/iothub/iotdevice"
	"github.com/amenzhinsky/iothub/iotdevice/transport/mqtt"
	"github.com/pkg/errors"

	"github.com/mchmarny/eventmaker/pkg/event"
)

// NewEventSender creates nee MetricProvider
func NewEventSender(ctx context.Context, l *log.Logger) (*EventSender, error) {
	connStr := os.Getenv("CONN_STR")
	if connStr == "" {
		return nil, errors.New("CONN_STR not defined")
	}
	c, err := iotdevice.NewFromConnectionString(mqtt.New(), connStr)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating client from: '%s'", connStr)
	}

	err = c.Connect(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error connecting")
	}
	l.Println("IoT Hub Client connected")

	return &EventSender{
		client: c,
		logger: l,
	}, nil
}

// EventSender sends events
type EventSender struct {
	client *iotdevice.Client
	logger *log.Logger
}

// Close closes the client connection
func (s *EventSender) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

// Publish sends provied events to IoT Hub
func (s *EventSender) Publish(ctx context.Context, e *event.MetricReading) error {
	data, _ := json.Marshal(e)
	opts := []iotdevice.SendOption{
		iotdevice.WithSendMessageID(e.ID),
		iotdevice.WithSendQoS(1),
		iotdevice.WithSendCorrelationID(e.SrcID),
		iotdevice.WithSendProperty("uom", e.Unit),
		iotdevice.WithSendProperty("src", e.SrcID),
	}

	if err := s.client.SendEvent(ctx, data, opts...); err != nil {
		if !errors.Is(err, context.Canceled) {
			errors.Wrapf(err, "error on publish: '%+v' with %v", e, opts)
		}
	}

	return nil
}
