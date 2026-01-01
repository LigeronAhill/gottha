package home

import (
	"context"

	"github.com/LigeronAhill/gottha/pkg/adaptor"
	"github.com/LigeronAhill/gottha/templates/components"
	"github.com/LigeronAhill/gottha/templates/pages"
	"github.com/gofiber/fiber/v2"
)

func Serve(ctx context.Context, app *fiber.App) {
	count := 0
	app.Post("/count", func(c *fiber.Ctx) error {
		count++
		return adaptor.Render(c, components.Counter(count))
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return adaptor.Render(c, pages.Home(count))
	})
}
