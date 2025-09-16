package ds

import "time"

// статусы
const (
	MedicalCardStatusDraft     = "черновик"
	MedicalCardStatusFormed    = "сформирован"
	MedicalCardStatusCompleted = "завершен"
	MedicalCardStatusRejected  = "отклонен"
	MedicalCardStatusDeleted   = "удален"
)

type MedicalCard struct {
	ID           uint      `gorm:"primaryKey"`
	Status       string    `gorm:"not null; default:'черновик'"` // дефолтный
	CreatedAt    time.Time `gorm:"not null; default:now()"`
	PatientName  string    `gorm:"not null"`
	FinalizedAt  *time.Time
	CompletedAt  *time.Time
	ModeratorID  *uint
	Moderator    User              `gorm:"foreignKey:ModeratorID; constraint:OnDelete:SET NULL"`
	TotalResult  float64           `gorm:"type:decimal(10,2)"`
	Calculations []CardCalculation `gorm:"foreignKey:MedicalCardID"` // Связь с формулами
}
