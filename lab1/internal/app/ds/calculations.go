package ds

type Calculation struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"` // Название формулы ("Мальчики 4-7 лет")
	Description string // Описание
	Formula     string `gorm:"not null"` // Сама формула
	ImageURL    string // Ссылка на изображение
	Category    string `gorm:"not null"`     // Категория ("дети", "взрослые")
	Gender      string `gorm:"not null"`     // Пол ("мужской", "женский")
	MinAge      int    `gorm:"not null"`     // Минимальный возраст
	MaxAge      int    `gorm:"not null"`     // Максимальный возраст
	IsActive    bool   `gorm:"default:true"` // Активна ли формула (вместо удаления)
}
