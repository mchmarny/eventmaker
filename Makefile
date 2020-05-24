GIT_COMMIT      ?=$(shell git rev-list -1 HEAD)
SERVICE_NAME    ?=eventmaker
RELEASE_VERSION ?=v0.1.2
RELEASE_COMMIT  ?=$(RELEASE_VERSION)-$(GIT_COMMIT)
DOCKER_USERNAME ?=$(DOCKER_USER)

all: test

.PHONY: mod
mod:
	go mod tidy
	go mod vendor

.PHONY: test
test: mod
	go test -v -count=1 -race ./...
	# go test -v -count=1 -run TestMakeCPUEvent ./...

.PHONY: run
run: mod
	go run cmd/*.go --file "conf/example.yaml"

.PHONY: send
send: mod
	go run cmd/*.go --file "conf/thermostat.yaml" --publisher iothub

.PHONY: build
build: mod
	CGO_ENABLED=0 go build -ldflags "-X main.Version=$(RELEASE_COMMIT)" \
    -mod vendor -o ./dist/$(SERVICE_NAME) ./cmd

.PHONY: image
image: mod
	docker build --build-arg VERSION=$(RELEASE_COMMIT) \
		-t "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)" .
	docker push "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)"

.PHONY: exec-image
exec-image:
	docker run -e CONN_STR=$(CONN_STR) -e DEV_NAME='test-run-1' \
						 -ti mchmarny/eventmaker:v0.1.1 /eventmaker --file conf/example.yaml

.PHONY: lint
lint:
	golangci-lint run --timeout=3m

.PHONY: clean
clean:
	go clean
	rm -fr ./bin/*



