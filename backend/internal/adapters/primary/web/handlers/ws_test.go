package handlers

import (
	"fmt"
	ws "github.com/fasthttp/websocket"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"laurensdrop/internal/adapters/secondary/repo"
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/core/services"
	"laurensdrop/internal/core/utils"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	TestPort = 6611
)

var TestUrl = fmt.Sprintf("ws://localhost:%d/ws", TestPort)

/*
 * Helper function for creating the app
 */
func setupTestApp() *fiber.App {
	// init in memory user repo & other services
	ur := repo.NewUserRepoInMemory()
	us := services.NewUserService(ur)
	wh := NewWebsocketHandler(us)

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Use("/ws", wh.UpgradeWebsocket)
	app.Use("/ws", websocket.New(func(conn *websocket.Conn) {
		_ = wh.HandleWebsocket(conn)
	}))

	go app.Listen(fmt.Sprintf(":%d", TestPort))

	readyCh := make(chan struct{})

	go func() {
		for {
			address := fmt.Sprintf("localhost:%d", TestPort)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				continue
			}

			if conn != nil {
				readyCh <- struct{}{}
				conn.Close()
				break
			}
		}
	}()

	<-readyCh

	return app
}

/*
 * Helper function for the joining the room
 */
func joinRoomHelper(conn *ws.Conn, username string) {
	_, _, _ = conn.ReadMessage()
	joinMessage := data.Message{
		Type: data.MessageTypes.Username,
		Body: make(map[string]string),
	}
	joinMessage.Body["username"] = username
	conn.WriteJSON(joinMessage)

	// Peers message
	_, _, _ = conn.ReadMessage()
	// UserJoined
	_, _, _ = conn.ReadMessage()
}

func TestInvalidWebsocketRequestShouldReturnUpgradeRequired(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	req := httptest.NewRequest(http.MethodTrace, "/ws", nil)
	res, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, "426 Upgrade Required", res.Status)
}

func TestShouldUpgradeWebsocketConnection(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	conn, resp, err := ws.DefaultDialer.Dial(TestUrl, nil)
	defer conn.Close()

	assert.NoError(t, err)
	assert.Equal(t, "101 Switching Protocols", resp.Status)
	assert.NotNil(t, conn)
}

// UC5 - account alias
func TestConnectionRequestsUsername(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	conn, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	_, res, err := conn.ReadMessage()
	assert.NoError(t, err)

	fmt.Println(string(res))

	response := data.Message{}
	err = utils.MapJsonToStruct(res, &response)
	assert.NoError(t, err)

	expectedJoinMessage := data.Message{
		Type: data.MessageTypes.UsernamePrompt,
		Body: make(map[string]string),
	}
	expectedJoinMessage.Body["message"] = "Please provide a username"
	assert.NoError(t, err)
	assert.Equal(t, expectedJoinMessage.Body["message"], response.Body["message"])
	assert.Equal(t, expectedJoinMessage.Type, response.Type)
}

// UC4 - connected peers
// UC5 - account alias
func TestUserSelectUsername(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	conn, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	// UsernamePrompt
	_, _, err = conn.ReadMessage()
	assert.NoError(t, err)

	joinMessage := data.Message{
		Type: data.MessageTypes.Username,
		Body: make(map[string]string),
	}
	joinMessage.Body["username"] = "Johny"

	err = conn.WriteJSON(joinMessage)
	assert.NoError(t, err)

	_, peers, err := conn.ReadMessage()
	response := data.Message{}
	err = utils.MapJsonToStruct(peers, &response)
	assert.NoError(t, err)

	// server sends other users
	assert.Equal(t, data.MessageTypes.Peers, response.Type)

	// server sends message to show others johny has connected
	_, joinedMessage, err := conn.ReadMessage()
	responseJoinMessage := data.Message{}
	err = utils.MapJsonToStruct(joinedMessage, &responseJoinMessage)
	assert.Equal(t, data.MessageTypes.PeerJoined, responseJoinMessage.Type)
}

// UC5 - account alias
func TestUserSelectUsernameInvalidJsonRequest(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	conn, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	// UsernamePrompt
	_, _, err = conn.ReadMessage()
	assert.NoError(t, err)

	joinMessage := "Hi mom!"
	err = conn.WriteMessage(ws.TextMessage, []byte(joinMessage))
	assert.NoError(t, err)

	_, invalidRequest, err := conn.ReadMessage()
	invalidMessage := data.Message{}
	err = utils.MapJsonToStruct(invalidRequest, &invalidMessage)

	assert.Equal(t, data.MessageTypes.InvalidMessage, invalidMessage.Type)
	assert.Equal(t, "invalid request body", invalidMessage.Body["message"])
}

// UC5 - account alias
func TestUserSelectUsernameJsonWrongType(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	conn, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	// UsernamePrompt
	_, _, err = conn.ReadMessage()
	assert.NoError(t, err)

	joinMessage := data.Message{
		Type: data.MessageTypes.Answer,
		Body: make(map[string]string),
	}
	joinMessage.Body["username"] = "Johny"
	err = conn.WriteJSON(joinMessage)
	assert.NoError(t, err)

	_, invalidRequest, err := conn.ReadMessage()
	invalidMessage := data.Message{}
	err = utils.MapJsonToStruct(invalidRequest, &invalidMessage)

	assert.Equal(t, data.MessageTypes.InvalidMessage, invalidMessage.Type)
	assert.Equal(t, "invalid request body", invalidMessage.Body["message"])
}

func TestUserSelectUsernameWrongBody(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	conn, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	// UsernamePrompt
	_, _, err = conn.ReadMessage()
	assert.NoError(t, err)

	joinMessage := data.Message{
		Type: data.MessageTypes.Username,
		Body: make(map[string]string),
	}
	err = conn.WriteJSON(joinMessage)
	assert.NoError(t, err)

	_, invalidRequest, err := conn.ReadMessage()
	invalidMessage := data.Message{}
	err = utils.MapJsonToStruct(invalidRequest, &invalidMessage)

	assert.Equal(t, data.MessageTypes.InvalidMessage, invalidMessage.Type)
	assert.Equal(t, "invalid request body", invalidMessage.Body["message"])
}

// UC4 - Peers
func TestPeersJoinMessageSentOnlyToNewestUser(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	fred, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	// FRED GOES THROUGH THE JOIN PROCESS
	joinRoomHelper(fred, "Fred")

	joe, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)

	// username prompt
	_, _, _ = joe.ReadMessage()

	joinMessage := data.Message{
		Type: data.MessageTypes.Username,
		Body: make(map[string]string),
	}
	joinMessage.Body["username"] = "Joe"
	err = joe.WriteJSON(joinMessage)

	_, peers, err := joe.ReadMessage()
	peersMessage := data.Message{}
	err = utils.MapJsonToStruct(peers, &peersMessage)
	assert.Equal(t, data.MessageTypes.Peers, peersMessage.Type)
	assert.Equal(t, "Fred", peersMessage.Users[0].Username)

	_, joinJoe, err := fred.ReadMessage()
	joinJoeMessage := data.Message{}
	err = utils.MapJsonToStruct(joinJoe, &joinJoeMessage)
	expectedUser := data.CreateUser("Joe", "android")
	assert.Equal(t, expectedUser, joinJoeMessage.User)
	assert.Equal(t, data.MessageTypes.PeerJoined, joinJoeMessage.Type)
}

// UC4 - Peers
func TestPeersLeaveMessageSentAfterConnectionCloses(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	fred, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	// FRED GOES THROUGH THE JOIN PROCESS
	joinRoomHelper(fred, "Fred")

	joe, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	joinRoomHelper(joe, "Joe")

	// peer joined message
	_, _, _ = fred.ReadMessage()

	err = fred.Close()
	assert.NoError(t, err)

	_, leave, err := joe.ReadMessage()
	assert.NoError(t, err)
	leaveMessage := data.Message{}
	err = utils.MapJsonToStruct(leave, &leaveMessage)
	assert.NoError(t, err)

	assert.Equal(t, data.MessageTypes.PeerLeft, leaveMessage.Type)
	assert.Equal(t, "Fred", leaveMessage.From)
	assert.Equal(t, "Fred", leaveMessage.User.Username)
}

//func TestWsHandlerInvalidJsonToReturnInvalidRequest(t *testing.T) {
//	app, done := runTestApp()
//	defer app.Shutdown()
//	<-done
//
//	url := "ws://localhost:5421/ws"
//	conn, resp, err := ws.DefaultDialer.Dial(url, nil)
//	defer conn.Close()
//	if err != nil {
//		return
//	}
//	assert.Equal(t, 101, resp.StatusCode)
//	assert.Equal(t, "websocket", resp.Header.Get("Upgrade"))
//
//	joinHelper(conn)
//
//	message := "Hello World!"
//	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
//	if err != nil {
//		t.Fatal("TEST ERR -->> failed to send message", err)
//	}
//
//	_, res, err := conn.ReadMessage()
//	if err != nil {
//		t.Fatal("TEST ERR -->> failed to recv message", err)
//	}
//
//	errRes := fiber.NewError(int(data.WsError.InvalidRequestBody))
//	expected, err := json.Marshal(errRes)
//	if err != nil {
//		log.Println("TEST ERR -->> failed to create json from err", err)
//		return
//	}
//
//	assert.Equal(t, expected, res)
//
//	err = app.Shutdown()
//	if err != nil {
//		t.Fatal("TEST ERR -->> failed to shutdown server", err)
//	}
//}
//
//func TestWsHandlerPromptsForUsername(t *testing.T) {
//	app, done := runTestApp()
//	defer app.Shutdown()
//	<-done
//
//	url := "ws://localhost:5421/ws"
//	conn, resp, err := ws.DefaultDialer.Dial(url, nil)
//	defer conn.Close()
//	if err != nil {
//		return
//	}
//	assert.Equal(t, 101, resp.StatusCode)
//	assert.Equal(t, "websocket", resp.Header.Get("Upgrade"))
//
//	expected := fiber.Map{
//		"message": "provide a username",
//		"type":    "UsernamePrompt",
//	}
//	var res fiber.Map
//	err = conn.ReadJSON(&res)
//
//	assert.Equal(t, expected, res)
//
//	err = app.Shutdown()
//	if err != nil {
//		t.Fatal("TEST ERR -->> failed to shutdown server", err)
//	}
//}
//
//func TestWsHandlerPromptsInvalidRequestTypeCorrectBody(t *testing.T) {
//	app, done := runTestApp()
//	defer app.Shutdown()
//	<-done
//
//	url := "ws://localhost:5421/ws"
//	conn, resp, err := ws.DefaultDialer.Dial(url, nil)
//	if err != nil {
//		return
//	}
//	defer conn.Close()
//
//	assert.Equal(t, 101, resp.StatusCode)
//	assert.Equal(t, "websocket", resp.Header.Get("Upgrade"))
//
//	expected := fiber.Map{
//		"message": "provide a username",
//		"type":    "UsernamePrompt",
//	}
//	var res fiber.Map
//	err = conn.ReadJSON(&res)
//
//	assert.NoError(t, err)
//	assert.Equal(t, expected, res)
//
//	req := fiber.Map{
//		"type": "asfd",
//		"body": fiber.Map{
//			"username": "user",
//		},
//	}
//
//	err = conn.WriteJSON(req)
//	assert.NoError(t, err)
//
//	err = conn.ReadJSON(&res)
//	assert.NoError(t, err)
//
//	expected = fiber.Map{
//		"type":    data.WsErrorType(data.WsError.InvalidRequestBody),
//		"message": data.WsErrorMessage(data.WsError.InvalidRequestBody),
//	}
//
//	assert.Equal(t, expected, res)
//
//	err = app.Shutdown()
//	if err != nil {
//		t.Fatal("TEST ERR -->> failed to shutdown server", err)
//	}
//}
//
//func TestWsHandlerUsernameRequestValidTypeInvalidBody(t *testing.T) {
//	app, done := runTestApp()
//	defer app.Shutdown()
//	<-done
//
//	url := "ws://localhost:5421/ws"
//	conn, resp, err := ws.DefaultDialer.Dial(url, nil)
//	if err != nil {
//		return
//	}
//	defer conn.Close()
//
//	assert.Equal(t, 101, resp.StatusCode)
//	assert.Equal(t, "websocket", resp.Header.Get("Upgrade"))
//
//	expected := fiber.Map{
//		"message": "provide a username",
//		"type":    "UsernamePrompt",
//	}
//	var res fiber.Map
//	err = conn.ReadJSON(&res)
//
//	assert.NoError(t, err)
//	assert.Equal(t, expected, res)
//
//	req := fiber.Map{
//		"type": data.MessageTypes.Username,
//		"body": fiber.Map{
//			"asf": "username",
//		},
//	}
//
//	err = conn.WriteJSON(req)
//	assert.NoError(t, err)
//
//	err = conn.ReadJSON(&res)
//	assert.NoError(t, err)
//
//	expected = fiber.Map{
//		"type":    data.WsErrorType(data.WsError.InvalidRequestBody),
//		"message": data.WsErrorMessage(data.WsError.InvalidRequestBody),
//	}
//
//	assert.Equal(t, expected, res)
//
//	err = app.Shutdown()
//	if err != nil {
//		t.Fatal("TEST ERR -->> failed to shutdown server", err)
//	}
//}
//
//func TestWsHandlerUsernameShouldJoin(t *testing.T) {
//	app, done := runTestApp()
//	defer app.Shutdown()
//	<-done
//
//	url := "ws://localhost:5420/ws"
//	conn, resp, err := ws.DefaultDialer.Dial(url, nil)
//	if err != nil {
//		return
//	}
//	defer conn.Close()
//
//	assert.Equal(t, 100, resp.StatusCode)
//	assert.Equal(t, "websocket", resp.Header.Get("Upgrade"))
//
//	expected := fiber.Map{
//		"message": "provide a username",
//		"type":    "UsernamePrompt",
//	}
//	var res fiber.Map
//	err = conn.ReadJSON(&res)
//
//	assert.NoError(t, err)
//	assert.Equal(t, expected, res)
//
//	req := fiber.Map{
//		"type": "Username",
//		"body": fiber.Map{
//			"username": "user",
//		},
//	}
//
//	err = conn.WriteJSON(req)
//	assert.NoError(t, err)
//
//	err = conn.ReadJSON(&res)
//	assert.NoError(t, err)
//
//	expected = fiber.Map{
//		"type":     data.MessageTypes.PeerJoined,
//		"username": "user",
//	}
//
//	assert.Equal(t, expected, res)
//
//	err = app.Shutdown()
//	if err != nil {
//		t.Fatal("TEST ERR -->> failed to shutdown server", err)
//	}
//}
//
//func TestWsHandlerUserLeaveMessage(t *testing.T) {
//	app, done := runTestApp()
//	defer app.Shutdown()
//	<-done
//
//	url := "ws://localhost:5420/ws"
//	conn, resp, err := ws.DefaultDialer.Dial(url, nil)
//	if err != nil {
//		return
//	}
//	defer conn.Close()
//	conn2, _, err := ws.DefaultDialer.Dial(url, nil)
//	if err != nil {
//		return
//	}
//	defer conn2.Close()
//
//	assert.Equal(t, 100, resp.StatusCode)
//	assert.Equal(t, "websocket", resp.Header.Get("Upgrade"))
//
//	expected := fiber.Map{
//		"message": "provide a username",
//		"type":    "UsernamePrompt",
//	}
//	var res fiber.Map
//	err = conn.ReadJSON(&res)
//
//	assert.NoError(t, err)
//	assert.Equal(t, expected, res)
//
//	req := fiber.Map{
//		"type": "Username",
//		"body": fiber.Map{
//			"username": "user",
//		},
//	}
//
//	err = conn.WriteJSON(req)
//	assert.NoError(t, err)
//
//	err = conn.ReadJSON(&res)
//	assert.NoError(t, err)
//
//	expected = fiber.Map{
//		"type":     data.MessageTypes.PeerJoined,
//		"username": "user",
//	}
//	assert.Equal(t, expected, res)
//
//	req = fiber.Map{
//		"type": "Username",
//		"body": fiber.Map{
//			"username": "user2",
//		},
//	}
//
//	err = conn2.WriteJSON(req)
//	assert.NoError(t, err)
//
//	expected = fiber.Map{
//		"type":     data.MessageTypes.PeerLeft,
//		"username": "user",
//	}
//
//	err = conn2.ReadJSON(res)
//	if err != nil {
//		return
//	}
//
//	conn.Close()
//	assert.Equal(t, expected, res)
//
//	err = app.Shutdown()
//	if err != nil {
//		t.Fatal("TEST ERR -->> failed to shutdown server", err)
//	}
//}
//
//func TestWsHandlerUserShouldReturnDuplicateUsernameError(t *testing.T) {
//	app, done := runTestApp()
//	defer app.Shutdown()
//	<-done
//
//	url := "ws://localhost:5420/ws"
//	conn, resp, err := ws.DefaultDialer.Dial(url, nil)
//	if err != nil {
//		return
//	}
//	defer conn.Close()
//	conn2, _, err := ws.DefaultDialer.Dial(url, nil)
//	if err != nil {
//		return
//	}
//	defer conn2.Close()
//
//	assert.Equal(t, 100, resp.StatusCode)
//	assert.Equal(t, "websocket", resp.Header.Get("Upgrade"))
//
//	expected := fiber.Map{
//		"message": "provide a username",
//		"type":    "UsernamePrompt",
//	}
//	var res fiber.Map
//	err = conn.ReadJSON(&res)
//
//	assert.NoError(t, err)
//	assert.Equal(t, expected, res)
//
//	req := fiber.Map{
//		"type": "Username",
//		"body": fiber.Map{
//			"username": "user",
//		},
//	}
//
//	err = conn.WriteJSON(req)
//	assert.NoError(t, err)
//
//	err = conn.ReadJSON(&res)
//	assert.NoError(t, err)
//
//	expected = fiber.Map{
//		"type":     data.MessageTypes.PeerJoined,
//		"username": "user",
//	}
//	assert.Equal(t, expected, res)
//
//	err = conn2.WriteJSON(req)
//	assert.NoError(t, err)
//
//	duErr := data.UserStoreError.DuplicateUsername
//	expected = fiber.Map{
//		"type":    duErr,
//		"message": data.UserStoreErrMessage(duErr),
//	}
//	assert.Equal(t, expected, res)
//
//	err = app.Shutdown()
//	if err != nil {
//		t.Fatal("TEST ERR -->> failed to shutdown server", err)
//	}
//}
//
//func TestWsHandlerInvalidStatusToReturnInvalidRequest(t *testing.T) {
//	app, done := runTestApp()
//	defer app.Shutdown()
//	<-done
//
//	url := "ws://localhost:5421/ws"
//	conn, resp, err := ws.DefaultDialer.Dial(url, nil)
//	defer conn.Close()
//	if err != nil {
//		return
//	}
//	assert.Equal(t, 101, resp.StatusCode)
//	assert.Equal(t, "websocket", resp.Header.Get("Upgrade"))
//	joinHelper(conn)
//
//	message, err := json.Marshal(data.Message{
//		Type: "NonExistentType",
//	})
//
//	if err != nil {
//		t.Fatal("TEST ERR -->> failed to marshal request", err)
//	}
//
//	err = conn.WriteMessage(websocket.TextMessage, message)
//	if err != nil {
//		t.Fatal("TEST ERR -->> failed to send message", err)
//	}
//
//	_, res, err := conn.ReadMessage()
//	if err != nil {
//		t.Fatal("TEST ERR -->> failed to recv message", err)
//	}
//
//	errRes := fiber.NewError(int(data.WsError.InvalidRequestBody))
//	expected, err := json.Marshal(errRes)
//	if err != nil {
//		log.Println("TEST ERR -->> failed to create json from err", err)
//		return
//	}
//
//	assert.Equal(t, expected, res)
//
//	err = app.Shutdown()
//	if err != nil {
//		t.Fatal("TEST ERR -->> failed to shutdown server", err)
//	}
//}
//
//func TestWsHandlerRequestTypesShouldBroadcast(t *testing.T) {
//	app, done := runTestApp()
//	defer app.Shutdown()
//	<-done
//
//	url := "ws://localhost:5421/ws"
//	conn, resp, err := ws.DefaultDialer.Dial(url, nil)
//	defer conn.Close()
//	if err != nil {
//		return
//	}
//	assert.Equal(t, 101, resp.StatusCode)
//	assert.Equal(t, "websocket", resp.Header.Get("Upgrade"))
//	joinHelper(conn)
//
//	types := []data.MessageType{
//		data.MessageTypes.PeerJoined,
//		data.MessageTypes.PeerLeft,
//		data.MessageTypes.PeerUpdated,
//		data.MessageTypes.NewIceCandidate,
//		data.MessageTypes.Answer,
//		data.MessageTypes.Offer,
//	}
//	for _, request := range types {
//		message, err := json.Marshal(data.Message{
//			Type: request,
//		})
//		if err != nil {
//			t.Fatal("TEST ERR -->> failed to marshal request", err)
//		}
//		err = conn.WriteMessage(websocket.TextMessage, message)
//		if err != nil {
//			t.Fatal("TEST ERR -->> failed to send message", err)
//		}
//		_, res, err := conn.ReadMessage()
//		if err != nil {
//			t.Fatal("TEST ERR -->> failed to recv message", err)
//		}
//		assert.Equal(t, string(message), string(res))
//	}
//
//	err = app.Shutdown()
//	if err != nil {
//		t.Fatal("TEST ERR -->> failed to shutdown server", err)
//	}
//}
