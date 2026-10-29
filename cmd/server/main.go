package main

import (
	"log"
	"mini-chat/internal/chat"
	"net/http"
)

func main() {
	hub := chat.NewHub()
	go hub.Run()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web"))) // 返回首页
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWS(hub, w, r)
	})

	log.Println("Serve started on :8765")
	if err := http.ListenAndServe(":8765", mux); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
