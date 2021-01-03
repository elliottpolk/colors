BIN=colr
PKG=github.com/elliottpolk/$(BIN)
CLI_VERSION=`cat .version`
COMPILED=`date +%s`
GIT_HASH=`git rev-parse --short HEAD`
GOOS?=linux
BUILD_DIR=./build/bin

M = $(shell printf "\033[34;1m◉\033[0m")

default: all ;                                              		@ ## defaulting to clean and build

.PHONY: all
all: clean build

.PHONY: clean
clean: ; $(info $(M) running clean ...)                             @ ## clean up the old build dir
	@rm -vrf build

.PHONY: test
test: unit-test;													@ ## wrapper to run all testing

.PHONY: unit-test
unit-test: ; $(info $(M) running unit tests...)                     @ ## run the unit tests
	@go get -v -u
	@go test -cover ./...

.PHONY: update
update: clean; $(info $(M) updating deps...)                        @ ## update the deps
	@GOOS=$(GOOS) go get -v -u

.PHONEY: build-dir
build-dir: ;
	@[ ! -d "${BUILD_DIR}" ] && mkdir -vp "${BUILD_DIR}" || true    @ ## generate the build dir

.PHONY: build
build: build-dir; $(info $(M) building ...)                         @ ## build the binary
	@GOOS=$(GOOS) go build \
		-ldflags "-X main.version=$(CLI_VERSION) -X main.compiled=$(COMPILED) -X main.githash=$(GIT_HASH)" \
		-o $(BUILD_DIR)/$(BIN) ./main.go

.PHONEY: install
install: clean build; $(info $(M) installing locally ...) 						@ ## install binary locally
	@cp -v $(BUILD_DIR)/$(BIN) $(GOPATH)/bin/$(BIN)

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

