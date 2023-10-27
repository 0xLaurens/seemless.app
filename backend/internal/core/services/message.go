package services

import (
	"errors"
	"github.com/gofiber/contrib/websocket"
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/core/utils"
	"laurensdrop/internal/ports"
	"log"
)

type MessageService struct {
	users ports.UserService
	conn  *websocket.Conn
}

func NewMessageService(us ports.UserService, conn *websocket.Conn) *MessageService {
	return &MessageService{
		users: us,
		conn:  conn,
	}
}

func (m *MessageService) Read() (*data.Message, error) {
	_, raw, err := m.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	message := &data.Message{}
	err = utils.MapJsonToStruct(raw, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (m *MessageService) Send(msg *data.Message) error {
	err := m.validateMessageOrigin(msg)
	if err != nil {
		log.Println("ERR", err)
		return err
	}
	return m.SendJSON(msg)
}

func (m *MessageService) SendJSON(json interface{}) error {
	err := m.conn.Conn.WriteJSON(json)
	if err != nil {
		return err
	}
	return nil
}

func (m *MessageService) validateMessageOrigin(msg *data.Message) error {
	if len(msg.From) < 1 {
		return nil
	}

	user, err := m.users.GetUserByConn(m.conn)
	if err != nil {
		return err
	}

	if user.Username != msg.From && len(msg.From) >= 1 {
		return errors.New("from message spoofed")
	}

	return nil
}

func (m *MessageService) Broadcast(msg *data.Message) error {
	err := m.validateMessageOrigin(msg)
	if err != nil {
		log.Println("ERR", err)
		return err
	}

	users, err := m.users.GetAllUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Username == msg.From {
			return nil
		}

		err = user.Connection.Conn.WriteJSON(msg)
		if err != nil {
			return err
		}

		return nil
	}
	return nil
}

func (m *MessageService) InvalidMessage(msg interface{}) error {
	return m.conn.WriteJSON(msg)
}
