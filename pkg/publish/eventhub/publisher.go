package eventhub

import (
	"context"
	"encoding/json"

	hub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/pkg/errors"

	"github.com/mchmarny/eventmaker/pkg/event"
)

// NewEventSender creates new EventSender
func NewEventSender(ctx context.Context, connStr string) (*EventSender, error) {
	if connStr == "" {
		return nil, errors.New("connStr not defined")
	}
	c, err := hub.NewHubFromConnectionString(connStr)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating client from: '%s'", connStr)
	}

	_, err = c.GetRuntimeInformation(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error connecting")
	}

	return &EventSender{
		client: c,
	}, nil
}

// EventSender sends events
type EventSender struct {
	client *hub.Hub
}

// Close closes the client connection
func (s *EventSender) Close() error {
	if s.client != nil {
		return s.client.Close(context.Background())
	}
	return nil
}

// Publish sends provied events to Event Hub
func (s *EventSender) Publish(ctx context.Context, e *event.MetricReading) error {
	data, _ := json.Marshal(e)

	ev := hub.NewEvent(data)
	ev.ID = e.ID
	ev.Properties = make(map[string]interface{})
	ev.Properties["uom"] = e.Unit
	ev.Properties["src"] = e.SrcID

	if err := s.client.Send(ctx, ev); err != nil {
		if !errors.Is(err, context.Canceled) {
			return errors.Wrapf(err, "error on publish: '%+v'", e)
		}
	}

	return nil
}
