package ds

type PvlcMedFormula struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"` // название формулы
	Description string // описание
	Formula     string `gorm:"not null"` // формула
	ImageURL    string // урл картинки
	Category    string `gorm:"not null"`
	Gender      string `gorm:"not null"`
	MinAge      int    `gorm:"not null"`
	MaxAge      int    `gorm:"not null"`
	IsActive    bool   `gorm:"default:true"`
}

// рост +
type PvlcMedFormulaWithHeight struct {
	PvlcMedFormula
	InputHeight float64 `json:"input_height"`
}
