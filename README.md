# WSCVenueBookingBackend

场馆预订系统后端 - 使用Go语言实现的标准化中型后端项目。

## 项目结构 / Project Structure

```
.
├── cmd/                    # 应用程序入口 / Application entry points
│   ├── main/              # 主服务器程序 / Main server application
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
├── configs/               # 配置文件 / Configuration files
│
├── docs/                  # 项目文档 / Project documentation
│
├── test/                  # 测试文件 / Test files
│   └── integration/      # 集成测试 / Integration tests
│
├── deployments/           # 部署配置 / Deployment configurations
|   ├── scripts/           # 构建和部署脚本 / Build and deployment scripts
│   └── docker/           # Docker配置 / Docker configurations
│
├── go.mod                 # Go模块定义 / Go module definition
├── Makefile              # Make构建脚本 / Make build script
└── docker-compose.yml.example  # Docker Compose示例 / Docker Compose example
```

## 开发指南 / Development Guide

### 环境要求 / Prerequisites
- Go 1.21+
- Make (可选 / optional)
- Docker (可选 / optional)

### 构建项目 / Build Project
参考 Makefile 中的构建命令注释。

Refer to build command comments in Makefile.

### 运行应用 / Run Application
参考 Makefile 中的运行命令注释。

Refer to run command comments in Makefile.

### 运行测试 / Run Tests
参考 Makefile 中的测试命令注释。

Refer to test command comments in Makefile.

### 代码检查 / Lint Code
待定

## 技术栈 / Tech Stack

- **语言 / Language**: Go 1.21+
- **数据库 / Database**: PostgreSQL (推荐 / 待定)
- **容器化 / Containerization**: Docker

## 许可证 / License

待定