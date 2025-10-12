// internal/app/ds/jwt_claims.go
package ds

import (
	"github.com/golang-jwt/jwt"
)

// JWTClaims - кастомные claims для JWT токена
// ДОБАВЛЕНО ДЛЯ ЛАБОРАТОРНОЙ РАБОТЫ 4
type JWTClaims struct {
	jwt.StandardClaims
	UserID      uint   `json:"user_id"`      // ID пользователя
	Login       string `json:"login"`        // Логин пользователя
	IsModerator bool   `json:"is_moderator"` // Флаг модератора
}
