.PHONY: all build test lint clean

all: lint test run

GIT_VERSION := $(shell git describe --abbrev=8 --dirty --always --tags)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date +%s)

run: build
	./go-translate

build:
	CGO_ENABLED=0 go build -v -ldflags "-X main.version=${GIT_VERSION} -X main.buildTime=${BUILD_TIME} -X main.commit=${GIT_COMMIT}" -o go-translate main.go

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -f go-translate
