package main

import (
	"database/sql"
	"log"

	"lab1/internal/api"
	"lab1/internal/app/config"
	"lab1/internal/app/dsn"
	"lab1/internal/app/repository"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("Application start!...")

	// подгружаем конфиг
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// бд подключение
	db, err := sql.Open("postgres", dsn.FromEnv())
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// репозиторий
	repo, err := repository.New(dsn.FromEnv())
	if err != nil {
		log.Fatalf("Ошибка создания репозитория: %v", err)
	}

	// запуск
	api.StartServer(cfg, repo, db)
	log.Println("Application terminated!")
}
