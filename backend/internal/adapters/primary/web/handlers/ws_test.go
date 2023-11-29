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
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

	app := fiber.New(fiber.Config{DisableStartupMessage: true, IdleTimeout: 600000})

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

func TestUserSelectUsernameWrongBodyAndWrongType(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	conn, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	// UsernamePrompt
	_, _, err = conn.ReadMessage()
	assert.NoError(t, err)

	joinMessage := data.Message{
		Type: data.MessageTypes.NewIceCandidate,
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
	expectedUser := data.CreateUser("Joe", "")
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

// UC5 - Username
func TestDuplicateUsernameErrorTwoUsersSameName(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	fred, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	// FRED GOES THROUGH THE JOIN PROCESS
	joinRoomHelper(fred, "Fred")

	fred2, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	_, _, _ = fred2.ReadMessage()

	joinMessage := data.Message{
		Type: data.MessageTypes.Username,
		Body: make(map[string]string),
	}
	joinMessage.Body["username"] = "Fred"
	err = fred2.WriteJSON(joinMessage)

	_, duplicate, err := fred2.ReadMessage()
	assert.NoError(t, err)
	duplicateMessage := data.Message{}
	err = utils.MapJsonToStruct(duplicate, &duplicateMessage)
	assert.NoError(t, err)

	assert.Equal(t, data.MessageTypes.DuplicateUsername, duplicateMessage.Type)
	assert.Equal(t, "username not unique", duplicateMessage.Body["message"])
}

// UC5 - Username
func TestDuplicateUsernameErrorTwoUsersDifferentCapitalization(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	fred, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	// FRED GOES THROUGH THE JOIN PROCESS
	joinRoomHelper(fred, "Fred")

	fred2, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	_, _, _ = fred2.ReadMessage()

	joinMessage := data.Message{
		Type: data.MessageTypes.Username,
		Body: make(map[string]string),
	}
	joinMessage.Body["username"] = "fred"
	err = fred2.WriteJSON(joinMessage)

	_, duplicate, err := fred2.ReadMessage()
	assert.NoError(t, err)
	duplicateMessage := data.Message{}
	err = utils.MapJsonToStruct(duplicate, &duplicateMessage)
	assert.NoError(t, err)

	assert.Equal(t, data.MessageTypes.DuplicateUsername, duplicateMessage.Type)
	assert.Equal(t, "username not unique", duplicateMessage.Body["message"])
}

func TestSelectiveForwardingToUser(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	fred, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	joinRoomHelper(fred, "Fred")

	joe, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	joinRoomHelper(joe, "Joe")

	harry, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	joinRoomHelper(harry, "Harry")

	mockIce := data.Message{
		Type: data.MessageTypes.NewIceCandidate,
		Candidate: &data.RTCIceCandidate{
			Candidate:        "SDP PROTOCOL",
			SdpMid:           "MID",
			SdpMLineIndex:    0,
			UsernameFragment: "u89432",
		},
		From:   "Fred",
		Target: "Harry",
	}
	err = fred.WriteJSON(mockIce)
	assert.NoError(t, err)

	_, iceOffer, err := harry.ReadMessage()
	iceOfferMessage := &data.Message{}
	err = utils.MapJsonToStruct(iceOffer, &iceOfferMessage)
	assert.Equal(t, mockIce, *iceOfferMessage)

	//peer joined
	_, _, _ = joe.ReadMessage()

	msgChannel := make(chan string)
	go func() {
		_, message, _ := joe.ReadMessage()
		msgChannel <- string(message)
	}()

	select {
	case <-time.After(500 * time.Millisecond):
		break
	case message := <-msgChannel:
		t.Errorf("Received %s whilst joe should have received anything", message)
	}
}

func TestFakeFromMessage(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	fred, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	joinRoomHelper(fred, "Fred")

	joe, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	joinRoomHelper(joe, "Joe")

	harry, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	joinRoomHelper(harry, "Harry")

	mockIce := data.Message{
		Type: data.MessageTypes.NewIceCandidate,
		Candidate: &data.RTCIceCandidate{
			Candidate:        "SDP PROTOCOL",
			SdpMid:           "MID",
			SdpMLineIndex:    0,
			UsernameFragment: "u89432",
		},
		From:   "Fred",
		Target: "Harry",
	}
	err = joe.WriteJSON(mockIce)
	assert.NoError(t, err)

	msgChannel := make(chan *data.Message)
	go func() {
		_, message, _ := harry.ReadMessage()
		msg := data.Message{}
		_ = utils.MapJsonToStruct(message, &msg)
		msgChannel <- &msg
	}()

	select {
	case <-time.After(500 * time.Millisecond):
		log.Println("DBG -->>", "Harry did not receive spoofed message")
		break
	case message := <-msgChannel:
		if message.Type == data.MessageTypes.NewIceCandidate {
			t.Errorf("Received %v whilst harry shouldn't receive a spoofed message", message)
		}
	}

	_, _, _ = joe.ReadMessage()
	_, message, err := joe.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	assert.Equal(t, "close 1008 (policy violation): Message Spoofing", string(message))
}