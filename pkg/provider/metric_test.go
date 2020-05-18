package provider

import (
	"testing"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/stretchr/testify/assert"
)

func TestRandomIntGen(t *testing.T) {
	arg := event.GenArg{Type: "int", Min: int64(0), Max: int64(100)}

	val, err := getRandomValue(arg)
	assert.Nil(t, err)
	assert.NotNil(t, val)

	valInt64 := val.(int64)
	assert.LessOrEqual(t, valInt64, int64(100))
	assert.GreaterOrEqual(t, valInt64, int64(0))
}

func TestRandomFloatGen(t *testing.T) {
	arg := event.GenArg{Type: "float", Min: float64(0), Max: float64(100)}

	val, err := getRandomValue(arg)
	assert.Nil(t, err)
	assert.NotNil(t, val)

	valFloat64 := val.(float64)
	assert.LessOrEqual(t, valFloat64, float64(100))
	assert.GreaterOrEqual(t, valFloat64, float64(0))
}

func TestRandomBoolGen(t *testing.T) {
	arg := event.GenArg{Type: "bool", Min: 0, Max: 1}

	val, err := getRandomValue(arg)
	assert.Nil(t, err)
	assert.NotNil(t, val)
}
