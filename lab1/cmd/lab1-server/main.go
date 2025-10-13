package main

import (
	"log"

	_ "lab1/docs" // Swagger docs - ДОБАВЛЕНО ДЛЯ ЛАБОРАТОРНОЙ РАБОТЫ 4
	"lab1/internal/api"
)

// @title Lung Capacity Calculation API
// @version 1.0
// @description API для расчета должной жизненной емкости легких (ДЖЕЛ)
// @description Лабораторная работа 4 - Добавление аутентификации и авторизации

// @contact.name API Support
// @contact.url http://localhost:8080
// @contact.email support@lungcalc.ru

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT токен для аутентификации. Формат: "Bearer {token}"

// main - точка входа приложения
// ОБНОВЛЕНО ДЛЯ ЛАБОРАТОРНОЙ РАБОТЫ 4 - добавлены Swagger аннотации
func main() {
	log.Println("Application start!")
	api.StartServer()
	log.Println("Application terminated!")
}
