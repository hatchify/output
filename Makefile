all:

install: export GO111MODULE=on
install:
	go install github.com/hatchify/output

lint:
	golangci-lint run --enable-all

test:
	go test
