package provider

import (
	"io/ioutil"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// LoadProviders loads user config
func LoadProviders(file string) ([]event.Provider, error) {
	rps, err := loadParamsFromConfig(file)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing file: %s", file)
	}

	ps := []event.Provider{}
	for _, rp := range rps {
		ps = append(ps, NewMetricProvider(&rp))
	}

	return ps, nil
}

func loadParamsFromConfig(file string) ([]event.ReadingParam, error) {
	if file == "" {
		return nil, errors.New("file argument required")
	}

	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading file: %s", file)
	}

	var c event.ParamConfig
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing content")
	}

	if c.Metrics == nil {
		return nil, errors.Wrapf(err, "invalid yaml format (nil metrics): %s", file)
	}

	return c.Metrics, nil
}
