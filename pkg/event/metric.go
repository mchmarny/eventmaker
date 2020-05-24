package event

import (
	"time"
)

// ConfigTemplate is the content of user config
type ConfigTemplate struct {
	Metrics []MetricTemplate `yaml:"metrics" json:"metrics"`
}

// MetricTemplate is readings parameter
type MetricTemplate struct {
	Label     string        `yaml:"label" json:"label"`
	Unit      string        `yaml:"unit" json:"unit"`
	Frequency time.Duration `yaml:"frequency" json:"frequency"`
	Template  ValueTemplate `yaml:"template" json:"template"`
}

// ValueTemplate defines the generation arguments
type ValueTemplate struct {
	Type string      `yaml:"type" json:"type"`
	Min  interface{} `yaml:"min" json:"min"`
	Max  interface{} `yaml:"max" json:"max"`
}

// MetricReading represents generic event
type MetricReading struct {
	ID    string      `json:"id"`
	SrcID string      `json:"src_id"`
	Time  int64       `json:"time"`
	Label string      `json:"label"`
	Data  interface{} `json:"data"`
	Unit  string      `json:"unit"`
}
