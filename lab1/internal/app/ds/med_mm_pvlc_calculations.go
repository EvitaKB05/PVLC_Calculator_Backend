// internal/app/ds/med_mm_pvlc_calculations.go
package ds

import "time"

type MedMmPvlcCalculation struct {
	PvlcMedCardID    uint           `gorm:"primaryKey"` // Вместе образуют составной ключ
	PvlcMedFormulaID uint           `gorm:"primaryKey"`
	PvlcMedCard      PvlcMedCard    `gorm:"foreignKey:PvlcMedCardID; constraint:OnDelete:RESTRICT"`
	PvlcMedFormula   PvlcMedFormula `gorm:"foreignKey:PvlcMedFormulaID; constraint:OnDelete:RESTRICT"`
	InputHeight      float64        `gorm:"not null; default:0"` // дефолтный рост
	FinalResult      float64        `gorm:"not null; default:0"` // Результат расчета по формуле
	CreatedAt        time.Time      `gorm:"autoCreateTime"`      // ДОБАВЛЕНО ДЛЯ ЛР4
	UpdatedAt        time.Time      `gorm:"autoUpdateTime"`      // ДОБАВЛЕНО ДЛЯ ЛР4
}
