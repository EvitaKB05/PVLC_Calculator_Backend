package main

import (
	"lab1/internal/app/ds"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Прямая строка подключения (временно для отладки)
	dsnString := "host=localhost port=5432 user=postgres password=password dbname=lung_capacity_db sslmode=disable"
	log.Printf("Подключаемся к БД: %s", dsnString)

	// Подключаемся к БД с явными настройками
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsnString,
		PreferSimpleProtocol: true, // важно для Windows
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Проверяем подключение
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Ошибка получения DB объекта: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Ошибка ping базы данных: %v", err)
	}

	log.Println("Успешно подключились к БД")

	// Выполняем миграции
	err = db.AutoMigrate(
		&ds.User{},
		&ds.Service{},
		&ds.Order{},
		&ds.OrderService{},
	)
	if err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v", err)
	}

	log.Println("Миграции успешно выполнены!")
}
