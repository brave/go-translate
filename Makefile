.PHONY: all build test lint clean

all: lint test build

build:
	go run main.go

test:
	go test -v ./...

lint:
	golangci-lint run -E gofmt -E golint --exclude-use-default=false

clean:
	rm -f go-translate
