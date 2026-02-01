# Makefile for WSCVenueBookingBackend
# WSCVenueBookingBackend 构建文件

.PHONY: help build run test clean lint

# help: 显示帮助信息
# Display this help screen
help:
	# 使用 grep 和 awk 解析 Makefile 中的注释，生成帮助信息

# build: 构建应用程序
# Build the application
build:
	# 编译主服务器: go build -o bin/main ./cmd/main
	# 编译worker服务: go build -o bin/worker ./cmd/worker
	# 编译migrate工具: go build -o bin/migrate ./cmd/migrate

# run: 运行应用程序
# Run the application
run:
	# 运行主服务器: go run ./cmd/main

# test: 运行测试
# Run tests
test:
	# 运行所有测试: go test -v ./...

# test-coverage: 运行测试并生成覆盖率报告
# Run tests with coverage
test-coverage:
	# 生成覆盖率报告: go test -v -coverprofile=coverage.out ./...
	# 生成HTML报告: go tool cover -html=coverage.out -o coverage.html

# lint: 运行代码检查
# Run linter
lint:
	# 运行 golangci-lint: golangci-lint run

# clean: 清理构建产物
# Clean build artifacts
clean:
	# 删除 bin/ 目录和覆盖率报告文件

# install-tools: 安装开发工具
# Install development tools
install-tools:
	# 安装 golangci-lint: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.DEFAULT_GOAL := help
