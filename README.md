# eventmaker

Utility to mock events with configurable format, random metric range, and frequency. The supported event publishers are 

* [x] `stdout` - prints events to the console 
* [x] `http` - posts events to specified URL ([how-to](doc/HTTP.md))
* [x] `iothub` - sends events to Azure IoT Hub ([how-to](doc/Azure-IoT-Hub.md))
* [x] `eventhub` - sends events to Azure Event Hub ([how-to](doc/Azure-Event-Hub.md))
* [ ] `iotcore` - sends events to GCP IoT Core 
* [x] `pubsub` - sends events to GCP Pub/Sub ([how-to](doc/GCP-PubSub.md))
* [ ] `sns` - sends events to AWS SNS
* [ ] `sqs` - sends events to AWS SQS


[![Test](https://github.com/mchmarny/eventmaker/workflows/Test/badge.svg)](https://github.com/mchmarny/eventmaker/actions?query=workflow%3ATest) [![Release](https://github.com/mchmarny/eventmaker/workflows/Release/badge.svg)](https://github.com/mchmarny/eventmaker/actions?query=workflow%3ARelease) [![Go Report Card](https://goreportcard.com/badge/github.com/mchmarny/eventmaker)](https://goreportcard.com/report/github.com/mchmarny/eventmaker) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mchmarny/eventmaker)

## usage 

To run `eventmaker` and publish events to the console run

```shell
./eventmaker stdout --file conf/thermostat.yaml
```

The file parameter can be local file path or a remote URL (e.g. [thermostat.yaml](https://raw.githubusercontent.com/mchmarny/eventmaker/master/conf/thermostat.yaml)). For more information about all the flags and commands supported by the `eventmaker` use the `--help` or `-h` flag 

```shell 
./eventmaker -h
```

For instructions how to publish events to other targets see the how-to links above

## events 

The mocked event look like this

```json
{
    "id": "fdf612b9-34a5-445e-9941-59c404ea9bef",
    "src_id": "device-1",
    "time": 1589745397,
    "label": "temp",
    "data": 70.79129651786347,
    "unit": "celsius"
}
```

* `id` - is a globally unique ID generated for each event 
* `src_id` - is the device ID or name (configured using the `--device` flag at launch)
* `time` - is the epoch (aka Unix time) of when the event was generated 

The `label`, `data`, and `unit` elements are based on the metrics defined in the template (see [metrics](#metrics) section below)

## metrics 

`eventmaker` dynamically configures metrics based on template where you can 

To create events, define one or more metrics in a template file where you configure the frequency in which that metric should be generated along with its type and range of value. For example, configuration for a virtual thermostat with temperature and humidity metrics would look like this

```yaml
--- 
metrics: 
- label: temperature
  frequency: "1s"
  unit: celsius
  template: 
    type: float
    min: 86.1
    max: 107.5
- label: humidity
  frequency: "1s"
  unit: percent
  template: 
    type: int
    min: 0
    max: 100
```

 That file where you define these metrics cab be either local (e.g. `--file conf/example.yaml`) or loaded from a remote URL (e.g. `--file https://raw.githubusercontent.com/mchmarny/eventmaker/master/conf/thermostat.yaml`)

## example 

To use `eventmaker` as a library, start by importing the package 

```shell
go get -u github.com/mchmarny/eventmaker
```

Then create some metric templates 

```go
metrics := []event.MetricTemplate{
  {
    Label:     "temperature",
    Unit:      "celsius",
    Frequency: time.Duration(1 * time.Second),
    Template: event.ValueTemplate{
      Type: "float",
      Min:  39.1,
      Max:  73.5,
    },
  },
  {
    Label:     "humidity",
    Unit:      "percent",
    Frequency: time.Duration(1 * time.Minute),
    Template: event.ValueTemplate{
      Type: "int",
      Min:  0,
      Max:  100,
    },
  },
}
```

Create a publisher target using one of the event senders (`stdout`, `iothub`, `http`)

```go
target, err := stdout.NewEventSender(ctx)
if err != nil {
  panic(err)
}
```

And finally create an instance of `mocker` with the previously created `metrics` and `target` and start it

```go
mocker, err := mock.New("demo-device-1", metrics, target)
if err != nil {
  panic(err)
}

mocker.Start(ctx)
```

A complete working example is available [here](example/main.go)

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License
This software is released under the [Apache v2 License](./LICENSE)


