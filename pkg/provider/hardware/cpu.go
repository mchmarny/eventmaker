package hardware

import (
	"time"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/cpu"
)

// NewCPUMetricProvider creates nee CPUMetricProvider
func NewCPUMetricProvider() *CPUMetricProvider {
	return &CPUMetricProvider{}
}

// CPUMetricProvider is a host CPU metric provider
type CPUMetricProvider struct{}

// Provide provides os process events
func (p *CPUMetricProvider) Provide(r *event.InvokerRequest, h func(e *event.Reading)) error {
	defer r.WaitGroup.Done()
	ticker := time.NewTicker(r.Frequency)

	for {
		select {
		case <-r.Context.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			e, err := getCPUMetric(r.Source, r.Frequency)
			if err != nil {
				return err
			}
			h(e)
		}
	}
}

func getCPUMetric(src string, d time.Duration) (e *event.Reading, err error) {
	if src == "" {
		return nil, errors.New("nil source ID")
	}
	mp, err := cpu.Percent(d, false)
	if err != nil {
		return nil, errors.Wrap(err, "error creating cpu utilization metric")
	}
	e = &event.Reading{
		ID:    uuid.NewV4().String(),
		SrcID: src,
		Time:  time.Now().UTC().Unix(),
		Label: "cpu utilization",
		Unit:  "percent",
		Data:  mp[0],
	}
	return
}
