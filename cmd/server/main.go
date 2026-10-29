package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// HTTP 路由管理器
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("mini-chat server is working!")
	})

	log.Println("mini-chat server is running on port 8765...")
	if err := http.ListenAndServe(":8765", mux); err != nil {
		log.Fatal(err)
	}
}
