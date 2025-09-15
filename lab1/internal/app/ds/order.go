package ds

import (
	"time"
)

// OrderStatus - тип для статусов заявок
type OrderStatus string

const (
	StatusDraft     OrderStatus = "черновик"
	StatusDeleted   OrderStatus = "удалён"
	StatusFormed    OrderStatus = "сформирован"
	StatusCompleted OrderStatus = "завершён"
	StatusRejected  OrderStatus = "отклонён"
)

// Order - модель заявки (расчета)
type Order struct {
	ID          uint        `gorm:"primaryKey"`
	Status      OrderStatus `gorm:"type:varchar(20);not null"`
	UserID      uint        `gorm:"not null"`
	CreatedAt   time.Time   `gorm:"not null"`
	FinalizedAt *time.Time  // Дата формирования (Nullable)
	CompletedAt *time.Time  // Дата завершения (Nullable)
	ModeratorID *uint       // ID модератора (Nullable)
	DoctorName  string      `gorm:"type:varchar(100)"` // Доп поле по теме: ФИО врача

	// Связи
	User      User           `gorm:"foreignKey:UserID"`
	Moderator *User          `gorm:"foreignKey:ModeratorID"`
	Services  []OrderService `gorm:"foreignKey:OrderID"`
}
