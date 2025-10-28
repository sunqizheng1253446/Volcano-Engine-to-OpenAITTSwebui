package api

import (
	"github.com/gin-gonic/gin"
	"volcano-engine-tts-monitor/cmd/handlers"
	"volcano-engine-tts-monitor/cmd/models"
	"volcano-engine-tts-monitor/config"
)

// SetupRouter 设置路由
func SetupRouter(cfg *config.Config, monitor *models.Monitor) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(cfg.GinMode)

	router := gin.Default()

	// 静态文件服务
	router.Static("/static", "./static")

	// 创建处理程序
	healthHandler := handlers.NewHealthHandler(monitor, cfg)
	websocketHandler := handlers.NewWebSocketHandler(monitor, cfg)

	// 健康检查相关路由
	router.GET("/health", healthHandler.HealthCheck)
	router.GET("/metrics", healthHandler.GetMetrics)
	router.GET("/errors", healthHandler.GetErrorRecords)

	// TTS WebSocket路由
	router.GET("/tts/websocket", websocketHandler.HandleWebSocket)

	// 监控页面路由
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	return router
}