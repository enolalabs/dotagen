.PHONY: build test lint install dev clean

build:
	go build -o dotagen ./cmd/dotagen

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/dotagen

dev:
	go run ./cmd/dotagen serve

clean:
	rm -rf bin/
