package repository

import (
	"lab1/internal/app/ds"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

// GetServices
func (r *Repository) GetServices() ([]ds.Service, error) {
	var services []ds.Service
	err := r.db.Where("is_active = ?", true).Find(&services).Error
	if err != nil {
		return nil, err
	}
	return services, nil
}

// GetServiceByID
func (r *Repository) GetServiceByID(id uint) (ds.Service, error) {
	var service ds.Service
	err := r.db.First(&service, id).Error
	if err != nil {
		return ds.Service{}, err
	}
	return service, nil
}

// GetServicesByTitle
func (r *Repository) GetServicesByTitle(title string) ([]ds.Service, error) {
	var services []ds.Service
	err := r.db.Where("title ILIKE ? AND is_active = ?", "%"+title+"%", true).Find(&services).Error
	if err != nil {
		return nil, err
	}
	return services, nil
}

// GetDraftOrder
func (r *Repository) GetDraftOrder(userID uint) (ds.Order, error) {
	var order ds.Order
	err := r.db.Where("user_id = ? AND status = ?", userID, ds.StatusDraft).First(&order).Error
	return order, err
}

// CreateOrder
func (r *Repository) CreateOrder(userID uint) (ds.Order, error) {
	order := ds.Order{
		Status:    ds.StatusDraft,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	err := r.db.Create(&order).Error
	return order, err
}

// AddServiceToOrder добавляет услугу в заявку
func (r *Repository) AddServiceToOrder(orderID, serviceID uint) error {
	orderService := ds.OrderService{
		OrderID:   orderID,
		ServiceID: serviceID,
		Quantity:  1,
	}

	return r.db.Create(&orderService).Error
}

// GetOrderServices возвращает услуги в заявке
func (r *Repository) GetOrderServices(orderID uint) ([]ds.OrderService, error) {
	var orderServices []ds.OrderService
	err := r.db.Preload("Service").Where("order_id = ?", orderID).Find(&orderServices).Error
	return orderServices, err
}

// GetOrdersCount возвращает количество расчетов (для иконки корзины)
func (r *Repository) GetOrdersCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&ds.Order{}).
		Where("user_id = ? AND status = ?", userID, ds.StatusDraft).
		Count(&count).Error
	return count, err
}

// DeleteOrderLogical логически удаляет заявку (меняет статус на "удалён")
func (r *Repository) DeleteOrderLogical(orderID uint) error {
	return r.db.Model(&ds.Order{}).
		Where("id = ?", orderID).
		Update("status", ds.StatusDeleted).Error
}
