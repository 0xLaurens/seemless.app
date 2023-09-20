package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"laurensdrop/handlers"
	"laurensdrop/store"
	"log"
)

func main() {
	// init in memory store and hub
	s := store.NewUserStoreInMemory()
	hub := handlers.CreateHub(s)
	go hub.Run()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World :)")
	})

	app.Use("/ws", handlers.WSUpgrader)
	app.Use("/ws", websocket.New(func(conn *websocket.Conn) {
		handlers.WSHandler(conn, hub)
	}))

	// MONITOR IS STILL IN BETA, API MIGHT CHANGE
	app.Use("/monitor", monitor.New())

	log.Fatal(app.Listen(":3000"))
}
