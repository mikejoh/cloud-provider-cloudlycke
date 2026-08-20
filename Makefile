# Use for ldflags to compile binary with more information
#DATE := $(shell date "+%Y-%m-%d@%H:%M:%S")
#GITBRANCH := $(shell git rev-parse --abbrev-ref HEAD)
#GITCOMMIT := $(shell git rev-parse HEAD)
VERSION := "v0.1.0"

PROJECTNAME := cloudlycke-cloud-controller-manager

# Go related vars
GOVER = 1.13
GOCMD = go
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
LDFLAGS_LINUX := -ldflags "-s -w -X=cloudlycke.io/cloudlycke/pkg/version.Version=$(VERSION)"

# Environment info
GOARCH=amd64
OS := linux

# Make make silent
MAKEFLAGS += --silent

# Docker vars
DOCKER_HUB_REPOSITORY := mikejoh
DEFAULT_VERSION_TAG := latest

# Commands
## test: Run all tests recursively, be verbose and output test coverage (in percent)
.PHONY: test
test:
	@$(GOCMD) test ./... -v -cover

## testcov: Run all tests and generate an HTML coverage report.
.PHONY: testcov
testcov:
	@$(GOCMD) test ./... -coverprofile=coverage.out
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html

## vet: Run go vet against the code.
.PHONY: vet
vet:
	@$(GOCMD) vet ./...

## lint: Run golangci-lint against the code.
.PHONY: lint
lint:
	@golangci-lint run -v --timeout=15m ./...

## dep: Verify and tidy module dependencies.
.PHONY: dep
dep:
	@$(GOCMD) mod tidy
	@$(GOCMD) mod verify

## build: Build the binary for CI (native OS/ARCH, no ldflags).
.PHONY: build
build:
	@$(GOCMD) build -o build/cloud-provider-cloudlycke ./cmd

## clean: Clean the working directory.
.PHONY: clean
clean:
	@echo " > Cleaning the working directory..."
	@rm -rf ./bin/ ./build/

## build-linux: Build Linux amd64 binary locally.
.PHONY: build-linux
build-linux:
	@echo " $(LDFLAGS_LINUX) "
	@echo "  >  Building Linux binary..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) $(GOCMD) build $(LDFLAGS_LINUX) -o ./bin/$(PROJECTNAME) ./cmd/$(wildcard *.go)

## build-linux-docker: Build Linux amd64 binary locally but through a Docker container.
.PHONY: build-linux-docker
build-linux-docker:
	@echo " $(LDFLAGS_LINUX) "
	@echo " > Build Linux binary in a Docker container..."
	@docker run --rm -it -v $(GOBASE):/$(PROJECTNAME) -w="/$(PROJECTNAME)" golang:$(GOVER)-alpine sh -c "CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) $(GOCMD) build $(LDFLAGS_LINUX) -o $(PROJECTNAME) $(wildcard *.go)"

## Build Docker image
.PHONY: docker-build
docker-build:
	@echo " > Building $(PROJECTNAME) Docker image"
	@docker build . -t $(DOCKER_HUB_REPOSITORY)/$(PROJECTNAME):$(DEFAULT_VERSION_TAG)

## Push Docker image to Docker registry
.PHONY: docker-push
docker-push:
	@echo " > Pushing $(PROJECTNAME) Docker image"
	@docker push $(DOCKER_HUB_REPOSITORY)/$(PROJECTNAME):$(DEFAULT_VERSION_TAG)

## Build and compile binary, build Docker image, push to Docker registry.
.PHONY: release
release: build-linux docker-build docker-push
	@echo " Releasing $(PROJECTNAME)"
