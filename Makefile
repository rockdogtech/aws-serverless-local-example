SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := $(shell echo $${PWD\#\#*/})
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 0.0.1
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean install uninstall fmt simplify check run

all: check install

# $(TARGET): $(SRC)
# 	@go build $(LDFLAGS) -o $(TARGET)

# build: $(TARGET)
# 	@true

clean:
	@rm -f $(TARGET)

# https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html
install:
	@brew install aws/tap/aws-sam-cli

uninstall: clean
	@brew remove aws/tap/aws-sam-cli

init:
	@sam init

build:
	@sam build --no-cached

invoke: build
	@sam local invoke --env-vars <(echo "{ \"Parameters\": `jq -n env` }")

# fmt:
# 	@gofmt -l -w $(SRC)
# 
# simplify:
# 	@gofmt -s -l -w $(SRC)
# 
# check:
# 	@test -z $(shell gofmt -l main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
# #	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done
# 	@go tool vet ${SRC}

# run: install
# 	@$(TARGET)
