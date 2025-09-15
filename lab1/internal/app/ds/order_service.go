package ds

// OrderService - модель связи многие-ко-многим между заявками и услугами
type OrderService struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint `gorm:"not null;uniqueIndex:idx_order_service"`
	ServiceID uint `gorm:"not null;uniqueIndex:idx_order_service"`
	Quantity  int  `gorm:"default:1"` // Количество (доп поле)

	// Связи
	Order   Order   `gorm:"foreignKey:OrderID"`
	Service Service `gorm:"foreignKey:ServiceID"`
}
