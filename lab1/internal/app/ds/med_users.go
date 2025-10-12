// internal/app/ds/med_users.go
package ds

import (
	"time"
)

type MedUser struct {
	ID          uint      `gorm:"primaryKey"`
	Login       string    `gorm:"unique; not null"`
	Password    string    `gorm:"not null"`
	IsModerator bool      `gorm:"default:false"`
	CreatedBy   uint      `gorm:"default:1"`      // Кто создал пользователя (для аудита)
	CreatedAt   time.Time `gorm:"autoCreateTime"` // ДОБАВЛЕНО ДЛЯ ЛР4
	UpdatedAt   time.Time `gorm:"autoUpdateTime"` // ДОБАВЛЕНО ДЛЯ ЛР4
}
