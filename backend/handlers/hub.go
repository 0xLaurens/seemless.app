package handlers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"laurensdrop/data"
	"laurensdrop/store"
	"log"
)

// Channels are for the different actions that are listened for
type Channels struct {
	broadcast      chan []byte
	register       chan *data.User
	unregister     chan *data.User
	invalidMessage chan *data.User
}

func createChannels() Channels {
	return Channels{
		broadcast:      make(chan []byte),
		register:       make(chan *data.User),
		unregister:     make(chan *data.User),
		invalidMessage: make(chan *data.User),
	}
}

// Hub handles the channels and connections
type Hub struct {
	users    store.UserStore
	channels Channels
}

func CreateHub(Store store.UserStore) *Hub {
	return &Hub{
		users:    Store,
		channels: createChannels(),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case msg := <-h.channels.broadcast:
			log.Println("DBG -->> broadcast")
			users, err := h.users.GetAllUsers()
			if err != nil {
				return
			}

			for _, u := range users {
				err := u.Connection.Conn.WriteMessage(1, msg)
				if err != nil {
					return
				}
			}
		case user := <-h.channels.register:
			log.Println("DBG -->> register", user)
			_, err := h.users.AddUser(user)
			if err != nil {
				log.Fatal("ERR -->> register failed", err)
				return
			}
		case user := <-h.channels.unregister:
			log.Println("DBG -->> unregister", user)
			_, err := h.users.RemoveUser(user.Username)
			if err != nil {
				log.Fatal("ERR -->> unregister failed", err)
				return
			}
		case user := <-h.channels.invalidMessage:
			log.Println("DBG -->> invalid message")
			errRes := fiber.NewError(int(data.WsError.InvalidRequestBody))
			res, err := json.Marshal(errRes)
			if err != nil {
				log.Println("ERR -->> failed to create json from err", err)
				return
			}

			err = user.Connection.Conn.WriteMessage(1, res)
			if err != nil {
				log.Println("ERR -->> failed to send message", err)
				return
			}

		}
	}
}
