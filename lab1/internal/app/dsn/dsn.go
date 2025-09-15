package dsn

import (
	"fmt"
	"os"
)

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
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "lung_capacity_db"
	}

	// Добавляем sslmode=disable и явно указываем md5
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable password_encryption=md5",
		host, port, user, password, dbname)
}
