package ds

type CardCalculation struct {
	MedicalCardID uint        `gorm:"primaryKey"` // Вместе образуют составной ключ
	CalculationID uint        `gorm:"primaryKey"`
	MedicalCard   MedicalCard `gorm:"foreignKey:MedicalCardID; constraint:OnDelete:RESTRICT"`
	Calculation   Calculation `gorm:"foreignKey:CalculationID; constraint:OnDelete:RESTRICT"`
	InputHeight   float64     `gorm:"not null"` // Рост, который ввел пользователь
	FinalResult   float64     `gorm:"not null"` // Результат расчета по формуле
}
