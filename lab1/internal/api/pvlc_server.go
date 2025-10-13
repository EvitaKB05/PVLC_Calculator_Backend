package api

import (
	"lab1/internal/app/handler"
	"lab1/internal/app/repository"
	"lab1/internal/auth"
	"lab1/internal/redis"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// StartServer запускает HTTP сервер
// ПОЛНОСТЬЮ ПЕРЕПИСАН ДЛЯ ЛАБОРАТОРНОЙ РАБОТЫ 4
func StartServer() {
	log.Println("Starting server")

	// Инициализация репозитория
	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Fatal("Ошибка инициализации репозитория: ", err)
	}

	// Инициализация Redis клиента
	redisClient, err := redis.NewRedisClient("localhost", 6379, "", 0)
	if err != nil {
		logrus.Fatal("Ошибка инициализации Redis: ", err)
	}
	defer redisClient.Close()

	// Инициализация аутентификации
	auth.InitAuth(redisClient)

	log.Println("Initializing MinIO bucket...")
	if err := repo.InitMinIOBucket(); err != nil {
		logrus.Warn("MinIO bucket initialization failed: ", err)
	} else {
		log.Println("MinIO bucket initialized successfully")
	}

	handler := handler.NewHandler(repo)

	// ИСПРАВЛЕНО: передаем Redis клиент при создании API
	api := NewAPI(repo, redisClient)

	r := gin.Default()

	// Swagger документация - ДОБАВЛЕНО ДЛЯ ЛР4
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// HTML routes (сохраняем существующие)
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	// GET routes для HTML интерфейса
	r.GET("/", handler.GetPvlcPatients)
	r.GET("/pvlc_patients", handler.GetPvlcPatients)
	r.GET("/pvlc_patient/:id", handler.GetPvlcPatient)
	r.GET("/pvlc_med_calc/:id", handler.GetPvlcMedCalc)

	// POST routes для HTML интерфейса
	r.POST("/pvlc_patient/:id/add", handler.AddPvlcMedFormulaToCart)
	r.POST("/pvlc_med_calc/delete", handler.DeletePvlcMedCart)

	// API Routes - ИСПРАВЛЕНО: убираем дублирование /api в группе
	apiGroup := r.Group("/api")
	{
		// Public routes (не требуют аутентификации)
		public := apiGroup.Group("")
		{
			public.POST("/auth/login", api.Login)                       // Аутентификация
			public.GET("/pvlc-med-formulas", api.GetPvlcMedFormulas)    // Список формул
			public.GET("/pvlc-med-formulas/:id", api.GetPvlcMedFormula) // Конкретная формула
		}

		// Auth required routes (требуют аутентификации)
		authRequired := apiGroup.Group("")
		authRequired.Use(auth.RequireAuth()) // Требует валидный JWT токен
		{
			// Auth routes
			authRequired.POST("/auth/logout", api.Logout)     // Выход
			authRequired.GET("/auth/profile", api.GetProfile) // Профиль

			// Cart routes
			authRequired.GET("/cart/icon", api.GetCartIcon) // Иконка корзины

			// Pvlc Med Formulas routes
			authRequired.POST("/pvlc-med-formulas/:id/add-to-cart", api.AddPvlcMedFormulaToCart) // Добавление в корзину

			// Pvlc Med Cards routes (user specific)
			authRequired.GET("/pvlc-med-cards", api.GetPvlcMedCards)                  // Список заявок пользователя
			authRequired.GET("/pvlc-med-cards/:id", api.GetPvlcMedCard)               // Конкретная заявка
			authRequired.PUT("/pvlc-med-cards/:id", api.UpdatePvlcMedCard)            // Обновление заявки
			authRequired.PUT("/pvlc-med-cards/:id/finalize", api.FinalizePvlcMedCard) // Формирование заявки
			authRequired.DELETE("/pvlc-med-cards/:id", api.DeletePvlcMedCard)         // Удаление заявки

			// Med Mm Pvlc Calculations routes
			authRequired.DELETE("/med-mm-pvlc-calculations", api.DeleteMedMmPvlcCalculation) // Удаление расчета
			authRequired.PUT("/med-mm-pvlc-calculations", api.UpdateMedMmPvlcCalculation)    // Обновление расчета

			// Med Users routes
			authRequired.GET("/med-users/profile", api.GetMedUserProfile)    // Профиль пользователя
			authRequired.PUT("/med-users/profile", api.UpdateMedUserProfile) // Обновление профиля
			authRequired.POST("/med-users/logout", api.LogoutMedUser)        // Выход (старый endpoint)
		}

		// Moderator only routes (требуют права модератора)
		moderator := apiGroup.Group("")
		moderator.Use(auth.RequireModerator()) // Требует права модератора
		{
			// Pvlc Med Formulas management
			moderator.POST("/pvlc-med-formulas", api.CreatePvlcMedFormula)                // Создание формулы
			moderator.PUT("/pvlc-med-formulas/:id", api.UpdatePvlcMedFormula)             // Обновление формулы
			moderator.DELETE("/pvlc-med-formulas/:id", api.DeletePvlcMedFormula)          // Удаление формулы
			moderator.POST("/pvlc-med-formulas/:id/image", api.UploadPvlcMedFormulaImage) // Загрузка изображения

			// Pvlc Med Cards moderation
			moderator.PUT("/pvlc-med-cards/:id/complete", api.CompletePvlcMedCard) // Завершение/отклонение заявки

			// User management
			moderator.POST("/med-users/register", api.RegisterMedUser) // Регистрация пользователя
		}
	}

	log.Println("Server starting on :8080")
	log.Println("Swagger UI available at: http://localhost:8080/swagger/index.html")
	r.Run(":8080")
	log.Println("Server down")
}
