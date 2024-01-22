package secondary

import (
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/ports"
)

type WebsocketMessageValidator struct {
	users ports.UserService
}

func NewWebsocketMessageValidator(users ports.UserService) *WebsocketMessageValidator {
	return &WebsocketMessageValidator{users: users}
}

func (v *WebsocketMessageValidator) ValidateMessageOrigin(msg *data.Message) error {
	//if msg.From == "" {
	//	return nil
	//}
	//
	//user, err := v.users.GetUserByName(msg.From)
	//if err != nil {
	//	return err
	//}
	//
	//if user == nil || msg.Conn == nil || user.RemoteAddr != msg.Conn.RemoteAddr().(*net.TCPAddr) {
	//	err := msg.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "Message Spoofing"))
	//	if err != nil {
	//		return err
	//	}
	//
	//	err = msg.Conn.Close()
	//	if err != nil {
	//		return err
	//	}
	//	log.Printf("DBG -->> User %s Spoofed a message\n", user.Username)
	//
	//	return errors.New("the message was spoofed")
	//}

	return nil
}
