package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseProviders(t *testing.T) {
	args := []string{
		"--metric",
		"temp|celsius|float|0:72.1|3s",
		"--metric",
		"speed|kmh|int|0:210|1s",
		"--metric",
		"friction|coefficient|float|0:1|1s",
	}

	list, err := ParseProviders(args)
	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Len(t, list, 3)

}

func TestParseFloatProvider(t *testing.T) {
	args := []string{
		"--metric",
		"temp|celsius|float|0:72.1|3s",
	}

	list, err := ParseProviders(args)
	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Len(t, list, 1)

	rp := list[0].GetParam()
	assert.NotNil(t, rp)
	assert.NotEmpty(t, rp.Label)
	assert.NotEmpty(t, rp.Unit)
	assert.NotEmpty(t, rp.Raw)

	assert.NotNil(t, rp.Frequency)

	assert.NotNil(t, rp.Template)
	assert.NotEmpty(t, rp.Template.Type)
	assert.NotNil(t, rp.Template.Min)
	assert.NotNil(t, rp.Template.Max)

}
