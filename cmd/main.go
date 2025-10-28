package main

import (
	"fmt"
	"log"

	"volcano-engine-tts-monitor/cmd/api"
	"volcano-engine-tts-monitor/cmd/models"
	"volcano-engine-tts-monitor/config"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建监控器
	monitor := models.NewMonitor()

	// 设置路由
	router := api.SetupRouter(cfg, monitor)

	// 启动服务器
	fmt.Printf("Server starting on %s\n", cfg.ListenAddr)
	if err := router.Run(cfg.ListenAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}