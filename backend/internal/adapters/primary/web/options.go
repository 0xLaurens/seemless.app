package web

import "github.com/gofiber/fiber/v2"

type AppOptions func(a *App)

func WithPort(port int) AppOptions {
	return func(a *App) {
		a.port = port
	}
}

func WithFiberConf(conf fiber.Config) AppOptions {
	return func(a *App) {
		a.fiber = fiber.New(conf)
	}
}
