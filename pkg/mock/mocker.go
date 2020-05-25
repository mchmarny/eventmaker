package mock

import (
	"context"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/mchmarny/eventmaker/pkg/provide"
	"github.com/pkg/errors"
)

// Make is a factory methood for EventMocker
func Make(ctx context.Context, src, file string, pub event.Publisher) (*EventMocker, error) {
	if file == "" {
		return nil, errors.New("file required")
	}

	if src == "" {
		return nil, errors.New("src required")
	}

	m := &EventMocker{
		source:   src,
		publsher: pub,
	}

	// providers
	providers, err := provide.LoadProviders(file)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing providers from file (%s)", file)
	}
	m.providers = providers

	return m, nil
}

// EventMocker represents instane of event mocker
type EventMocker struct {
	source    string
	providers []provide.MetricProvider
	publsher  event.Publisher
}

// Start initiates mocking of events
func (m *EventMocker) Start(ctx context.Context) (context.CancelFunc, <-chan error) {
	ctx, cancel := context.WithCancel(ctx)
	errs := make(chan error, 1)
	for i := range m.providers {
		m.run(ctx, i, errs)
	} // for providers
	return cancel, errs
}

func (m *EventMocker) run(ctx context.Context, i int, er chan error) {
	go func() {
		p := m.providers[i]
		err := p.Provide(ctx, m.source, func(e *event.MetricReading) {
			if err := m.publsher.Publish(ctx, e); err != nil {
				if !errors.Is(err, context.Canceled) {
					er <- errors.Wrapf(err, "error publishing event [%+v]", e)
				}
			}
		})
		if err != nil {
			er <- errors.Wrapf(err, "error starting provider [%+v]", p.GetParam())
		}
	}()
}
