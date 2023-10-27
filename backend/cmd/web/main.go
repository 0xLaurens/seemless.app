package main

import (
	"laurensdrop/internal/adapters/primary/web"
	"laurensdrop/internal/adapters/primary/web/handlers"
	"laurensdrop/internal/adapters/secondary/repo"
	"laurensdrop/internal/core/services"
	"log"
)

func main() {
	ur := repo.NewUserRepoInMemory()
	us := services.NewUserService(ur)
	wh := handlers.NewWebsocketHandler(us)

	app := web.NewApp(wh)

	log.Fatal(app.Run())
}
