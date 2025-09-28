package api

import (
	"lab1/internal/app/handler"
	"lab1/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Fatal("Ошибка инициализации репозитория: ", err)
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()

	// шаблоны
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	// маршруты GET запросы
	r.GET("/", handler.GetServices)
	r.GET("/services", handler.GetServices)
	r.GET("/service/:id", handler.GetService)
	r.GET("/calculation/:id", handler.GetCalculation) // ИЗМЕНЯЕМ НА С ID

	// маршруты POST запросы
	r.POST("/service/:id/add", handler.AddToCart)     // Добавление в корзину
	r.POST("/calculation/delete", handler.DeleteCart) // Удаление корзины

	r.Run()
	log.Println("Server down")
}
