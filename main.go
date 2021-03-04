package main

import (
	"log"
	"net/http"

	"github.com/alejogs4/golangchat/chat"
	"github.com/gorilla/websocket"
)

func main() {
	wsConnection := chat.NewWebSocketConnection()
	go wsConnection.ListenChatEvents()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebsocket(&wsConnection, w, r)
	})

	log.Println("Initializing ws server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWebsocket(wsConnection *chat.WebSocketConnection, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go wsConnection.NewChatSubscriber(conn)
}
