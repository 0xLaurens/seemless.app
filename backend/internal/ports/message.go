package ports

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"laurensdrop/internal/core/data"
)

type MessageService interface {
	Read() (*data.Message, error)
	Send(msg *data.Message) error
	SendTargeted(msg *data.Message, user *data.User) error
	SendJSON(json interface{}) error
	Broadcast(msg *data.Message, roomId uuid.UUID) error
	InvalidMessage(msg interface{}) error
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

type MessageValidator interface {
	ValidateMessageOrigin(msg *data.Message) error
}
