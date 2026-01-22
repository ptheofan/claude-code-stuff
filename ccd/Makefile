.PHONY: build install clean test run

BINARY_NAME=ccd
VERSION=1.0.0
BUILD_DIR=./build

build:
	@echo "Building $(BINARY_NAME)..."
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME) ./cmd/ccd
	@echo "Moving to parent folder as 'deploy'..."
	mv $(BINARY_NAME) ../deploy
	chmod +x ../deploy

build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/ccd
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/ccd
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/ccd

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BINARY_NAME) /usr/local/bin/

install-local: build
	@echo "Installing $(BINARY_NAME) to ~/bin..."
	@mkdir -p ~/bin
	cp $(BINARY_NAME) ~/bin/

clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f ../deploy
	rm -rf $(BUILD_DIR)

test:
	go test -v ./...

run: build
	./$(BINARY_NAME)

run-dry: build
	./$(BINARY_NAME) --dry-run

deps:
	go mod download
	go mod tidy

fmt:
	go fmt ./...

lint:
	golangci-lint run

help:
	@echo "Available targets:"
	@echo "  build        - Build the binary and move to ../deploy"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  install      - Install to /usr/local/bin (requires sudo)"
	@echo "  install-local- Install to ~/bin"
	@echo "  clean        - Remove built binaries"
	@echo "  test         - Run tests"
	@echo "  run          - Build and run"
	@echo "  run-dry      - Build and run in dry-run mode"
	@echo "  deps         - Download dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
