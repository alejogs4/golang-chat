package chat

import (
	"strings"

	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	rooms  map[string]*room
	events chan *event
}

// NewWebSocketConnection build a new instance of webSocketConnection struct, giving to rooms and expected events an initial value
func NewWebSocketConnection() WebSocketConnection {
	return WebSocketConnection{
		events: make(chan *event),
		rooms:  make(map[string]*room),
	}
}

func (wsc *WebSocketConnection) NewChatSubscriber(connection *websocket.Conn) {
	newSubscriber := subscriber{
		nickname:   "Anonymous",
		connection: connection,
		event:      wsc.events,
	}

	newSubscriber.ListenSubscriberMessages()
}

func (wsc *WebSocketConnection) ListenChatEvents() {
	for gotEvent := range wsc.events {
		switch gotEvent.id {
		case joinRoom:
			wsc.joinRoom(gotEvent.args[0], gotEvent)
		case sendMessage:
			wsc.sendAllMessage(gotEvent)
		case leave:
			wsc.quitCurrentRoom(gotEvent)
		case changeNickname:
			wsc.changeSubscriberNickname(gotEvent)
		}
	}
}

func (wsc *WebSocketConnection) joinRoom(newRoom string, event *event) {
	currentRoom, ok := wsc.rooms[strings.TrimSpace(newRoom)]
	if !ok {
		currentRoom = &room{
			name:    newRoom,
			members: make(map[*websocket.Conn]*subscriber),
		}
	}

	currentRoom.members[event.subscriber.connection] = event.subscriber
	wsc.rooms[strings.TrimSpace(newRoom)] = currentRoom
	wsc.quitCurrentRoom(event)

	currentRoom.BroadcastMessage(event, event.subscriber.connection)
}

func (wsc *WebSocketConnection) sendAllMessage(event *event) {
	if event.subscriber.room != nil {
		event.subscriber.connection.WriteJSON(map[string]string{"id": event.id, "message": event.message})
		event.subscriber.room.BroadcastMessage(event, event.subscriber.connection)
	}
}

func (wsc *WebSocketConnection) changeSubscriberNickname(incommingEvent *event) {
	newNickname := incommingEvent.args[0]
	incommingEvent.subscriber.nickname = newNickname

	wsc.sendAllMessage(&event{
		id:         incommingEvent.id,
		message:    incommingEvent.subscriber.nickname + " has changed nickname to " + newNickname,
		subscriber: incommingEvent.subscriber,
	})
}

func (wsc *WebSocketConnection) quitCurrentRoom(event *event) {
	subs := event.subscriber

	if subs.room != nil {
		oldRoom := event.subscriber.room
		subs.room.RemoveMember(subs.connection)
		oldRoom.BroadcastMessage(event, event.subscriber.connection)
	}
}
