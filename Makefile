MAKEFILE := $(lastword $(MAKEFILE_LIST))

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

BINARY_NAME=scheduler
BINARY_DARWIN=$(BINARY_NAME)_darwin
BINARY_UNIX=$(BINARY_NAME)_unix
BUILD_PACKAGE=cmd/scheduler/main.go

DOCKERFILE=Dockerfile
DOCKERFILE_CONTEXT=.
IMAGE_REPO=tioxy/scheduler
IMAGE_TAG=latest

AWS_DEFAULT_REGION=us-west-2

CLOUDFORMATION_FOLDER=infra/cloudformation
CLOUDFORMATION_TASKCAT_FILE=$(CLOUDFORMATION_FOLDER)/ci/taskcat.yml

PACKER_FOLDER=infra/packer
PACKER_DEFAULT_DISTRO=debian
PACKER_AWS_ACCESS_KEY=YOURACCESSKEY
PACKER_AWS_SECRET_KEY=YOURSECRETKEY

all: test build

build: 
		$(GOBUILD) -o $(BINARY_NAME) -v $(BUILD_PACKAGE)
build-darwin:
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(BUILD_PACKAGE)
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(BUILD_PACKAGE)
build-image:
		docker build \
			-t "$(IMAGE_REPO):$(IMAGE_TAG)" \
			-f $(DOCKERFILE) $(DOCKERFILE_CONTEXT) \
			--build-arg IN_BINARY=$(BINARY_UNIX)
build-ami:
		packer build \
			-var 'aws_access_key=$(PACKER_AWS_ACCESS_KEY)' \
			-var 'aws_secret_key=$(PACKER_AWS_SECRET_KEY)' \
			-var 'aws_region=$(AWS_DEFAULT_REGION)' \
			-var 'instance_type=t3.micro' \
			$(PACKER_FOLDER)/$(PACKER_DEFAULT_DISTRO)-base.json

clean: 
		$(GOCLEAN)
		$(GOCLEAN) -testcache
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
		rm -f $(BINARY_DARWIN)
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

test:
		$(GOTEST) -count=1 -v ./...
test-cf:
		taskcat -c $(CLOUDFORMATION_TASKCAT_FILE)
