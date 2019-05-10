# iot-event maker

IoT Core event generator. In this demo we will create and configure Google Cloud IOT Core registry, configure a device on that gateway, and send mocked events to that device which will persist them into CLoud PubSub topic which we will create for this demo.

## Configure IoT Core

First, lets configure the IOT Core registry and a device as well as the Cloud PubSub topic for our events.

### Generate TLS certificates

Before we start, we need to generate both your private and public keys which our client and IOT Core registry will use to secure its communications.

First, download the root CA from Google

```shell
curl https://pki.google.com/roots.pem > ./root-ca.pem
```

Then create private and public device keys

```shell
openssl req -x509 -nodes -newkey rsa:2048 \
            -keyout device.key.pem \
            -out device.crt.pem \
            -days 365 \
            -subj "/CN=demo"
```

To make sure everything worked, make sure you have these 3 files in your directory:

```shell
device.crt.pem
device.key.pem
root-ca.pem
```

### Create PubSub topic

IOT Core will drain all published events to specified topic. You can use one that already exists or create a brand new:

```shell
gcloud pubsub topics create demo-iot-events
```

### Create Registry

Now that you have the topic and all the necessary keys, you can create a IOT Core Registry

```shell
gcloud iot registries create demo-reg \
		--project=${GCP_PROJECT} \
		--region=us-central1 \
		--event-notification-config=topic=demo-iot-events
```

### Configure Device on the Registry

To add new device to the above created registry:

```shell
gcloud iot devices create demo-device-1 \
		--project=${GCP_PROJECT} \
		--region=us-central1 \
		--registry=demo-reg \
		--public-key path=device.crt.pem,type=rs256
```

> Note, this assumes your certificates you created above are in current directory

## Run

To send data to IoT Core, run

```shell
go mod tidy
go run *.go --project=${GCP_PROJECT} --region=us-central1 --registry=demo-reg \
			--device=demo-device-1 --ca="root-ca.pem" --key="device.key.pem" \
			--src="demo-client" --metric="my-metric" --range="0.01-10.00" --freq="3s"
```

Besides references to the IOT Core resources we created above, there are a few parameters worth explaining:

* `--src` - Unique name of the device from which you are sending events. This is used to identify the specific sender of events in case you are running multiple clients (e.g. `my-laptop`)
* `--metric` - Name of the metric to generate that will be used as `label` in the sent event (e.g. `friction`)
* `--range` - Range of the random data points that will be generated for the above defined metric (e.g. `0.01-10.00` which means floats between `0.01` and `10.00`)
* `--freq` - Frequency in which these events will be sent to IoT Core (e.g. `3s` which means every 3 sec.)

Most of these have reasonable defaults so if you are not sure, just omit it.


# Cleanup

To delete all the resources created on GCP (topic, registry, and devices) run:

```shell
gcloud iot devices delete demo-device-1 \
    --project=${GCP_PROJECT} \
    --registry=demo-reg \
    --region=us-central1

gcloud iot registries delete demo-reg \
    --project=${GCP_PROJECT} \
    --region=us-central1

gcloud pubsub topics delete demo-iot-events
```