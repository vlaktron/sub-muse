BINARY := sub-muse
PKG    := ./...
INTERNAL := ./internal/...

.PHONY: all build run test coverage lint fmt vet tidy install-tools clean

all: build

## Build

build:
	go build -o ./bin/$(BINARY) .

run:
	go run .

## Dependencies

tidy:
	go mod tidy

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	mise install

## Quality

fmt:
	gofmt -w .

vet:
	go vet $(PKG)

lint:
	golangci-lint run --timeout 3m $(PKG)

## Tests

test:
	go test -v $(INTERNAL)

test-short:
	go test -short $(INTERNAL)

coverage:
	go test -v $(INTERNAL) -coverprofile=coverage.out
	go tool cover -html=coverage.out

## Cleanup

clean:
	rm -f $(BINARY) coverage.out
