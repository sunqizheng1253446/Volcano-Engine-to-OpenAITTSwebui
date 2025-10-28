package models

import (
	"sync"
	"time"
)

// TargetHealthStatus 被监控服务的健康状态数据结构
type TargetHealthStatus struct {
	Status            string  `json:"status"`             // 服务状态: ok, warning, error
	UptimeSeconds     int64   `json:"uptime_seconds"`     // 服务运行时间（秒）
	ActiveConnections int     `json:"active_connections"` // 活动连接数
	CurrentCalls      int     `json:"current_calls"`      // 当前并发调用数
	AvgResponseTime   int     `json:"avg_response_time"`  // 平均响应时间（毫秒）
	SuccessRate       float64 `json:"success_rate"`       // 成功率（百分比）
	ErrorCount        int     `json:"error_count"`        // 错误数量
	LastCheckTime     string  `json:"last_check_time"`    // 最后检查时间
	CPUUsage          float64 `json:"cpu_usage"`          // CPU使用率（百分比）
	MemoryUsage       float64 `json:"memory_usage"`       // 内存使用率（百分比）
	RequestCount      int     `json:"request_count"`      // 总请求数
	Version           string  `json:"version"`            // 服务版本
}

// Monitor 监控器结构
type Monitor struct {
	sync.RWMutex
	targetHealthStatus TargetHealthStatus
	lastUpdate         time.Time
	errorRecords       []ErrorRecord
}

// ErrorRecord 错误记录结构
type ErrorRecord struct {
	Timestamp string `json:"timestamp"`
	ErrorType string `json:"error_type"`
	Message   string `json:"message"`
}

// NewMonitor 创建新的监控器实例
func NewMonitor() *Monitor {
	return &Monitor{
		targetHealthStatus: TargetHealthStatus{
			Status:        "unknown",
			LastCheckTime: time.Now().Format(time.RFC3339),
			Version:       "unknown",
		},
		lastUpdate:   time.Now(),
		errorRecords: make([]ErrorRecord, 0, 50),
	}
}

// UpdateTargetHealthStatus 更新被监控服务的健康状态
func (m *Monitor) UpdateTargetHealthStatus(status TargetHealthStatus) {
	m.Lock()
	defer m.Unlock()
	m.targetHealthStatus = status
	m.lastUpdate = time.Now()
}

// GetTargetHealthStatus 获取被监控服务的健康状态
func (m *Monitor) GetTargetHealthStatus() TargetHealthStatus {
	m.RLock()
	defer m.RUnlock()
	return m.targetHealthStatus
}

// RecordError 记录错误信息
func (m *Monitor) RecordError(errorType, message string) {
	m.Lock()
	defer m.Unlock()

	// 记录错误信息，最多保存50条
	record := ErrorRecord{
		Timestamp: time.Now().Format(time.RFC3339),
		ErrorType: errorType,
		Message:   message,
	}

	m.errorRecords = append(m.errorRecords, record)
	if len(m.errorRecords) > 50 {
		m.errorRecords = m.errorRecords[1:]
	}
}

// GetErrorRecords 获取错误记录
func (m *Monitor) GetErrorRecords() []ErrorRecord {
	m.RLock()
	defer m.RUnlock()
	// 返回副本以避免并发问题
	records := make([]ErrorRecord, len(m.errorRecords))
	copy(records, m.errorRecords)
	return records
}

// GetLastUpdate 获取最后更新时间
func (m *Monitor) GetLastUpdate() time.Time {
	m.RLock()
	defer m.RUnlock()
	return m.lastUpdate
}
