package dsn

import (
	"fmt"
	"os"
)

// FromEnv собирает строку подключения к PostgreSQL из переменных окружения
func FromEnv() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "simpleuser"
	}

	pass := os.Getenv("DB_PASS")
	if pass == "" {
		pass = "12345"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "lung_db"
	}

	// ИСПРАВЛЕННАЯ СТРОКА: убрали лишние параметры
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname)
}
