VERSION := 0.0.1
PROJECT_NAME := $(shell pwd | awk -F/ '{ print $$NF }')
DOCKER_REPOSITORY := docker.io
DOCKER_USERNAME := tennashi
DOCKER_IMAGE_TAG := $(DOCKER_REPOSITORY)/$(DOCKER_USERNAME)/$(PROJECT_NAME):$(VERSION)
GAUGE_PATH := gauge

.PHONY: e2e
e2e:
	$(GAUGE_PATH) run -d e2e

.PHONY: test
test:
	go test -v ./...

.PHONY: update-golden
update-golden:
	go test -update ./...

.PHONY: lint
lint:
	go vet ./...
	golangci-lint run ./...

.PHONY: build
build:
	go build -o bin/pj .

.PHONY: image-build
image-build:
	docker build -t $(DOCKER_IMAGE_TAG) -f contrib/docker/Dockerfile .
