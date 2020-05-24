package stdout

import (
	"context"
	"log"

	"github.com/mchmarny/eventmaker/pkg/event"
)

// NewEventSender creates nee MetricProvider
func NewEventSender(ctx context.Context, l *log.Logger) (*EventSender, error) {
	return &EventSender{
		logger: l,
	}, nil
}

// EventSender sends events
type EventSender struct {
	logger *log.Logger
}

// Close closes the client connection
func (s *EventSender) Close() error {
	return nil
}

// Publish sends provied events to stdout
func (s *EventSender) Publish(ctx context.Context, e *event.MetricReading) error {
	s.logger.Printf("%+v", e)
	return nil
}
