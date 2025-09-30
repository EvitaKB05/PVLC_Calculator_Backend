package ds

type MedUser struct {
	ID          uint   `gorm:"primaryKey"`
	Login       string `gorm:"unique; not null"`
	Password    string `gorm:"not null"`
	IsModerator bool   `gorm:"default:false"`
}
