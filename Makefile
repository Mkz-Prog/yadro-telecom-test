.PHONY: build run test clean fmt lint

build: fmt lint
	go build -o bin/dungeon-challenge ./cmd/app/main.go

run:
	go run ./cmd/app/main.go -config config.json -events events

test:
	go test -v ./...

clean:
	rm -rf bin/

fmt:
	go fmt ./...
	go vet ./...

lint:
	golangci-lint run