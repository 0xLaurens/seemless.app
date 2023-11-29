package web

import (
	"errors"
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/stretchr/testify/assert"
	"laurensdrop/internal/adapters/primary/web/handlers"
	"laurensdrop/internal/adapters/secondary/repo"
	"laurensdrop/internal/core/services"
	"net/http"
	"testing"
	"time"
)

func setupTestApp() *App {
	ur := repo.NewUserRepoInMemory()
	us := services.NewUserService(ur)
	wh := handlers.NewWebsocketHandler(us)

	port := WithPort(4543)
	app := NewApp(wh, port)
	return app
}

func runTestApp() (*App, chan struct{}) {
	app := setupTestApp()
	done := make(chan struct{})

	go func() {
		err := app.Run()
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

func TestRouteRoot(t *testing.T) {
	app, done := runTestApp()
	defer app.Close()
	<-done

	url := "http://127.0.0.1:4543/"
	get, err := http.Get(url)
	if err != nil {
		return
	}

	assert.Equal(t, "200 OK", get.Status)
}

func TestWsShouldUpgradeRequest(t *testing.T) {
	app, done := runTestApp()
	defer app.Close()
	<-done

	url := "ws://127.0.0.1:4543/ws"
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	assert.Equal(t, err, nil)
	defer conn.Close()
	if err != nil {
		return
	}
	fmt.Println(resp)

	assert.Equal(t, 101, resp.StatusCode)
	assert.Equal(t, "websocket", resp.Header.Get("Upgrade"))
}

func TestHttpWsShouldNotUpgradeRequest(t *testing.T) {
	app, done := runTestApp()
	defer app.Close()
	<-done

	url := "http://127.0.0.1:4543/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		assert.Equal(t, errors.New("malformed ws or wss URL"), err)
	}
}
