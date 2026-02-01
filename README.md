# WSCVenueBookingBackend

场馆预订系统后端 - 使用Go语言实现的标准化中型后端项目。

## 项目结构 / Project Structure

```
.
├── cmd/                    # 应用程序入口 / Application entry points
│   ├── api/               # API服务器主程序 / Main API server application
│   ├── worker/            # 后台工作进程 / Background worker application
│   └── migrate/           # 数据库迁移工具 / Database migration tool
│
├── internal/              # 私有应用代码 / Private application code
│   ├── api/              # HTTP处理器和路由 / HTTP handlers and routing
│   ├── config/           # 配置管理 / Configuration management
│   ├── database/         # 数据库连接和迁移 / Database connection and migrations
│   ├── middleware/       # HTTP中间件 / HTTP middleware (auth, logging, etc.)
│   ├── models/           # 数据模型和领域实体 / Data models and domain entities
│   ├── repository/       # 数据访问层 / Data access layer
│   ├── service/          # 业务逻辑层 / Business logic layer
│   └── utils/            # 内部工具函数 / Internal utility functions
│
├── pkg/                   # 可复用的公共代码 / Public reusable code
│   ├── logger/           # 日志包 / Logging package
│   ├── validator/        # 验证工具 / Validation utilities
│   └── crypto/           # 加密工具 / Cryptography utilities
│
├── api/                   # API定义 / API definitions
│   ├── proto/            # Protocol Buffer定义 / Protocol Buffer definitions
│   └── swagger/          # OpenAPI/Swagger文档 / OpenAPI/Swagger documentation
│
├── configs/               # 配置文件 / Configuration files
│   └── config.example.yml # 配置示例 / Configuration example
│
├── scripts/               # 构建和部署脚本 / Build and deployment scripts
│
├── docs/                  # 项目文档 / Project documentation
│
├── test/                  # 测试文件 / Test files
│   ├── integration/      # 集成测试 / Integration tests
│   └── e2e/              # 端到端测试 / End-to-end tests
│
├── deployments/           # 部署配置 / Deployment configurations
│   ├── docker/           # Docker配置 / Docker configurations
│   └── kubernetes/       # Kubernetes清单 / Kubernetes manifests
│
├── go.mod                 # Go模块定义 / Go module definition
├── Makefile              # Make构建脚本 / Make build script
├── Dockerfile.example    # Docker构建示例 / Dockerfile example
└── docker-compose.yml.example  # Docker Compose示例 / Docker Compose example
```

## 目录说明 / Directory Descriptions

### cmd/
应用程序的入口点。每个子目录代表一个可执行程序。

Application entry points. Each subdirectory represents an executable program.

### internal/
私有应用程序代码。这些代码不应被其他项目导入。

Private application code. This code should not be imported by other projects.

### pkg/
可以被外部项目导入使用的公共库代码。

Public library code that can be imported by external projects.

### api/
API协议定义文件，如Protocol Buffers定义、OpenAPI规范等。

API protocol definition files, such as Protocol Buffers definitions, OpenAPI specifications, etc.

### configs/
配置文件模板和示例。

Configuration file templates and examples.

### scripts/
执行各种构建、安装、分析等操作的脚本。

Scripts for performing various build, install, analysis operations.

### test/
额外的外部测试应用程序和测试数据。

Additional external test applications and test data.

### deployments/
系统和容器编排部署配置和模板。

System and container orchestration deployment configurations and templates.

## 开发指南 / Development Guide

### 构建项目 / Build Project
```bash
make build
```

### 运行应用 / Run Application
```bash
make run
```

### 运行测试 / Run Tests
```bash
make test
```

### 代码检查 / Lint Code
```bash
make lint
```

### 清理构建产物 / Clean Build Artifacts
```bash
make clean
```

## 技术栈 / Tech Stack

- **语言 / Language**: Go 1.21+
- **数据库 / Database**: PostgreSQL (推荐 / recommended)
- **容器化 / Containerization**: Docker
- **编排 / Orchestration**: Kubernetes (可选 / optional)

## 贡献指南 / Contributing

欢迎贡献！请先阅读 [CONTRIBUTING.md](./docs/CONTRIBUTING.md) 了解详细信息。

Contributions are welcome! Please read [CONTRIBUTING.md](./docs/CONTRIBUTING.md) for details.

## 许可证 / License

[MIT License](LICENSE)