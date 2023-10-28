package secondary

import (
	"github.com/gofiber/contrib/websocket"
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/core/utils"
)

type WebsocketMsgNotifier struct {
	conn *websocket.Conn
}

func NewWebsocketMsgNotifier(conn *websocket.Conn) *WebsocketMsgNotifier {
	return &WebsocketMsgNotifier{
		conn: conn,
	}
}

func (m *WebsocketMsgNotifier) Read() (*data.Message, error) {
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

func (m *WebsocketMsgNotifier) Send(msg *data.Message) error {
	return m.SendJSON(msg)
}

func (m *WebsocketMsgNotifier) SendTargeted(msg *data.Message, target *data.User) error {
	err := target.Connection.Conn.WriteJSON(msg)
	if err != nil {
		return err
	}
	return nil
}

func (m *WebsocketMsgNotifier) SendJSON(json interface{}) error {
	err := m.conn.Conn.WriteJSON(json)
	if err != nil {
		return err
	}
	return nil
}

func (m *WebsocketMsgNotifier) InvalidMessage(msg interface{}) error {
	return m.conn.WriteJSON(msg)
}
