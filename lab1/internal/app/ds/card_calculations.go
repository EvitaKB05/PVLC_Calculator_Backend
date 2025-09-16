package ds

type CardCalculation struct {
	MedicalCardID uint        `gorm:"primaryKey"` // ключ составной!
	CalculationID uint        `gorm:"primaryKey"`
	MedicalCard   MedicalCard `gorm:"foreignKey:MedicalCardID; constraint:OnDelete:RESTRICT"`
	Calculation   Calculation `gorm:"foreignKey:CalculationID; constraint:OnDelete:RESTRICT"`
	InputHeight   float64     `gorm:"not null"` // вводимый рост
	FinalResult   float64     `gorm:"not null"` // результат
}
