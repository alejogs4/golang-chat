package main

import (
	"log"
	"net/http"

	"github.com/alejogs4/golangchat/chat"
	"github.com/alejogs4/golangchat/handlers"
)

func main() {
	wsConnection := chat.NewWebSocketConnection()
	go wsConnection.ListenChatEvents()

	http.HandleFunc("/ws", handlers.ServeWebsocket(&wsConnection))
	http.HandleFunc("/rooms", handlers.GetAllCurrentChatRooms(&wsConnection))

	log.Println("Initializing ws server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
