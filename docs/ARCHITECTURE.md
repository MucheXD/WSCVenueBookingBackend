# Architecture

## 系统架构 / System Architecture

本项目采用分层架构设计，遵循Clean Architecture原则。

### 架构图

```
┌─────────────────────────────────────────────────┐
│                   Presentation Layer            │
│              (cmd/main, internal/api)           │
└─────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────┐
│                   Business Logic Layer          │
│               (internal/service)                │
└─────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────┐
│                   Data Access Layer             │
│              (internal/repository)              │
└─────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────┐
│                   Database                      │
└─────────────────────────────────────────────────┘
```

## 层次说明

### Presentation Layer (表现层)
- 对应模型：定义在 Controller 中的 DTO 模型，可有 `json` 反射字
- 处理HTTP请求和响应
- 参数验证
- 认证和授权
- 中间件处理

### Business Logic Layer (业务逻辑层)
- 亦称：Service / Domain
- 对应模型：定义在 Models 中的业务模型，一般无反射字
- 核心业务规则
- 业务流程编排
- 领域模型操作

### Data Access Layer (数据访问层)
- 亦称：Repository
- 对应模型：定义在 Repository 中的 Entity 模型，可有 `gorm` 反射字
- 数据库操作
- 数据持久化
- 查询优化

### 模型转换

Repository 依赖于 Domain，负责 Models 和 Entity 的互相转换

Presentation 依赖于 Domain，由 Controller 负责 DTO 向 Models 的单向转换

## 错误处理

表现层的统一 API 错误与业务层错误分离，在 Controller 中实现转换。

表现层错误：统一定义在 utils-apiException 中，通过 errorHandler 中间件实现统一错误处理与转换成 JSON 格式返回。

业务层错误：哨兵错误统一定义在对应业务包中的 exceptions.go 中，以便外部通过 errors.Is() 进行错误类型判断。如果需要传递来自数据访问层等的错误，可以使用多重错误包装的方式向外传递。

业务层错误到表现层错误为多对一关系，避免过多的信息泄露给外部。原始业务层错误通过 apiException.AbortWithException 方法的第三个参数传递给 errorHandler，通过日志打印。

数据访问层为方便起见，可直接提出原始未包装错误。