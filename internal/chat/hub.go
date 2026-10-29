package chat

import "log"

// 管理所有在线客户端
type Hub struct {
	clients    map[*Client]struct{}
	register   chan *Client // 新客户端连接，通知 Hub 注册
	unregister chan *Client // 客户端断开连接，通知 Hub 注销
	broadcast  chan []byte  // 广播消息，通知 Hub 转发给所有客户端
}

// NewHub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	// 无限循环 select 负责处理 channel 内的客户端注册、注销和消息广播
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
			log.Printf("Client connected, total clients: %d", len(h.clients))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected, total clients: %d", len(h.clients))
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
					// 消息进入客户端的发送 channel，表示消息已成功发送
					log.Printf("Message sent to client, total clients: %d", len(h.clients))
				default:
					// 如果客户端的发送 channel 已满，说明客户端处理消息过慢，直接关闭连接
					delete(h.clients, client)
					close(client.send)
					client.conn.Close()
					log.Printf("Client disconnected, total clients: %d", len(h.clients))
				}
			}
		}
	}
}
