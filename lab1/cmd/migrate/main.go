package main

import (
	"lab1/internal/app/ds"
	"lab1/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Загружаем переменные окружения из файла .env
	_ = godotenv.Load()

	// Получаем DSN строку из переменных окружения
	dsnString := dsn.FromEnv()

	// Подключаемся к БД
	db, err := gorm.Open(postgres.Open(dsnString), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Выполняем миграции (создание таблиц)
	err = db.AutoMigrate(
		&ds.MedUser{},
		&ds.PvlcMedFormula{},
		&ds.PvlcMedCard{},
		&ds.MedMmPvlcCalculation{},
	)
	if err != nil {
		panic("cant migrate db: " + err.Error())
	}

	println("Миграции успешно выполнены!")
}
