package hardware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadMetricProvider(t *testing.T) {
	e, err := getLoadMetric("test")
	assert.Nil(t, err)
	runEventTest(t, e)
}
