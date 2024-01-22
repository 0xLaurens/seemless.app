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
	rr := repo.NewRoomRepoInMemory()
	cr := repo.NewCodeRepoInMemory()
	cs := services.NewCodeService(cr)
	rs := services.NewRoomService(rr, cs)
	wh := NewWebsocketHandler(us, rs)

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
func TestConnectionProvidesUserWithDisplayName(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	conn, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	_, res, err := conn.ReadMessage()
	assert.NoError(t, err)

	response := data.Message{}
	err = utils.MapJsonToStruct(res, &response)
	assert.NoError(t, err)

	assert.Equal(t, data.DisplayName, response.Type)
}

// UC4 - connected peers
// UC5 - account alias
func TestUserShouldReceiveUniqueAlias(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	conn, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	_, displayName, err := conn.ReadMessage()
	resName := data.Message{}
	err = utils.MapJsonToStruct(displayName, &resName)
	assert.NoError(t, err)

	assert.NotNil(t, resName.User.Username)
}

// UC4 - connected peers
// UC5 - account alias
func TestUserCanChangeAlias(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	conn, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	// user receives default username
	_, _, err = conn.ReadMessage()
	assert.NoError(t, err)

	// user receives room code
	_, _, err = conn.ReadMessage()
	assert.NoError(t, err)

	// user receives peers in room
	_, _, err = conn.ReadMessage()
	assert.NoError(t, err)

	// user receives peer joined message in room
	_, _, err = conn.ReadMessage()
	assert.NoError(t, err)

	// user send new alias
	alias := "newAlias"
	newAlias := data.Message{
		Type:        data.ChangeDisplayName,
		DisplayName: alias,
	}
	err = conn.WriteJSON(newAlias)

	_, displayName, err := conn.ReadMessage()
	resName := data.Message{}
	err = utils.MapJsonToStruct(displayName, &resName)
	assert.NoError(t, err)

	assert.Equal(t, data.ChangeDisplayName, newAlias.Type)
	assert.Equal(t, alias, resName.User.Username)
}

// UC4 - Peers
func TestPeersJoinMessageSentOnlyToNewestUser(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	user1, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	// user1 display name
	_, _, _ = user1.ReadMessage()

	user2, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	// user2 display name
	_, _, _ = user2.ReadMessage()
	_, _, _ = user2.ReadMessage()
	_, peers, err := user2.ReadMessage()
	assert.NoError(t, err)
	peersMessage := data.Message{}
	err = utils.MapJsonToStruct(peers, &peersMessage)
	assert.NoError(t, err)
	assert.Equal(t, data.Peers, peersMessage.Type)
}

// UC4 - Peers
func TestPeersLeaveMessageSentAfterConnectionCloses(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	user1, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	// user1 display name
	_, displayName, err := user1.ReadMessage()
	assert.NoError(t, err)
	displayNameMessage := data.Message{}
	err = utils.MapJsonToStruct(displayName, &displayNameMessage)
	assert.NoError(t, err)

	// user1 room code
	_, roomCode, err := user1.ReadMessage()
	assert.NoError(t, err)
	roomCodeMessage := data.Message{}
	err = utils.MapJsonToStruct(roomCode, &roomCodeMessage)
	assert.NoError(t, err)

	// user1 peers
	_, _, _ = user1.ReadMessage()

	// user1 joined
	_, _, _ = user1.ReadMessage()

	user2, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)

	err = user2.WriteJSON(data.Message{Type: data.RoomJoin, RoomCode: roomCodeMessage.RoomCode})
	assert.NoError(t, err)

	_, _, err = user1.ReadMessage()

	time.Sleep(100 * time.Millisecond)
	err = user2.Close()
	assert.NoError(t, err)

	_, leave, err := user1.ReadMessage()
	assert.NoError(t, err)
	leaveMessage := data.Message{}
	err = utils.MapJsonToStruct(leave, &leaveMessage)
	assert.NoError(t, err)

	assert.Equal(t, data.PeerLeft, leaveMessage.Type)
}

func TestSelectiveForwardingToUser(t *testing.T) {
	app := setupTestApp()
	defer app.Shutdown()

	fred, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)
	_, fredDisplayMessage, _ := fred.ReadMessage()
	fredDisplayNameMessage := &data.Message{}
	utils.MapJsonToStruct(fredDisplayMessage, &fredDisplayNameMessage)

	_, _, _ = fred.ReadMessage()
	_, _, _ = fred.ReadMessage()

	joe, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)

	_, _, _ = joe.ReadMessage()
	_, _, _ = joe.ReadMessage()
	_, _, _ = joe.ReadMessage()

	harry, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
	assert.NoError(t, err)

	_, harryDisplayName, _ := harry.ReadMessage()
	harryDisplayNameMessage := &data.Message{}
	utils.MapJsonToStruct(harryDisplayName, &harryDisplayNameMessage)
	_, _, _ = harry.ReadMessage()
	_, _, _ = harry.ReadMessage()
	_, _, _ = harry.ReadMessage()

	mockIce := data.Message{
		Type: data.NewIceCandidate,
		Candidate: &data.RTCIceCandidate{
			Candidate:        "SDP PROTOCOL",
			SdpMid:           "MID",
			SdpMLineIndex:    0,
			UsernameFragment: "u89432",
		},
		From:   fredDisplayNameMessage.User.Username,
		Target: harryDisplayNameMessage.User.Username,
	}
	err = fred.WriteJSON(mockIce)
	assert.NoError(t, err)

	_, iceOffer, err := harry.ReadMessage()
	iceOfferMessage := &data.Message{}
	err = utils.MapJsonToStruct(iceOffer, &iceOfferMessage)
	assert.Equal(t, mockIce, *iceOfferMessage)

	//harry joined
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

//func TestFakeFromMessage(t *testing.T) {
//	app := setupTestApp()
//	defer app.Shutdown()
//
//	fred, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
//	assert.NoError(t, err)
//	_, fredDisplayName, _ := fred.ReadMessage()
//	fredDisplayNameMessage := &data.Message{}
//	utils.MapJsonToStruct(fredDisplayName, &fredDisplayNameMessage)
//	_, _, _ = fred.ReadMessage()
//	_, _, _ = fred.ReadMessage()
//
//	joe, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
//	assert.NoError(t, err)
//
//	_, joeDisplayName, _ := joe.ReadMessage()
//	joeDisplayNameMessage := &data.Message{}
//	utils.MapJsonToStruct(joeDisplayName, &joeDisplayNameMessage)
//
//	_, _, _ = joe.ReadMessage()
//	_, _, _ = joe.ReadMessage()
//
//	harry, _, err := ws.DefaultDialer.Dial(TestUrl, nil)
//	assert.NoError(t, err)
//
//	_, harryDisplayName, _ := harry.ReadMessage()
//	harryDisplayNameMessage := &data.Message{}
//	utils.MapJsonToStruct(harryDisplayName, &harryDisplayNameMessage)
//	_, _, _ = harry.ReadMessage()
//
//	mockIce := data.Message{
//		Type: data.NewIceCandidate,
//		Candidate: &data.RTCIceCandidate{
//			Candidate:        "SDP PROTOCOL",
//			SdpMid:           "MID",
//			SdpMLineIndex:    0,
//			UsernameFragment: "u89432",
//		},
//		From:   joeDisplayNameMessage.User.Username,
//		Target: harryDisplayNameMessage.User.Username,
//	}
//	err = fred.WriteJSON(mockIce)
//	assert.NoError(t, err)
//
//	msgChannel := make(chan *data.Message)
//	go func() {
//		_, message, _ := harry.ReadMessage()
//		msg := data.Message{}
//		_ = utils.MapJsonToStruct(message, &msg)
//		msgChannel <- &msg
//	}()
//
//	select {
//	case <-time.After(500 * time.Millisecond):
//		log.Println("DBG -->>", "Harry did not receive spoofed message")
//		break
//	case message := <-msgChannel:
//		if message.Type == data.NewIceCandidate {
//			t.Errorf("Received %v whilst harry shouldn't receive a spoofed message", message)
//		}
//	}
//
//	_, _, _ = fred.ReadMessage()
//
//	_, message, err := fred.ReadMessage()
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	assert.Equal(t, "close 1008 (policy violation): Message Spoofing", string(message))
//}
