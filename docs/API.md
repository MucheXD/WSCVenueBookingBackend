# API Documentation

## Overview

本文档描述了场馆预订系统的API接口。

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

大多数API端点需要JWT令牌认证。

在请求头中包含：
```
Authorization: Bearer <token>
```

## Endpoints

### Health Check

```
GET /health
```

检查服务健康状态。

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Error Responses

所有错误响应遵循统一格式：

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message",
    "details": {}
  }
}
```

## Rate Limiting

API限流规则：
- 未认证用户: 100 请求/小时
- 已认证用户: 1000 请求/小时
