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

func WSHandler(c *websocket.Conn) {
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read: ", err)
			break
		}
		log.Printf("recv: %s", msg)
		RequestMatcher(msg)

		if err = c.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func RequestMatcher(msg []byte) {
	req := data.Request{}
	err := utils.MapJsonToStruct(msg, &req)
	if err != nil {
		log.Println("Json error", err)
		return
	}

	switch req.Type {
	case data.Offer:
		log.Println("Offer")
		break
	case data.Answer:
		log.Println("Answer")
		break
	case data.PeerLeft:
		log.Println("PeerLeft")
		break
	case data.PeerJoined:
		log.Println("PeerJoined")
		break
	case data.PeerUpdated:
		log.Println("PeerUpdated")
		break
	case data.NewIceCandidate:
		log.Println("NewIceCandidate")
	default:
		log.Println("InvalidRequest")
	}
}
