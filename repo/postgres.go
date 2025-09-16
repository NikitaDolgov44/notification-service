package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
	"notification-service/config"
)

func NewPostgresDB(cfg *config.PostgreSQLConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&connect_timeout=%d",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.ConnectTimeout)

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(time.Duration(cfg.IdleTimeout) * time.Second)

	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return db, nil
}

