package main

import (
	"lab1/internal/app/ds"
	"lab1/internal/app/dsn" // не забудь!

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// подгружаем из энвиромент
	_ = godotenv.Load()

	// гет стринг из энв
	dsnString := dsn.FromEnv()

	// подключение к бдшке
	db, err := gorm.Open(postgres.Open(dsnString), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// создаем миграции
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
