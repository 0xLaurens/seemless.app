package services

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/ports"
	"log"
)

type MessageService struct {
	users     ports.UserService
	room      ports.RoomService
	notifier  ports.MessageNotifier
	validator ports.MessageValidator
}

func NewMessageService(us ports.UserService, room ports.RoomService, notifier ports.MessageNotifier, validator ports.MessageValidator) *MessageService {
	return &MessageService{
		users:     us,
		room:      room,
		notifier:  notifier,
		validator: validator,
	}
}

var _ ports.MessageService = (*MessageService)(nil)

func (m *MessageService) SetWebsocketMsgNotifierConn(conn *websocket.Conn) {
	m.notifier.SetWebsocketMsgNotifierConn(conn)
}

func (m *MessageService) Read() (*data.Message, error) {
	return m.notifier.Read()
}

func (m *MessageService) Send(msg *data.Message) error {
	return m.notifier.Send(msg)
}

func (m *MessageService) SendTargeted(msg *data.Message, user *data.User) error {
	err := m.validator.ValidateMessageOrigin(msg)
	if err != nil {
		return err
	}
	return m.notifier.SendTargeted(msg, user)
}

func (m *MessageService) SendJSON(json interface{}) error {
	return m.notifier.SendJSON(json)
}

func (m *MessageService) Broadcast(msg *data.Message, roomId uuid.UUID) error {
	room, err := m.room.GetRoomById(roomId)
	if err != nil {
		return err
	}

	users := room.GetClients()
	log.Println("rooms users", users, len(users))
	if len(users) < 1 {
		return nil
	}

	if msg.Target != "" {
		target, err := m.users.GetUserByName(msg.Target)
		if err != nil {
			return err
		}
		return m.notifier.SendTargeted(msg, target)
	}

	for _, user := range users {
		err = m.notifier.SendTargeted(msg, user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MessageService) InvalidMessage(msg interface{}) error {
	if msg == nil {
		invalid := data.Message{
			Type: data.InvalidMessage,
			Body: make(map[string]string),
		}

		invalid.Body["message"] = data.InvalidRequestBody.Error().Error()
		return m.notifier.InvalidMessage(invalid)
	} else {
		return m.notifier.InvalidMessage(msg)
	}

}
