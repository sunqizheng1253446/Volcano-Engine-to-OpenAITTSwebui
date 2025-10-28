# Volcano Engine TTS 健康监测系统

这是一个通用的健康监测系统，用于监控其他服务的运行状态、性能指标和错误信息。系统采用前后端分离架构，后端使用 Go 语言开发，前端页面独立存放，便于热更新。

## 功能特性

- **实时健康监控**：监控目标服务的运行状态、连接数、并发调用等关键指标
- **性能分析**：展示响应时间、成功率、资源使用等性能指标
- **可视化展示**：通过图表直观展示性能趋势和系统状态
- **错误记录**：记录并展示系统错误，便于问题排查
- **配置管理**：通过环境变量管理服务配置，支持自定义参数
- **热更新支持**：前端页面与后端分离，支持独立更新

## 项目结构

```
Volcano-Engine-to-OpenAITTSwebui/
├── cmd/                     # 应用程序入口
│   ├── main.go              # 主程序入口
│   ├── api/                 # API 路由定义
│   │   └── router.go        # 路由设置
│   ├── handlers/            # 请求处理器
│   │   ├── health_handler.go  # 健康检查处理器
│   │   └── websocket_handler.go # WebSocket 处理器
│   ├── models/              # 数据模型
│   │   └── monitor.go       # 监控相关模型
│   └── internal/            # 内部包
├── config/                  # 配置相关
│   └── config.go            # 配置加载和管理
├── static/                  # 静态资源（前端文件）
│   └── index.html           # 健康监测页面
├── templates/               # 模板文件（预留）
├── .env                     # 环境变量配置
├── .env.example             # 环境变量示例
├── go.mod                   # Go 模块定义
├── Dockerfile               # Docker构建文件
├── .dockerignore            # Docker忽略文件
└── README.md                # 项目说明
```

## 安装与使用

### 前置条件

- Go 1.20 或更高版本
- 要监控的目标服务API
- Docker（可选，用于容器化部署）

### 快速开始

1. **配置环境变量**

   复制 `.env.example` 文件为 `.env` 并填写相关配置：

   ```bash
   cp .env.example .env
   # 编辑 .env 文件，填写配置信息
   ```

2. **启动服务**

   方式一：直接运行（需要安装Go）：
   ```bash
   go run cmd/main.go
   ```

   方式二：编译后运行：
   ```bash
   go build -o tts-monitor cmd/main.go
   ./tts-monitor
   ```

   方式三：使用Docker运行：
   ```bash
   docker build -t health-monitor .
   docker run -p 8080:8080 --env-file .env health-monitor
   ```

3. **访问监控页面**

   打开浏览器，访问 `http://localhost:8080` 查看健康监测页面。

## API 端点

- **`GET /health`** - 获取目标服务的健康状态信息
- **`GET /metrics`** - 获取目标服务的性能指标
- **`GET /errors`** - 获取监控系统自身的错误记录
- **`GET /`** - 访问健康监测页面

## 环境变量配置

| 环境变量 | 说明 | 默认值 | 必需 |
|---------|------|-------|------|
| LISTEN_ADDR | 监控服务监听地址和端口 | :8080 | 否 |
| MAX_CONNECTIONS | 监控系统最大并发连接数 | 100 | 否 |
| MAX_CONCURRENT_CALLS | 监控系统最大并发调用数 | 10 | 否 |
| CHECK_INTERVAL | 检查间隔（秒） | 5 | 否 |
| TARGET_API_URL | 被监控服务的健康检查API地址 | http://localhost:8081/api/health | 是 |
| LOG_LEVEL | 日志级别 | info | 否 |
| GIN_MODE | Gin 框架模式 | release | 否 |

## 被监控服务API数据格式要求

被监控的服务需要提供一个健康检查API端点，返回以下JSON格式的数据：

```json
{
  "status": "ok",
  "uptime_seconds": 3600,
  "active_connections": 10,
  "current_calls": 5,
  "avg_response_time": 120,
  "success_rate": 99.5,
  "error_count": 2,
  "cpu_usage": 45.5,
  "memory_usage": 60.2,
  "request_count": 12345,
  "version": "1.2.3"
}
```

### 字段说明

| 字段名 | 类型 | 说明 |
|--------|------|------|
| status | string | 服务状态 ("ok", "warning", "error") |
| uptime_seconds | int64 | 服务运行时间（秒） |
| active_connections | int | 活动连接数 |
| current_calls | int | 当前并发调用数 |
| avg_response_time | int | 平均响应时间（毫秒） |
| success_rate | float64 | 成功率（百分比） |
| error_count | int | 错误数量 |
| cpu_usage | float64 | CPU使用率（百分比） |
| memory_usage | float64 | 内存使用率（百分比） |
| request_count | int | 总请求数 |
| version | string | 服务版本号 |

## Docker部署

### 构建镜像
```bash
docker build -t health-monitor .
```

### 运行容器
```bash
docker run -d -p 8080:8080 --name health-monitor health-monitor
```

### 使用环境变量文件运行
```bash
docker run -d -p 8080:8080 --env-file .env --name health-monitor health-monitor
```

### 查看日志
```bash
docker logs health-monitor
```

## 前端页面热更新

由于前端页面（`static/index.html`）与后端代码完全分离，您可以随时更新页面内容而无需重启后端服务：

1. 修改 `static/index.html` 文件
2. 刷新浏览器页面即可看到更新后的内容

## 自定义监控指标

您可以通过修改 `cmd/models/monitor.go` 和 `cmd/handlers/health_handler.go` 文件来添加或修改监控指标，以满足特定需求。

## 许可证

MIT License