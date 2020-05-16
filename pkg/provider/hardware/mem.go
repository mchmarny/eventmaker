package hardware

import (
	"context"
	"sync"
	"time"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/mem"
)

// NewRAMMetricProvider creates nee RAMMetricProvider
func NewRAMMetricProvider() *RAMMetricProvider {
	return &RAMMetricProvider{}
}

// RAMMetricProvider is a host RAM metric provider
type RAMMetricProvider struct{}

// Describe provides provider info
func (p *RAMMetricProvider) Describe() *event.MetricInfo {
	return &event.MetricInfo{
		Metric: "swap memory",
		Type:   "float",
		Unit:   "percent",
	}
}

// Provide provides os ram memory metrics at duration interval
func (p *RAMMetricProvider) Provide(ctx context.Context, wg *sync.WaitGroup, src string, d time.Duration, h func(e *event.SimpleEvent)) error {
	defer wg.Done()
	ticker := time.NewTicker(d)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			e, err := getRAMMetric(src)
			if err != nil {
				return err
			}
			h(e)
		}
	}
}

func getRAMMetric(src string) (e *event.SimpleEvent, err error) {
	if src == "" {
		return nil, errors.New("nil source ID")
	}
	mp, err := mem.SwapMemory()
	if err != nil {
		return nil, errors.Wrap(err, "error creating swap memory metric")
	}
	e = &event.SimpleEvent{
		ID:    uuid.NewV4().String(),
		SrcID: src,
		Time:  time.Now().UTC().Unix(),
		Label: "swap memory",
		Unit:  "percent",
		Value: mp.UsedPercent,
	}
	return
}
