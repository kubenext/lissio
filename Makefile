# Copyright (c) 2019 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

SHELL=/bin/bash
BUILD_TIME=$(shell date -u +%Y-%m-%dT%T%z)
GIT_COMMIT=$(shell git rev-parse --short HEAD)

LD_FLAGS= '-X "main.buildTime=$(BUILD_TIME)" -X main.gitCommit=$(GIT_COMMIT)'
GO_FLAGS= -ldflags=$(LD_FLAGS)
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

VERSION ?= v0.7.0

ifdef XDG_CONFIG_HOME
	LISSIO_PLUGINSTUB_DIR ?= ${XDG_CONFIG_HOME}/lissio/plugins
# Determine in on windows
else ifeq ($(OS),Windows_NT)
	LISSIO_PLUGINSTUB_DIR ?= ${LOCALAPPDATA}/lissio/plugins
else
	LISSIO_PLUGINSTUB_DIR ?= ${HOME}/.config/lissio/plugins
endif

.PHONY: version
version:
	@echo $(VERSION)

# Run all tests
.PHONY: test
test: generate
	@echo "-> $@"
	@env go test -v ./internal/... ./pkg/...

# Run govet
.PHONY: vet
vet:
	@echo "-> $@"
	@env go vet ./internal/... ./pkg/...

lissio-dev:
	@mkdir -p ./build
	@env $(GOBUILD) -o build/lissio $(GO_FLAGS) -v ./cmd/lissio

lissio-docker:
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o /lissio $(GO_FLAGS) -v ./cmd/lissio

generate:
	@echo "-> $@"
	@find pkg internal -name fake -type d | xargs rm -rf
	@go generate -v ./pkg/... ./internal/...

go-install:
	@env GO111MODULE=on $(GOINSTALL) github.com/GeertJohan/go.rice
	@env GO111MODULE=on $(GOINSTALL) github.com/GeertJohan/go.rice/rice
	@env GO111MODULE=on $(GOINSTALL) github.com/golang/mock/gomock
	@env GO111MODULE=on $(GOINSTALL) github.com/golang/mock/mockgen
	@env GO111MODULE=on $(GOINSTALL) github.com/golang/protobuf/protoc-gen-go

# Remove all generated go files
.PHONY: clean
clean:
	@find pkg internal -name fake -type d | xargs rm -rf
	@rm ./pkg/icon/rice-box.go

web-deps:
	@cd web; npm ci

web-build: web-deps
	@cd web; npm run build
	@go generate ./web

web-test: web-deps
	@cd web; npm run test:headless

ui-server:
	LISSIO_DISABLE_OPEN_BROWSER=1 LISSIO_LISTENER_ADDR=localhost:7777 $(GOCMD) run ./cmd/lissio/main.go $(LISSIO_FLAGS)

ui-client:
	@cd web; API_BASE=http://localhost:7777 npm run start

gen-electron:
	@GOCACHE=${HOME}/cache/go-build astilectron-bundler -v -c configs/electron/bundler.json

.PHONY: changelogs
changelogs:
	hacks/changelogs.sh

.PHONY: release
release:
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push --follow-tags

.PHONY: ci
ci: test vet web-test web-build lissio-dev

.PHONY: ci-quick
ci-quick:
	@cd web; npm run build
	@go generate ./web
	make lissio-dev

install-test-plugin:
	@echo $(LISSIO_PLUGINSTUB_DIR)
	mkdir -p $(LISSIO_PLUGINSTUB_DIR)
	go build -o $(LISSIO_PLUGINSTUB_DIR)/lissio-sample-plugin github.com/kubenext/lissio/cmd/lissio-sample-plugin

.PHONY:
build-deps:
