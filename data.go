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

var (
	memStats  *mem.SwapMemoryStat
	cpuStats  cpu.InfoStat
	loadStats *load.AvgStat
)

func init() {

	ms, err := mem.SwapMemory()
	failOnErr(err)
	memStats = ms

	cs, err := cpu.Info()
	failOnErr(err)
	cpuStats = cs[0]

	ls, err := load.Avg()
	failOnErr(err)
	loadStats = ls
}

func makeEvent(min, max float64) string {

	event := struct {
		SourceID    string    `json:"source_id"`
		EventID     string    `json:"event_id"`
		EventTs     time.Time `json:"event_ts"`
		Label       string    `json:"label"`
		MemFree     float64   `json:"memFree"`
		CPUFree     float64   `json:"cpuFree"`
		LoadAvg1    float64   `json:"loadAvg1"`
		LoadAvg5    float64   `json:"loadAvg5"`
		LoadAvg15   float64   `json:"loadAvg15"`
		RandomValue float64   `json:"randomValue"`
	}{
		SourceID:    *eventSrc,
		EventID:     fmt.Sprintf("%s-%s", idPrefix, uuid.NewV4().String()),
		EventTs:     time.Now().UTC(),
		Label:       *metricLabel,
		RandomValue: min + rand.Float64()*(max-min),
		MemFree:     memStats.UsedPercent,
		CPUFree:     cpuStats.Mhz,
		LoadAvg1:    loadStats.Load1,
		LoadAvg5:    loadStats.Load5,
		LoadAvg15:   loadStats.Load15,
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
