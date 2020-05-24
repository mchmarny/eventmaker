package event

import "context"

const (
	// StdoutPublisher is the default publsher to stdout
	StdoutPublisher = "stdout"
	// AzureIoTHubPublsher is the IoT Hub Publisher
	AzureIoTHubPublsher = "iothub"
)

// Publisher defines the interface for event sender
type Publisher interface {
	Send(ctx context.Context, e *Reading) error
	Close() error
}
