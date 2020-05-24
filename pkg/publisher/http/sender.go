package http

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/pkg/errors"
)

// NewEventSender creates nee MetricProvider
func NewEventSender(ctx context.Context, l *log.Logger) (*EventSender, error) {
	t := os.Getenv("TARGET_URL")
	if t == "" {
		return nil, errors.New("TARGET_URL not defined")
	}
	return &EventSender{
		logger: l,
		url:    t,
	}, nil
}

// EventSender sends events
type EventSender struct {
	logger *log.Logger
	url    string
}

// Close closes the client connection
func (s *EventSender) Close() error {
	return nil
}

// Send sends provied events to stdout
func (s *EventSender) Send(ctx context.Context, e *event.Reading) error {
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
		return errors.Wrapf(err, "invalid response code to %s: %d", s.url, resp.Status)
	}

	return nil
}
