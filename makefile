.PHONY: build debug clean

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BIN_DIR := bin
BIN_NAME := aurumServer

build_all: build_amd64 build_arm64 build_386 build_arm build_ppc64 build_ppc64le build_s390x build_mips64 build_windows build_windows_386

build:
ifeq ($(GOOS),windows)
	$(MAKE) build_windows
else
	$(MAKE) build_single
endif

build_amd64:
	GOARCH=amd64 $(MAKE) build_single

build_arm64:
	GOARCH=arm64 $(MAKE) build_single

build_386:
	GOARCH=386 $(MAKE) build_single

build_arm:
	GOARCH=arm GOARM=7 $(MAKE) build_single

build_ppc64:
	GOARCH=ppc64 $(MAKE) build_single

build_ppc64le:
	GOARCH=ppc64le $(MAKE) build_single

build_s390x:
	GOARCH=s390x $(MAKE) build_single

build_mips64:
	GOARCH=mips64 $(MAKE) build_single

build_windows:
	$(MAKE) build_single_windows GOOS=windows GOARCH=amd64

build_windows_386:
	$(MAKE) build_single_windows GOOS=windows GOARCH=386

build_single_windows:
	go build -o $(BIN_DIR)/$(BIN_NAME)-$(GOOS)-$(GOARCH).exe

build_single:
	go build -o $(BIN_DIR)/$(BIN_NAME)-$(GOOS)-$(GOARCH)

debug:
	go build -tags debug -gcflags "all=-N -l" -o $(BIN_DIR)/$(BIN_NAME)-$(GOOS)-$(GOARCH)-debug

clean:
	@rm -rf $(BIN_DIR)

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protos/*.proto