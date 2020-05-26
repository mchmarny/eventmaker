package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/mchmarny/eventmaker/pkg/mock"
	"github.com/mchmarny/eventmaker/pkg/publish/http"
	"github.com/mchmarny/eventmaker/pkg/publish/iothub"
	"github.com/mchmarny/eventmaker/pkg/publish/stdout"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	logger = log.New(os.Stdout, "", 0)

	// Version will be overritten during build
	Version = "v0.0.1-default"
)

func main() {
	ctx := context.Background()

	// common flags
	deviceFlag := &cli.StringFlag{
		Name:        "device",
		Aliases:     []string{"d"},
		Value:       "device-0",
		DefaultText: "device-0",
		Usage:       "event source device name",
		EnvVars:     []string{"DEVICE_NAME", "DEVICE_ID"},
	}

	fileFlag := &cli.StringFlag{
		Name:     "file",
		Aliases:  []string{"f"},
		Usage:    "metric template file path or URL",
		Required: true,
	}

	// commands
	stdoutCmd := &cli.Command{
		Name:  "stdout",
		Usage: "Mocks events and prints them in console",
		Flags: []cli.Flag{deviceFlag, fileFlag},
		Action: func(c *cli.Context) error {
			pub, err := stdout.NewEventSender(ctx)
			if err != nil {
				return errors.Wrapf(err, "error creating stdout publisher")
			}
			return execute(c, pub)
		},
	}

	httpCmd := &cli.Command{
		Name:  "http",
		Usage: "Mocks events and POSTs them to URL",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Aliases:  []string{"u"},
				Usage:    "URL to which the events will be posted",
				Required: true,
			},
			deviceFlag,
			fileFlag,
		},
		Action: func(c *cli.Context) error {
			postURL := c.String("url")
			if postURL == "" {
				return errors.New("url required")
			}
			pub, err := http.NewEventSender(ctx, postURL)
			if err != nil {
				return errors.Wrapf(err, "error creating http publisher")
			}
			return execute(c, pub)
		},
	}

	iothubCmd := &cli.Command{
		Name:  "iothub",
		Usage: "Mocks events and sends them to Azure IoT Hub",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "connect",
				Usage:    "connection string to IoT Hub",
				EnvVars:  []string{"CONN_STR", "IOTHUB_CONN_STR"},
				Required: true,
			},
			deviceFlag,
			fileFlag,
		},
		Action: func(c *cli.Context) error {
			connStr := c.String("connect")
			if connStr == "" {
				return errors.New("connect required")
			}
			pub, err := iothub.NewEventSender(ctx, connStr)
			if err != nil {
				return errors.Wrapf(err, "error creating iot hub publisher")
			}
			return execute(c, pub)
		},
	}

	// app
	app := &cli.App{
		Name:    "eventmaker",
		Version: Version,
		Usage:   "Mocks events with configurable format, metric range, and frequency",
		Action: func(c *cli.Context) error {
			logger.Printf("supported publishers: %v", event.ListPublishers())
			logger.Println("for more help use --help or -h flag")
			return nil
		},
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			stdoutCmd,
			httpCmd,
			iothubCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}

func execute(c *cli.Context, pub event.Publisher) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	ctx := context.Background()

	deviceID := c.String("device")
	file := c.String("file")

	em, err := mock.Make(ctx, deviceID, file, pub)
	if err != nil {
		return errors.Wrap(err, "error creating mocker")
	}

	cancel, errCh := em.Start(ctx)

	for {
		select {
		case e := <-errCh:
			cancel()
			return e
		case <-sigChan:
			cancel()
			return nil
		}
	}
}
