# API数据项说明

本文档详细列出了Volcano Engine TTS健康监测系统所需的API数据项，用于与外部系统对接。

## 1. 健康状态接口 (/health)

### 响应数据项

| 字段名 | 类型 | 说明 |
|--------|------|------|
| status | string | 服务状态 ("ok", "warning", "error") |
| uptime_seconds | int64 | 服务运行时间（秒） |
| active_connections | int | 活动连接数 |
| current_calls | int | 当前并发调用数 |
| avg_response_time | int | 平均响应时间（毫秒） |
| success_rate | float64 | 成功率 (%) |
| error_count | int | 错误数量 |
| cpu_usage | float64 | CPU使用率 (%) |
| memory_usage | float64 | 内存使用率 (%) |
| request_count | int | 总请求数 |
| version | string | 服务版本号 |
| last_check_time | string | 最后检查时间 (RFC3339格式) |

## 2. 错误记录接口 (/errors)

### 响应数据项

| 字段名 | 类型 | 说明 |
|--------|------|------|
| error_records | array | 错误记录数组 |
| count | int | 错误记录总数 |

### 错误记录项 (error_records数组元素)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| timestamp | string | 错误时间戳 (RFC3339格式) |
| error_type | string | 错误类型 |
| message | string | 错误消息 |

## 3. 性能指标接口 (/metrics)

### 响应数据项

| 字段名 | 类型 | 说明 |
|--------|------|------|
| avg_response_time | int | 平均响应时间（毫秒） |
| success_rate | float64 | 成功率 (%) |
| uptime_seconds | int64 | 服务运行时间（秒） |
| active_connections | int | 活动连接数 |
| current_calls | int | 当前并发调用数 |
| cpu_usage | float64 | CPU使用率 (%) |
| memory_usage | float64 | 内存使用率 (%) |
| request_count | int | 总请求数 |

## 4. WebSocket接口 (/api/websocket)

### 说明
用于实时数据推送的WebSocket连接

## 5. 环境变量配置项

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| LISTEN_ADDR | :8080 | 监控服务监听地址和端口 |
| MAX_CONNECTIONS | 100 | 监控系统最大并发连接数 |
| MAX_CONCURRENT_CALLS | 10 | 监控系统最大并发调用数 |
| CHECK_INTERVAL | 5 | 检查间隔（秒） |
| TARGET_API_URL | http://localhost:8081/api/health | 被监控服务的健康检查API地址 |
| LOG_LEVEL | info | 日志级别 |
| GIN_MODE | release | Gin 框架模式 |