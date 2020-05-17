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

```shell
make run
``` 

Your events will be published to Azure Service Bus. To find the connection string

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


