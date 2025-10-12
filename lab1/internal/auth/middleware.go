// internal/auth/middleware.go
package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"lab1/internal/app/ds"
	"lab1/internal/redis"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

// Константы для работы с аутентификацией
const AuthUserKey = "user"                  // Ключ для хранения пользователя в контексте
const jwtPrefix = "Bearer "                 // Префикc JWT токена
const jwtSecret = "your-secret-key-for-jwt" // Секретный ключ для JWT

var redisClient *redis.Client

// InitAuth инициализирует модуль аутентификации
// ДОБАВЛЕНО ДЛЯ ЛАБОРАТОРНОЙ РАБОТЫ 4
func InitAuth(redis *redis.Client) {
	redisClient = redis
}

// GenerateToken создает JWT токен для пользователя
func GenerateToken(userID uint, login string, isModerator bool) (string, error) {
	claims := ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Токен действует 24 часа
			IssuedAt:  time.Now().Unix(),
			Issuer:    "lung-capacity-app",
		},
		UserID:      userID,
		Login:       login,
		IsModerator: isModerator,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// ParseToken проверяет и парсит JWT токен
func ParseToken(tokenString string) (*ds.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*ds.JWTClaims); ok && token.Valid {
		// Проверяем, не в черном списке ли токен
		inBlacklist, err := redisClient.CheckJWTInBlacklist(context.Background(), tokenString)
		if err != nil {
			return nil, err
		}
		if inBlacklist {
			return nil, errors.New("token revoked")
		}
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// AuthMiddleware - middleware для проверки JWT токена
// ДОБАВЛЕНО ДЛЯ ЛАБОРАТОРНОЙ РАБОТЫ 4
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Если заголовка нет, продолжаем без аутентификации
			c.Next()
			return
		}

		// Проверяем формат заголовка
		if !strings.HasPrefix(authHeader, jwtPrefix) {
			logrus.Warn("Invalid authorization header format")
			c.Next()
			return
		}

		// Извлекаем токен
		tokenString := authHeader[len(jwtPrefix):]
		claims, err := ParseToken(tokenString)
		if err != nil {
			logrus.Warn("Invalid token: ", err)
			c.Next()
			return
		}

		// Сохраняем claims в контексте
		c.Set(AuthUserKey, claims)
		c.Next()
	}
}

// RequireAuth - middleware, требующий аутентификации
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ИСПРАВЛЕНА ОШИБКА: используем пустой идентификатор для переменной, которая не используется
		if _, exists := c.Get(AuthUserKey); !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireModerator - middleware, требующий прав модератора
func RequireModerator() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get(AuthUserKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
			c.Abort()
			return
		}

		claims, ok := user.(*ds.JWTClaims)
		if !ok || !claims.IsModerator {
			c.JSON(http.StatusForbidden, gin.H{"error": "Требуются права модератора"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(c *gin.Context) *ds.JWTClaims {
	user, exists := c.Get(AuthUserKey)
	if !exists {
		return nil
	}
	claims, ok := user.(*ds.JWTClaims)
	if !ok {
		return nil
	}
	return claims
}
