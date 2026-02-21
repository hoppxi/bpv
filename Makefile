.PHONY: build dev clean web-build web-dev help

VERSION := $(shell grep 'Version:' cmd/bpv/root.go | head -1 | sed 's/.*"\(.*\)".*/\1/')
BINARY := bpv
BUILD_FLAGS := -ldflags="-s -w"

build:
	CGO_ENABLED=1 go build $(BUILD_FLAGS) -o $(BINARY) ./cmd/bpv
	CGO_ENABLED=1 go build $(BUILD_FLAGS) -o $(BINARY)d ./cmd/bpvd

dev: build
	./$(BINARY) $(MUSIC_DIR)

clean:
	rm -f $(BINARY) $(BINARY)d
	go clean

web-build:
	cd web && npm install && npm run build

web-dev:
	cd web && npm run dev

tidy:
	go mod tidy

help:
	@echo "BPV - Build Targets"
	@echo ""
	@echo "  make build       Build the BPV binary"
	@echo "  make dev         Build and run TUI (set MUSIC_DIR=~/Music)"
	@echo "  make clean       Remove build artifacts"
	@echo "  make web-build   Build Vue frontend"
	@echo "  make web-dev     Run Vue dev server"
	@echo "  make tidy        Run go mod tidy"
