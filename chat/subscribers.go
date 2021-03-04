package chat

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type subscriber struct {
	nickname   string
	room       *room
	connection *websocket.Conn
	event      chan<- *event
}

func (s *subscriber) ListenSubscriberMessages() {
	defer func() {
		s.connection.Close()
	}()

	for {
		_, message, err := s.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)

			}
			break
		}

		var incommingMessage struct {
			Event    string `json:"event"`
			Message  string `json:"message,omitempty"`
			Room     string `json:"room,omitempty"`
			Nickname string `json:"nickname,omitempty"`
		}

		messageBuffer := bytes.NewBuffer(message)
		err = json.NewDecoder(messageBuffer).Decode(&incommingMessage)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		switch incommingMessage.Event {
		case joinRoom:
			s.event <- &event{
				id:         incommingMessage.Event,
				message:    s.nickname + " has joined to the room",
				args:       []string{incommingMessage.Room},
				subscriber: s,
			}
		case sendMessage:
			s.event <- &event{
				id:         incommingMessage.Event,
				message:    incommingMessage.Message,
				subscriber: s,
			}
		case leave:
			s.event <- &event{
				id:         incommingMessage.Event,
				message:    incommingMessage.Message,
				subscriber: s,
			}
		case changeNickname:
			s.event <- &event{
				id:         incommingMessage.Event,
				message:    incommingMessage.Message,
				subscriber: s,
				args:       []string{incommingMessage.Nickname},
			}
		}
	}
}

func (s *subscriber) SendMessage(message *event) {
	s.connection.WriteJSON(map[string]string{
		"event":   message.id,
		"message": message.message,
	})
}
