package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

func makeEvent(min, max float64) string {

	ms, err := mem.SwapMemory()
	failOnErr(err)

	d, err := time.ParseDuration("1s")
	failOnErr(err)
	cs, err := cpu.Percent(d, false)
	failOnErr(err)

	ls, err := load.Avg()
	failOnErr(err)

	event := struct {
		SourceID     string    `json:"source_id"`
		EventID      string    `json:"event_id"`
		EventTs      time.Time `json:"event_ts"`
		Label        string    `json:"label"`
		MemUsed      float64   `json:"mem_used"`
		CPUUsed      float64   `json:"cpu_used"`
		Load1        float64   `json:"load_1"`
		Load5        float64   `json:"load_5"`
		Load15       float64   `json:"load_15"`
		RandomMetric float64   `json:"random_metric"`
	}{
		SourceID:     *eventSrc,
		EventID:      fmt.Sprintf("%s-%s", idPrefix, uuid.NewV4().String()),
		EventTs:      time.Now().UTC(),
		Label:        *metricLabel,
		RandomMetric: min + rand.Float64()*(max-min),
		MemUsed:      ms.UsedPercent,
		CPUUsed:      cs[0],
		Load1:        ls.Load1,
		Load5:        ls.Load5,
		Load15:       ls.Load15,
	}

	data, _ := json.Marshal(event)

	return string(data)

}

func mustParseRange(r string) (min, max float64) {

	rangeParts := strings.Split(r, "-")
	if len(rangeParts) != 2 {
		log.Fatal(errorInvalidMetricRange)
	}

	min, minErr := strconv.ParseFloat(rangeParts[0], 64)
	max, maxErr := strconv.ParseFloat(rangeParts[1], 64)
	if minErr != nil || maxErr != nil {
		log.Fatal(errorInvalidMetricRange)
	}

	return min, max

}
