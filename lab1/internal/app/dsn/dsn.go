package dsn

import (
	"fmt"
	"os"
)

// FromEnv собирает строку подключения к PostgreSQL из переменных окружения
func FromEnv() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	// Формат: host=localhost port=5432 user=user password=pass dbname=dbname sslmode=disable
	return fmt.Sprintf("host=localhost port=5432 user=simpleuser password=12345 dbname=lung_db sslmode=disable",
		host, port, user, pass, dbname)
}
