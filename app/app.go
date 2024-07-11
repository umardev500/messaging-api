package app

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/routes"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) error {
	app := fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
		},
	)

	router := routes.NewRouter(app)
	router.Initialize()

	ch := make(chan error, 1)
	go func() {
		port := os.Getenv("PORT")
		addr := ":" + port
		log.Info().Msgf("🔥 Listening on %s", port)
		ch <- app.Listen(addr)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		app.Shutdown()
	}

	return nil
}
