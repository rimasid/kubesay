.PHONY: build test lint

build:
	go build -o bin/kubesay ./cmd/kubesay

test:
	go test ./...

lint:
	golangci-lint run
