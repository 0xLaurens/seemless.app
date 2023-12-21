package handlers

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mssola/user_agent"
	"laurensdrop/internal/adapters/secondary"
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/core/services"
	"laurensdrop/internal/core/utils"
	"laurensdrop/internal/ports"
	"log"
	"strings"
)

type WebsocketHandler struct {
	us   ports.UserService
	room ports.RoomService
	msg  ports.MessageService
}

func NewWebsocketHandler(us ports.UserRepo, room ports.RoomService) *WebsocketHandler {
	return &WebsocketHandler{
		us: services.NewUserService(us),
		msg: services.NewMessageService(us,
			room,
			secondary.NewWebsocketMsgNotifier(),
			secondary.NewWebsocketMessageValidator(us),
		),
		room: room,
	}
}

func (wh *WebsocketHandler) UpgradeWebsocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (wh *WebsocketHandler) HandleWebsocket(c *websocket.Conn) error {
	ua := user_agent.New(c.Headers("User-Agent"))
	os := ua.OSInfo().Name
	wh.msg.SetWebsocketMsgNotifierConn(c)

	user := data.CreateUser(os, data.WithConnection(c))
	_, err := wh.us.AddUser(user)
	if err != nil {
		return err
	}

	room, err := wh.room.JoinLocalRoom(user)
	if err != nil {
		return err
	}
	defer wh.wsDefer(user, c)

	displayName := &data.Message{Type: data.DisplayName, User: user}
	err = wh.msg.Send(displayName)
	if err != nil {
		log.Println("ERR -->> display name", err)
		return err
	}

	err = wh.msg.SendTargeted(&data.Message{Type: data.Peers, Users: room.GetClients()}, user)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = wh.msg.Broadcast(&data.Message{Type: data.PeerJoined, User: user, From: user.Username}, room.GetId())
	if err != nil {
		fmt.Println(err)
		return err
	}

	log.Printf("DBG -->> created user: %v", user.Username)
	for {
		msg, err := ReadMessage(c)
		if err != nil {
			log.Println("ERR -->> read message", err)
		}
		err = wh.WsRequestHandler(msg, user)
		if err != nil {
			log.Println("ERR -->> read loop", err)
		}
	}
}

func ReadMessage(conn *websocket.Conn) (*data.Message, error) {
	_, raw, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	message := &data.Message{}
	err = utils.MapJsonToStruct(raw, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (wh *WebsocketHandler) wsDefer(user *data.User, conn *websocket.Conn) {
	log.Println("DBG", "defer", user.Username)
	_, err := wh.us.RemoveUser(user.Username)
	if err != nil {
		return
	}
	_ = conn.Close() // attempt to close ignore if it's not successful

	if user.LocalRoom != uuid.Nil {
		localRoom, err := wh.room.GetRoomById(user.LocalRoom)
		if err != nil {
			log.Println("err ->> local room get by id", err)
			return
		}

		localRoom.RemoveClient(user)
		err = wh.msg.Broadcast(&data.Message{Type: data.PeerLeft, User: user, From: user.Username}, user.LocalRoom)
		if err != nil {
			log.Println("ERR -->> local room", err)
			return
		}
	}

	if user.PublicRoom != uuid.Nil {
		err = wh.msg.Broadcast(&data.Message{Type: data.PublicRoomLeft, User: user, From: user.Username}, user.PublicRoom)
		if err != nil {
			log.Println("ERR -->> public room", err)
			return
		}
	}
}

func (wh *WebsocketHandler) WsRequestHandler(msg *data.Message, user *data.User) error {
	log.Println("WsRequestHandler", msg, user.LocalRoom, user.PublicRoom)
	switch msg.Type {
	case data.Offer,
		data.Answer,
		data.PeerLeft,
		data.PeerUpdated,
		data.PeerJoined,
		data.NewIceCandidate:
		msg.Conn = user.Connection
		if user.LocalRoom != uuid.Nil {
			err := wh.msg.Broadcast(msg, user.LocalRoom)
			if err != nil {
				log.Println("ERR -->> local room broadcast", err)
				return err
			}
		}
		if user.PublicRoom != uuid.Nil {
			err := wh.msg.Broadcast(msg, user.PublicRoom)
			if err != nil {
				log.Println("ERR -->> public room broadcast", err)
				return err
			}
		}
	case data.PublicRoomCreate:
		if user.PublicRoom != uuid.Nil {
			log.Println("You already made a room")
			break
		}

		room, err := wh.room.CreatePublicRoom()
		if err != nil {
			fmt.Println("Failed to create room")
			return err
		}

		err = wh.room.JoinPublicRoom(room.GetCode(), user)
		if err != nil {
			log.Println(err)
			return err
		}

		room, err = wh.room.GetRoomByCode(room.GetCode())
		if err != nil {
			return err
		}

		err = wh.msg.SendTargeted(&data.Message{Type: data.PublicRoomCreated, RoomCode: room.GetCode()}, user)
		if err != nil {
			log.Println(err)
			return err
		}
		err = wh.msg.Broadcast(&data.Message{Type: data.PublicRoomJoin, User: user}, user.PublicRoom)
		if err != nil {
			log.Println(err)
			return err
		}

	case data.PublicRoomJoin:
		log.Printf("PUBLIC ROOM %s JOIN REQUEST %s\n", msg.RoomCode, user.Username)
		room, err := wh.room.GetRoomByCode(data.RoomCode(strings.ToUpper(string(msg.RoomCode))))
		if err != nil {
			log.Println(err)
			_ = wh.msg.SendTargeted(&data.Message{Type: data.PublicRoomIdInvalid}, user)
		}

		err = wh.msg.SendTargeted(&data.Message{Type: data.PublicRoomPeers, Users: room.GetClients(), RoomCode: room.GetCode()}, user)
		if err != nil {
			log.Println(err)
			return err
		}

		err = wh.room.JoinPublicRoom(room.GetCode(), user)
		if err != nil {
			log.Println(err)
			return err
		}

		err = wh.msg.Broadcast(&data.Message{Type: data.PublicRoomJoin, User: user}, room.GetId())
		if err != nil {
			log.Println(err)
			return err
		}
	default:
		log.Println("ERR -->> invalid request")
		err := wh.msg.InvalidMessage(nil)
		if err != nil {
			return err
		}
	}
	return nil
}
