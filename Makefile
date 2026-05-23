BINARY := stackdiff
CMD     := ./cmd/stackdiff
GOFLAGS := -trimpath

.PHONY: all build test lint clean

all: build

build:
	go build $(GOFLAGS) -o bin/$(BINARY) $(CMD)

test:
	go test ./...

test-verbose:
	go test -v ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/

run: build
	@echo "Usage: bin/$(BINARY) <state-a.json> <state-b.json>"
