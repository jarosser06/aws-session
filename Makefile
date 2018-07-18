INSTALLPRE = /usr/local
BINARY_NAME = aws-session
SUPPORTED_SYSTEMS = linux darwin
RELEASE = $(shell git describe --always --tags)
GOLIST = $(shell go list ./... | grep -v /vendor/)

all: build

build: bin/aws-session
bin/aws-session:
		@echo "Building binary..."
		go build -o bin/${BINARY_NAME}

test:
		@test -z "$(gofmt -s -l . | tee /dev/stderr)"
		@test -z "$(golint $(GOLIST) | tee /dev/stderr)"
		@go test -v -race $(GOLIST)
		@go vet $(GOLIST)

clean:
		@echo "Cleaning up..."
		rm -rf bin

rebuild: clean build

install:
		@echo "Installing..."
		cp bin/aws-session $(INSTALLPRE)/bin/

cross_compile: linux darwin

linux: bin/linux/aws-session
bin/linux/aws-session:
		@mkdir -p bin/linux
		GOOS=linux go build -v -o bin/linux/aws-session

darwin: bin/darwin/aws-session
bin/darwin/aws-session:
		@mkdir -p bin/aws-session
		GOOS=darwin go build -v -o bin/darwin/aws-session

.PHONY: linux darwin cross_compile rebuild
