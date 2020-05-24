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

// Publisher defines the interface for event sender
type Publisher interface {
	Send(ctx context.Context, e *Reading) error
	Close() error
}
