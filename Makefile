BINARY := stackdiff
CMD     := ./cmd/stackdiff
GOFLAGS := -trimpath

.PHONY: all build test lint clean install

all: build

build:
	go build $(GOFLAGS) -o bin/$(BINARY) $(CMD)

test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ coverage.out coverage.html

install: build
	cp bin/$(BINARY) $(GOPATH)/bin/$(BINARY)

run: build
	@echo "Usage: bin/$(BINARY) <state-a.json> <state-b.json>"
