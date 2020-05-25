package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/mchmarny/eventmaker/pkg/mock"
	"github.com/mchmarny/gcputil/env"
)

var (
	logger = log.New(os.Stdout, "", 0)

	// Version will be overritten during build
	Version = "v0.0.1-default"

	deviceID = env.MustGetEnvVar("DEV_NAME", "eventmakerdev-0")

	file    string
	pubType string
)

func main() {
	logger.Printf("version: %s", Version)

	flag.StringVar(&file, "file", "", "metric template file path")
	flag.StringVar(&pubType, "publisher", "stdout", "event publisher (stdout, iothub, http)")
	flag.Parse()

	if file == "" {
		log.Printf("args: %v", os.Args)
		log.Fatalln("--file flag required")
	}

	// signals
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)

	ctx := context.Background()
	em, err := mock.Make(ctx, deviceID, file, pubType)
	if err != nil {
		logger.Fatalf("error creating mocker: %v", err)
	}

	cancel, errCh := em.Start(ctx)

	for {
		select {
		case e := <-errCh:
			logger.Println(e)
			return
		case <-sigChan:
			return
		}
	}

	cancel()
}
