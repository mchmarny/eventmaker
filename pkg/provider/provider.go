package provider

import (
	"context"
	"sync"
	"time"

	"github.com/mchmarny/eventmaker/pkg/event"
)

// Provider defines the interface for data provider
type Provider interface {
	// Provide generates events
	Provide(ctx context.Context, wg *sync.WaitGroup, src string, d time.Duration, h func(e *event.SimpleEvent)) error
	// Describe provides provider info
	Describe() *event.MetricInfo
}
