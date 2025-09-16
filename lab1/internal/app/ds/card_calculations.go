package ds

type CardCalculation struct {
	MedicalCardID uint        `gorm:"primaryKey"` // Вместе образуют составной ключ
	CalculationID uint        `gorm:"primaryKey"`
	MedicalCard   MedicalCard `gorm:"foreignKey:MedicalCardID; constraint:OnDelete:RESTRICT"`
	Calculation   Calculation `gorm:"foreignKey:CalculationID; constraint:OnDelete:RESTRICT"`
	InputHeight   float64     `gorm:"not null; default:0"` // дефолтный рост
	FinalResult   float64     `gorm:"not null; default:0"` // Результат расчета по формуле
}
