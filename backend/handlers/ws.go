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
	username := rand.Intn(100)
	user := data.CreateUser(string(rune(username)), "android", data.WithConnection(c))
	defer func() {
		log.Println("DBG -->> defer ws handler")
		hub.channels.unregister <- user
		err := c.Close()
		if err != nil {
			log.Println("ERR -->> failed to close connection")
			return
		}
	}()
	hub.channels.register <- user
	log.Printf("DBG -->> created user: %v", user.Username)
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("ERR -->> read message", err)
			return
		}
		log.Printf("DBG -->> recv: %s", msg)
		WsRequestHandler(msg, hub, user)
	}
}

func WsRequestHandler(msg []byte, hub *Hub, user *data.User) {
	log.Println("DBG -->> ws request handler")
	req := data.Request{}
	err := utils.MapJsonToStruct(msg, &req)
	if err != nil {
		log.Println("ERR -->> json matching", err)
		hub.channels.invalidMessage <- user
		return
	}

	switch req.Type {
	case data.RequestTypes.Offer,
		data.RequestTypes.Answer,
		data.RequestTypes.PeerLeft,
		data.RequestTypes.PeerUpdated,
		data.RequestTypes.PeerJoined,
		data.RequestTypes.NewIceCandidate:
		log.Println("")
		hub.channels.broadcast <- msg
	default:
		log.Println("ERR -->> invalid request")
		hub.channels.invalidMessage <- user
	}
}
