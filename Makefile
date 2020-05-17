GIT_COMMIT      =$(shell git rev-list -1 HEAD)
SERVICE_NAME    =eventmaker
RELEASE_VERSION =v0.1.1
RELEASE_COMMIT  =$(RELEASE_VERSION)-$(GIT_COMMIT)
DOCKER_USERNAME =$(DOCKER_USER)

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
	go run cmd/main.go

.PHONY: build
build: mod
	GOOS=darwin GOARCH=amd64 \
    go build -ldflags "-X main.Version=$(RELEASE_COMMIT)" \
    -mod vendor -o ./bin/$(SERVICE_NAME) ./cmd

.PHONY: exec
exec:
	bin/eventmaker

.PHONY: image
image: mod
	docker build --build-arg VERSION=$(RELEASE_COMMIT) \
		-t "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)" .
	docker push "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)"

.PHONY: lint
lint:
	golangci-lint run --timeout=3m

.PHONY: clean
clean:
	go clean
	rm -fr ./bin/*



