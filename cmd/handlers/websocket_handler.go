package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"volcano-engine-tts-monitor/cmd/models"
	"volcano-engine-tts-monitor/config"
)

type WebSocketHandler struct {
	monitor *models.Monitor
	config  *config.Config
	upgrader websocket.Upgrader
}

func NewWebSocketHandler(monitor *models.Monitor, config *config.Config) *WebSocketHandler {
	return &WebSocketHandler{
		monitor: monitor,
		config:  config,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// 在生产环境中应该检查Origin
				return true
			},
		},
	}
}

// HandleWebSocket 处理WebSocket连接
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// 检查连接数限制
	if h.monitor.GetActiveConnections() >= h.config.MaxConnections {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Connection limit exceeded",
		})
		h.monitor.RecordError("connection_limit", "Connection rejected due to limit")
		return
	}

	// 升级HTTP连接为WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		h.monitor.RecordError("websocket_upgrade", err.Error())
		return
	}
	// 增加连接计数
	h.monitor.IncrementConnection()

	// 确保在函数退出时清理资源
	defer func() {
		// 关闭WebSocket连接并处理可能的错误
		if err := conn.Close(); err != nil {
			// 记录关闭错误但不中断程序流程
			log.Printf("WebSocket连接关闭错误: %v", err)
		}
		// 减少连接计数
		h.monitor.DecrementConnection()
	}()

	// 记录连接成功
	h.monitor.RecordSuccess()

	// 记录连接开始时间，用于计算连接时长
	startTime := time.Now()
	defer func() {
		h.monitor.RecordConnectionTime(time.Since(startTime))
	}()

	// 模拟处理TTS请求
	h.handleTTSRequest(conn)
}

// handleTTSRequest 处理TTS请求的WebSocket消息
func (h *WebSocketHandler) handleTTSRequest(conn *websocket.Conn) {
	// 增加并发调用计数
	h.monitor.IncrementCall()
	defer h.monitor.DecrementCall()

	// 模拟处理时间
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		h.monitor.RecordResponseTime(int(elapsed.Milliseconds()))
	}()

	// 读取客户端消息
	messageType, message, err := conn.ReadMessage()
	if err != nil {
		log.Printf("WebSocket read error: %v", err)
		h.monitor.RecordError("websocket_read", err.Error())
		return
	}

	log.Printf("Received TTS request: %s", message)

	// 模拟处理延迟
	time.Sleep(500 * time.Millisecond)

	// 模拟发送音频数据（在实际应用中，这里应该是与火山引擎TTS服务交互）
	responseMessage := []byte("Simulated audio data chunk")
	err = conn.WriteMessage(messageType, responseMessage)
	if err != nil {
		log.Printf("WebSocket write error: %v", err)
		h.monitor.RecordError("websocket_write", err.Error())
		return
	}

	// 模拟发送多个音频数据块
	for i := 0; i < 3; i++ {
		time.Sleep(200 * time.Millisecond)
		chunk := []byte("Audio chunk " + string(i+'0'))
		err = conn.WriteMessage(messageType, chunk)
		if err != nil {
			log.Printf("WebSocket write error: %v", err)
			h.monitor.RecordError("websocket_write", err.Error())
			return
		}
	}

	// 记录成功处理
	h.monitor.RecordSuccess()
}