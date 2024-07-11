package routes

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/umardev500/messaging-api/app/handler"
	"github.com/umardev500/messaging-api/app/service"
)

func (r *Routes) Api() {
	app := r.app.Group("api")
	app.Use(cors.New())

	chatService := service.NewChatService()
	chat := handler.NewChatHandler(chatService)
	chatRoute := app.Group("chat")
	chatRoute.Get("/list", chat.WsChatList())
	chatRoute.Get("/:room", chat.WsChat())
	chatRoute.Post("/new", chat.PushNewChat)
}
