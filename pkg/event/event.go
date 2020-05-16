package event

// SimpleEvent represents generic event
type SimpleEvent struct {
	ID    string      `json:"id"`
	SrcID string      `json:"src_id"`
	Time  int64       `json:"time"`
	Label string      `json:"label"`
	Value interface{} `json:"value"`
	Unit  string      `json:"unit"`
}

// MetricInfo describes the type of event the provider can provide
type MetricInfo struct {
	Metric string `json:"metric"`
	Unit   string `json:"unit"`
	Type   string `json:"type"`
}
