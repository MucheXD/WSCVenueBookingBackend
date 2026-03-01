# WSCVenueBookingBackend

场馆预订系统后端

## 项目结构 / Project Structure

```
.
├───.github # Action 自动部署工作流
├───cmd # 可运行程序与功能入口
│   ├───batchRegister # 提供批量注册功能
│   ├───deploy # 项目自动化部署相关
│   ├───main # 主进程启动位置
│   ├───migrate # SQL 自动建表相关
│   └───worker # 后台工作进程
├───configs # 外部配置
│   ├───secret # 服务器配置
│   └───sql # 数据库配置
├───docs # 文档
│   ├───prompt # 智能体辅助开发文档
│   └───svc # 业务开发文档
└───internal # 内部业务实现
    ├───config # 内部模块配置
    │   ├───database # 数据库管理
    │   ├───logger # 日志管理
    │   └───server # 服务器、路由管理
    ├───controllers # 请求控制器实现
    │   ├───accessCtrl # 场地权限请求控制器
    │   ├───applicationCtrl # 申请单请求控制器
    │   ├───fileCtrl # 文件请求控制器
    │   ├───notificationCtrl # 公告与站内信请求控制器
    │   ├───userCtrl # 用户与账号请求控制器
    │   └───venueCtrl # 场地请求控制器
    ├───middlewares # 中间件
    ├───models # 业务模型
    ├───repository # 存储层实现
    ├───service # 业务层实现
    │   ├───accessSvc # 场地权限业务
    │   ├───applicationSvc # 申请单业务
    │   ├───fileStorage # 文件仓储业务
    │   ├───fileSvc # 文件出入库业务
    │   ├───notificationSvc # 公告与站内信业务
    │   ├───roleSvc # 场地角色业务
    │   ├───userSvc # 用户与账号业务
    │   └───venueSvc # 场地业务
    └───utils # 通用组件
        ├───apiException # 业务异常处理组件
        ├───systemPermission # 系统权限组件
        ├───venuePermission # 场地权限组件
        └───webtoken # 用户令牌组件
```

## 开发指南 / Development Guide

### 环境要求 / Prerequisites
- Go 1.24+
- Docker (可选 / optional)

### 构建项目 / Build Project
参考 Dockerfile 中的构建命令注释。

Refer to build command comments in Makefile.

### 运行应用 / Run Application
参考 Dockerfile 中的运行命令注释。

Refer to run command comments in Makefile.

## 技术栈 / Tech Stack

- **语言 / Language**: Go 1.24+
- **数据库 / Database**: MySQL
- **容器化 / Containerization**: Docker