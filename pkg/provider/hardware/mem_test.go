package hardware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRAMMetricProvider(t *testing.T) {
	e, err := getRAMMetric("test")
	assert.Nil(t, err)
	runEventTest(t, e)
}
