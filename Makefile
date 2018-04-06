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
LDFLAGS 			:=-ldflags "-s -w -X damo/pkg/config.version=$(VERSION_LABEL)"

AWS_DEFAULT_REGION 	:=ap-southeast-2

DOCKER_BUILD_IMAGE := damo-build
DOCKER_DEPLOY_IMAGE := damo-deploy
DOCKER_BUILD := docker run --rm -it -p 3000:3000 -v $(PWD):/go/src/damo -e ENVIRONMENT=$(ENVIRONMENT) $(DOCKER_BUILD_IMAGE)
DOCKER_DEPLOY := docker run --rm  -it -v $(PWD):/src -e ENVIRONMENT=$(ENVIRONMENT) \
	-e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) \
	-e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) \
	-e AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION) \
	$(DOCKER_DEPLOY_IMAGE)

include config/$(ENVIRONMENT)-env.sh

all: run

.main:
	@echo "Building package"
	@GOOS=linux go build -a $(LDFLAGS) -o $(RELEASE_DIR)/$(BINARY) $(BINARY).go

.run:
	@echo "Running locally"
	@source config/$(ENVIRONMENT)-env.sh && go run $(LDFLAGS) $(BINARY).go -l

.deploy:
	@echo "Deploying functions"
	@source config/$(ENVIRONMENT)-env.sh && sls deploy --stage $(ENVIRONMENT) | tee .deploy.txt
	@echo Invoking health check endpoint...
	@APIURL=$$(cat .deploy.txt | sed -n -e "s/.*\(https:\/\/.*execute-api\..*\.amazonaws\.com\).*/\1/gp"); \
	curl $${APIURL}/$(ENVIRONMENT)/health
	@rm .deploy.txt

.go-deps:
	@echo "Installing go dependencies" 
	@dep ensure -v

.node-deps:
	@echo "Installing node dependencies"
	@npm i

.deps: .go-deps .node-deps

.test:
	@echo "Running tests"
	@golint $(TEST_DIR)
	@go vet $(TEST_DIR)
	@go test -cover $(TEST_DIR)

docker:
	@echo "Building deploy container"
	@docker build -t $(DOCKER_DEPLOY_IMAGE) -f Dockerfile.deploy .
	@echo "Building build container"
	@docker build -t $(DOCKER_BUILD_IMAGE) -f Dockerfile.build .

main:
	@$(DOCKER_BUILD) make .main

run: 
	@$(DOCKER_BUILD) make .run

deploy: main test
	@$(DOCKER_DEPLOY) make .deploy

deps: 
	@$(DOCKER_BUILD) make .go-deps
	@$(DOCKER_DEPLOY) make .node-deps

test:
	@$(DOCKER_BUILD) make .test

