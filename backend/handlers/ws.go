package handlers

import (
	"encoding/json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"laurensdrop/data"
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

func ParseMessageToRequest(msg []byte) (data.Request, error) {
	req := data.Request{}
	err := json.Unmarshal(msg, &req)
	if err != nil {
		return data.Request{}, err
	}

	return req, nil
}

func RequestMatcher(msg []byte) {
	req, err := ParseMessageToRequest(msg)
	if err != nil {
		log.Println("NotJSON")
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
