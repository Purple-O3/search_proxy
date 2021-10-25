BINARY=search_proxy
GODIR=../cmd

all: build 

build:
	go build -o ${BINARY} ${GODIR}/main.go

run:
	go run ${GODIR}/main.go

clean:
	rm -f ${BINARY}

.PHONY: all build install run clean
