// internal/api/auth_handlers.go
package api

import (
	"net/http"
	"time"

	"lab1/internal/app/ds"
	"lab1/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoginRequest - структура запроса для аутентификации
// ДОБАВЛЕНО ДЛЯ ЛАБОРАТОРНОЙ РАБОТЫ 4
type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse - структура ответа при успешной аутентификации
type LoginResponse struct {
	Token     string             `json:"token"`
	User      ds.MedUserResponse `json:"user"`
	ExpiresAt time.Time          `json:"expires_at"`
}

// LogoutRequest - структура запроса для выхода из системы
type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Выполняет вход пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для входа"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
func (a *API) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	// Ищем пользователя в базе данных
	user, err := a.repo.GetMedUserByLogin(req.Login)
	if err != nil {
		logrus.Warn("Login failed - user not found: ", req.Login)
		a.errorResponse(c, http.StatusUnauthorized, "Неверный логин или пароль")
		return
	}

	// Проверяем пароль (В РЕАЛЬНОМ ПРИЛОЖЕНИИ НУЖНО ХЕШИРОВАТЬ!)
	if user.Password != req.Password {
		logrus.Warn("Login failed - invalid password for user: ", req.Login)
		a.errorResponse(c, http.StatusUnauthorized, "Неверный логин или пароль")
		return
	}

	// Генерируем JWT токен
	token, err := auth.GenerateToken(user.ID, user.Login, user.IsModerator)
	if err != nil {
		logrus.Error("Error generating token: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка генерации токена")
		return
	}

	// Создаем ответ
	response := LoginResponse{
		Token: token,
		User: ds.MedUserResponse{
			ID:          user.ID,
			Login:       user.Login,
			IsModerator: user.IsModerator,
		},
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	a.successResponse(c, response)
}

// Logout godoc
// @Summary Выход пользователя
// @Description Добавляет JWT токен в черный список
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LogoutRequest true "Токен для выхода"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/auth/logout [post]
// @Security BearerAuth
func (a *API) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	// Добавляем токен в черный список на 24 часа
	err := a.redis.WriteJWTToBlacklist(c.Request.Context(), req.Token, 24*time.Hour)
	if err != nil {
		logrus.Error("Error adding token to blacklist: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка выхода из системы")
		return
	}

	a.successResponse(c, gin.H{"message": "Успешный выход из системы"})
}

// GetProfile godoc
// @Summary Получение профиля пользователя
// @Description Возвращает информацию о текущем пользователе
// @Tags auth
// @Produce json
// @Success 200 {object} ds.MedUserResponse
// @Failure 401 {object} map[string]string
// @Router /api/auth/profile [get]
// @Security BearerAuth
func (a *API) GetProfile(c *gin.Context) {
	claims := auth.GetUserFromContext(c)
	if claims == nil {
		a.errorResponse(c, http.StatusUnauthorized, "Требуется аутентификация")
		return
	}

	user, err := a.repo.GetMedUserByID(claims.UserID)
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
