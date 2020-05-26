GIT_COMMIT       =$(shell git rev-list -1 HEAD)
SERVICE_NAME     =eventmaker
WINDOWS          =$(SERVICE_NAME)-windows.exe
LINUX            =$(SERVICE_NAME)-linux
DARWIN           =$(SERVICE_NAME)
RELEASE_VERSION  =v0.4.8
RELEASE_COMMIT   =$(RELEASE_VERSION)-$(GIT_COMMIT)
DOCKER_USERNAME ?=$(DOCKER_USER)

.PHONY: mod test run send build exec image imagerun lint clean, tag
all: test

mod: ## Updates the go modules and vendors all dependancies 
	go mod tidy
	go mod vendor

test: mod ## Tests the entire project 
	go test -v -count=1 -race ./...
	# go test -v -count=1 -run TestMakeCPUEvent ./...

run: mod ## Runs the uncompiled code with stdout publisher 
	go run cmd/*.go stdout --file conf/example.yaml

build: mod ## Build local release binary
		env CGO_ENABLED=0 go build -ldflags "-X main.Version=$(RELEASE_COMMIT)" \
    	-mod vendor -o ./dist/$(SERVICE_NAME) ./cmd

exec: build ## Builds binaries and executes it 
	dist/eventmaker stdout --file conf/example.yaml

image: mod ## Builds docker iamge 
	docker build --build-arg VERSION=$(RELEASE_COMMIT) \
		-t "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)" .
	docker push "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)"

imagerun: ## Runs the pre-built docker image 
	docker run -e DEV_NAME="docker-1" \
		-ti "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)" \
		stdout --file https://raw.githubusercontent.com/mchmarny/eventmaker/master/conf/example.yaml

lint: ## Lints the entire project 
	golangci-lint run --timeout=3m

tag: ## Creates release tag 
	git tag $(RELEASE_VERSION)
	git push origin $(RELEASE_VERSION)

clean: ## Cleans dist directory
	go clean
	rm -fr ./dist/*

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
