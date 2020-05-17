package hardware

import (
	"context"
	"sync"
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

// Describe provides provider info
func (p *LoadMetricProvider) Describe() *event.MetricInfo {
	return &event.MetricInfo{
		Metric: "system load",
		Type:   "float",
		Unit:   "percent",
	}
}

// Provide provides os ram memory metrics at duration interval
func (p *LoadMetricProvider) Provide(ctx context.Context, wg *sync.WaitGroup, src string, d time.Duration, h func(e *event.SimpleEvent)) error {
	defer wg.Done()
	ticker := time.NewTicker(d)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			e, err := getLoadMetric(src)
			if err != nil {
				return err
			}
			h(e)
		}
	}
}

func getLoadMetric(src string) (e *event.SimpleEvent, err error) {
	if src == "" {
		return nil, errors.New("nil source ID")
	}
	mp, err := load.Avg()
	if err != nil {
		return nil, errors.Wrap(err, "error creating avg system load metric")
	}
	e = &event.SimpleEvent{
		ID:    uuid.NewV4().String(),
		SrcID: src,
		Time:  time.Now().UTC().Unix(),
		Label: "system load",
		Unit:  "percent",
		Value: mp.Load1,
	}
	return
}
