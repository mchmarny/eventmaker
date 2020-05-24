package event

import "context"

const (
	// StdoutPublisher is the default publsher to stdout
	StdoutPublisher = "stdout"
	// AzureIoTHubPublsher is the IoT Hub publisher
	AzureIoTHubPublsher = "iothub"
	// HTTPPublsher is a basic HTTP publisher
	HTTPPublsher = "http"
)

// Publisher defines the interface for event publisher
type Publisher interface {
	Publish(ctx context.Context, e *MetricReading) error
	Close() error
}
