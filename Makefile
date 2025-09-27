# Makefile for Git Server

.PHONY: build run test clean

# Build the server
build:
	go build -o ethglobal main.go

# Run the server
run: build
	./ethglobal

# Run tests
test:
	./test.sh

# Test Lighthouse functionality
test-lighthouse:
	./test-lighthouse.sh

# Clean build artifacts
clean:
	rm -f ethglobal

# Install dependencies
deps:
	go mod tidy

# Run with custom port
run-port:
	@read -p "Enter port (default 8080): " port; \
	PORT=$${port:-8080} ./ethglobal

# Help
help:
	@echo "Available targets:"
	@echo "  build     - Build the Git server"
	@echo "  run       - Build and run the server"
	@echo "  test      - Run test script"
	@echo "  clean     - Remove build artifacts"
	@echo "  deps      - Install dependencies"
	@echo "  run-port  - Run server on custom port"
	@echo "  help      - Show this help"
