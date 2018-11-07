# iot-code-client

Knative IoT Core Demo

## Setup

You will need to created a self-signed cert for each of the devices. You can follow the IoT Core documentation or just run the include `make` command

```shell
make certs
```

This will generate both your private and public keys as well as download the IoT Core CA key.

When done, you can navigate to the device configuration page in IoT Core and upload the public key (`device.crt.key`) as `RS256_X509` format.


## Build

To compile the client, execute

```shell
make build
```

## Run

To execute the sample run the following command

```shell
bin/iot-code-client --project "your project id" \
                    --region "registry region e.g. us-central1" \
                    --registry "your registry name" \
		    --device "your device ID" \
                    --ca "CA key [root-ca.pem]" \
                    --key "Your private key [device.key.pem]" \
		    --events "Number of events you want to send [3]" \
                    --src "Source of the events e.g. knative-client"

```
