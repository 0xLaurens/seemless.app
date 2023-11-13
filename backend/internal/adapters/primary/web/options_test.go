package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"laurensdrop/internal/adapters/primary/web/handlers"
	"laurensdrop/internal/adapters/secondary/repo"
	"laurensdrop/internal/core/services"
	"testing"
)

func TestWithPortToEqual3333(t *testing.T) {
	port := 3333
	withPort := WithPort(port)
	ur := repo.NewUserRepoInMemory()
	us := services.NewUserService(ur)
	wh := handlers.NewWebsocketHandler(us)

	app := NewApp(wh, withPort)
	assert.Equal(t, port, app.port)
}

func TestWithPortToEqual4343(t *testing.T) {
	port := 4343
	withPort := WithPort(port)
	ur := repo.NewUserRepoInMemory()
	us := services.NewUserService(ur)
	wh := handlers.NewWebsocketHandler(us)

	app := NewApp(wh, withPort)
	assert.Equal(t, port, app.port)
}

func TestWithPortDefaultPort3000(t *testing.T) {
	port := 3000
	ur := repo.NewUserRepoInMemory()
	us := services.NewUserService(ur)
	wh := handlers.NewWebsocketHandler(us)

	app := NewApp(wh)
	assert.Equal(t, port, app.port)
}

func TestWithFiberConfBodyLimit(t *testing.T) {
	ur := repo.NewUserRepoInMemory()
	us := services.NewUserService(ur)
	wh := handlers.NewWebsocketHandler(us)

	conf := fiber.Config{BodyLimit: 3999}
	withFiberConf := WithFiberConf(conf)

	app := NewApp(wh, withFiberConf)
	assert.Equal(t, conf.BodyLimit, app.fiber.Config().BodyLimit)
}

func TestWithFiberConfETag(t *testing.T) {
	ur := repo.NewUserRepoInMemory()
	us := services.NewUserService(ur)
	wh := handlers.NewWebsocketHandler(us)

	conf := fiber.Config{ETag: false}
	withFiberConf := WithFiberConf(conf)

	app := NewApp(wh, withFiberConf)
	assert.Equal(t, conf.ETag, app.fiber.Config().ETag)
}
