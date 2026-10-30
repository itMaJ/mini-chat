package chat

import "time"

// 客户端发送给服务器的消息
type ChatRequest struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// 服务端根据客户端信息补充生成的消息
type ChatMessage struct {
	Type     string    `json:"type"`
	Content  string    `json:"content"`
	Username string    `json:"username"` // 后续可以考虑改为整个用户信息或者用户id
	SentAt   time.Time `json:"sent_at"`
}
