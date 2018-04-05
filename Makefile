PROJECT     := damo
VERSION		:= 0.0.1
RELEASE_DIR	:= bin
BINARY 		:= main
TEST_DIR 	:= ./pkg/...

ENVIRONMENT	?=dev
SHELL		:=/usr/bin/env bash

BUILD_VCS_NUMBER ?= $(shell git rev-parse HEAD)
BUILD_NUMBER     ?= $(shell hostname)
VERSION_LABEL    ?= $(VERSION)-$(BUILD_NUMBER).$(shell SHA=$(BUILD_VCS_NUMBER) && echo $${SHA:0:7})

AWS_DEFAULT_REGION 	:=ap-southeast-2

DOCKER_IMAGE : damo-build
DOCKER: docker run --rm -e ENVIRONMENT=$(ENVIRONMENT) \
	-e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) \
	-e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) \
	-e AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION) \
	$(DAMO_BUILD)

include config/$(ENVIRONMENT)-env.sh

all: run

.main:
	GOOS=linux go build -a -ldflags "-s -w -X main.version=$(VERSION_LABEL)" -o $(RELEASE_DIR)/$(BINARY) $(BINARY).go

.run:
	source config/$(ENVIRONMENT)-env.sh && go run $(BINARY).go -l

.deploy:
	source config/$(ENVIRONMENT)-env.sh && sls deploy --stage $(ENVIRONMENT) | tee .deploy.txt
	@echo Invoking health check endpoint...
	APIURL=$$(cat .deploy.txt | sed -n -e "s/.*\(https:\/\/.*execute-api\..*\.amazonaws\.com\).*/\1/gp"); \
	curl $${APIURL}/$(ENVIRONMENT)/health
	@rm .deploy.txt

.deps:
	dep ensure -v

.test:
	golint $(TEST_DIR)
	go vet $(TEST_DIR)
	go test -cover $(TEST_DIR)

docker:
	docker build -t $(DAMO_BUILD) .

main:
	$(DOCKER) make .main

run: 
	$(DOCKER) make run

deploy: docker
	$(DOCKER) make .deploy
