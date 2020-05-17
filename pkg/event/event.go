package event

// Reading represents generic event
type Reading struct {
	ID    string      `json:"id"`
	SrcID string      `json:"src_id"`
	Time  int64       `json:"time"`
	Label string      `json:"label"`
	Data  interface{} `json:"data"`
	Unit  string      `json:"unit"`
}
