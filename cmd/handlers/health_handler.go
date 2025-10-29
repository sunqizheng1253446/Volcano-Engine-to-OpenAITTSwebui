package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"volcano-engine-tts-monitor/cmd/models"
	"volcano-engine-tts-monitor/config"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	monitor *models.Monitor
	config  *config.Config
	client  *http.Client
}

func NewHealthHandler(monitor *models.Monitor, config *config.Config) *HealthHandler {
	return &HealthHandler{
		monitor: monitor,
		config:  config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// HealthCheck 处理健康检查请求
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// 从目标API获取健康数据
	targetHealth, err := h.fetchTargetHealth()
	if err != nil {
		// 记录错误
		h.monitor.RecordError("fetch_target_health", err.Error())

		// 返回错误状态
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Failed to fetch target service health data",
			"message": err.Error(),
		})
		return
	}

	// 更新监控器中的数据
	h.monitor.UpdateTargetHealthStatus(targetHealth)

	// 返回健康状态
	c.JSON(http.StatusOK, targetHealth)
}

// fetchTargetHealth 从目标API获取健康数据
func (h *HealthHandler) fetchTargetHealth() (models.TargetHealthStatus, error) {
	var targetHealth models.TargetHealthStatus

	// 构造完整的API URL，添加 /api/health 路径
	apiURL := h.config.TargetAPIURL + "/api/health"

	// 发起HTTP请求获取目标服务健康数据
	resp, err := h.client.Get(apiURL)
	if err != nil {
		return targetHealth, fmt.Errorf("failed to connect to target API: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态
	if resp.StatusCode != http.StatusOK {
		return targetHealth, fmt.Errorf("target API returned status code: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return targetHealth, fmt.Errorf("failed to read response body: %w", err)
	}

	// 解析JSON数据
	if err := json.Unmarshal(body, &targetHealth); err != nil {
		return targetHealth, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// 更新最后检查时间
	targetHealth.LastCheckTime = time.Now().Format(time.RFC3339)

	return targetHealth, nil
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
	targetHealth := h.monitor.GetTargetHealthStatus()

	metrics := gin.H{
		"status":             targetHealth.Status,
		"uptime_seconds":     targetHealth.UptimeSeconds,
		"active_connections": targetHealth.ActiveConnections,
		"current_calls":      targetHealth.CurrentCalls,
		"avg_response_time":  targetHealth.AvgResponseTime,
		"success_rate":       targetHealth.SuccessRate,
		"error_count":        targetHealth.ErrorCount,
		"cpu_usage":          targetHealth.CPUUsage,
		"memory_usage":       targetHealth.MemoryUsage,
		"request_count":      targetHealth.RequestCount,
		"last_check_time":    targetHealth.LastCheckTime,
	}

	c.JSON(http.StatusOK, metrics)
}
