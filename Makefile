MAKEFILE := $(lastword $(MAKEFILE_LIST))

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

BINARY_NAME=scheduler
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_DARWIN=$(BINARY_NAME)_darwin
BUILD_PACKAGE=cmd/scheduler/main.go

IMAGE_REPO=tioxy/scheduler
IMAGE_TAG=latest
DOCKERFILE=Dockerfile
DOCKERFILE_CONTEXT=.

all: test build

build: 
		$(GOBUILD) -o $(BINARY_NAME) -v $(BUILD_PACKAGE)
test: 
		$(GOTEST) -count=1 -v ./...
clean: 
		$(GOCLEAN)
		$(GOCLEAN) -testcache
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
		rm -f $(BINARY_DARWIN)

build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(BUILD_PACKAGE)
build-darwin:
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(BUILD_PACKAGE)
build-image:
		docker build \
			-t "$(IMAGE_REPO):$(IMAGE_TAG)" \
			-f $(DOCKERFILE) $(DOCKERFILE_CONTEXT) \
			--build-arg IN_BINARY=$(BINARY_UNIX)

clean-image:
		docker image rm -f "$(IMAGE_REPO):$(IMAGE_TAG)"
		docker image prune -f

gen-image:
		$(MAKE) -f $(MAKEFILE) test
		$(MAKE) -f $(MAKEFILE) build-linux
		$(MAKE) -f $(MAKEFILE) build-image
		$(MAKE) -f $(MAKEFILE) clean
push-image:
		docker image push "$(IMAGE_REPO):$(IMAGE_TAG)"
