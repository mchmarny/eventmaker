# eventmaker

Creates and configures Azure IoT Hub with virtual devices and sends mocked events to them.

> Note, this is still work in progress

## setup 

Export the name of your Azure IoT Hub and an example device

```shell
export HUB_NAME="cloudylabs"
export DEV_NAME="device-1"
```

## hub 

> assumes you have resource group and location defaults set 

Create a IoT Hub with a standard pricing tier

```shell
az iot hub create --name $HUB_NAME --sku S1
```

## device 

Create the device in the identity registry 

```shell
az iot hub device-identity create \
  --hub-name $HUB_NAME \
  --device-id $DEV_NAME
```

Retrieve device connection string

```shell
export CONN_STR=$(az iot hub device-identity show-connection-string \
  --device-id $DEV_NAME \
  --hub-name $HUB_NAME \
  -o tsv)
```

## run 

Build the executable 

```shell
make build
``` 

And start sending 

```shell
bin/eventmaker --metric "temp|celsius|float|0:72.1|3s"
```

Where `--metric` is the content of the metric you want to send. For example, the above metric would generate an event similar to this one every 3 seconds

```json
{
    "id":"fdf612b9-34a5-445e-9941-59c404ea9bef",
    "src_id":"client-1",
    "time":1589745397,
    "label":"temp",
    "data":70.79129651786347,
    "unit":"celsius"
}
```

The format of metrics are as follow 

`<label>|<unit>|<type of content in data field>|<range of data to generate>|<frequency>`

The supported types are `int`, `float`, and `bool` as well as most common derivates of these (e.g. `int64` or `float32`).

The `ranges` follow `min:max` format. So int in between 0 and 100 would be formatted as `0:100`. This way you can include negative numbers. 

Finally, the `frequency` follows standard go `time.Duration` format (e.g. `1s` for every second, `2m` for every 2 minutes, or `3h` for every 3 hours)

The one defaults you set using environment variables is the device name (`DEV_NAME`) which is the device ID associated with this client (default: `device-1`)

To find the Azure Service Bus here these events will be published:

```shell
az iot hub show \
  --name $HUB_NAME \
  --query "properties.eventHubEndpoints.events.endpoint" \
  -o tsv
```


## cleanup 

Delete hub and all of it's devices 

```shell
az iot hub delete --name $HUB_NAME
```


## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License
This software is released under the [Apache v2 License](./LICENSE)


