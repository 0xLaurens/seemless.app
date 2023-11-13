package handlers

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/mssola/user_agent"
	"laurensdrop/internal/adapters/secondary"
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/core/services"
	"laurensdrop/internal/core/utils"
	"laurensdrop/internal/ports"
	"log"
)

type WebsocketHandler struct {
	us  ports.UserService
	msg ports.MessageService
}

func NewWebsocketHandler(us ports.UserRepo) *WebsocketHandler {
	return &WebsocketHandler{
		us: services.NewUserService(us),
		msg: services.NewMessageService(us,
			secondary.NewWebsocketMsgNotifier(),
			secondary.NewWebsocketMessageValidator(us),
		),
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
	username := ""
	usernamePrompt := data.Message{
		Type: data.MessageTypes.UsernamePrompt,
		Body: make(map[string]string),
	}
	usernamePrompt.Body["message"] = "Please provide a username"

	err := wh.msg.SendJSON(usernamePrompt)
	if err != nil {
		log.Println("ERR -->> write JSON error")
		return err
	}

	for username == "" {
		msg, err := ReadMessage(c)
		if err != nil {
			log.Println("ERR -->> read message", err)
			_ = wh.msg.InvalidMessage(nil)
			return err
		}

		username = msg.Body["username"]
		if username == "" || msg.Type != data.MessageTypes.Username {
			_ = wh.msg.InvalidMessage(nil)
			username = ""
		}

		if u, _ := wh.us.GetUserByName(username); u != nil || msg.Type != data.MessageTypes.Username {
			err := data.UserStoreError.DuplicateUsername
			msg := &data.Message{
				Type: data.MessageTypes.DuplicateUsername,
				Body: make(map[string]string),
			}
			msg.Body["message"] = data.UserStoreErrMessage(err)

			_ = wh.msg.InvalidMessage(msg)
			username = ""
		}
	}

	users, err := wh.us.GetAllUsers()
	if err != nil {
		return err
	}
	log.Printf("DBG -->> users: %v\n", users)

	peers := &data.Message{Type: data.MessageTypes.Peers, Users: users}
	err = wh.msg.Send(peers)
	if err != nil {
		return err
	}

	user := data.CreateUser(username, os, data.WithConnection(c))
	_, err = wh.us.AddUser(user)
	if err != nil {
		return err
	}

	defer wh.wsDefer(user, c)

	err = wh.msg.Broadcast(&data.Message{Type: data.MessageTypes.PeerJoined, User: user, From: username})
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	log.Printf("DBG -->> created user: %v", user.Username)
	for {
		msg, err := ReadMessage(c)
		if err != nil {
			log.Println("ERR -->> read message", err)
			return err
		}
		log.Println("DBG -->>", msg.Type, msg.SDP)
		err = wh.WsRequestHandler(msg, c)
		if err != nil {
			log.Println("ERR -->> readloop", err)
			return err
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
	log.Println("DBG", "defer")
	_, err := wh.us.RemoveUser(user.Username)
	if err != nil {
		return
	}
	err = conn.Close()
	if err != nil {
		return
	}

	err = wh.msg.Broadcast(&data.Message{Type: data.MessageTypes.PeerLeft, User: user, From: user.Username})
	if err != nil {
		log.Println("ERR", err)
		return
	}
}

func (wh *WebsocketHandler) WsRequestHandler(msg *data.Message, conn *websocket.Conn) error {
	switch msg.Type {
	case data.MessageTypes.Offer,
		data.MessageTypes.Answer,
		data.MessageTypes.PeerLeft,
		data.MessageTypes.PeerUpdated,
		data.MessageTypes.PeerJoined,
		data.MessageTypes.NewIceCandidate:
		msg.Conn = conn
		err := wh.msg.Broadcast(msg)
		if err != nil {
			log.Println("ERR -->> ws handler", err)
			return err
		}
		return nil
	default:
		log.Println("ERR -->> invalid request")
		err := wh.msg.InvalidMessage(nil)
		if err != nil {
			return err
		}
		return nil
	}
}
