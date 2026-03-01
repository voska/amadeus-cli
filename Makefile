VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X main.version=$(VERSION)
BIN     := bin/amadeus

.PHONY: build test lint clean install generate

build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/amadeus

test:
	go test -race ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ dist/

install: build
	cp $(BIN) $(GOPATH)/bin/amadeus 2>/dev/null || cp $(BIN) ~/go/bin/amadeus

generate:
	go generate ./...
