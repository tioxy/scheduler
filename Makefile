MAKEFILE := $(lastword $(MAKEFILE_LIST))

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

BINARY_NAME=scheduler
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_DARWIN=$(BINARY_NAME)_darwin
BUILD_PACKAGE=internal/app/main.go

CONTAINER_IMAGE=tioxy/scheduler
CONTAINER_TAG=latest
CONTAINER_FILE=Dockerfile
CONTAINER_CONTEXT=.

all: test build

build: 
		$(GOBUILD) -o $(BINARY_NAME) -v $(BUILD_PACKAGE)
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
		rm -f $(BINARY_DARWIN)

build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(BUILD_PACKAGE)
build-darwin:
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(BUILD_PACKAGE)

docker-build:
		docker build \
			-t "$(CONTAINER_IMAGE):$(CONTAINER_TAG)" \
			-f $(CONTAINER_FILE) $(CONTAINER_CONTEXT) \
			--build-arg IN_BINARY=$(BINARY_UNIX)
docker-push:
		docker image push "$(CONTAINER_IMAGE):$(CONTAINER_TAG)"

gen-image:
		$(MAKE) -f $(MAKEFILE) test
		$(MAKE) -f $(MAKEFILE) build-linux
		$(MAKE) -f $(MAKEFILE) docker-build
		$(MAKE) -f $(MAKEFILE) clean
