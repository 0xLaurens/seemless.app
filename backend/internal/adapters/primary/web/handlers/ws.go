package handlers

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
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

	room, err := wh.room.CreateRoom()
	if err != nil {
		return err
	}
	defer wh.wsDefer(user, room, c)

	displayName := &data.Message{Type: data.DisplayName, User: user}
	err = wh.msg.Send(displayName)
	if err != nil {
		log.Println("ERR -->> display name", err)
		return err
	}

	err = wh.msg.SendTargeted(&data.Message{Type: data.RoomCreated, RoomCode: room.GetCode()}, user)
	if err != nil {
		return err
	}

	_, err = wh.room.JoinRoom(room.GetCode(), user)
	if err != nil {
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
			break
		}
		err = wh.WsRequestHandler(msg, user)
		if err != nil {
			log.Println("ERR -->> read loop", err)
			break
		}
	}
	return nil
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

func (wh *WebsocketHandler) wsDefer(user *data.User, room *data.Room, conn *websocket.Conn) {
	log.Println("DBG", "defer", user.Username)

	log.Println(user.Username, "disconnected")
	latestUser, err := wh.us.GetUserById(user.GetId())
	if err != nil {
		log.Println("ERR -->> get room by id", err)
	}

	usersRoom, _ := wh.room.GetRoomById(latestUser.RoomID)
	log.Println("DBG", "usersRoom", usersRoom.GetCode())

	err = wh.msg.Broadcast(&data.Message{Type: data.PeerLeft, User: user}, latestUser.RoomID)
	if err != nil {
		log.Println("ERR -->> broadcast peer left", err)
	}

	err = wh.room.LeaveRoom(usersRoom.GetId(), user)
	if err != nil {
		log.Println("ERR -->> leave room", err)
	}
	_, err = wh.us.RemoveUser(user.GetId())
	if err != nil {
		log.Println("ERR -->> remove user", err)
	}
}

func (wh *WebsocketHandler) WsRequestHandler(msg *data.Message, user *data.User) error {
	switch msg.Type {
	case data.Offer,
		data.Answer:
		target, err := wh.us.GetUserByName(msg.Target)
		err = wh.msg.SendTargeted(msg, target)
		if err != nil {
			return err
		}
		break
	case data.PeerLeft,
		data.PeerUpdated,
		data.PeerJoined,
		data.NewIceCandidate:
		msg.Conn = user.Connection
		err := wh.msg.Broadcast(msg, user.RoomID)
		if err != nil {
			return err
		}
	case data.RoomJoin:
		log.Printf("PUBLIC ROOM %s JOIN REQUEST %s %s\n", msg.RoomCode, user.Username)
		room, err := wh.room.GetRoomByCode(data.RoomCode(strings.ToUpper(string(msg.RoomCode))))
		if err != nil {
			log.Println(err)
			_ = wh.msg.SendTargeted(&data.Message{Type: data.RoomCodeInvalid}, user)
		}

		user.SetRoom(room.GetId())

		err = wh.msg.SendTargeted(&data.Message{Type: data.RoomJoined, RoomCode: room.GetCode()}, user)
		if err != nil {
			return err
		}

		err = wh.msg.SendTargeted(&data.Message{Type: data.Peers, Users: room.GetClients(), RoomCode: room.GetCode()}, user)
		if err != nil {
			log.Println(err)
			return err
		}

		_, err = wh.room.JoinRoom(room.GetCode(), user)
		if err != nil {
			log.Println(err)
			return err
		}

		err = wh.msg.Broadcast(&data.Message{Type: data.PeerJoined, User: user}, room.GetId())
		if err != nil {
			log.Println(err)
			return err
		}
	case data.ChangeDisplayName:
		err := wh.room.ChangeDisplayName(user.RoomID, user, msg.DisplayName)
		if err != nil {
			err := wh.msg.SendTargeted(&data.Message{Type: data.DuplicateDisplayName}, user)
			if err != nil {
				return err
			}
		}
		err = wh.msg.Broadcast(&data.Message{Type: data.PeerUpdated, User: user, From: msg.From}, user.RoomID)
	default:
		log.Println("ERR -->> invalid request")
		err := wh.msg.InvalidMessage(nil)
		if err != nil {
			return err
		}
	}
	return nil
}
