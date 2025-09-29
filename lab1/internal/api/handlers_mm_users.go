package api

import (
	"lab1/internal/app/ds"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Домен: М-М (Card Calculations)

// DELETE /api/card-calculations - удаление из заявки
func (a *API) DeleteCardCalculation(c *gin.Context) {
	var request struct {
		CardID        uint `json:"card_id" binding:"required"`
		CalculationID uint `json:"calculation_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	// Проверяем что заявка существует и это черновик
	card, err := a.repo.GetMedicalCardByID(request.CardID)
	if err != nil || card.Status != ds.MedicalCardStatusDraft {
		a.errorResponse(c, http.StatusBadRequest, "Неверная заявка")
		return
	}

	if err := a.repo.DeleteCardCalculation(request.CardID, request.CalculationID); err != nil {
		logrus.Error("Error deleting card calculation: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка удаления расчета")
		return
	}

	a.successResponse(c, gin.H{"message": "Расчет удален из заявки"})
}

// PUT /api/card-calculations - изменение м-м
func (a *API) UpdateCardCalculation(c *gin.Context) {
	var request struct {
		CardID        uint                            `json:"card_id" binding:"required"`
		CalculationID uint                            `json:"calculation_id" binding:"required"`
		Data          ds.UpdateCardCalculationRequest `json:"data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	// Проверяем что заявка существует и это черновик
	card, err := a.repo.GetMedicalCardByID(request.CardID)
	if err != nil || card.Status != ds.MedicalCardStatusDraft {
		a.errorResponse(c, http.StatusBadRequest, "Неверная заявка")
		return
	}

	if err := a.repo.UpdateCardCalculation(request.CardID, request.CalculationID, request.Data.InputHeight); err != nil {
		logrus.Error("Error updating card calculation: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка обновления расчета")
		return
	}

	a.successResponse(c, gin.H{"message": "Расчет успешно обновлен"})
}

// Домен: Пользователи

// POST /api/users/register - регистрация
func (a *API) RegisterUser(c *gin.Context) {
	var request ds.UserRegistrationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	// Проверяем что логин не занят
	existing, _ := a.repo.GetUserByLogin(request.Login)
	if existing != nil {
		a.errorResponse(c, http.StatusBadRequest, "Пользователь с таким логином уже существует")
		return
	}

	user := ds.User{
		Login:       request.Login,
		Password:    request.Password, // В реальном приложении нужно хэшировать!
		IsModerator: request.IsModerator,
	}

	if err := a.repo.CreateUser(&user); err != nil {
		logrus.Error("Error creating user: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка регистрации")
		return
	}

	a.successResponse(c, gin.H{
		"message": "Пользователь успешно зарегистрирован",
		"user_id": user.ID,
	})
}

// GET /api/users/profile - профиль пользователя
func (a *API) GetUserProfile(c *gin.Context) {
	// Фиксированный пользователь для демонстрации
	userID := uint(1)

	user, err := a.repo.GetUserByID(userID)
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Пользователь не найден")
		return
	}

	response := ds.UserResponse{
		ID:          user.ID,
		Login:       user.Login,
		IsModerator: user.IsModerator,
	}

	a.successResponse(c, response)
}

// PUT /api/users/profile - обновление профиля
func (a *API) UpdateUserProfile(c *gin.Context) {
	// Фиксированный пользователь для демонстрации
	userID := uint(1)

	var request ds.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	user, err := a.repo.GetUserByID(userID)
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Пользователь не найден")
		return
	}

	if request.Password != "" {
		user.Password = request.Password // В реальном приложении нужно хэшировать!
	}

	if err := a.repo.UpdateUser(&user); err != nil {
		logrus.Error("Error updating user: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка обновления профиля")
		return
	}

	a.successResponse(c, gin.H{"message": "Профиль успешно обновлен"})
}

// POST /api/users/login - аутентификация
func (a *API) Login(c *gin.Context) {
	var request struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	user, err := a.repo.GetUserByLogin(request.Login)
	if err != nil || user.Password != request.Password { // В реальном приложении сравнивать хэши!
		a.errorResponse(c, http.StatusUnauthorized, "Неверный логин или пароль")
		return
	}

	response := ds.UserResponse{
		ID:          user.ID,
		Login:       user.Login,
		IsModerator: user.IsModerator,
	}

	a.successResponse(c, gin.H{
		"message": "Успешная аутентификация",
		"user":    response,
	})
}

// POST /api/users/logout - деавторизация
func (a *API) Logout(c *gin.Context) {
	// В реальном приложении здесь инвалидируем токен
	a.successResponse(c, gin.H{"message": "Успешный выход из системы"})
}
