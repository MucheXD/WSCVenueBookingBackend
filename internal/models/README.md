# Models

Data models and domain entities.

通用数据模型，专用模型与模块翻译中间层请置于模块内。

注意：此处定义的是 `Domain Model`，因此不应该存在任何反射标签。

- 含有 `gorm` 标签的为 `Entity Model`，请定义在 `repository` 中；

- 含有 `json` 标签的为 `DTO Layer`，请定义在 `api` 中；