package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"laurensdrop/handlers"
	"laurensdrop/store"
	"log"
)

func main() {
	// init in memory store and hub
	s := store.NewUserStoreInMemory()
	hub := handlers.CreateHub(s)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World :)")
	})

	app.Use("/ws", handlers.WSUpgrader)
	app.Use("/ws", websocket.New(func(conn *websocket.Conn) {
		handlers.WSHandler(conn, hub)
	}))

	log.Fatal(app.Listen(":3000"))
}
