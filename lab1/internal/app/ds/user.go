package ds

// Users - модель пользователя (используем встроенную, но для примера создаем свою)
type User struct {
	ID          uint   `gorm:"primaryKey"`
	Login       string `gorm:"type:varchar(25);unique;not null"`
	Password    string `gorm:"type:varchar(100);not null"`
	IsModerator bool   `gorm:"type:boolean;default:false"`
}
