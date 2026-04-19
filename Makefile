GO ?= /usr/local/go/bin/go

.PHONY: build test lint run tidy

build:
	$(GO) build -o bin/bot ./cmd/bot

run:
	$(GO) run ./cmd/bot

test:
	$(GO) test ./...

lint:
	golangci-lint run

tidy:
	$(GO) mod tidy
