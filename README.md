# Volcano Engine TTS 健康监测系统

这是一个为 Volcano Engine TTS 服务设计的健康监测系统，提供实时监控、性能分析和可视化展示功能。系统采用前后端分离架构，后端使用 Go 语言开发，前端页面独立存放，便于热更新。

## 功能特性

- **实时健康监控**：监控 TTS 服务的运行状态、连接数、并发调用等关键指标
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
├── start.sh                 # Linux/Mac 启动脚本
├── start.bat                # Windows 启动脚本
└── README.md                # 项目说明
```

## 安装与使用

### 前置条件

- Go 1.20 或更高版本
- 火山引擎 TTS 服务的访问凭证（可选，用于实际运行）

### 快速开始

1. **配置环境变量**

   复制 `.env.example` 文件为 `.env` 并填写相关配置：

   ```bash
   cp .env.example .env
   # 编辑 .env 文件，填写配置信息
   ```

2. **启动服务**

   - 在 Linux/Mac 环境：
     ```bash
     chmod +x start.sh
     ./start.sh
     ```

   - 在 Windows 环境：
     ```cmd
     start.bat
     ```

3. **访问监控页面**

   打开浏览器，访问 `http://localhost:8080` 查看健康监测页面。

## API 端点

- **`GET /health`** - 获取健康状态信息
- **`GET /metrics`** - 获取性能指标
- **`GET /errors`** - 获取错误记录
- **`GET /tts/websocket`** - TTS WebSocket 服务端点
- **`GET /`** - 访问健康监测页面

## 环境变量配置

| 环境变量 | 说明 | 默认值 | 必需 |
|---------|------|-------|------|
| LISTEN_ADDR | 服务监听地址和端口 | :8080 | 否 |
| MAX_CONNECTIONS | 最大并发连接数 | 100 | 否 |
| MAX_CONCURRENT_CALLS | 最大并发调用数 | 10 | 否 |
| CHECK_INTERVAL | 检查间隔（秒） | 5 | 否 |
| BYTEDANCE_TTS_APP_ID | 火山引擎 App ID | - | 否（实际使用时需要） |
| BYTEDANCE_TTS_BEARER_TOKEN | 火山引擎认证令牌 | - | 否（实际使用时需要） |
| BYTEDANCE_TTS_CLUSTER | 火山引擎集群名称 | - | 否（实际使用时需要） |
| BYTEDANCE_TTS_VOICE_TYPE | 火山引擎语音类型 | - | 否（实际使用时需要） |
| OPENAI_TTS_API_KEY | OpenAI TTS API 密钥 | - | 否 |
| LOG_LEVEL | 日志级别 | info | 否 |
| GIN_MODE | Gin 框架模式 | release | 否 |

## 前端页面热更新

由于前端页面（`static/index.html`）与后端代码完全分离，您可以随时更新页面内容而无需重启后端服务：

1. 修改 `static/index.html` 文件
2. 刷新浏览器页面即可看到更新后的内容

## 自定义监控指标

您可以通过修改 `cmd/models/monitor.go` 和 `cmd/handlers/health_handler.go` 文件来添加或修改监控指标，以满足特定需求。

## 注意事项

- 本项目中的 WebSocket 处理部分是模拟实现，在实际使用时需要集成火山引擎的 TTS WebSocket API
- 环境变量中的敏感信息（如 API 密钥）应妥善保管，避免泄露
- 在生产环境中，建议配置适当的连接和调用限制，以防止资源耗尽

## 许可证

MIT License