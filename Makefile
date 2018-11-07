
# Go parameters
BINARY_NAME="knative-iot-client"
PROJECT="s9-demo"
REGISSTRY="knative-demo"
DEVICE="knative-demo-client"
REGION="us-central1"
TOPIC_DATA="knative-iot-demo"
TOPIC_DEVICE="knative-iot-demo-device"
CA_KEY="root-ca.pem"
DEVICE_KEY="device.key.pem"
NUMBER_OF_EVENTS_TO_SEND=3
EVENT_SRC="next18-demo-client"

all: test
build:
	go build -o ./bin/$(BINARY_NAME) -v

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(BINARY_NAME)

clean:
	go clean
	rm -f ./bin/$(BINARY_NAME)

run:
	go run *.go --project $(PROJECT) --region $(REGION) --registry $(REGISSTRY) \
				--device $(DEVICE) --ca $(CA_KEY) --key $(DEVICE_KEY) \
				--events $(NUMBER_OF_EVENTS_TO_SEND) --src $(EVENT_SRC)

deps:
	go mod tidy

certs:
	openssl req -x509 -nodes -newkey rsa:2048 \
				-keyout device.key.pem \
				-out device.crt.pem \
				-days 365 \
				-subj "/CN=unused"
	curl https://pki.google.com/roots.pem > ./root-ca.pem