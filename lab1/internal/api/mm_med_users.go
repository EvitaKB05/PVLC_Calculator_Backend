package api

import (
	"lab1/internal/app/ds"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Домен: М-М расчеты (MedMmPvlcCalculations)

// DELETE /api/med-mm-pvlc-calculations - удаление из заявки
func (a *API) DeleteMedMmPvlcCalculation(c *gin.Context) {
	var request struct {
		CardID           uint `json:"card_id" binding:"required"`
		PvlcMedFormulaID uint `json:"pvlc_med_formula_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	// Проверяем что заявка существует и это черновик
	card, err := a.repo.GetPvlcMedCardByID(request.CardID)
	if err != nil || card.Status != ds.PvlcMedCardStatusDraft {
		a.errorResponse(c, http.StatusBadRequest, "Неверная заявка")
		return
	}

	if err := a.repo.DeleteMedMmPvlcCalculation(request.CardID, request.PvlcMedFormulaID); err != nil {
		logrus.Error("Error deleting med mm pvlc calculation: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка удаления расчета")
		return
	}

	a.successResponse(c, gin.H{"message": "Расчет удален из заявки"})
}

// PUT /api/med-mm-pvlc-calculations - изменение м-м
func (a *API) UpdateMedMmPvlcCalculation(c *gin.Context) {
	var request struct {
		CardID           uint                                 `json:"card_id" binding:"required"`
		PvlcMedFormulaID uint                                 `json:"pvlc_med_formula_id" binding:"required"`
		Data             ds.UpdateMedMmPvlcCalculationRequest `json:"data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	// Проверяем что заявка существует и это черновик
	card, err := a.repo.GetPvlcMedCardByID(request.CardID)
	if err != nil || card.Status != ds.PvlcMedCardStatusDraft {
		a.errorResponse(c, http.StatusBadRequest, "Неверная заявка")
		return
	}

	if err := a.repo.UpdateMedMmPvlcCalculation(request.CardID, request.PvlcMedFormulaID, request.Data.InputHeight); err != nil {
		logrus.Error("Error updating med mm pvlc calculation: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка обновления расчета")
		return
	}

	a.successResponse(c, gin.H{"message": "Расчет успешно обновлен"})
}

// Домен: Пользователи (MedUsers)

// POST /api/med-users/register - регистрация
func (a *API) RegisterMedUser(c *gin.Context) {
	var request ds.MedUserRegistrationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	// Проверяем что логин не занят
	existing, _ := a.repo.GetMedUserByLogin(request.Login)
	if existing != nil {
		a.errorResponse(c, http.StatusBadRequest, "Пользователь с таким логином уже существует")
		return
	}

	user := ds.MedUser{
		Login:       request.Login,
		Password:    request.Password, // В реальном приложении нужно хэшировать!
		IsModerator: request.IsModerator,
	}

	if err := a.repo.CreateMedUser(&user); err != nil {
		logrus.Error("Error creating med user: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка регистрации")
		return
	}

	a.successResponse(c, gin.H{
		"message": "Пользователь успешно зарегистрирован",
		"user_id": user.ID,
	})
}

// GET /api/med-users/profile - профиль пользователя
func (a *API) GetMedUserProfile(c *gin.Context) {
	// Фиксированный пользователь для демонстрации
	userID := uint(1)

	user, err := a.repo.GetMedUserByID(userID)
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Пользователь не найден")
		return
	}

	response := ds.MedUserResponse{
		ID:          user.ID,
		Login:       user.Login,
		IsModerator: user.IsModerator,
	}

	a.successResponse(c, response)
}

// PUT /api/med-users/profile - обновление профиля
func (a *API) UpdateMedUserProfile(c *gin.Context) {
	// Фиксированный пользователь для демонстрации
	userID := uint(1)

	var request ds.UpdateMedUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	user, err := a.repo.GetMedUserByID(userID)
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Пользователь не найден")
		return
	}

	if request.Password != "" {
		user.Password = request.Password // В реальном приложении нужно хэшировать!
	}

	if err := a.repo.UpdateMedUser(&user); err != nil {
		logrus.Error("Error updating med user: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка обновления профиля")
		return
	}

	a.successResponse(c, gin.H{"message": "Профиль успешно обновлен"})
}

// POST /api/med-users/login - аутентификация
func (a *API) LoginMedUser(c *gin.Context) {
	var request struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	user, err := a.repo.GetMedUserByLogin(request.Login)
	if err != nil || user.Password != request.Password { // В реальном приложении сравнивать хэши!
		a.errorResponse(c, http.StatusUnauthorized, "Неверный логин или пароль")
		return
	}

	response := ds.MedUserResponse{
		ID:          user.ID,
		Login:       user.Login,
		IsModerator: user.IsModerator,
	}

	a.successResponse(c, gin.H{
		"message": "Успешная аутентификация",
		"user":    response,
	})
}

// POST /api/med-users/logout - деавторизация
func (a *API) LogoutMedUser(c *gin.Context) {
	// В реальном приложении здесь инвалидируем токен
	a.successResponse(c, gin.H{"message": "Успешный выход из системы"})
}
