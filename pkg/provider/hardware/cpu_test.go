package hardware

import (
	"testing"
	"time"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/stretchr/testify/assert"
)

func runEventTest(t *testing.T, e *event.SimpleEvent) {
	assert.NotNil(t, e)
	assert.Equal(t, "test", e.SrcID)
	assert.NotEmpty(t, e.ID)
	assert.NotEmpty(t, e.Label)
	assert.NotEmpty(t, e.Unit)
	assert.NotNil(t, e.Value)
	assert.NotZero(t, e.Time)
}

func TestCPUMetricProvider(t *testing.T) {
	e, err := getCPUMetric("test", time.Duration(1*time.Second))
	assert.Nil(t, err)
	runEventTest(t, e)
}
