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

// Send sends provied events to stdout
func (s *EventSender) Send(ctx context.Context, src string, e *event.Reading) error {
	fmt.Printf("%s[%+v]\n", src, e)
	return nil
}
