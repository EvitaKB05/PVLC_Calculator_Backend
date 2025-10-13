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
	// ИСПРАВЛЕНО: добавляем проверку на nil redisClient
	if redisClient == nil {
		logrus.Warn("Redis client is nil in ParseToken")
		return nil, errors.New("redis client not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// ИСПРАВЛЕНО: проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		logrus.Warn("JWT parsing error:", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*ds.JWTClaims); ok && token.Valid {
		// Проверяем, не в черном списке ли токен
		inBlacklist, err := redisClient.CheckJWTInBlacklist(context.Background(), tokenString)
		if err != nil {
			logrus.Warn("Error checking JWT blacklist:", err)
			return nil, err
		}
		if inBlacklist {
			logrus.Warn("Token is in blacklist")
			return nil, errors.New("token revoked")
		}
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// extractTokenFromHeader извлекает токен из заголовка Authorization
// ИСПРАВЛЕНО: обрабатываем разные форматы заголовков
func extractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is empty")
	}

	// Убираем лишние пробелы
	authHeader = strings.TrimSpace(authHeader)

	// Проверяем разные форматы
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[len("Bearer "):], nil
	} else if strings.HasPrefix(authHeader, "bearer ") {
		return authHeader[len("bearer "):], nil
	} else if strings.HasPrefix(authHeader, "Token ") {
		return authHeader[len("Token "):], nil
	} else if strings.HasPrefix(authHeader, "token ") {
		return authHeader[len("token "):], nil
	} else {
		// Если нет префикса, возможно токен передан без префикса
		// Проверяем, похож ли на JWT (содержит точки)
		if strings.Count(authHeader, ".") == 2 {
			return authHeader, nil
		}
		return "", errors.New("invalid authorization header format")
	}
}

// AuthMiddleware - middleware для проверки JWT токена
// ДОБАВЛЕНО ДЛЯ ЛАБОРАТОРНОЙ РАБОТЫ 4
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Debug("AuthMiddleware started for path: ", c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")
		logrus.Debug("Raw Authorization header: '", authHeader, "'")

		if authHeader == "" {
			logrus.Debug("No Authorization header found")
			// Если заголовка нет, продолжаем без аутентификации
			c.Next()
			return
		}

		// Извлекаем токен из заголовка
		tokenString, err := extractTokenFromHeader(authHeader)
		if err != nil {
			logrus.Warn("Error extracting token from header: ", err)
			c.Next()
			return
		}

		logrus.Debug("Token extracted: ", tokenString[:10]+"...") // Логируем только начало токена

		claims, err := ParseToken(tokenString)
		if err != nil {
			logrus.Warn("Invalid token: ", err)
			c.Next()
			return
		}

		// Сохраняем claims в контексте
		c.Set(AuthUserKey, claims)
		logrus.Debug("User claims set in context: ", claims.UserID, " ", claims.Login)

		c.Next()
	}
}

// RequireAuth - middleware, требующий аутентификации
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Debug("RequireAuth started for path: ", c.Request.URL.Path)

		user, exists := c.Get(AuthUserKey)
		if !exists {
			logrus.Warn("No user found in context for protected route")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
			c.Abort()
			return
		}

		claims, ok := user.(*ds.JWTClaims)
		if !ok {
			logrus.Warn("Invalid user type in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат аутентификации"})
			c.Abort()
			return
		}

		logrus.Debug("User authenticated: ", claims.UserID, " ", claims.Login)
		c.Next()
	}
}

// RequireModerator - middleware, требующий прав модератора
func RequireModerator() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Debug("RequireModerator started for path: ", c.Request.URL.Path)

		user, exists := c.Get(AuthUserKey)
		if !exists {
			logrus.Warn("No user found in context for moderator route")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
			c.Abort()
			return
		}

		claims, ok := user.(*ds.JWTClaims)
		if !ok || !claims.IsModerator {
			logrus.Warn("User is not moderator: ", claims.Login)
			c.JSON(http.StatusForbidden, gin.H{"error": "Требуются права модератора"})
			c.Abort()
			return
		}

		logrus.Debug("Moderator authenticated: ", claims.UserID, " ", claims.Login)
		c.Next()
	}
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(c *gin.Context) *ds.JWTClaims {
	user, exists := c.Get(AuthUserKey)
	if !exists {
		logrus.Debug("No user found in context")
		return nil
	}
	claims, ok := user.(*ds.JWTClaims)
	if !ok {
		logrus.Warn("Invalid user type in context")
		return nil
	}
	return claims
}
