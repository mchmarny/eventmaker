package provider

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/pkg/errors"
)

const (
	metricPrefix           = "--metric"
	metricPartsLength      = 5
	metricPartsDeliminator = "|"
	metricRangeDeliminator = ":"
	metricRangePartsLength = 2
)

var (
	logger = log.New(os.Stdout, "", 0)
)

// ParseProvider parses provider from commandline string
// Expected format '--metric|temp|celsius|float|0:72.1|3s'
func ParseProvider(arg string) (event.Provider, error) {
	if arg == "" {
		return nil, errors.New("empty metric arguments")
	}

	argParts := strings.Split(arg, metricPartsDeliminator)
	//logger.Printf("arg parts: %v", argParts)

	if len(argParts) != metricPartsLength {
		return nil, errors.New("invalid metric format (number of parts)")
	}

	// 0: temp
	// 1: celsius
	// 2: float
	// 3: 0:72.1
	// 4: 3s

	rp := &event.ReadingParam{
		Raw:   arg,
		Label: argParts[0],
		Unit:  argParts[1],
	}

	d, e := time.ParseDuration(argParts[4])
	if e != nil {
		return nil, errors.Wrapf(e, "error parsing frequency arg: %s", argParts[4])
	}
	rp.Frequency = d

	argType := strings.ToLower(argParts[2])
	argVal := argParts[3]

	switch argType {
	case "int", "int8", "int32", "int64":
		t, e := getIntArg(argVal)
		if e != nil {
			return nil, errors.Wrapf(e, "error parsing int arg: %s", argVal)
		}
		rp.Template = *t
	case "float", "float32", "float64":
		t, e := getFloatArg(argVal)
		if e != nil {
			return nil, errors.Wrapf(e, "error parsing float arg: %s", argVal)
		}
		rp.Template = *t
	case "bool":
		t, e := getBoolArg(argVal)
		if e != nil {
			return nil, errors.Wrapf(e, "error parsing float arg: %s", argVal)
		}
		rp.Template = *t
	default:
		return nil, errors.New("invalid data type in template")
	}

	//logger.Printf("parsed reading param: %+v", rp.Template)

	return NewMetricProvider(rp), nil
}

func getIntArg(s string) (*event.GenArg, error) {
	parts, err := validateRange(s)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing int range")
	}
	a := &event.GenArg{
		Type: "int",
	}

	min, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing int range, min: %s", parts[0])
	}
	a.Min = min

	max, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing int range, max: %s", parts[1])
	}
	a.Max = max

	return a, nil
}

func getFloatArg(s string) (*event.GenArg, error) {
	parts, err := validateRange(s)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing float range")
	}
	a := &event.GenArg{
		Type: "float",
	}

	min, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing float range, min: %s", parts[0])
	}
	a.Min = min

	max, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing float range, max: %s", parts[1])
	}
	a.Max = max

	return a, nil
}

func getBoolArg(s string) (*event.GenArg, error) {
	_, err := validateRange(s)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing bool range")
	}
	a := &event.GenArg{
		Type: "float",
		Min:  false,
		Max:  true,
	}

	return a, nil
}

func validateRange(s string) ([]string, error) {
	if s == "" {
		return nil, errors.New("empty range")
	}
	//logger.Printf("parsing range: %s", s)

	rangeParts := strings.Split(s, metricRangeDeliminator)
	//logger.Printf("range parts: %v", rangeParts)

	if len(rangeParts) != metricRangePartsLength {
		return nil, errors.New("invalid range format (want 2 parts)")
	}

	return rangeParts, nil

}
