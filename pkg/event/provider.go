package event

import (
	"context"
	"sync"
	"time"
)

// ReadingParam is readings parameter
type ReadingParam struct {
	Raw       string
	Label     string
	Unit      string
	Frequency time.Duration
	Template  *GenArg
}

// GenArg defines the generation arguments
type GenArg struct {
	Type string
	Min  interface{}
	Max  interface{}
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
	GetParam() *ReadingParam
}
