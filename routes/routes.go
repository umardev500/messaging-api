package routes

import "github.com/gofiber/fiber/v2"

type Routes struct {
	app *fiber.App
}

func NewRouter(app *fiber.App) *Routes {
	return &Routes{
		app: app,
	}
}

func (r *Routes) Initialize() {
	r.Api()
}
