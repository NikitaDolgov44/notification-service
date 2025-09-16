package main

import (
	"context"
	"log"

	"notification-service/config"
	"notification-service/repo"
)

func main() {
	config.Init("./config", "config") // ваш уже готовый loader

	db, err := repo.NewPostgresDB(config.GlobalConfig.Postgres)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer db.Close()

	repository := repo.NewNotificationRepo(db)

	nn, err := repository.FindAllByPage(context.Background(), repo.Page{Offset: 0, Limit: 10})
	if err != nil {
		log.Fatalf("repo: %v", err)
	}
	log.Printf("notifications: %+v", nn)
}