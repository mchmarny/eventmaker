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

	file   string
	metric string
)

func main() {
	logger.Printf("version: %s", Version)

	flag.StringVar(&file, "file", "", "metric template file path")
	flag.StringVar(&metric, "metric", "", "individual template (e.g. temp|celsius|float|68.9:72.1|3s)")
	flag.Parse()

	logger.Printf("file: %v", file)
	logger.Printf("metric: %v", metric)

	if file == "" && metric == "" {
		log.Fatalln("either --file or --metric arguments required")
	}

	// providers
	ps := []event.Provider{}
	if metric != "" {
		pr, err := provider.ParseProvider(metric)
		if err != nil {
			log.Fatalf("error parsing provider from metric (%s): %v", metric, err)
		}
		ps = append(ps, pr)
	}
	if file != "" {
		prs, err := provider.LoadProviders(file)
		if err != nil {
			log.Fatalf("error parsing provider from file (%s): %v", file, err)
		}
		ps = append(ps, prs...)
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
