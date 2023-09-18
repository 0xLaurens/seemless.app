package handlers

import (
	"laurensdrop/data"
	"laurensdrop/store"
	"log"
)

// Channels are for the different actions that are listened for
type Channels struct {
	broadcast  chan []byte
	register   chan *data.User
	unregister chan *data.User
}

func createChannels() Channels {
	return Channels{
		broadcast:  make(chan []byte),
		register:   make(chan *data.User),
		unregister: make(chan *data.User),
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

func (h *Hub) run() {
	for {
		select {
		case msg := <-h.channels.broadcast:
			log.Println("DBG -->> broadcast")
			users, err := h.users.GetAllUsers()
			log.Println(users)
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
			log.Println("DBG -->> register", user)
			_, err := h.users.RemoveUser(user.Username)
			if err != nil {
				log.Fatal("ERR -->> unregister failed", err)
				return
			}
		}
	}
}
