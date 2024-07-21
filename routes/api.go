package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/umardev500/messaging-api/config"
	"github.com/umardev500/messaging-api/injections"
	"github.com/umardev500/messaging-api/middlewares"
)

func (r *Routes) Api() {
	app := r.app.Group("api")
	app.Use(cors.New())
	app.Use("/static", filesystem.New(filesystem.Config{
		Root: http.Dir(config.GetConfig().Upload.Path),
	}))

	// Connection
	conn := config.NewPgx()
	var injects = injections.NewInject(conn, app)

	// Chat
	injects.Chat()

	// Auth
	injects.Auth()

	// Message
	injects.Message()

	// Upload
	uploadRoute := app.Group("upload")
	uploadRoute.Post("/", middlewares.UploadMiddleware(config.GetConfig().Upload.Path))
	uploadRoute.Put("/", middlewares.UpdateUploadMiddleware(config.GetConfig().Upload.Path))
}
