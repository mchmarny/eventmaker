package hardware

import (
	"time"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/load"
)

// NewLoadMetricProvider creates nee LoadMetricProvider
func NewLoadMetricProvider() *LoadMetricProvider {
	return &LoadMetricProvider{}
}

// LoadMetricProvider is a host system load metric provider
type LoadMetricProvider struct{}

// Provide provides os ram memory metrics at duration interval
func (p *LoadMetricProvider) Provide(r *event.InvokerRequest, h func(e *event.Reading)) error {
	defer r.WaitGroup.Done()
	ticker := time.NewTicker(r.Frequency)

	for {
		select {
		case <-r.Context.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			e, err := getLoadMetric(r.Source)
			if err != nil {
				return err
			}
			h(e)
		}
	}
}

func getLoadMetric(src string) (e *event.Reading, err error) {
	if src == "" {
		return nil, errors.New("nil source ID")
	}
	mp, err := load.Avg()
	if err != nil {
		return nil, errors.Wrap(err, "error creating avg system load metric")
	}
	e = &event.Reading{
		ID:    uuid.NewV4().String(),
		SrcID: src,
		Time:  time.Now().UTC().Unix(),
		Label: "system load",
		Unit:  "percent",
		Data:  mp.Load1,
	}
	return
}
