package ds

import "time"

// Статусы медицинской карты
const (
	MedicalCardStatusDraft     = "черновик"    // Пользователь еще выбирает формулы
	MedicalCardStatusFormed    = "сформирован" // Пользователь отправил на расчет
	MedicalCardStatusCompleted = "завершен"    // Модератор проверил и завершил
	MedicalCardStatusRejected  = "отклонен"    // Модератор отклонил
	MedicalCardStatusDeleted   = "удален"      // Пользователь "удалил" карту
)

type MedicalCard struct {
	ID           uint              `gorm:"primaryKey"`
	Status       string            `gorm:"not null; default:'черновик'"` // Статус из констант
	CreatedAt    time.Time         `gorm:"not null; default:now()"`
	PatientName  string            `gorm:"not null"` // ФИО пациента (заменил "создателя")
	FinalizedAt  *time.Time        // Когда отправили на расчет (аналог "дата формирования")
	CompletedAt  *time.Time        // Когда модератор завершил
	ModeratorID  *uint             // ID модератора, который проверил
	Moderator    User              `gorm:"foreignKey:ModeratorID; constraint:OnDelete:SET NULL"`
	TotalResult  float64           `gorm:"type:decimal(10,2)"`       // Суммарный результат (если нужно)
	Calculations []CardCalculation `gorm:"foreignKey:MedicalCardID"` // Связь с формулами
}
