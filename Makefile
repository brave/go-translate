.PHONY: all build test lint clean

all: lint test build

build:
	go run main.go

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -f go-translate
