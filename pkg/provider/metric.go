package provider

import (
	"math/rand"
	"time"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// NewMetricProvider creates nee MetricProvider
func NewMetricProvider(param *event.ReadingParam) *MetricProvider {
	return &MetricProvider{
		param: param,
	}
}

// MetricProvider generates metric readers based on dynamic value
type MetricProvider struct {
	param *event.ReadingParam
}

// GetParam returns local param
func (p *MetricProvider) GetParam() *event.ReadingParam {
	return p.param
}

// Provide provides os process events
func (p *MetricProvider) Provide(r *event.InvokerRequest, h func(e *event.Reading)) error {
	defer r.WaitGroup.Done()
	ticker := time.NewTicker(r.Frequency)

	for {
		select {
		case <-r.Context.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			e, err := makeMetric(r.Source, p.param)
			if err != nil {
				return err
			}
			h(e)
		}
	}
}

func makeMetric(src string, rp *event.ReadingParam) (e *event.Reading, err error) {
	if rp == nil {
		return nil, errors.New("nil reading param")
	}

	v, ge := getRandomValue(rp.Template)
	if ge != nil {
		return nil, errors.Wrap(ge, "error generating rundom value")
	}

	e = &event.Reading{
		ID:    uuid.NewV4().String(),
		SrcID: src,
		Time:  time.Now().UTC().Unix(),
		Label: rp.Label,
		Unit:  rp.Unit,
		Data:  v,
	}
	return
}

func getRandomValue(arg event.GenArg) (val interface{}, err error) {
	switch arg.Type {
	case "int", "int8", "int32", "int64":
		return getRandomIntValue(arg.Min.(int64), arg.Max.(int64)), nil
	case "float", "float32", "float64":
		return getRandomFloatValue(arg.Min.(float64), arg.Max.(float64)), nil
	case "bool":
		return getRandomBoolValue(), nil
	default:
		return nil, errors.New("invalid data type in template")
	}
}

func getRandomIntValue(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(int(max)-int(min)) + int(min))
}

func getRandomFloatValue(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

func getRandomBoolValue() bool {
	return getRandomIntValue(0, 100) < 50
}
