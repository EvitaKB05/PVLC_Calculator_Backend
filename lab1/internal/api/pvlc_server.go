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
	r.GET("/", handler.GetPvlcPatients)
	r.GET("/pvlc_patients", handler.GetPvlcPatients)
	r.GET("/pvlc_patient/:id", handler.GetPvlcPatient)
	r.GET("/pvlc_med_calc/:id", handler.GetPvlcMedCalc)

	// POST
	r.POST("/pvlc_patient/:id/add", handler.AddPvlcMedFormulaToCart)
	r.POST("/pvlc_med_calc/delete", handler.DeletePvlcMedCart)
	// API Routes
	apiGroup := r.Group("/api")
	{
		// Pvlc Med Formulas domain
		formulas := apiGroup.Group("/pvlc-med-formulas")
		{
			formulas.GET("", api.GetPvlcMedFormulas)
			formulas.GET("/:id", api.GetPvlcMedFormula)
			formulas.POST("", api.CreatePvlcMedFormula)
			formulas.PUT("/:id", api.UpdatePvlcMedFormula)
			formulas.DELETE("/:id", api.DeletePvlcMedFormula)
			formulas.POST("/:id/image", api.UploadPvlcMedFormulaImage)
			formulas.POST("/:id/add-to-cart", api.AddPvlcMedFormulaToCart)
		}

		// Pvlc Med Cards domain
		cards := apiGroup.Group("/pvlc-med-cards")
		{
			cards.GET("", api.GetPvlcMedCards)
			cards.GET("/:id", api.GetPvlcMedCard)
			cards.PUT("/:id", api.UpdatePvlcMedCard)
			cards.PUT("/:id/finalize", api.FinalizePvlcMedCard)
			cards.PUT("/:id/complete", api.CompletePvlcMedCard)
			cards.DELETE("/:id", api.DeletePvlcMedCard)
		}

		// Cart domain
		cart := apiGroup.Group("/cart")
		{
			cart.GET("/icon", api.GetCartIcon)
		}

		// Med Mm Pvlc Calculations domain (M-M)
		calculations := apiGroup.Group("/med-mm-pvlc-calculations")
		{
			calculations.DELETE("", api.DeleteMedMmPvlcCalculation)
			calculations.PUT("", api.UpdateMedMmPvlcCalculation)
		}

		// Med Users domain
		users := apiGroup.Group("/med-users")
		{
			users.POST("/register", api.RegisterMedUser)
			users.GET("/profile", api.GetMedUserProfile)
			users.PUT("/profile", api.UpdateMedUserProfile)
			users.POST("/login", api.LoginMedUser)
			users.POST("/logout", api.LogoutMedUser)
		}
	}

	r.Run()
	log.Println("Server down")
}
