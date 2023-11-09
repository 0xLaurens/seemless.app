package ports

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"laurensdrop/internal/core/data"
)

type MessageService interface {
	Read() (*data.Message, error)
	Send(msg *data.Message) error
	SendJSON(json interface{}) error
	Broadcast(msg *data.Message) error
	InvalidMessage() error
	SetWebsocketMsgNotifierConn(conn *websocket.Conn)
}

type MessageHandler interface {
	UpgradeWebsocket(c *fiber.Ctx) error
	HandleWebsocket(c *websocket.Conn) error
}

type MessageNotifier interface {
	Read() (*data.Message, error)
	Send(msg *data.Message) error
	SendTargeted(msg *data.Message, target *data.User) error
	SendJSON(json interface{}) error
	InvalidMessage(msg interface{}) error
	SetWebsocketMsgNotifierConn(conn *websocket.Conn)
}
