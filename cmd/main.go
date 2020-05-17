package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/mchmarny/eventmaker/pkg/provider"
	"github.com/mchmarny/eventmaker/pkg/provider/hardware"
	"github.com/mchmarny/gcputil/env"
)

var (
	logger = log.New(os.Stdout, "", 0)

	// Version will be overritten during build
	Version = "v0.0.1-default"

	freqStr       = env.MustGetEnvVar("SEND_FREQ", "1s")
	clientID      = env.MustGetEnvVar("CLIENT_ID", "client-1")
	clientConnStr = env.MustGetEnvVar("CONN_STR", "")
)

func main() {
	logger.Printf("version: %s", Version)

	freq, err := time.ParseDuration(freqStr)
	if err != nil {
		log.Fatalf("invalid send frequency value (SEND_FREQ): %s", freqStr)
	}

	// providers
	ps := []provider.Provider{
		hardware.NewCPUMetricProvider(),
		hardware.NewLoadMetricProvider(),
		hardware.NewRAMMetricProvider(),
	}
	logger.Printf("loaded %d metric proviers", len(ps))

	// client
	c, err := newClient()
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}
	defer c.Close()

	// signals
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// send
	for _, p := range ps {
		logger.Printf("metric %v", p.Describe())
		wg.Add(1)
		j := &job{
			clientID:  clientID,
			ctx:       ctx,
			wg:        wg,
			client:    c,
			provider:  p,
			frequency: freq,
		}
		go send(j)
	}

	// wait
	<-sigChan
	cancel()
	logger.Println("\nwaiting for providers to drain existing work...")
	wg.Wait()
}
