.PHONY: all build run test lint

APP_NAME := server

all: build

build:
	go build -o $(APP_NAME) ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v ./...

lint:
	go vet ./...
