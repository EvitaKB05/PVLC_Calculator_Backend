package ds

import "time"

// Статусы медицинской карты
const (
	MedicalCardStatusDraft     = "черновик"    //
	MedicalCardStatusFormed    = "сформирован" //
	MedicalCardStatusCompleted = "завершен"    //
	MedicalCardStatusRejected  = "отклонен"    //
	MedicalCardStatusDeleted   = "удален"      //
)

type MedicalCard struct {
	ID           uint              `gorm:"primaryKey"`
	Status       string            `gorm:"not null; default:'черновик'"` //
	CreatedAt    time.Time         `gorm:"not null; default:now()"`
	PatientName  string            `gorm:"not null"`                                 //
	DoctorName   string            `gorm:"type:varchar(100); default:'Иванов И.И.'"` // ДОБАВЛЯЕМ ПОЛЕ ДЛЯ ФИО ВРАЧА
	FinalizedAt  *time.Time        //
	CompletedAt  *time.Time        //
	ModeratorID  *uint             //
	Moderator    User              `gorm:"foreignKey:ModeratorID; constraint:OnDelete:SET NULL"`
	TotalResult  float64           `gorm:"type:decimal(10,2)"`       // Суммарный результат (если нужно)
	Calculations []CardCalculation `gorm:"foreignKey:MedicalCardID"` // Связь с формулами
}
