package main

import (
	"fmt"
	"log"
	"mini-chat/internal/chat"
	"net/http"
)

func main() {
	hub := chat.NewHub()
	go hub.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "mini-chat server is running!")
	})
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWS(hub, w, r)
	})

	log.Println("Serve started on :8765")
	if err := http.ListenAndServe(":8765", mux); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
