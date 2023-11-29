package web

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"laurensdrop/internal/ports"
)

type App struct {
	fiber *fiber.App
	port  int
	ws    ports.MessageHandler
}

func NewApp(ws ports.MessageHandler, opts ...AppOptions) *App {
	s := &App{
		fiber: fiber.New(),
		port:  3000,
	}

	for _, applyOption := range opts {
		applyOption(s)
	}

	s.initAppRoutes(ws)

	return s
}

func (a *App) Run() error {
	return a.fiber.Listen(fmt.Sprintf(":%d", a.port))
}

func (a *App) Close() error {
	return a.fiber.Shutdown()
}
