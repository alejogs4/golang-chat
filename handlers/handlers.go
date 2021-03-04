package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alejogs4/golangchat/chat"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ServeWebsocket register a new chat subscriber making it possible make use of the chat features
func ServeWebsocket(wsConnection *chat.WebSocketConnection) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")

		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		go wsConnection.NewChatSubscriber(conn)
	}
}

// GetAllCurrentChatRooms get the rooms recently created by users so far
func GetAllCurrentChatRooms(wsConnection *chat.WebSocketConnection) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")

		rooms := wsConnection.GetChatRooms()
		roomsAsJSON, err := json.Marshal(map[string]interface{}{"data": rooms})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(`{message: "Error getting rooms"}`))
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write(roomsAsJSON)
	}
}
