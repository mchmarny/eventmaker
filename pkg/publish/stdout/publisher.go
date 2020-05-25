package stdout

import (
	"context"
	"fmt"

	"github.com/mchmarny/eventmaker/pkg/event"
)

// NewEventSender creates nee MetricProvider
func NewEventSender(ctx context.Context) (*EventSender, error) {
	return &EventSender{}, nil
}

// EventSender sends events
type EventSender struct{}

// Close closes the client connection
func (s *EventSender) Close() error {
	return nil
}

// Publish sends provied events to stdout
func (s *EventSender) Publish(ctx context.Context, e *event.MetricReading) error {
	fmt.Printf("%+v\n", e)
	return nil
}
