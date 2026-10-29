package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	log.Printf("Client connected addr: %s", r.RemoteAddr)

	// 持续接收客户端发送的消息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("Received message from %s: %s", r.RemoteAddr, string(message))

		// 将消息发送回客户端
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

func main() {

	mux := http.NewServeMux()

	// 普通的 Http 路由
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "mini-chat server is working!")
	})

	// WebSocket 路由
	mux.HandleFunc("/ws", wsHandler)

	log.Println("mini-chat server is running on port 8765...")
	if err := http.ListenAndServe(":8765", mux); err != nil {
		log.Fatal(err)
	}
}
