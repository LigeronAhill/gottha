package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func GetPool(ctx context.Context, cfg *viper.Viper) (*pgxpool.Pool, error) {
	dbURL := cfg.GetString("database.url")
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}
	var greeting string
	err = pool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Запрос к БД провален: %v\n", err)
		os.Exit(1)
	} else {
		slog.Info("Соединение с БД установлено")
	}
	return pool, nil
}
func Migrate(pool *pgxpool.Pool) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	db := sql.OpenDB(stdlib.GetPoolConnector(pool))
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	return nil
}
