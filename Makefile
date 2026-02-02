.PHONY: build test lint run clean

build:
	go build -o visor ./cmd/visor

test:
	go test ./...

lint:
	golangci-lint run

run:
	go run ./cmd/visor

clean:
	rm -f visor
