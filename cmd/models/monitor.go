package models

import (
	"sync"
	"time"
)

type HealthStatus struct {
	Status              string `json:"status"`              // ok, warning, error
	ActiveConnections   int    `json:"active_connections"`   // 活动WebSocket连接数
	MaxConnections      int    `json:"max_connections"`      // 最大连接数
	CurrentCalls        int    `json:"current_calls"`        // 当前并发调用数
	MaxConcurrentCalls  int    `json:"max_concurrent_calls"` // 最大并发调用数
	UptimeSeconds       int64  `json:"uptime_seconds"`       // 服务运行时间（秒）
	LastCheckTime       string `json:"last_check_time"`      // 最后检查时间
	AvgResponseTime     int    `json:"avg_response_time"`    // 平均响应时间（毫秒）
	TodayRequestCount   int    `json:"today_request_count"`  // 今日请求总数
	SuccessRate         float64 `json:"success_rate"`        // 成功率
	ErrorCount          int    `json:"error_count"`          // 错误数量
	Config              ServiceConfig `json:"config"`        // 配置信息
	ResourceUsage       ResourceUsage `json:"resource_usage"` // 资源使用情况
	ConnectionStats     ConnectionStats `json:"connection_stats"` // 连接统计
}

type ServiceConfig struct {
	ByteDanceAppID       string `json:"byte_dance_app_id"`       // 火山引擎App ID
	ByteDanceCluster     string `json:"byte_dance_cluster"`     // 火山引擎集群
	ByteDanceVoiceType   string `json:"byte_dance_voice_type"`   // 语音类型
	ListenAddr           string `json:"listen_addr"`            // 监听地址
	MaxConnections       int    `json:"max_connections"`        // 最大连接数
	MaxConcurrentCalls   int    `json:"max_concurrent_calls"`    // 最大并发调用数
}

type ResourceUsage struct {
	CPUUsage    float64 `json:"cpu_usage"`    // CPU使用率
	MemoryUsage float64 `json:"memory_usage"` // 内存使用率
	MemoryMB    int64   `json:"memory_mb"`    // 内存使用量（MB）
	NetworkIn   int64   `json:"network_in"`   // 网络入站流量（KB）
	NetworkOut  int64   `json:"network_out"`  // 网络出站流量（KB）
}

type ConnectionStats struct {
	TodayPeak           int    `json:"today_peak"`           // 今日峰值连接数
	ConnectionSuccess   int    `json:"connection_success"`   // 成功连接数
	ConnectionFailure   int    `json:"connection_failure"`   // 失败连接数
	AvgConnectionTime   int    `json:"avg_connection_time"`   // 平均连接时间（秒）
}

type ErrorRecord struct {
	Timestamp string `json:"timestamp"`
	ErrorType string `json:"error_type"`
	Message   string `json:"message"`
}

type Monitor struct {
	sync.RWMutex
	startTime           time.Time
	activeConnections   int
	currentCalls        int
	todayRequestCount   int
	errorCount          int
	successCount        int
	todayPeak           int
	connectionSuccess   int
	connectionFailure   int
	responseTimes       []int
	avgConnectionTime   float64
	connectionTotalTime int64
	connectionsHandled  int
	errorRecords        []ErrorRecord
}

func NewMonitor() *Monitor {
	return &Monitor{
		startTime:           time.Now(),
		activeConnections:   0,
		currentCalls:        0,
		todayRequestCount:   0,
		errorCount:          0,
		successCount:        0,
		todayPeak:           0,
		connectionSuccess:   0,
		connectionFailure:   0,
		responseTimes:       make([]int, 0, 100),
		avgConnectionTime:   0,
		connectionTotalTime: 0,
		connectionsHandled:  0,
		errorRecords:        make([]ErrorRecord, 0, 50),
	}
}

// 更新各种指标的方法
func (m *Monitor) IncrementConnection() {
	m.Lock()
	defer m.Unlock()
	m.activeConnections++
	if m.activeConnections > m.todayPeak {
		m.todayPeak = m.activeConnections
	}
}

func (m *Monitor) DecrementConnection() {
	m.Lock()
	defer m.Unlock()
	if m.activeConnections > 0 {
		m.activeConnections--
	}
}

func (m *Monitor) IncrementCall() {
	m.Lock()
	defer m.Unlock()
	m.currentCalls++
	m.todayRequestCount++
}

func (m *Monitor) DecrementCall() {
	m.Lock()
	defer m.Unlock()
	if m.currentCalls > 0 {
		m.currentCalls--
	}
}

func (m *Monitor) RecordSuccess() {
	m.Lock()
	defer m.Unlock()
	m.successCount++
}

func (m *Monitor) RecordError(errorType, message string) {
	m.Lock()
	defer m.Unlock()
	m.errorCount++
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

func (m *Monitor) RecordResponseTime(timeMs int) {
	m.Lock()
	defer m.Unlock()
	m.responseTimes = append(m.responseTimes, timeMs)
	// 只保留最近100条记录用于计算平均值
	if len(m.responseTimes) > 100 {
		m.responseTimes = m.responseTimes[1:]
	}
}

func (m *Monitor) RecordConnectionTime(duration time.Duration) {
	m.Lock()
	defer m.Unlock()
	m.connectionTotalTime += int64(duration.Seconds())
	m.connectionsHandled++
	if m.connectionsHandled > 0 {
		m.avgConnectionTime = float64(m.connectionTotalTime) / float64(m.connectionsHandled)
	}
}

func (m *Monitor) GetSuccessRate() float64 {
	m.RLock()
	defer m.RUnlock()
	total := m.successCount + m.errorCount
	if total == 0 {
		return 100.0
	}
	return float64(m.successCount) / float64(total) * 100.0
}

func (m *Monitor) GetAvgResponseTime() int {
	m.RLock()
	defer m.RUnlock()
	if len(m.responseTimes) == 0 {
		return 0
	}
	total := 0
	for _, t := range m.responseTimes {
		total += t
	}
	return total / len(m.responseTimes)
}

func (m *Monitor) GetUptimeSeconds() int64 {
	return int64(time.Since(m.startTime).Seconds())
}

func (m *Monitor) GetActiveConnections() int {
	m.RLock()
	defer m.RUnlock()
	return m.activeConnections
}

func (m *Monitor) GetCurrentCalls() int {
	m.RLock()
	defer m.RUnlock()
	return m.currentCalls
}

func (m *Monitor) GetErrorRecords() []ErrorRecord {
	m.RLock()
	defer m.RUnlock()
	// 返回副本以避免并发问题
	records := make([]ErrorRecord, len(m.errorRecords))
	copy(records, m.errorRecords)
	return records
}