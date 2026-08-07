.PHONY: build check clean fmt run test vet

GOCACHE ?= /tmp/bifunctional-go-cache
CONFIG ?= config/repos.json

build:
	mkdir -p bin
	GOCACHE=$(GOCACHE) go build -trimpath -ldflags="-s -w" -o bin/bifunctional-mcp ./cmd/bifunctional-mcp

run:
	GOCACHE=$(GOCACHE) go run ./cmd/bifunctional-mcp --config $(CONFIG)

fmt:
	gofmt -w .

test:
	GOCACHE=$(GOCACHE) go test -race -cover ./...

vet:
	GOCACHE=$(GOCACHE) go vet ./...

check: vet test build

clean:
	rm -rf bin coverage.out
