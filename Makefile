# Copyright (c) 2021 Red Hat, Inc.
# Copyright Contributors to the Open Cluster Management project

BINDIR ?= output

include build/Configfile

USE_VENDORIZED_BUILD_HARNESS ?=

ifndef USE_VENDORIZED_BUILD_HARNESS
-include $(shell curl -s -H 'Authorization: token ${GITHUB_TOKEN}' -H 'Accept: application/vnd.github.v4.raw' -L https://api.github.com/repos/stolostron/build-harness-extensions/contents/templates/Makefile.build-harness-bootstrap -o .build-harness-bootstrap; echo .build-harness-bootstrap)
else
-include vbh/.build-harness-vendorized
endif

default::
	@echo "Build Harness Bootstrapped"

.PHONY: deps
deps:
	go mod tidy

.PHONY: insights-metrics
insights-metrics:
	 CGO_ENABLED=1 go build -a -v -i -installsuffix cgo -ldflags '-s -w' -o $(BINDIR)/insights-metrics ./

.PHONY: build
build: insights-metrics

.PHONY: build-linux
build-linux:
	make insights-metrics GOOS=linux

.PHONY: lint
lint:
	GOPATH=$(go env GOPATH)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${GOPATH}/bin" v1.64.6
	CGO_ENABLED=1 GOGC=25 golangci-lint run --timeout=3m
	
run:
	 go run main.go

.PHONY: test
test:
	 go test ./... -v -coverprofile cover.out

.PHONY: coverage
coverage:
	 go tool cover -html=cover.out -o=cover.html


.PHONY: clean
clean::
	go clean
	rm -f cover*
	rm -rf ./$(BINDIR)
