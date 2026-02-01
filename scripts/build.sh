#!/bin/bash

# Build script for the project

set -e

echo "Building WSCVenueBookingBackend..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Build API server
echo "Building API server..."
go build -o bin/api ./cmd/api

# Build worker
echo "Building worker..."
go build -o bin/worker ./cmd/worker

# Build migrate tool
echo "Building migrate tool..."
go build -o bin/migrate ./cmd/migrate

echo "Build completed successfully!"
echo "Binaries are available in ./bin/"
