#!/bin/bash

# 启动脚本 - 用于Linux/Mac环境
echo "Starting Volcano Engine TTS Health Monitor..."

# 确保依赖已安装
go mod tidy

# 构建并运行项目
go build -o tts-monitor ./cmd/main.go

# 检查构建是否成功
if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Starting server..."
    ./tts-monitor
else
    echo "Build failed!"
    exit 1
fi