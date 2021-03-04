package chat

import "github.com/gorilla/websocket"

type room struct {
	name    string
	members map[*websocket.Conn]*subscriber
}

func (r *room) RemoveMember(connection *websocket.Conn) {
	delete(r.members, connection)
}

func (r *room) BroadcastMessage(message *event, sender *websocket.Conn) {
	for member, subscriber := range r.members {
		if member != sender {
			subscriber.SendMessage(message)
		}
	}
}
