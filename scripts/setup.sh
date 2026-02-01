#!/bin/bash

# Setup development environment

set -e

echo "Setting up development environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

echo "Go version: $(go version)"

# Initialize Go modules if needed
if [ ! -f "go.mod" ]; then
    echo "Initializing Go modules..."
    go mod init github.com/MucheXD/WSCVenueBookingBackend
fi

# Download dependencies
echo "Downloading dependencies..."
go mod download

# Install development tools
echo "Installing development tools..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

echo "Development environment setup completed!"
echo ""
echo "Next steps:"
echo "  1. Copy configs/config.example.yml to configs/config.yml and update it"
echo "  2. Run 'make build' to build the project"
echo "  3. Run 'make test' to run tests"
echo "  4. Run 'make run' to start the application"
