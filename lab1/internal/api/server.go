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
	// ДОБАВЬТЕ ЭТОТ БЛОК - инициализация MinIO
	log.Println("Initializing MinIO bucket...")
	if err := repo.InitMinIOBucket(); err != nil {
		logrus.Warn("MinIO bucket initialization failed: ", err)
	} else {
		log.Println("MinIO bucket initialized successfully")
	}

	handler := handler.NewHandler(repo)
	api := NewAPI(repo)

	r := gin.Default()

	// HTML routes (сохраняем существующие)
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	// GET
	r.GET("/", handler.GetDjelPatients)
	r.GET("/djel_patients", handler.GetDjelPatients)
	r.GET("/djel_patient/:id", handler.GetDjelPatient)
	r.GET("/djel_request/:id", handler.GetDjelRequest)

	// POST
	r.POST("/djel_patient/:id/add", handler.AddToCart)
	r.POST("/djel_request/delete", handler.DeleteCart)

	// API Routes (добавляем новые)
	apiGroup := r.Group("/api")
	{
		// Calculations domain
		calculations := apiGroup.Group("/calculations")
		{
			calculations.GET("", api.GetCalculations)
			calculations.GET("/:id", api.GetCalculation)
			calculations.POST("", api.CreateCalculation)
			calculations.PUT("/:id", api.UpdateCalculation)
			calculations.DELETE("/:id", api.DeleteCalculation)
			calculations.POST("/:id/image", api.UploadCalculationImage)
			calculations.POST("/:id/add-to-cart", api.AddToCart)
		}

		// Medical Cards domain
		medicalCards := apiGroup.Group("/medical-cards")
		{
			medicalCards.GET("", api.GetMedicalCards)
			medicalCards.GET("/:id", api.GetMedicalCard)
			medicalCards.PUT("/:id", api.UpdateMedicalCard)
			medicalCards.PUT("/:id/finalize", api.FinalizeMedicalCard)
			medicalCards.PUT("/:id/complete", api.CompleteMedicalCard)
			medicalCards.DELETE("/:id", api.DeleteMedicalCard)
		}

		// Cart domain
		cart := apiGroup.Group("/cart")
		{
			cart.GET("/icon", api.GetCartIcon)
		}

		// Card Calculations domain (M-M)
		cardCalculations := apiGroup.Group("/card-calculations")
		{
			cardCalculations.DELETE("", api.DeleteCardCalculation)
			cardCalculations.PUT("", api.UpdateCardCalculation)
		}

		// Users domain
		users := apiGroup.Group("/users")
		{
			users.POST("/register", api.RegisterUser)
			users.GET("/profile", api.GetUserProfile)
			users.PUT("/profile", api.UpdateUserProfile)
			users.POST("/login", api.Login)
			users.POST("/logout", api.Logout)
		}
	}

	r.Run()
	log.Println("Server down")
}
