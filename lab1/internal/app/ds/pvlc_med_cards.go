// internal/app/ds/pvlc_med_cards.go
package ds

import "time"

// Статусы медицинской карты
const (
	PvlcMedCardStatusDraft     = "черновик"
	PvlcMedCardStatusFormed    = "сформирован"
	PvlcMedCardStatusCompleted = "завершен"
	PvlcMedCardStatusRejected  = "отклонен"
	PvlcMedCardStatusDeleted   = "удален"
)

type PvlcMedCard struct {
	ID           uint      `gorm:"primaryKey"`
	Status       string    `gorm:"not null; default:'черновик'"`
	CreatedAt    time.Time `gorm:"autoCreateTime"` // ИЗМЕНЕНО ДЛЯ ЛР4
	UpdatedAt    time.Time `gorm:"autoUpdateTime"` // ДОБАВЛЕНО ДЛЯ ЛР4
	PatientName  string    `gorm:"not null"`
	DoctorName   string    `gorm:"type:varchar(100); default:'Иванов И.И.'"`
	FinalizedAt  *time.Time
	CompletedAt  *time.Time
	ModeratorID  *uint
	Moderator    MedUser                `gorm:"foreignKey:ModeratorID; constraint:OnDelete:SET NULL"`
	TotalResult  float64                `gorm:"type:decimal(10,2)"`
	Calculations []MedMmPvlcCalculation `gorm:"foreignKey:PvlcMedCardID"`
	UserID       uint                   `gorm:"not null; default:1"` // Владелец заявки - ДОБАВЛЕНО ДЛЯ ЛР4
	User         MedUser                `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE"`
}
