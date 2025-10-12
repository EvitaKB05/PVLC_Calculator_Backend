// internal/app/ds/pvlc_med_formulas.go
package ds

import "time"

type PvlcMedFormula struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"not null"` // название формулы
	Description string    // описание
	Formula     string    `gorm:"not null"` // формула
	ImageURL    string    // урл картинки
	Category    string    `gorm:"not null"`
	Gender      string    `gorm:"not null"`
	MinAge      int       `gorm:"not null"`
	MaxAge      int       `gorm:"not null"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time `gorm:"autoCreateTime"` // ДОБАВЛЕНО ДЛЯ ЛР4
	UpdatedAt   time.Time `gorm:"autoUpdateTime"` // ДОБАВЛЕНО ДЛЯ ЛР4
}

// рост +
type PvlcMedFormulaWithHeight struct {
	PvlcMedFormula
	InputHeight float64 `json:"input_height"`
}
