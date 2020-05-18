package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadingConfig(t *testing.T) {
	rps, err := loadParamsFromConfig("../../config/example.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, rps)
	assert.GreaterOrEqual(t, len(rps), 1)

	rp := rps[0]
	assert.NotNil(t, rp)
	assert.NotEmpty(t, rp.Label)
	assert.NotEmpty(t, rp.Unit)

	assert.NotNil(t, rp.Frequency)

	assert.NotNil(t, rp.Template)
	assert.NotEmpty(t, rp.Template.Type)
	assert.NotNil(t, rp.Template.Min)
	assert.NotNil(t, rp.Template.Max)
}
