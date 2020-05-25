package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/pkg/errors"
)

// NewEventSender creates nee MetricProvider
func NewEventSender(ctx context.Context, url string) (*EventSender, error) {
	if url == "" {
		return nil, errors.New("url not defined")
	}
	return &EventSender{
		url: url,
	}, nil
}

// EventSender sends events
type EventSender struct {
	url string
}

// Close closes the client connection
func (s *EventSender) Close() error {
	return nil
}

// Publish sends provied events to stdout
func (s *EventSender) Publish(ctx context.Context, e *event.MetricReading) error {
	data, _ := json.Marshal(e)

	req, err := http.NewRequest(http.MethodPost, s.url, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrapf(err, "error creating post request: %s", s.url)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "error executing post request: %s", s.url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(err, "invalid response code to %s: %d (%s)",
			s.url, resp.StatusCode, resp.Status)
	}

	return nil
}
