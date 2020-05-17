package hardware

import (
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

// Provide provides os ram memory metrics at duration interval
func (p *RAMMetricProvider) Provide(r *event.InvokerRequest, h func(e *event.Reading)) error {
	defer r.WaitGroup.Done()
	ticker := time.NewTicker(r.Frequency)

	for {
		select {
		case <-r.Context.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			e, err := getRAMMetric(r.Source)
			if err != nil {
				return err
			}
			h(e)
		}
	}
}

func getRAMMetric(src string) (e *event.Reading, err error) {
	if src == "" {
		return nil, errors.New("nil source ID")
	}
	mp, err := mem.SwapMemory()
	if err != nil {
		return nil, errors.Wrap(err, "error creating swap memory metric")
	}
	e = &event.Reading{
		ID:    uuid.NewV4().String(),
		SrcID: src,
		Time:  time.Now().UTC().Unix(),
		Label: "swap memory",
		Unit:  "percent",
		Data:  mp.UsedPercent,
	}
	return
}
