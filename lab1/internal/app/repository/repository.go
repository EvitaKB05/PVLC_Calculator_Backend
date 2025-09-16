package repository

import (
	"lab1/internal/app/ds"
	"lab1/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() (*Repository, error) {
	// Загружаем переменные окружения
	// (этот код уже есть в функции, оставляем его)
	_ = godotenv.Load()

	// Получаем DSN строку
	dsnString := dsn.FromEnv()

	// Подключаемся к БД
	db, err := gorm.Open(postgres.Open(dsnString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

// Получаем все активные расчеты из БД
func (r *Repository) GetServices() ([]ds.Calculation, error) {
	var calculations []ds.Calculation
	err := r.db.Where("is_active = ?", true).Find(&calculations).Error
	if err != nil {
		return nil, err
	}
	return calculations, nil
}

// Поиск расчетов по названию
func (r *Repository) GetServicesByTitle(title string) ([]ds.Calculation, error) {
	var calculations []ds.Calculation
	err := r.db.Where("is_active = ? AND title ILIKE ?", true, "%"+title+"%").Find(&calculations).Error
	if err != nil {
		return nil, err
	}
	return calculations, nil
}

// Получаем один расчет по ID
func (r *Repository) GetService(id int) (ds.Calculation, error) {
	var calculation ds.Calculation
	err := r.db.Where("is_active = ? AND id = ?", true, id).First(&calculation).Error
	if err != nil {
		return ds.Calculation{}, err
	}
	return calculation, nil
}

// Получаем количество расчетов в корзине (временная заглушка)
func (r *Repository) GetCalculationsCount() (int, error) {
	// Пока возвращаем 0, потом реализуем
	return 0, nil
}

// Получаем расчеты для корзины (временная заглушка)
func (r *Repository) GetCalculation() ([]ds.Calculation, error) {
	// Пока возвращаем пустой массив, потом реализуем
	return []ds.Calculation{}, nil
}
