package ds

type User struct {
	ID          uint   `gorm:"primaryKey"`
	Login       string `gorm:"unique; not null"` // Логин должен быть уникальным
	Password    string `gorm:"not null"`         // Пароль (будем хранить хэш)
	IsModerator bool   `gorm:"default:false"`    // Является ли пользователь модератором
}
