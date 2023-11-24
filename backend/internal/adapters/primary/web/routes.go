package web

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"laurensdrop/internal/ports"
)

func (a *App) initAppRoutes(ws ports.MessageHandler) {
	a.fiber.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World :)")
	})

	a.fiber.Use("/ws", ws.UpgradeWebsocket)

	a.fiber.Use("/ws", websocket.New(func(conn *websocket.Conn) {
		err := ws.HandleWebsocket(conn)
		if err != nil {
			return
		}
	}))

	a.fiber.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
}
