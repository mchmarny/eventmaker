package provider

import (
	"testing"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/stretchr/testify/assert"
)

func validateConfig(t *testing.T, rp event.ReadingParam) {
	assert.NotNil(t, rp)
	assert.NotEmpty(t, rp.Label)
	assert.NotEmpty(t, rp.Unit)

	assert.NotNil(t, rp.Frequency)

	assert.NotNil(t, rp.Template)
	assert.NotEmpty(t, rp.Template.Type)
	assert.NotNil(t, rp.Template.Min)
	assert.NotNil(t, rp.Template.Max)
}

func TestLoadingConfigFromFile(t *testing.T) {
	rps, err := loadParamsFromConfig("../../conf/example.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, rps)
	assert.GreaterOrEqual(t, len(rps), 1)

	rp := rps[0]
	validateConfig(t, rp)
}

func TestLoadingConfigFromURL(t *testing.T) {
	u := "https://raw.githubusercontent.com/mchmarny/eventmaker/master/conf/example.yaml"
	rps, err := loadParamsFromConfig(u)
	assert.Nil(t, err)
	assert.NotNil(t, rps)
	assert.GreaterOrEqual(t, len(rps), 1)

	rp := rps[0]
	validateConfig(t, rp)
}
