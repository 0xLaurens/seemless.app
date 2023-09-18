package handlers

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"laurensdrop/data"
	"laurensdrop/utils"
	"log"
	"math/rand"
)

func WSUpgrader(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func WSHandler(c *websocket.Conn, hub *Hub) {
	go hub.run()
	username := rand.Intn(100)
	user := data.CreateUser(string(rune(username)), "andoird", data.WithConnection(c))
	hub.channels.register <- user
	log.Printf("created user: %v", user)
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Fatal("ERR -->>", err)
		}
		log.Printf("DBG -->> recv: %s", msg)
		RequestMatcher(msg, hub)
	}
}

func RequestMatcher(msg []byte, hub *Hub) {
	log.Println("DBG -->> request matcher")
	req := data.Request{}
	err := utils.MapJsonToStruct(msg, &req)
	if err != nil {
		log.Fatal("ERR -->> json matching", err)
		return
	}

	switch req.Type {
	case data.Offer,
		data.Answer,
		data.PeerLeft,
		data.PeerUpdated,
		data.PeerJoined,
		data.NewIceCandidate:
		log.Println("")
		hub.channels.broadcast <- msg
		break
	default:
		log.Fatal("ERR -->> invalid request")
	}
}
