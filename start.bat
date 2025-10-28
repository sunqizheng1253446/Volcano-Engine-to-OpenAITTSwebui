@echo off

echo Starting Volcano Engine TTS Health Monitor...

REM 确保依赖已安装
echo Installing dependencies...
go mod tidy

REM 构建并运行项目
echo Building project...
go build -o tts-monitor.exe .\cmd\main.go

REM 检查构建是否成功
if %errorlevel% equ 0 (
    echo Build successful!
    echo Starting server...
    tts-monitor.exe
) else (
    echo Build failed!
    pause
    exit /b 1
)