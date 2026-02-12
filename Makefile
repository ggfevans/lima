BINARY_NAME=linkedin-tui
BUILD_DIR=bin
GO=go

.PHONY: build run clean fmt vet

build:
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/linkedin-tui/

run:
	$(GO) run ./cmd/linkedin-tui/

clean:
	rm -rf $(BUILD_DIR)

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...
