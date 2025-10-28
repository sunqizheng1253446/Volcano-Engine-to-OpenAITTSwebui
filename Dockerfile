# 使用官方Go镜像作为构建环境
FROM golang:1.20-alpine AS builder

# 安装git（用于go mod download）
RUN apk add --no-cache git

# 创建源代码目录
RUN mkdir /src
WORKDIR /src

# 复制go mod和sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o ./bin/server ./cmd/main.go

# 使用alpine作为运行环境
FROM alpine:latest

# 安装ca证书（用于HTTPS请求）
RUN apk --no-cache add ca-certificates

# 创建非root用户
RUN adduser -D -s /bin/sh appuser

# 设置工作目录
WORKDIR /app

# 从builder阶段复制构建好的二进制文件
COPY --from=builder /src/bin/server .

# 复制静态文件
COPY --from=builder /src/static ./static

# 复制环境变量文件
COPY --from=builder /src/.env.example .env

# 更改文件所有者
RUN chown -R appuser:appuser /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./server"]