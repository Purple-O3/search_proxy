all: build

build:
	go build -ldflags "-s -w" -o scripts/search_proxy  cmd/main.go

run:
	cd scripts && ./search_proxy

clean:
	cd scripts && rm ./search_proxy

.PHONY: all build install run clean