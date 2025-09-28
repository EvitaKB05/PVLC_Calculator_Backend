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

	// GET
	r.GET("/", handler.GetDjelPatients)
	r.GET("/djel_patients", handler.GetDjelPatients)
	r.GET("/djel_patient/:id", handler.GetDjelPatient)
	r.GET("/djel_request/:id", handler.GetDjelRequest)

	//  POST
	r.POST("/djel_patient/:id/add", handler.AddToCart)
	r.POST("/djel_request/delete", handler.DeleteCart)

	r.Run()
	log.Println("Server down")
}
