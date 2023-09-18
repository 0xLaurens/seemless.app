package handlers

import (
	"encoding/json"
	"fmt"
	ws "github.com/fasthttp/websocket"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"laurensdrop/data"
	"laurensdrop/store"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupTestApp() *fiber.App {
	// init in memory store and hub
	s := store.NewUserStoreInMemory()
	hub := CreateHub(s)
	go hub.Run()

	app := fiber.New()

	app.Use("/ws", WSUpgrader)
	app.Use("/ws", websocket.New(func(conn *websocket.Conn) {
		WSHandler(conn, hub)
	}))

	return app
}

func runTestApp() (*fiber.App, chan struct{}) {
	app := setupTestApp()
	done := make(chan struct{})
	go func() {
		err := app.Listen(":5421")
		if err != nil {
			return
		}
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		close(done)
	}()

	return app, done
}

func TestInvalidWebsocketRequestShouldReturnUpgradeRequired(t *testing.T) {
	app := setupTestApp()
	req := httptest.NewRequest(http.MethodTrace, "/ws", nil)
	resp, err := app.Test(req)
	if err != nil {
		return
	}

	if resp.StatusCode != fiber.StatusUpgradeRequired {
		t.Errorf("got %v : expected %v", resp.StatusCode, 426)
	}
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	if string(body) == "Upgrade required" {
		t.Errorf("got %s : expected %s", string(body), "Upgrade Required")
	}
}

func TestWSUpgraderShouldUpgradeToWSConnection(t *testing.T) {
	app, done := runTestApp()
	defer app.Shutdown()
	<-done

	url := "ws://localhost:5421/ws"
	conn, resp, err := ws.DefaultDialer.Dial(url, nil)
	defer conn.Close()
	if err != nil {
		return
	}
	fmt.Println(resp)

	assert.Equal(t, 101, resp.StatusCode)
	assert.Equal(t, "websocket", resp.Header.Get("Upgrade"))

	err = app.Shutdown()
	if err != nil {
		t.Fatal("TEST ERR -->> failed to shutdown server", err)
	}
}

func TestWsHandlerInvalidJsonToReturnInvalidRequest(t *testing.T) {
	app, done := runTestApp()
	defer app.Shutdown()
	<-done

	url := "ws://localhost:5421/ws"
	conn, resp, err := ws.DefaultDialer.Dial(url, nil)
	defer conn.Close()
	if err != nil {
		return
	}
	assert.Equal(t, 101, resp.StatusCode)
	assert.Equal(t, "websocket", resp.Header.Get("Upgrade"))

	message := "Hello World!"
	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		t.Fatal("TEST ERR -->> failed to send message", err)
	}

	_, res, err := conn.ReadMessage()
	if err != nil {
		t.Fatal("TEST ERR -->> failed to recv message", err)
	}

	errRes := fiber.NewError(int(data.WsError.InvalidRequestBody))
	expected, err := json.Marshal(errRes)
	if err != nil {
		log.Println("TEST ERR -->> failed to create json from err", err)
		return
	}

	assert.Equal(t, expected, res)

	err = app.Shutdown()
	if err != nil {
		t.Fatal("TEST ERR -->> failed to shutdown server", err)
	}
}

func TestWsHandlerInvalidStatusToReturnInvalidRequest(t *testing.T) {
	app, done := runTestApp()
	defer app.Shutdown()
	<-done

	url := "ws://localhost:5421/ws"
	conn, resp, err := ws.DefaultDialer.Dial(url, nil)
	defer conn.Close()
	if err != nil {
		return
	}
	assert.Equal(t, 101, resp.StatusCode)
	assert.Equal(t, "websocket", resp.Header.Get("Upgrade"))

	message, err := json.Marshal(data.Request{
		Type: "NonExistentType",
	})

	if err != nil {
		t.Fatal("TEST ERR -->> failed to marshal request", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		t.Fatal("TEST ERR -->> failed to send message", err)
	}

	_, res, err := conn.ReadMessage()
	if err != nil {
		t.Fatal("TEST ERR -->> failed to recv message", err)
	}

	errRes := fiber.NewError(int(data.WsError.InvalidRequestBody))
	expected, err := json.Marshal(errRes)
	if err != nil {
		log.Println("TEST ERR -->> failed to create json from err", err)
		return
	}

	assert.Equal(t, expected, res)

	err = app.Shutdown()
	if err != nil {
		t.Fatal("TEST ERR -->> failed to shutdown server", err)
	}
}
