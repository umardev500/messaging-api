package routes

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/umardev500/messaging-api/app/handler"
	"github.com/umardev500/messaging-api/app/repository"
	"github.com/umardev500/messaging-api/app/service"
	"github.com/umardev500/messaging-api/config"
)

func (r *Routes) Api() {
	app := r.app.Group("api")
	app.Use(cors.New())

	// Connection
	conn := config.NewPgx()

	participantRepository := repository.NewParticipantRepository(conn)
	chatRepository := repository.NewChatRepository(conn)
	chatService := service.NewChatService(chatRepository, participantRepository, conn)
	chat := handler.NewChatHandler(chatService)
	chatRoute := app.Group("chat")
	chatRoute.Get("/list", chat.WsChatList())
	chatRoute.Get("/:room", chat.WsChat())
	chatRoute.Post("/new", chat.PushNewChat)
}
