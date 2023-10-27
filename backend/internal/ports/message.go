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
	InvalidMessage(msg interface{}) error
}

type MessageHandler interface {
	UpgradeWebsocket(c *fiber.Ctx) error
	HandleWebsocket(c *websocket.Conn) error
}
