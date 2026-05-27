# Gin Basic 项目脚手架

这是一个基于 Gin 框架的 Go Web 项目脚手架，包含了常用的项目结构和功能模块。

## 项目结构

```
gin-basic/
├── main.go               # 主入口文件
├── Makefile              # Make命令集合
├── Dockerfile            # Docker配置
├── go.mod
├── go.sum
├── config/               # 配置文件
│   ├── local.yaml        # 本地开发配置
│   ├── dev.yaml          # 开发环境配置
│   └── production.yaml   # 生产环境配置
├── handler/              # 控制器层
│   └── user_handler.go   # 用户相关控制器
├── middleware/           # 中间件
│   └── locales_middleware.go # 国际化中间件
├── models/               # 数据模型
│   └── user.go           # 用户模型
├── response/             # 响应处理
│   ├── code.go           # 响应码定义
│   └── response.go       # 响应格式
├── router/               # 路由配置
│   ├── router.go         # 主路由
│   └── user_router.go    # 用户路由
├── service/              # 业务逻辑层
│   └── user_service.go   # 用户服务
├── settings/             # 配置处理
│   └── setting.go        # 配置加载
├── logger/               # 日志处理
│   └── logger.go         # 日志初始化
├── internal/             # 内部工具
│   └── local.go          # 国际化工具
├── locales/              # 国际化文件
│   ├── active.en.toml    # 英文翻译
│   └── active.zh.toml    # 中文翻译
└── README.md
```

## 功能特性

- Gin Web 框架
- 配置管理（支持多环境）
- 日志记录（支持文件轮转）
- 国际化支持
- 统一响应格式
- 用户注册/登录功能
- Docker 部署支持
- 结构化日志记录

## 快速开始

### 本地开发

1. 安装依赖：
```bash
make deps
```

2. 运行项目：
```bash
make run
```

或者直接使用 Go：
```bash
go run main.go
```

### 使用 Make 命令

```bash
# 构建项目
make build

# 运行项目
make run

# 清理构建文件
make clean

# 运行测试
make test

# 构建 Docker 镜像
make docker-build

# 安装依赖
make deps
```

### Docker 部署

1. 构建镜像：
```bash
make docker-build
```

2. 运行容器：
```bash
make docker-run
```

或者直接使用 Docker 命令：
```bash
docker build -t gin-basic .
docker run -p 8080:8080 gin-basic
```

## API 接口

- `POST /v1/user/register` - 用户注册
- `POST /v1/user/login` - 用户登录
- `GET /v1/user/info` - 获取用户信息
- `GET /health` - 健康检查

## 环境变量

- `ENVIRONMENT` - 设置运行环境（local/dev/production），默认为 local

## 依赖库

- [Gin](https://github.com/gin-gonic/gin) - Web 框架
- [Zap](https://github.com/uber-go/zap) - 日志库
- [Lumberjack](https://github.com/natefinch/lumberjack) - 日志轮转
- [YAML](https://github.com/go-yaml/yaml) - 配置文件解析
- [go-i18n](https://github.com/nicksnyder/go-i18n) - 国际化支持
- [Toml](https://github.com/BurntSushi/toml) - TOML 解析
- [Golang.org/x/text](https://pkg.go.dev/golang.org/x/text) - 文本处理
- [Golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) - 加密算法

## 项目特点

### 1. 分层架构
- Handler 层：处理 HTTP 请求和响应
- Service 层：封装业务逻辑
- Model 层：定义数据结构
- Router 层：处理路由映射

### 2. 统一日志
- 使用 Zap 进行结构化日志记录
- 支持日志轮转，防止日志文件无限增长
- 在关键节点记录日志，便于问题排查

### 3. 国际化支持
- 支持多语言（中英文）
- 通过中间件自动处理语言选择
- 可轻松扩展支持更多语言

### 4. 统一响应格式
- 预定义的标准响应结构
- 易于前端处理和解析
- 包含错误码和消息定义

## 扩展指南

要添加新功能，请按照以下步骤：

1. 在 `models` 目录下创建数据模型
2. 在 `service` 目录下创建业务逻辑
3. 在 `handler` 目录下创建处理器
4. 在 `router` 目录下注册路由
5. 在 `response` 中添加响应码（如需要）

## 部署建议

- 在生产环境中使用 Docker 部署
- 配置外部日志收集系统
- 使用反向代理（如 Nginx）处理 SSL 终止
- 设置适当的资源限制和健康检查