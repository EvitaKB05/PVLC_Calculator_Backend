package ds

import "time"

// Статусы медицинской карты
const (
	PvlcMedCardStatusDraft     = "черновик"    //
	PvlcMedCardStatusFormed    = "сформирован" //
	PvlcMedCardStatusCompleted = "завершен"    //
	PvlcMedCardStatusRejected  = "отклонен"    //
	PvlcMedCardStatusDeleted   = "удален"      //
)

type PvlcMedCard struct {
	ID           uint                   `gorm:"primaryKey"`
	Status       string                 `gorm:"not null; default:'черновик'"` //
	CreatedAt    time.Time              `gorm:"not null; default:now()"`
	PatientName  string                 `gorm:"not null"`                                 //
	DoctorName   string                 `gorm:"type:varchar(100); default:'Иванов И.И.'"` // ДОБАВЛЯЕМ ПОЛЕ ДЛЯ ФИО ВРАЧА
	FinalizedAt  *time.Time             //
	CompletedAt  *time.Time             //
	ModeratorID  *uint                  //
	Moderator    MedUser                `gorm:"foreignKey:ModeratorID; constraint:OnDelete:SET NULL"`
	TotalResult  float64                `gorm:"type:decimal(10,2)"`       // Суммарный результат (если нужно)
	Calculations []MedMmPvlcCalculation `gorm:"foreignKey:PvlcMedCardID"` // Связь с формулами
}
