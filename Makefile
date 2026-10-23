BINARY := bin/server
MOCKERY ?= $(shell command -v mockery 2>/dev/null || echo $(HOME)/go/bin/mockery)

IMAGE ?= condominium-manager

.PHONY: build test mocks docker

docker:
	docker build -f build/docker/Dockerfile -t $(IMAGE) .

build:
	go build -o $(BINARY) ./cmd/server

test:
	go test ./...

mocks:
	$(MOCKERY)
