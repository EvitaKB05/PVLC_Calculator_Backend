package api

import (
	"database/sql"
	"lab1/internal/app/config"
	"lab1/internal/app/handler"
	"lab1/internal/app/repository"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StartServer(cfg *config.Config, repo *repository.Repository, db *sql.DB) {
	log.Println("Запуск сервера...")

	//создаем хэндлер
	handler := handler.NewHandler(repo, db)

	// гин
	r := gin.Default()

	// шаблоны и статика
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	// руты
	registerRoutes(r, handler)

	// сервер запуск
	serverAddress := cfg.ServiceHost + ":" + strconv.Itoa(cfg.ServicePort)
	if err := r.Run(serverAddress); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func registerRoutes(r *gin.Engine, handler *handler.Handler) {
	// GET
	r.GET("/", handler.GetServices)
	r.GET("/services", handler.GetServices)
	r.GET("/service/:id", handler.GetService)
	r.GET("/order", handler.GetOrder) // корзина

	// POST
	r.POST("/order/add", handler.AddServiceToOrder)  // добавить
	r.POST("/order/:id/delete", handler.DeleteOrder) // удалить
}
