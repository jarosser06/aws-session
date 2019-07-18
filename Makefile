COMMAND_NAME = aws-session
SUPPORTED_SYSTEMS = linux darwin windows
RELEASE = $(shell git describe --always --tags)
GOLIST = $(shell go list ./... | grep -v /vendor/)
BUILDS_DIR = bin
LDFLAGS = "-X main.Version=${RELEASE}"

ifeq ($(OS), Windows_NT)
	BINARY_NAME = $(COMMAND_NAME).exe
	INSTALLPRE = "c:\aws-session\bin"
else
	BINARY_NAME = $(COMMAND_NAME)
	INSTALL_PRE = /usr/local/bin
endif

all: build

build: ${BUILDS_DIR}/${BINARY_NAME}
${BUILDS_DIR}/${BINARY_NAME}:
		@echo "Building binary..."
		go build -ldflags ${LDFLAGS} -o ${BUILDS_DIR}/${BINARY_NAME}

test:
		@test -z "$(gofmt -s -l . | tee /dev/stderr)"
		@test -z "$(golint $(GOLIST) | tee /dev/stderr)"
		@go test -v -race $(GOLIST)
		@go vet $(GOLIST)

clean:
		@echo "Cleaning up..."
		rm -rf bin
		rm -rf releases

rebuild: clean build

install:
		@echo "Installing..."
		cp ${BUILDS_DIR}/${BINARY_NAME} $(INSTALLPRE)/

cross_compile: linux darwin windows

linux: ${BUILDS_DIR}/linux/${COMMAND_NAME}
${BUILDS_DIR}/linux/${COMMAND_NAME}:
		@mkdir -p ${BUILDS_DIR}/linux
		GOOS=linux go build -v -ldflags ${LDFLAGS} -o bin/linux/${COMMAND_NAME}

darwin: ${BUILDS_DIR}/darwin/${COMMAND_NAME}
${BUILDS_DIR}/darwin/${COMMAND_NAME}:
		@mkdir -p ${BUILDS_DIR}/darwin
		GOOS=darwin go build -v -ldflags ${LDFLAGS} -o ${BUILDS_DIR}/darwin/${COMMAND_NAME}

windows: ${BUILDS_DIR}/windows/${COMMAND_NAME}.exe
${BUILDS_DIR}/windows/${COMMAND_NAME}.exe:
		@mkdir -p bin/windows
		GOOS=windows go build -v -ldflags ${LDFLAGS} -o ${BUILDS_DIR}/windows/${COMMAND_NAME}.exe

release: cross_compile
		@mkdir -p releases
		./.release.sh

.PHONY: linux darwin windows cross_compile rebuild
