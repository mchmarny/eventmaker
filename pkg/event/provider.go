package event

import (
	"context"
	"sync"
	"time"
)

// ParamConfig is the content of user config
type ParamConfig struct {
	Metrics []ReadingParam `yaml:"metrics"`
}

// ReadingParam is readings parameter
type ReadingParam struct {
	Raw       string        `yaml:"-" json:"-"`
	Label     string        `yaml:"label" json:"label"`
	Unit      string        `yaml:"unit" json:"unit"`
	Frequency time.Duration `yaml:"frequency" json:"frequency"`
	Template  GenArg        `yaml:"template" json:"template"`
}

// GenArg defines the generation arguments
type GenArg struct {
	Type string      `yaml:"type" json:"type"`
	Min  interface{} `yaml:"min" json:"min"`
	Max  interface{} `yaml:"max" json:"max"`
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
