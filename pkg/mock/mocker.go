package mock

import (
	"context"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/mchmarny/eventmaker/pkg/provide"
	"github.com/mchmarny/eventmaker/pkg/publish/http"
	"github.com/mchmarny/eventmaker/pkg/publish/iothub"
	"github.com/mchmarny/eventmaker/pkg/publish/stdout"
	"github.com/pkg/errors"
)

// Make is a factory methood for EventMocker
func Make(ctx context.Context, src, file, pub string) (*EventMocker, error) {
	if file == "" {
		return nil, errors.New("nil file")
	}

	if pub == "" {
		return nil, errors.New("nil pub")
	}

	m := &EventMocker{
		source: src,
	}

	// providers
	providers, err := provide.LoadProviders(file)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing providers from file (%s)", file)
	}
	m.providers = providers

	// publisher
	switch pub {
	case event.StdoutPublisher:
		stdoutPublisher, err := stdout.NewEventSender(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating stdout publisher")
		}
		m.publsher = stdoutPublisher
	case event.AzureIoTHubPublsher:
		iotHubPublisher, err := iothub.NewEventSender(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating iot hub publisher")
		}
		m.publsher = iotHubPublisher
	case event.HTTPPublsher:
		httpPublisher, err := http.NewEventSender(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating http publisher")
		}
		m.publsher = httpPublisher
	default:
		return nil, errors.Wrapf(err, "invalid publisher type (%s)", pub)
	}

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
