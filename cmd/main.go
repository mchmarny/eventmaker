package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/mchmarny/eventmaker/pkg/provider"
	"github.com/mchmarny/gcputil/env"
)

var (
	logger = log.New(os.Stdout, "", 0)

	// Version will be overritten during build
	Version = "v0.0.1-default"

	deviceID      = env.MustGetEnvVar("DEV_NAME", "device-1")
	clientConnStr = env.MustGetEnvVar("CONN_STR", "")

	file string
)

func main() {
	logger.Printf("version: %s", Version)

	flag.StringVar(&file, "file", "", "metric template file path")
	flag.Parse()

	if file == "" {
		log.Fatalln("--file flag required")
	}

	// providers
	prs, err := provider.LoadProviders(file)
	if err != nil {
		log.Fatalf("error parsing provider from file (%s): %v", file, err)
	}
	logger.Printf("loaded %d metric proviers", len(prs))

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
	for _, p := range prs {
		wg.Add(1)
		j := &job{
			client:   c,
			provider: p,
			request: &event.InvokerRequest{
				Source:    deviceID,
				Context:   ctx,
				WaitGroup: wg,
				Frequency: p.GetParam().Frequency,
			},
		}
		go send(j)
	}

	// wait
	<-sigChan
	cancel()
	logger.Println("\nwaiting for providers to drain existing work...")
	wg.Wait()
}
