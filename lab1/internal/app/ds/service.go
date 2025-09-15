package ds

// Service - модель услуги (формулы расчета ДЖЕЛ)
type Service struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:text"`
	Formula     string `gorm:"type:text;not null"`
	Image       string `gorm:"type:varchar(255)"` // Nullable
	Category    string `gorm:"type:varchar(50)"`
	Gender      string `gorm:"type:varchar(20)"`
	MinAge      int
	MaxAge      int
	IsActive    bool `gorm:"type:boolean;default:true"` // Статус удален/действует
}
