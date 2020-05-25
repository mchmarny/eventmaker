package provide

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// LoadProviders loads user config
func LoadProviders(file string) ([]MetricProvider, error) {
	templates, err := loadParamsFromConfig(file)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing metric template file: %s", file)
	}

	providers := make([]MetricProvider, 0)
	for _, t := range templates {
		providers = append(providers, NewMetricProvider(t))
	}

	return providers, nil
}

func loadParamsFromConfig(file string) ([]event.MetricTemplate, error) {
	if file == "" {
		return nil, errors.New("file argument required")
	}

	var content []byte

	// load only https files
	if strings.HasPrefix(file, "https://") {
		b, err := getContentFromURL(file)
		if err != nil {
			return nil, errors.Wrapf(err, "error reading url: %s", file)
		}
		content = b
	} else {
		f, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, errors.Wrapf(err, "error reading file: %s", file)
		}
		content = f
	}

	var c event.ConfigTemplate
	err := yaml.Unmarshal(content, &c)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing content")
	}

	if c.Metrics == nil {
		return nil, errors.Wrapf(err, "invalid yaml format (nil metrics): %s", file)
	}

	return c.Metrics, nil
}

func getContentFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
