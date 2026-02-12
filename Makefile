BINARY_NAME=li-cli
BUILD_DIR=bin
GO=go

.PHONY: build run clean fmt vet

build:
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/li-cli/

run:
	$(GO) run ./cmd/li-cli/

clean:
	rm -rf $(BUILD_DIR)

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...
