package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// 管理单个客户端，客户端中包含对应用户信息
type Client struct {
	hub      *Hub            // 所属的 Hub
	conn     *websocket.Conn // WebSocket 连接
	send     chan []byte     // 发送消息的 channel
	username string          // client对应用户的相关信息
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}

	username := strings.TrimSpace(r.URL.Query().Get("nickname"))

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256), // 发送消息的 channel，缓冲区大小为 256
		username: username,               // 从请求中获取用户名
	}

	// 注册客户端到 Hub
	client.hub.register <- client

	// 启动 goroutine 分别处理读和写
	go client.writePump()
	client.readPump()
}

// readPump 持续读取客户端发送的消息，并将消息转发给 Hub 的广播 channel
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// 单条消息的最大字节数限制为 4KB，防止客户端发送过大的消息导致服务端内存占用过高
	c.conn.SetReadLimit(4 * 1024)

	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		if messageType != websocket.TextMessage {
			log.Printf("Ignoring non-text message of type %d", messageType)
			continue
		}

		// 解析 JSON，并组装成服务端消息，发送给 Hub 的广播 channel，作为字符串进行传输
		var request ChatRequest
		if err := json.Unmarshal(message, &request); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}
		if request.Type != "chat" {
			log.Printf("Unsupported message type %s", request.Type)
			continue
		}
		content := strings.TrimSpace(request.Content)
		if content == "" {
			log.Printf("Ignoring empty chat message")
			continue
		}

		response := ChatMessage{
			Type:     "chat",
			Username: c.username,
			Content:  content,
			SentAt:   time.Now().UTC(),
		}

		data, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling message: %v", err)
			continue
		}
		c.hub.broadcast <- data
	}
}

// writePump 持续从 Hub 的广播 channel 中读取消息，并将消息发送给客户端
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
}
