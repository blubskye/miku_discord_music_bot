.PHONY: build run clean install deps test

BINARY_NAME=miku_bot
BUILD_DIR=bin
MAIN_PATH=./cmd/bot

build:
	@echo "Building Miku Bot..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux"

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows.exe $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)-windows.exe"

build-mac:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-macos $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)-macos"

build-all: build-linux build-windows build-mac
	@echo "All builds complete!"

run: build
	@echo "Starting Miku Bot..."
	./$(BUILD_DIR)/$(BINARY_NAME)

dev:
	@echo "Running in development mode..."
	go run $(MAIN_PATH)

clean:
	@echo "Cleaning build files..."
	@rm -rf $(BUILD_DIR)
	@rm -f *.db *.db-shm *.db-wal
	@echo "Clean complete!"

deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed!"

test:
	@echo "Running tests..."
	go test -v ./...

install: build
	@echo "Installing Miku Bot..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete! Run with: $(BINARY_NAME)"

help:
	@echo "Miku Bot Makefile Commands:"
	@echo "  make build         - Build the bot for current platform"
	@echo "  make build-linux   - Build for Linux (amd64)"
	@echo "  make build-windows - Build for Windows (amd64)"
	@echo "  make build-mac     - Build for macOS (arm64)"
	@echo "  make build-all     - Build for all platforms"
	@echo "  make run           - Build and run the bot"
	@echo "  make dev           - Run in development mode (no build)"
	@echo "  make clean         - Remove build files and databases"
	@echo "  make deps          - Install/update dependencies"
	@echo "  make test          - Run tests"
	@echo "  make install       - Install to /usr/local/bin"
	@echo "  make help          - Show this help message"
