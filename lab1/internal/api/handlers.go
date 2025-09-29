package api

import (
	"lab1/internal/app/ds"
	"lab1/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type API struct {
	repo *repository.Repository
}

func NewAPI(repo *repository.Repository) *API {
	return &API{repo: repo}
}

// Вспомогательные функции
func (a *API) successResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ds.APIResponse{
		Status: "success",
		Data:   data,
	})
}

func (a *API) errorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ds.APIResponse{
		Status:  "error",
		Message: message,
	})
}

// Домен: Услуги (Calculations)

// GET /api/calculations - список услуг с фильтрацией
func (a *API) GetCalculations(c *gin.Context) {
	var filter ds.CalculationFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные параметры фильтрации")
		return
	}

	calculations, err := a.repo.GetCalculationsWithFilter(filter)
	if err != nil {
		logrus.Error("Error getting calculations: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка получения услуг")
		return
	}

	var response []ds.CalculationResponse
	for _, calc := range calculations {
		response = append(response, ds.CalculationResponse{
			ID:          calc.ID,
			Title:       calc.Title,
			Description: calc.Description,
			Formula:     calc.Formula,
			ImageURL:    calc.ImageURL,
			Category:    calc.Category,
			Gender:      calc.Gender,
			MinAge:      calc.MinAge,
			MaxAge:      calc.MaxAge,
			IsActive:    calc.IsActive,
		})
	}

	a.successResponse(c, response)
}

// GET /api/calculations/:id - одна услуга
func (a *API) GetCalculation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID услуги")
		return
	}

	calculation, err := a.repo.GetCalculationByID(uint(id))
	if err != nil {
		logrus.Error("Error getting calculation: ", err)
		a.errorResponse(c, http.StatusNotFound, "Услуга не найдена")
		return
	}

	response := ds.CalculationResponse{
		ID:          calculation.ID,
		Title:       calculation.Title,
		Description: calculation.Description,
		Formula:     calculation.Formula,
		ImageURL:    calculation.ImageURL,
		Category:    calculation.Category,
		Gender:      calculation.Gender,
		MinAge:      calculation.MinAge,
		MaxAge:      calculation.MaxAge,
		IsActive:    calculation.IsActive,
	}

	a.successResponse(c, response)
}

// POST /api/calculations - добавление услуги
func (a *API) CreateCalculation(c *gin.Context) {
	var request ds.CreateCalculationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	calculation := ds.Calculation{
		Title:       request.Title,
		Description: request.Description,
		Formula:     request.Formula,
		Category:    request.Category,
		Gender:      request.Gender,
		MinAge:      request.MinAge,
		MaxAge:      request.MaxAge,
		IsActive:    true,
	}

	if err := a.repo.CreateCalculation(&calculation); err != nil {
		logrus.Error("Error creating calculation: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка создания услуги")
		return
	}

	a.successResponse(c, gin.H{"id": calculation.ID})
}

// PUT /api/calculations/:id - изменение услуги
func (a *API) UpdateCalculation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID услуги")
		return
	}

	var request ds.UpdateCalculationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	calculation, err := a.repo.GetCalculationByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Услуга не найдена")
		return
	}

	// Обновляем только переданные поля
	if request.Title != "" {
		calculation.Title = request.Title
	}
	if request.Description != "" {
		calculation.Description = request.Description
	}
	if request.Formula != "" {
		calculation.Formula = request.Formula
	}
	if request.Category != "" {
		calculation.Category = request.Category
	}
	if request.Gender != "" {
		calculation.Gender = request.Gender
	}
	if request.MinAge > 0 {
		calculation.MinAge = request.MinAge
	}
	if request.MaxAge > 0 {
		calculation.MaxAge = request.MaxAge
	}
	if request.IsActive != nil {
		calculation.IsActive = *request.IsActive
	}

	if err := a.repo.UpdateCalculation(&calculation); err != nil {
		logrus.Error("Error updating calculation: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка обновления услуги")
		return
	}

	a.successResponse(c, gin.H{"message": "Услуга успешно обновлена"})
}

// DELETE /api/calculations/:id - удаление услуги
func (a *API) DeleteCalculation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID услуги")
		return
	}

	// Получаем услугу для проверки существования и получения image_url
	calculation, err := a.repo.GetCalculationByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Услуга не найдена")
		return
	}

	// Удаляем изображение из MinIO если оно есть
	if calculation.ImageURL != "" {
		if err := a.repo.DeleteImageFromMinIO(calculation.ImageURL); err != nil {
			logrus.Warn("Error deleting image from MinIO: ", err)
		}
	}

	// Удаляем услугу из БД
	if err := a.repo.DeleteCalculation(uint(id)); err != nil {
		logrus.Error("Error deleting calculation: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка удаления услуги")
		return
	}

	a.successResponse(c, gin.H{"message": "Услуга успешно удалена"})
}

// POST /api/calculations/:id/image - добавление изображения
func (a *API) UploadCalculationImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID услуги")
		return
	}

	// Проверяем существование услуги
	calculation, err := a.repo.GetCalculationByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Услуга не найдена")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Файл изображения обязателен")
		return
	}

	// Загружаем изображение в MinIO
	imageURL, err := a.repo.UploadImageToMinIO(file, uint(id))
	if err != nil {
		logrus.Error("Error uploading image to MinIO: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка загрузки изображения")
		return
	}

	// Удаляем старое изображение если было
	if calculation.ImageURL != "" {
		if err := a.repo.DeleteImageFromMinIO(calculation.ImageURL); err != nil {
			logrus.Warn("Error deleting old image from MinIO: ", err)
		}
	}

	// Обновляем image_url в БД
	calculation.ImageURL = imageURL
	if err := a.repo.UpdateCalculation(&calculation); err != nil {
		logrus.Error("Error updating calculation image: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка обновления услуги")
		return
	}

	a.successResponse(c, gin.H{
		"message":   "Изображение успешно загружено",
		"image_url": imageURL,
	})
}

// POST /api/calculations/:id/add-to-cart - добавление в заявку-черновик
func (a *API) AddToCart(c *gin.Context) {
	idStr := c.Param("id")
	calculationID, err := strconv.Atoi(idStr)
	if err != nil || calculationID <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID услуги")
		return
	}

	// Фиксированный пользователь (как требуется в задании)
	userID := uint(1)

	// Создаем или получаем черновик
	card, err := a.repo.GetOrCreateDraftCard(userID)
	if err != nil {
		logrus.Error("Error getting draft card: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка доступа к корзине")
		return
	}

	// Добавляем услугу в заявку
	if err := a.repo.AddCalculationToCard(card.ID, uint(calculationID)); err != nil {
		logrus.Error("Error adding to cart: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка добавления в корзину")
		return
	}

	a.successResponse(c, gin.H{
		"message": "Услуга добавлена в заявку",
		"card_id": card.ID,
	})
}
