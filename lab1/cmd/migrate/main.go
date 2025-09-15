package main

import (
	"lab1/internal/app/ds"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Загружаем переменные окружения из файла .env
	_ = godotenv.Load()

	// Подключаемся к БД
	db, err := gorm.Open(postgres.Open("host=localhost port=5432 user=simpleuser password=12345 dbname=lung_db sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Выполняем миграции (создание таблиц)
	err = db.AutoMigrate(
		&ds.User{},
		&ds.Calculation{},
		&ds.MedicalCard{},
		&ds.CardCalculation{},
	)
	if err != nil {
		panic("cant migrate db: " + err.Error())
	}

	println("Миграции успешно выполнены!")
}
