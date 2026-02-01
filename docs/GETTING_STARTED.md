# Getting Started

欢迎使用 WSCVenueBookingBackend！

## 快速开始 / Quick Start

### 1. 环境要求 / Prerequisites

- Go 1.21 或更高版本 / Go 1.21 or higher
- Make (可选但推荐 / Optional but recommended)
- Docker (可选 / Optional)

### 2. 初始化项目 / Initialize Project

```bash
# 克隆项目
git clone https://github.com/MucheXD/WSCVenueBookingBackend.git
cd WSCVenueBookingBackend

# 运行设置脚本
./scripts/setup.sh
```

### 3. 配置 / Configuration

```bash
# 复制配置示例
cp configs/config.example.yml configs/config.yml

# 编辑配置文件
vim configs/config.yml
```

### 4. 构建项目 / Build

```bash
make build
```

或者使用构建脚本 / Or use build script:
```bash
./scripts/build.sh
```

### 5. 运行应用 / Run Application

```bash
# 运行 API 服务器
make run

# 或直接运行二进制文件
./bin/api
```

### 6. 运行测试 / Run Tests

```bash
make test
```

### 7. 代码检查 / Linting

```bash
make lint
```

## 项目结构说明 / Project Structure Guide

### cmd/
存放应用程序的主入口点。每个子目录应该包含一个 `main.go` 文件。

Place application entry points here. Each subdirectory should contain a `main.go` file.

### internal/
存放私有应用程序代码。这些代码不应被外部项目导入。

Place private application code here. This code should not be imported by external projects.

- **api/**: HTTP 路由和处理器
- **config/**: 配置加载和管理
- **database/**: 数据库连接和迁移
- **middleware/**: HTTP 中间件
- **models/**: 数据模型
- **repository/**: 数据访问层
- **service/**: 业务逻辑
- **utils/**: 工具函数

### pkg/
存放可以被外部项目导入的公共库代码。

Place public library code that can be imported by external projects.

### api/
存放 API 规范和文档。

Place API specifications and documentation.

### configs/
存放配置文件。

Place configuration files.

### scripts/
存放构建、部署和其他脚本。

Place build, deployment, and other scripts.

### docs/
存放项目文档。

Place project documentation.

### test/
存放测试代码。

Place test code.

### deployments/
存放部署相关的配置。

Place deployment-related configurations.

## 下一步 / Next Steps

1. 在 `cmd/` 目录下实现应用程序入口
2. 在 `internal/models/` 定义数据模型
3. 在 `internal/repository/` 实现数据访问层
4. 在 `internal/service/` 实现业务逻辑
5. 在 `internal/api/` 实现 HTTP 处理器
6. 添加测试用例

## 需要帮助？ / Need Help?

查看以下文档：

- [API Documentation](./API.md)
- [Architecture](./ARCHITECTURE.md)
- [Contributing Guide](./CONTRIBUTING.md)
