package event

import (
	"context"
	"sync"
	"time"
)

// Reading represents generic event
type Reading struct {
	ID    string      `json:"id"`
	SrcID string      `json:"src_id"`
	Time  int64       `json:"time"`
	Label string      `json:"label"`
	Data  interface{} `json:"data"`
	Unit  string      `json:"unit"`
}

// InvokerRequest is the context of the provider invoker
type InvokerRequest struct {
	Source    string
	Context   context.Context
	WaitGroup *sync.WaitGroup
	Frequency time.Duration
}

// Provider defines the interface for data provider
type Provider interface {
	// Provide generates events and invokes handler
	Provide(r *InvokerRequest, h func(e *Reading)) error
}
