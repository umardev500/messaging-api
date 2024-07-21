package injections

import (
	"github.com/gofiber/fiber/v2"
	"github.com/umardev500/messaging-api/app/handler"
	"github.com/umardev500/messaging-api/app/repository"
	"github.com/umardev500/messaging-api/app/service"
	"github.com/umardev500/messaging-api/config"
	"github.com/umardev500/messaging-api/middlewares"
)

type Inject struct {
	conn *config.PgxConfig
	app  fiber.Router
}

func NewInject(conn *config.PgxConfig, app fiber.Router) *Inject {
	return &Inject{
		conn: conn,
		app:  app,
	}
}

func (i *Inject) Auth() {
	authRoute := i.app.Group("auth")
	userRepository := repository.NewuserRepository(i.conn)
	authService := service.NewAuthService(userRepository)
	auth := handler.NewAuthHandler(authService)
	authRoute.Post("/login", auth.Login)
}

func (i *Inject) Chat() {
	participantRepository := repository.NewParticipantRepository(i.conn)
	chatRepository := repository.NewChatRepository(i.conn)
	chatService := service.NewChatService(chatRepository, participantRepository, i.conn)
	chat := handler.NewChatHandler(chatService)
	chatRoute := i.app.Group("chat")
	chatRoute.Get("/list", chat.WsChatList()) // ws
	chatRoute.Get("/chat_list", middlewares.CheckAuth, chat.GetChatList)
	chatRoute.Get("/:room", chat.WsChat()) // ws
	chatRoute.Post("/new", middlewares.CheckAuth, chat.CreateNewChat)
}

func (i *Inject) Message() {
	messageRoute := i.app.Group("message")
	messageRepository := repository.NewMessageRepository(i.conn)
	messageService := service.NewMessageService(messageRepository)
	message := handler.NewMessageHandler(messageService)
	messageRoute.Post("/:room", middlewares.CheckAuth, message.Create)
	messageRoute.Get("/:room", middlewares.CheckAuth, message.GetMessage)
}
