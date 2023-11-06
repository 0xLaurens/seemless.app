package handlers

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
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
		us:  services.NewUserService(us),
		msg: services.NewMessageService(us, secondary.NewWebsocketMsgNotifier()),
	}
}

func (wh *WebsocketHandler) UpgradeWebsocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (wh *WebsocketHandler) HandleWebsocket(c *websocket.Conn) error {
	wh.msg.SetWebsocketMsgNotifierConn(c)
	username := ""
	err := wh.msg.SendJSON(fiber.Map{
		"type":    "UsernamePrompt",
		"message": "provide a username",
	})
	if err != nil {
		log.Println("ERR -->> write JSON error")
		return err
	}

	for username == "" {
		msg, err := ReadMessage(c)
		if err != nil {
			log.Println("ERR -->> read message", err)
			return err
		}

		username = msg.Body["username"]
		if username == "" || msg.Type != data.MessageTypes.Username {
			err := data.WsError.InvalidRequestBody
			msg := fiber.Map{
				"type":    data.WsErrorType(err),
				"message": data.WsErrorMessage(err),
			}
			_ = wh.msg.InvalidMessage(msg)
			username = ""
		}

		if u, _ := wh.us.GetUserByName(username); u != nil || msg.Type != data.MessageTypes.Username {
			err := data.UserStoreError.DuplicateUsername
			msg := fiber.Map{
				"type":    err,
				"message": data.UserStoreErrMessage(err),
			}
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

	user := data.CreateUser(username, "android", data.WithConnection(c))
	_, err = wh.us.AddUser(user)
	if err != nil {
		return err
	}

	defer wh.wsDefer(user)

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
		wh.WsRequestHandler(msg)
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

func (wh *WebsocketHandler) wsDefer(user *data.User) {
	log.Println("DBG", "defer")
	_, err := wh.us.RemoveUser(user.Username)
	if err != nil {
		return
	}

	err = wh.msg.Broadcast(&data.Message{Type: data.MessageTypes.PeerLeft, User: user, From: user.Username})
	if err != nil {
		log.Println("ERR", err)
		return
	}
}

func (wh *WebsocketHandler) WsRequestHandler(msg *data.Message) {
	switch msg.Type {
	case data.MessageTypes.Offer,
		data.MessageTypes.Answer,
		data.MessageTypes.PeerLeft,
		data.MessageTypes.PeerUpdated,
		data.MessageTypes.PeerJoined,
		data.MessageTypes.NewIceCandidate:
		err := wh.msg.Broadcast(msg)
		if err != nil {
			log.Println("ERR ws handler", err)
			return
		}
		break
	default:
		log.Println("ERR -->> invalid request")
		err := wh.msg.InvalidMessage(fiber.Map{})
		if err != nil {
			return
		}
		break
	}
}
