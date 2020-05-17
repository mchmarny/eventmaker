package main

import (
	"context"
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

	args = os.Args[1:]
)

func main() {
	logger.Printf("version: %s", Version)
	logger.Printf("args: %v", args)

	// providers
	ps, err := provider.ParseProviders(args)
	if err != nil {
		log.Fatalf("error parsing providers: %v", err)
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
