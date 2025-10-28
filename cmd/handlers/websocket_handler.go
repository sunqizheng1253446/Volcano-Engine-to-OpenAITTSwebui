package handlers

import (
	"log"
	"net/http"

	"volcano-engine-tts-monitor/cmd/models"
	"volcano-engine-tts-monitor/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	monitor  *models.Monitor
	config   *config.Config
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
	// 升级HTTP连接为WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		h.monitor.RecordError("websocket_upgrade", err.Error())
		return
	}
	defer func() {
		// 关闭WebSocket连接并处理可能的错误
		if err := conn.Close(); err != nil {
			// 记录关闭错误但不中断程序流程
			log.Printf("WebSocket连接关闭错误: %v", err)
		}
	}()

	// 发送欢迎消息
	welcomeMsg := []byte("Connected to health monitoring system")
	if err := conn.WriteMessage(websocket.TextMessage, welcomeMsg); err != nil {
		log.Printf("WebSocket write error: %v", err)
		h.monitor.RecordError("websocket_write", err.Error())
		return
	}

	// 这里可以实现与被监控服务的WebSocket连接逻辑
	// 目前保持连接打开，等待客户端消息
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			// 连接可能已关闭
			break
		}

		// 回显收到的消息
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("WebSocket write error: %v", err)
			h.monitor.RecordError("websocket_write", err.Error())
			break
		}
	}
}
