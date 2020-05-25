package event

import "context"

// ListPublishers list supported publishers
func ListPublishers() []string {
	return []string{
		"stdout",
		"iothub",
		"http",
	}
}

// Publisher defines the interface for event publisher
type Publisher interface {
	Publish(ctx context.Context, e *MetricReading) error
	Close() error
}
