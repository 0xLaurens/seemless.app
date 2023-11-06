package services

import (
	"github.com/gofiber/contrib/websocket"
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/ports"
)

type MessageService struct {
	users    ports.UserService
	notifier ports.MessageNotifier
}

func NewMessageService(us ports.UserService, notifier ports.MessageNotifier) *MessageService {
	return &MessageService{
		users:    us,
		notifier: notifier,
	}
}

func (m *MessageService) SetWebsocketMsgNotifierConn(conn *websocket.Conn) {
	m.notifier.SetWebsocketMsgNotifierConn(conn)
}

func (m *MessageService) Read() (*data.Message, error) {
	return m.notifier.Read()
}

func (m *MessageService) Send(msg *data.Message) error {
	return m.notifier.Send(msg)
}

func (m *MessageService) SendJSON(json interface{}) error {
	return m.notifier.SendJSON(json)
}

func (m *MessageService) Broadcast(msg *data.Message) error {
	users, err := m.users.GetAllUsers()
	if err != nil {
		return err
	}
	if msg.Target != "" {
		target, err := m.users.GetUserByName(msg.Target)
		if err != nil {
			return err
		}

		return m.notifier.SendTargeted(msg, target)
	}

	for _, user := range users {
		err := m.notifier.SendTargeted(msg, user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MessageService) InvalidMessage(msg interface{}) error {
	return m.notifier.InvalidMessage(msg)
}
