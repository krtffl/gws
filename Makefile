REVISION := $(shell git rev-parse --short=8 HEAD || echo unknown)
VERSION ?= $(shell git describe --tags --always)
BUILT := $(shell date +%Y-%m-%dT%H:%M:%S)
BRANCH := $(shell git show-ref | grep "$(REVISION)" | grep -v HEAD | awk '{print $$2}' | sed 's|refs/remotes/origin/||' | sed 's|refs/heads/||' | sort | head -n 1)
MODULE_NAME := $(shell cat go.mod | head -n 1 | cut -d " " -f2- | cut -d "/" -f2-)
BASE_DIR := $(PWD)

GO_LDFLAGS ?= -X github.com/krtffl/get-well-soon/internal/version.Version=$(VERSION) \
			-X github.com/krtffl/get-well-soon/internal/version.Branch=$(BRANCH) \
			-X github.com/krtffl/get-well-soon/internal/version.Revision=$(REVISION) \
			-X github.com/krtffl/get-well-soon/internal/version.Built=$(BUILT) \
			-s \
			-w

BUILD_DIR=build
DIST_DIR=out


run:
	go run cmd/server/main.go

build: clean
	CGO_ENABLED=0 go build -a -ldflags "$(GO_LDFLAGS)" -o="$(BUILD_DIR)/gws" ./cmd/server

dist: clean
	CGO_ENABLED=0 go build -a -ldflags "$(GO_LDFLAGS)" -o="$(DIST_DIR)/gws" ./cmd/server
	@mkdir $(DIST_DIR)/config && cp ./config/config.default.yaml $(DIST_DIR)/config

dist-arm64: clean
	CGO_ENABLED=0 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -a -ldflags "$(GO_LDFLAGS)" -o="$(DIST_DIR)/gws" ./cmd/server
	@mkdir $(DIST_DIR)/config && cp ./config/config.default.yaml $(DIST_DIR)/config

clean:
	@rm -rf ./${BUILD_DIR} ./${DIST_DIR}


