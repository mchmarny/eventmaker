RELEASE_VERSION = v0.1.1
GIT_COMMIT      = $(shell git rev-list -1 HEAD)
GIT_VERSION     = $(shell git describe --always --abbrev=7 --dirty)

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
	go build -o ./bin/eventmaker ./cmd
	# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/eventmaker-linux

.PHONY: lint
lint:
	golangci-lint run --timeout=3m

.PHONY: clean
clean:
	go clean
	rm -fr ./bin/*



