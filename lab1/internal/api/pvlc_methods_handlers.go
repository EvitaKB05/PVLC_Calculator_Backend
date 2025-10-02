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
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (a *API) errorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}

// Домен: Формулы ДЖЕЛ (PvlcMedFormulas)

// GET /api/pvlc-med-formulas - список формул с фильтрацией
func (a *API) GetPvlcMedFormulas(c *gin.Context) {
	var filter ds.PvlcMedFormulaFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные параметры фильтрации")
		return
	}

	formulas, err := a.repo.GetPvlcMedFormulasWithFilter(filter)
	if err != nil {
		logrus.Error("Error getting pvlc med formulas: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка получения формул")
		return
	}

	var response []ds.PvlcMedFormulaResponse
	for _, formula := range formulas {
		response = append(response, ds.PvlcMedFormulaResponse{
			ID:          formula.ID,
			Title:       formula.Title,
			Description: formula.Description,
			Formula:     formula.Formula,
			ImageURL:    formula.ImageURL,
			Category:    formula.Category,
			Gender:      formula.Gender,
			MinAge:      formula.MinAge,
			MaxAge:      formula.MaxAge,
			IsActive:    formula.IsActive,
		})
	}

	a.successResponse(c, response)
}

// GET /api/pvlc-med-formulas/:id - одна формула
func (a *API) GetPvlcMedFormula(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID формулы")
		return
	}

	formula, err := a.repo.GetPvlcMedFormulaByID(uint(id))
	if err != nil {
		logrus.Error("Error getting pvlc med formula: ", err)
		a.errorResponse(c, http.StatusNotFound, "Формула не найдена")
		return
	}

	response := ds.PvlcMedFormulaResponse{
		ID:          formula.ID,
		Title:       formula.Title,
		Description: formula.Description,
		Formula:     formula.Formula,
		ImageURL:    formula.ImageURL,
		Category:    formula.Category,
		Gender:      formula.Gender,
		MinAge:      formula.MinAge,
		MaxAge:      formula.MaxAge,
		IsActive:    formula.IsActive,
	}

	a.successResponse(c, response)
}

// POST /api/pvlc-med-formulas - добавление формулы
func (a *API) CreatePvlcMedFormula(c *gin.Context) {
	var request ds.CreatePvlcMedFormulaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	formula := ds.PvlcMedFormula{
		Title:       request.Title,
		Description: request.Description,
		Formula:     request.Formula,
		Category:    request.Category,
		Gender:      request.Gender,
		MinAge:      request.MinAge,
		MaxAge:      request.MaxAge,
		IsActive:    true,
	}

	if err := a.repo.CreatePvlcMedFormula(&formula); err != nil {
		logrus.Error("Error creating pvlc med formula: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка создания формулы")
		return
	}

	a.successResponse(c, gin.H{"id": formula.ID})
}

// PUT /api/pvlc-med-formulas/:id - изменение формулы
func (a *API) UpdatePvlcMedFormula(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID формулы")
		return
	}

	var request ds.UpdatePvlcMedFormulaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	formula, err := a.repo.GetPvlcMedFormulaByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Формула не найдена")
		return
	}

	// Обновляем только переданные поля
	if request.Title != "" {
		formula.Title = request.Title
	}
	if request.Description != "" {
		formula.Description = request.Description
	}
	if request.Formula != "" {
		formula.Formula = request.Formula
	}
	if request.Category != "" {
		formula.Category = request.Category
	}
	if request.Gender != "" {
		formula.Gender = request.Gender
	}
	if request.MinAge > 0 {
		formula.MinAge = request.MinAge
	}
	if request.MaxAge > 0 {
		formula.MaxAge = request.MaxAge
	}
	if request.IsActive != nil {
		formula.IsActive = *request.IsActive
	}

	if err := a.repo.UpdatePvlcMedFormula(&formula); err != nil {
		logrus.Error("Error updating pvlc med formula: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка обновления формулы")
		return
	}

	a.successResponse(c, gin.H{"message": "Формула успешно обновлена"})
}

// DELETE /api/pvlc-med-formulas/:id - удаление формулы
func (a *API) DeletePvlcMedFormula(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID формулы")
		return
	}

	// Получаем формулу для проверки существования и получения image_url
	formula, err := a.repo.GetPvlcMedFormulaByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Формула не найдена")
		return
	}

	// Удаляем изображение из MinIO если оно есть
	if formula.ImageURL != "" {
		if err := a.repo.DeleteImageFromMinIO(formula.ImageURL); err != nil {
			logrus.Warn("Error deleting image from MinIO: ", err)
		}
	}

	// Удаляем формулу из БД
	if err := a.repo.DeletePvlcMedFormula(uint(id)); err != nil {
		logrus.Error("Error deleting pvlc med formula: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка удаления формулы")
		return
	}

	a.successResponse(c, gin.H{"message": "Формула успешно удалена"})
}

// POST /api/pvlc-med-formulas/:id/image - добавление изображения
func (a *API) UploadPvlcMedFormulaImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID формулы")
		return
	}

	// Проверяем существование формулы
	formula, err := a.repo.GetPvlcMedFormulaByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Формула не найдена")
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
	if formula.ImageURL != "" {
		if err := a.repo.DeleteImageFromMinIO(formula.ImageURL); err != nil {
			logrus.Warn("Error deleting old image from MinIO: ", err)
		}
	}

	// Обновляем image_url в БД
	formula.ImageURL = imageURL
	if err := a.repo.UpdatePvlcMedFormula(&formula); err != nil {
		logrus.Error("Error updating pvlc med formula image: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка обновления формулы")
		return
	}

	a.successResponse(c, gin.H{
		"message":   "Изображение успешно загружено",
		"image_url": imageURL,
	})
}

// POST /api/pvlc-med-formulas/:id/add-to-cart - добавление в заявку-черновик
func (a *API) AddPvlcMedFormulaToCart(c *gin.Context) {
	idStr := c.Param("id")
	formulaID, err := strconv.Atoi(idStr)
	if err != nil || formulaID <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID формулы")
		return
	}

	// Фиксированный пользователь (как требуется в задании)
	userID := uint(1)

	// Создаем или получаем черновик
	card, err := a.repo.GetOrCreateDraftPvlcMedCard(userID)
	if err != nil {
		logrus.Error("Error getting draft pvlc med card: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка доступа к корзине")
		return
	}

	// Добавляем формулу в заявку
	if err := a.repo.AddPvlcMedFormulaToCard(card.ID, uint(formulaID)); err != nil {
		logrus.Error("Error adding to cart: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка добавления в корзину")
		return
	}

	a.successResponse(c, gin.H{
		"message":     "Формула добавлена в заявку",
		"med_card_id": card.ID,
	})
}
