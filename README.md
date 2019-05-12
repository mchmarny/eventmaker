# iot-event-maker

IoT Core event generator. In this demo we will create and configure Google Cloud IOT Core registry, configure its device, and send mocked events based on host's CPU/Mem/Load stats to that device which will persist them into CLoud PubSub topic which we will create for this demo.

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

> if you don't have `go` installed locally you can run the latest release binaries published here https://github.com/mchmarny/iot-event-maker/releases

After few lines of configuration output, you should see `eventmaker` posting to IOT Core

```shell
2019/05/12 14:43:26 Publishing: {"source_id":"demo-client","event_id":"eid-6ae3ab0d-a4d1-40a7-803c-f1e5158fe2b9","event_ts":"2019-05-12T21:43:26.303646Z","label":"my-metric","memFree":73.10791015625,"cpuFree":3500,"loadAvg1":2.65,"loadAvg5":2.88,"loadAvg15":3.46,"randomValue":1.2132739730794428}
```

The JSON payload in each one of these events looks something like this:

```json
{
    "source_id":"demo-client",
    "event_id":"eid-6ae3ab0d-a4d1-40a7-803c-f1e5158fe2b9",
    "event_ts":"2019-05-12T21:43:26.303646Z",
    "label":"my-metric",
    "memFree":73.10791015625,
    "cpuFree":3500,
    "loadAvg1":2.65,
    "loadAvg5":2.88,
    "loadAvg15":3.46,
    "randomValue":1.2132739730794428
}
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