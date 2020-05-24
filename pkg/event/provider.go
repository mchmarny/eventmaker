package event

import (
	"context"
	"sync"
	"time"
)

// ProviderRequest is the context of the provider invoker
type ProviderRequest struct {
	Source    string
	Context   context.Context
	WaitGroup *sync.WaitGroup
	Frequency time.Duration
}

// Provider defines the interface for data provider
type Provider interface {
	// Provide generates events and invokes handler
	Provide(r *ProviderRequest, h func(e *MetricReading)) error
	GetParam() *MetricTemplate
}
