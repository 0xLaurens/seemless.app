package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"laurensdrop/handlers"
	"log"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World :)")
	})

	app.Use("/ws", handlers.WSUpgrader)
	app.Use("/ws", websocket.New(handlers.WSHandler))

	log.Fatal(app.Listen(":3000"))
}
