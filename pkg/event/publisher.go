package event

import "context"

// ListPublishers list supported publishers
func ListPublishers() []string {
	return []string{
		"eventhub",
		"http",
		"iothub",
		"pubsub",
		"stdout",
	}
}

// Publisher defines the interface for event publisher
type Publisher interface {
	Publish(ctx context.Context, e *MetricReading) error
	Close() error
}
