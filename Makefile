BRANCH := $(shell git branch --show-current)
COMMIT := $(shell git log -1 --format='%h')
DATE := $(shell date +%Y%m%d%H%M%S)
IMAGE_TAG := ${BRANCH}-${DATE}-${COMMIT}

all: build

build:
	go build -ldflags "-s -w" -o scripts/search_proxy  cmd/main.go

run:
	cd scripts && ./search_proxy

clean:
	cd scripts && rm ./search_proxy

image-builder:
	docker build --file scripts/Dockerfile --tag search_proxy:${IMAGE_TAG} .

.PHONY: all build install run clean