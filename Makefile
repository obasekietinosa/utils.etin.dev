.PHONY: all build run test lint docker-build

APP_NAME := server
DOCKER_IMAGE := utils-etin-dev

all: build

build:
	go build -o $(APP_NAME) ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v ./...

lint:
	go vet ./...

docker-build:
	docker build -t $(DOCKER_IMAGE) .
