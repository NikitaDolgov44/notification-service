package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"notification-service/config"
	"notification-service/internal/kafka"
	"notification-service/internal/service"
	"notification-service/repo"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	config.Init("../config", "application")

	// 1. Открываем *sql.DB (goose требует именно *sql.DB)
	sqlDB, err := sql.Open("postgres", buildDSN(config.GlobalConfig.Postgres))
	if err != nil {
		log.Fatalf("open postgres: %v", err)
	}
	defer sqlDB.Close()

	// 2. Накатываем миграции
	if err := migrate(sqlDB); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// 3. Оборачиваем в sqlx
	db := sqlx.NewDb(sqlDB, "postgres")

	// 4. Обычный код сервиса
	notificationRepo := repo.NewNotificationRepo(db)
	notificationService := service.NewNotificationService(notificationRepo)

	consumer := kafka.NewNotificationConsumer(
		[]string{"localhost:9092"},
		"notification-service-group",
		notificationService,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := consumer.Consume(ctx); err != nil {
			log.Printf("consumer stopped: %v", err)
		}
		cancel()
	}()

	<-ctx.Done()
	_ = consumer.Close()
	log.Println("shutdown complete")
}

// migrate запускает goose Up
func migrate(db *sql.DB) error {
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, "migrations")
}

// buildDSN собирает строку подключения из конфига
func buildDSN(cfg *config.PostgreSQLConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}