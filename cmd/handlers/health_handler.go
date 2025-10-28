package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"volcano-engine-tts-monitor/cmd/models"
	"volcano-engine-tts-monitor/config"
)

type HealthHandler struct {
	monitor *models.Monitor
	config  *config.Config
}

func NewHealthHandler(monitor *models.Monitor, config *config.Config) *HealthHandler {
	return &HealthHandler{
		monitor: monitor,
		config:  config,
	}
}

// HealthCheck 处理健康检查请求
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// 模拟资源使用情况（在实际应用中，应该使用系统库获取真实资源使用情况）
	resourceUsage := models.ResourceUsage{
		CPUUsage:    35.5, // 模拟CPU使用率
		MemoryUsage: 42.3, // 模拟内存使用率
		MemoryMB:    512,  // 模拟内存使用量（MB）
		NetworkIn:   1024, // 模拟网络入站流量（KB）
		NetworkOut:  2048, // 模拟网络出站流量（KB）
	}

	// 连接统计信息
	connectionStats := models.ConnectionStats{
		TodayPeak:         h.monitor.GetActiveConnections(),
		ConnectionSuccess: 1234,
		ConnectionFailure: 23,
		AvgConnectionTime: int(h.monitor.GetAvgResponseTime() / 1000), // 转换为秒
	}

	// 服务配置信息（注意：不返回敏感信息如token）
	serviceConfig := models.ServiceConfig{
		ByteDanceAppID:     maskString(h.config.ByteDanceAppID),
		ByteDanceCluster:   h.config.ByteDanceCluster,
		ByteDanceVoiceType: h.config.ByteDanceVoiceType,
		ListenAddr:         h.config.ListenAddr,
		MaxConnections:     h.config.MaxConnections,
		MaxConcurrentCalls: h.config.MaxConcurrentCalls,
	}

	// 确定服务状态
	status := "ok"
	if h.monitor.GetActiveConnections() > int(float64(h.config.MaxConnections)*0.8) ||
		h.monitor.GetCurrentCalls() > int(float64(h.config.MaxConcurrentCalls)*0.8) {
		status = "warning"
	}
	if h.monitor.GetSuccessRate() < 95.0 {
		status = "error"
	}

	// 构建健康状态响应
	healthStatus := models.HealthStatus{
		Status:             status,
		ActiveConnections:  h.monitor.GetActiveConnections(),
		MaxConnections:     h.config.MaxConnections,
		CurrentCalls:       h.monitor.GetCurrentCalls(),
		MaxConcurrentCalls: h.config.MaxConcurrentCalls,
		UptimeSeconds:      h.monitor.GetUptimeSeconds(),
		LastCheckTime:      time.Now().Format(time.RFC3339),
		AvgResponseTime:    h.monitor.GetAvgResponseTime(),
		TodayRequestCount:  4567, // 模拟值
		SuccessRate:        h.monitor.GetSuccessRate(),
		ErrorCount:         23,   // 模拟值
		Config:             serviceConfig,
		ResourceUsage:      resourceUsage,
		ConnectionStats:    connectionStats,
	}

	c.JSON(http.StatusOK, healthStatus)
}

// GetErrorRecords 获取错误记录
func (h *HealthHandler) GetErrorRecords(c *gin.Context) {
	errorRecords := h.monitor.GetErrorRecords()
	c.JSON(http.StatusOK, gin.H{
		"error_records": errorRecords,
		"count":         len(errorRecords),
	})
}

// GetMetrics 获取性能指标
func (h *HealthHandler) GetMetrics(c *gin.Context) {
	metrics := gin.H{
		"avg_response_time": h.monitor.GetAvgResponseTime(),
		"success_rate":      h.monitor.GetSuccessRate(),
		"uptime_seconds":    h.monitor.GetUptimeSeconds(),
		"active_connections": h.monitor.GetActiveConnections(),
		"current_calls":     h.monitor.GetCurrentCalls(),
	}
	c.JSON(http.StatusOK, metrics)
}

// 辅助函数：遮蔽敏感字符串，只显示前几位和后几位
func maskString(s string) string {
	if len(s) <= 8 {
		return s
	}
	return s[:4] + "***" + s[len(s)-4:]
}