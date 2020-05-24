package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/mchmarny/eventmaker/pkg/provider"
	"github.com/mchmarny/eventmaker/pkg/publisher/iothub"
	"github.com/mchmarny/eventmaker/pkg/publisher/stdout"
	"github.com/mchmarny/gcputil/env"
)

var (
	logger = log.New(os.Stdout, "", 0)

	// Version will be overritten during build
	Version = "v0.0.1-default"

	deviceID = env.MustGetEnvVar("DEV_NAME", "device-1")

	file    string
	pubType string
)

func main() {
	logger.Printf("version: %s", Version)

	flag.StringVar(&file, "file", "", "metric template file path")
	flag.StringVar(&pubType, "publisher", "stdout", "event publisher (stdout, iothub)")
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

	// signals
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// publisher
	var publisher event.Publisher

	switch pubType {
	case event.StdoutPublisher:
		stdoutPublisher, err := stdout.NewEventSender(ctx)
		if err != nil {
			log.Fatalf("error creating stdout publisher: %v", err)
		}
		defer stdoutPublisher.Close()
		publisher = stdoutPublisher
	case event.AzureIoTHubPublsher:
		iotHubPublisher, err := iothub.NewEventSender(ctx)
		if err != nil {
			log.Fatalf("error creating iot hub publisher: %v", err)
		}
		defer iotHubPublisher.Close()
		publisher = iotHubPublisher
	default:
		log.Fatalf("invalid publisher type (%s)", pubType)
	}

	// process
	for _, p := range prs {
		wg.Add(1)
		r := event.ProviderRequest{
			Source:    deviceID,
			Context:   ctx,
			WaitGroup: wg,
			Frequency: p.GetParam().Frequency,
		}
		go run(ctx, p, r, publisher)
	}

	// wait
	<-sigChan
	cancel()
	logger.Println("\nwaiting for providers to drain existing work...")
	wg.Wait()
}

// run executes provider with invoker request
func run(ctx context.Context, p event.Provider, r event.ProviderRequest, s event.Publisher) {
	err := p.Provide(&r, func(e *event.Reading) {
		if err := s.Send(ctx, r.Source, e); err != nil {
			if !errors.Is(err, context.Canceled) {
				logger.Printf("error sending: '%+v'", e)
			}
		}
	})

	if err != nil {
		logger.Fatalf("error initializing provide: %v", err)
	}
}
