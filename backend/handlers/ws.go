package handlers

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"laurensdrop/data"
	"laurensdrop/utils"
	"log"
)

func WSUpgrader(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func WSHandler(c *websocket.Conn, hub *Hub) {
	username := ""
	err := c.WriteJSON(fiber.Map{
		"type":    "UsernamePrompt",
		"message": "provide a username",
	})
	if err != nil {
		log.Println("ERR -->> write JSON error")
		return
	}

	for username == "" {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("ERR -->> read message", err)
			return
		}

		req := data.Request{}
		err = utils.MapJsonToStruct(msg, &req)
		if err != nil {
			log.Println("ERR -->> json matching", err)
			return
		}

		username = req.Body["username"]
		if username == "" || req.Type != data.RequestTypes.Username {
			writeMSG(c, fiber.Map{
				"type":    data.WsError.InvalidRequestBody,
				"message": "Please Send the right format for the username request",
			})
		}

		if u, _ := hub.users.GetUserByName(username); u != nil {
			writeMSG(c, fiber.Map{
				"type":    data.UserStoreError.DuplicateUsername,
				"message": "duplicate username error",
			})
			username = ""
		}
	}

	writeMSG(c, fiber.Map{
		"type":     data.RequestTypes.PeerJoined,
		"username": username,
	})
	user := data.CreateUser(username, "android", data.WithConnection(c))
	hub.channels.register <- user

	defer wsDefer(hub, user, c)

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

func wsDefer(hub *Hub, user *data.User, c *websocket.Conn) {
	log.Println("DBG -->> defer ws handler")
	hub.channels.unregister <- user
	err := c.Close()
	if err != nil {
		log.Println("ERR -->> failed to close connection")
		return
	}
}

func writeMSG(c *websocket.Conn, json interface{}) {
	err := c.WriteJSON(json)
	if err != nil {
		log.Println("ERR -->> failed to writeMSG", err)
		return
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
