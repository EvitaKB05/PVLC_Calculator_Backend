package repository

import (
	"context"
	"fmt"
	"lab1/internal/app/ds"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Методы для API

// GetCalculationsWithFilter - получение услуг с фильтрацией
func (r *Repository) GetCalculationsWithFilter(filter ds.CalculationFilter) ([]ds.Calculation, error) {
	var calculations []ds.Calculation
	query := r.db.Model(&ds.Calculation{})

	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Gender != "" {
		query = query.Where("gender = ?", filter.Gender)
	}
	if filter.MinAge > 0 {
		query = query.Where("min_age >= ?", filter.MinAge)
	}
	if filter.MaxAge > 0 {
		query = query.Where("max_age <= ?", filter.MaxAge)
	}
	if filter.Active != nil {
		query = query.Where("is_active = ?", *filter.Active)
	}

	err := query.Find(&calculations).Error
	return calculations, err
}

// GetCalculationByID - получение услуги по ID
func (r *Repository) GetCalculationByID(id uint) (ds.Calculation, error) {
	var calculation ds.Calculation
	err := r.db.First(&calculation, id).Error
	return calculation, err
}

// CreateCalculation - создание услуги
func (r *Repository) CreateCalculation(calculation *ds.Calculation) error {
	return r.db.Create(calculation).Error
}

// UpdateCalculation - обновление услуги
func (r *Repository) UpdateCalculation(calculation *ds.Calculation) error {
	return r.db.Save(calculation).Error
}

// DeleteCalculation - удаление услуги
func (r *Repository) DeleteCalculation(id uint) error {
	return r.db.Delete(&ds.Calculation{}, id).Error
}

// GetMedicalCardsWithFilter - получение заявок с фильтрацией
func (r *Repository) GetMedicalCardsWithFilter(filter ds.MedicalCardFilter) ([]ds.MedicalCard, error) {
	var cards []ds.MedicalCard
	query := r.db.Where("status != ? AND status != ?",
		ds.MedicalCardStatusDeleted, ds.MedicalCardStatusDraft)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.DateFrom != "" {
		if dateFrom, err := time.Parse("2006-01-02", filter.DateFrom); err == nil {
			query = query.Where("created_at >= ?", dateFrom)
		}
	}
	if filter.DateTo != "" {
		if dateTo, err := time.Parse("2006-01-02", filter.DateTo); err == nil {
			query = query.Where("created_at <= ?", dateTo.AddDate(0, 0, 1)) // включая весь день
		}
	}

	err := query.Preload("Moderator").Find(&cards).Error
	return cards, err
}

// GetMedicalCardByID - получение заявки по ID
func (r *Repository) GetMedicalCardByID(id uint) (ds.MedicalCard, error) {
	var card ds.MedicalCard
	err := r.db.Preload("Moderator").First(&card, id).Error
	return card, err
}

// UpdateMedicalCard - обновление заявки
func (r *Repository) UpdateMedicalCard(card *ds.MedicalCard) error {
	return r.db.Save(card).Error
}

// GetDraftCardByUserID - получение черновика по ID пользователя
func (r *Repository) GetDraftCardByUserID(userID uint) (*ds.MedicalCard, error) {
	var card ds.MedicalCard
	err := r.db.Where("status = ?", ds.MedicalCardStatusDraft).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// GetCalculationsCountByCardID - количество расчетов в заявке
func (r *Repository) GetCalculationsCountByCardID(cardID uint) (int, error) {
	var count int64
	err := r.db.Model(&ds.CardCalculation{}).Where("medical_card_id = ?", cardID).Count(&count).Error
	return int(count), err
}

// GetCardCalculationsByCardID - получение расчетов заявки
func (r *Repository) GetCardCalculationsByCardID(cardID uint) ([]ds.CardCalculation, error) {
	var calculations []ds.CardCalculation
	err := r.db.Where("medical_card_id = ?", cardID).
		Preload("Calculation").
		Find(&calculations).Error
	return calculations, err
}

// DeleteCardCalculation - удаление расчета из заявки
func (r *Repository) DeleteCardCalculation(cardID, calculationID uint) error {
	return r.db.Where("medical_card_id = ? AND calculation_id = ?", cardID, calculationID).
		Delete(&ds.CardCalculation{}).Error
}

// UpdateCardCalculation - обновление расчета в заявке
func (r *Repository) UpdateCardCalculation(cardID, calculationID uint, inputHeight float64) error {
	return r.db.Model(&ds.CardCalculation{}).
		Where("medical_card_id = ? AND calculation_id = ?", cardID, calculationID).
		Update("input_height", inputHeight).Error
}

// CalculateTotalDjel - вычисление общего ДЖЕЛ для заявки
func (r *Repository) CalculateTotalDjel(cardID uint) (float64, error) {
	var calculations []ds.CardCalculation
	err := r.db.Where("medical_card_id = ?", cardID).
		Preload("Calculation").
		Find(&calculations).Error
	if err != nil {
		return 0, err
	}

	total := 0.0
	for _, cc := range calculations {
		// Здесь должна быть реализация расчета ДЖЕЛ по формуле
		// Для демонстрации используем простой расчет
		if cc.InputHeight > 0 {
			// Пример расчета: используем рост как базовый параметр
			result := cc.InputHeight * 0.05 // Упрощенная формула
			cc.FinalResult = result
			total += result

			// Сохраняем результат расчета
			r.db.Save(&cc)
		}
	}

	return total, nil
}

// User methods
func (r *Repository) GetUserByID(id uint) (ds.User, error) {
	var user ds.User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *Repository) GetUserByLogin(login string) (*ds.User, error) {
	var user ds.User
	err := r.db.Where("login = ?", login).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateUser(user *ds.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) UpdateUser(user *ds.User) error {
	return r.db.Save(user).Error
}

// MinIO methods
func (r *Repository) UploadImageToMinIO(file *multipart.FileHeader, calculationID uint) (string, error) {
	// Создаем клиент MinIO
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minio", "minio124", ""),
		Secure: false,
	})
	if err != nil {
		return "", err
	}

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Генерируем имя файла
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("calculation_%d%s", calculationID, ext)

	// Загружаем в MinIO
	_, err = minioClient.PutObject(context.Background(), "pics", filename, src, file.Size, minio.PutObjectOptions{
		ContentType: "image/" + strings.TrimPrefix(ext, "."),
	})
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (r *Repository) DeleteImageFromMinIO(imageURL string) error {
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minio", "minio124", ""),
		Secure: false,
	})
	if err != nil {
		return err
	}

	err = minioClient.RemoveObject(context.Background(), "pics", imageURL, minio.RemoveObjectOptions{})
	return err
}

// Добавьте этот метод в конец api_repository.go если отсутствует
func (r *Repository) InitMinIOBucket() error {
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minio", "minio124", ""),
		Secure: false,
	})
	if err != nil {
		return fmt.Errorf("failed to create MinIO client: %v", err)
	}

	exists, err := minioClient.BucketExists(context.Background(), "pics")
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %v", err)
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), "pics", minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
		fmt.Println("MinIO bucket 'pics' created successfully")
	} else {
		fmt.Println("MinIO bucket 'pics' already exists")
	}
	return nil
}
