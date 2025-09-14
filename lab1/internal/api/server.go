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
		logrus.Error("Ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()

	// шаблоны
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	// маршруты запросы
	r.GET("/", handler.GetServices)
	r.GET("/services", handler.GetServices)
	r.GET("/service/:id", handler.GetService)
	r.GET("/calculation", handler.GetCalculation)

	r.Run()
	log.Println("Server down")
}
