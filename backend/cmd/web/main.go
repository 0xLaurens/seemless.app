package main

import (
	"laurensdrop/internal/adapters/primary/web"
	"laurensdrop/internal/adapters/primary/web/handlers"
	"laurensdrop/internal/adapters/secondary/repo"
	"laurensdrop/internal/core/services"
	"log"
)

func main() {
	userRepo := repo.NewUserRepoInMemory()
	userService := services.NewUserService(userRepo)

	codeRepo := repo.NewCodeRepoInMemory()
	codeService := services.NewCodeService(codeRepo)

	roomRepo := repo.NewRoomRepoInMemory()
	roomService := services.NewRoomService(roomRepo, codeService)

	wh := handlers.NewWebsocketHandler(userService, roomService)

	app := web.NewApp(wh)

	log.Fatal(app.Run())
}
