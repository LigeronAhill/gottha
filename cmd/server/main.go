package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/LigeronAhill/gottha/internal/controllers/home"
	"github.com/LigeronAhill/gottha/internal/controllers/sitemap"
	"github.com/LigeronAhill/gottha/internal/middleware"
	"github.com/LigeronAhill/gottha/pkg/config"
	"github.com/LigeronAhill/gottha/pkg/database"
	"github.com/LigeronAhill/gottha/pkg/db"
	"github.com/LigeronAhill/gottha/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
	slogfiber "github.com/samber/slog-fiber"
)

type Counter struct {
	Count int
}

func main() {
	// Контекст
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Журнал
	customLogger := logger.Init(slog.LevelDebug)

	// Конфигурация
	cfg, err := config.New("", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Пул соединений с базой данных
	pool, err := database.GetPool(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	// Миграции базы данных
	if err = database.Migrate(pool); err != nil {
		log.Fatal(err)
	}

	// Сгенерированные запросы к базе данных
	querier := db.New(pool)
	version, err := querier.GetVersion(ctx)
	if err != nil {
		slog.Error("Не получилось получить версию")
	}
	slog.Info("Версия", slog.String("tag", version))

	sessionsDB := postgres.New(postgres.Config{
		DB:         pool,
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	})

	sessionsStorage := session.New(session.Config{
		Storage: sessionsDB,
	})

	host := cfg.GetString("host")
	port := cfg.GetInt("port")
	if port == 0 {
		port = 3000
	}
	addr := fmt.Sprintf("%s:%d", host, port)

	// Сервер fiber
	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// Промежуточное програмное обеспечение
	app.Use(requestid.New())
	app.Use(slogfiber.New(customLogger))
	app.Use(recover.New())
	app.Use(middleware.Auth(sessionsStorage))

	// Маршруты
	// Статичные файлы
	app.Static("/", "./public")

	// Основные маршруты
	home.Serve(ctx, app)
	sitemap.Serve(ctx, app, addr)

	// Запуск сервера
	slog.Info(fmt.Sprintf("🚀 Сервер запущен на %s", addr))
	log.Fatal(app.Listen(addr))
}
