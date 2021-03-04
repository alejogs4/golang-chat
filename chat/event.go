package chat

const (
	sendMessage    = "SEND_MESSAGE"
	joinRoom       = "JOIN_ROOM"
	leave          = "LEAVE"
	changeNickname = "CHANGE_NICKNAME"
)

type event struct {
	id         string
	message    string
	args       []string
	subscriber *subscriber
}
