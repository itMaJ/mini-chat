package main

import (
	"log"
	"mini-chat/internal/chat"
	"mini-chat/internal/server"
)

func main() {
	hub := chat.NewHub()
	go hub.Run()

	// 创建 Gin 路由
	router := server.NewRouter(hub)

	log.Println("Serve started on :8765")
	if err := router.Run(":8765"); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
