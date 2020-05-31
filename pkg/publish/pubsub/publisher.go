package pubsub

import (
	"context"
	"encoding/json"
	"time"

	gcp "cloud.google.com/go/pubsub"
	"github.com/pkg/errors"

	"github.com/mchmarny/eventmaker/pkg/event"
)

// NewEventSender creates new EventSender
func NewEventSender(ctx context.Context, projectID, topicName string) (*EventSender, error) {
	if projectID == "" {
		return nil, errors.New("projectID not defined")
	}
	if topicName == "" {
		return nil, errors.New("topicName not defined")
	}
	c, err := gcp.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating client: project='%s', topic='%s'",
			projectID, topicName)
	}

	t := c.Topic(topicName)
	topicExists, err := t.Exists(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error connecting")
	}

	if !topicExists {
		t, err = c.CreateTopic(ctx, topicName)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating topic '%s'", topicName)
		}
	}

	return &EventSender{
		client: c,
		topic:  t,
	}, nil
}

// EventSender sends events
type EventSender struct {
	client *gcp.Client
	topic  *gcp.Topic
}

// Close closes the client connection
func (s *EventSender) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

// Publish sends provied events to GCP PubSub
func (s *EventSender) Publish(ctx context.Context, e *event.MetricReading) error {
	data, _ := json.Marshal(e)

	msg := &gcp.Message{
		ID:          e.ID,
		Data:        data,
		PublishTime: time.Now(),
		Attributes: map[string]string{
			"uom": e.Unit,
			"src": e.SrcID,
		},
	}

	result := s.topic.Publish(ctx, msg)
	_, err := result.Get(ctx)
	if !errors.Is(err, context.Canceled) {
		return errors.Wrapf(err, "error on publish: '%+v'", e)
	}

	return nil
}
