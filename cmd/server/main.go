package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/LigeronAhill/gottha/internal/controllers/home"
	"github.com/LigeronAhill/gottha/internal/controllers/sitemap"
	"github.com/LigeronAhill/gottha/pkg/config"
	"github.com/LigeronAhill/gottha/pkg/database"
	"github.com/LigeronAhill/gottha/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"
)

type Counter struct {
	Count int
}

func main() {
	ctx := context.Background()
	customLogger := logger.Init(slog.LevelDebug)

	cfg, err := config.New("", nil)
	if err != nil {
		log.Fatal(err)
	}
	pool, err := database.GetPool(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err = database.Migrate(pool); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Use(slogfiber.New(customLogger))
	app.Use(recover.New())
	app.Static("/", "./public")

	home.Serve(ctx, app)

	sitemap.Serve(ctx, app, cfg.GetString("host"))
	log.Fatal(app.Listen(":3000"))
}
