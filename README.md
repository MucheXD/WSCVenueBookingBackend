> [!NOTE]
> **此项目是一个练习项目，不适用于生产。**

> [!NOTE]
> **此项目已经结束，将不再进行后续开发与维护。**

# WSCVenueBookingBackend

使用 Go 编写的“场馆预订系统”后端。

与本项目关连的前端仓库为 [WSCVenueBookingFrontend](https://github.com/Promise2679/WSCVenueBookingFrontend) 

## 项目结构 / Project Structure

```
.
├───.github # Action 自动部署工作流
├───cmd # 可运行程序与功能入口
│   ├───batchRegister # 提供批量注册功能
│   ├───deploy # 项目自动化部署相关
│   └───main # 主进程启动位置
├───configs # 外部配置
│   ├───secret # 服务器配置 (gitignore)
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
    │   ├───applicationCtrl # 申请单请求控制器
    │   ├───fileCtrl # 文件请求控制器
    │   ├───notificationCtrl # 公告与站内信请求控制器
    │   ├───userCtrl # 用户与账号请求控制器
    │   └───venueCtrl # 场地请求控制器
    ├───middlewares # 中间件
    ├───models # 业务模型
    ├───repository # 存储层实现
    ├───service # 业务层实现
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

## 技术栈 / Tech Stack

- **语言**: Go 1.25.1
- **主要依赖**: Gin, Gorm
- **数据库**: MySQL
- **容器化**: Docker
