.PHONY: build test lint install dev clean

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build:
	go build -ldflags="-X github.com/enolalabs/dotagen/v2/internal/cli.version=$(VERSION)" -o dotagen ./cmd/dotagen

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/dotagen

dev:
	go run ./cmd/dotagen serve

clean:
	rm -rf bin/ dotagen
