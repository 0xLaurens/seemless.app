package handlers

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"laurensdrop/internal/adapters/secondary/repo"
	"laurensdrop/internal/core/services"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

const TestPort = 6611

func setupTestApp() *fiber.App {
	// init in memory user repo & other services
	ur := repo.NewUserRepoInMemory()
	us := services.NewUserService(ur)
	ws := NewWebsocketHandler(us)

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hi mom!")
	})

	app.Use("/ws", ws.UpgradeWebsocket)
	app.Use("/ws", websocket.New(func(conn *websocket.Conn) {
		_ = ws.HandleWebsocket(conn)
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

func TestInvalidWebsocketRequestShouldReturnUpgradeRequired(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	req := httptest.NewRequest(http.MethodTrace, "/ws", nil)
	resp, err := app.Test(req)
	fmt.Println(resp, err)
	assert.NoError(t, err)

	assert.Equal(t, "426 Upgrade Required", resp.Status)
}

//func TestInvalidWebsocketRequestShouldReturnUpgradeRequired(t *testing.T) {
//	app := setupTestApp()
//	req := httptest.NewRequest(http.MethodTrace, "/ws", nil)
//	resp, err := app.Test(req)
//	if err != nil {
//		return
//	}
//
//	if resp.StatusCode != fiber.StatusUpgradeRequired {
//		t.Errorf("got %v : expected %v", resp.StatusCode, 426)
//	}
//	body, err := io.ReadAll(resp.Body)
//	fmt.Println(string(body))
//	if string(body) == "Upgrade required" {
//		t.Errorf("got %s : expected %s", string(body), "Upgrade Required")
//	}
//}
//
//func TestWSUpgraderShouldUpgradeToWSConnection(t *testing.T) {
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
//	fmt.Println(resp)
//
//	assert.Equal(t, 101, resp.StatusCode)
//	assert.Equal(t, "websocket", resp.Header.Get("Upgrade"))
//
//	err = app.Shutdown()
//	if err != nil {
//		t.Fatal("TEST ERR -->> failed to shutdown server", err)
//	}
//}
//
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
