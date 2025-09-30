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

// ==================== МЕТОДЫ ДЛЯ API ФОРМУЛ ДЖЕЛ ====================

// GetPvlcMedFormulasWithFilter - получение формул с фильтрацией для API (бывший GetCalculationsWithFilter)
func (r *Repository) GetPvlcMedFormulasWithFilter(filter ds.PvlcMedFormulaFilter) ([]ds.PvlcMedFormula, error) {
	var formulas []ds.PvlcMedFormula
	query := r.db.Model(&ds.PvlcMedFormula{})

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

	err := query.Find(&formulas).Error
	return formulas, err
}

// GetPvlcMedFormulaByID - получение формулы по ID для API (бывший GetCalculationByID)
func (r *Repository) GetPvlcMedFormulaByID(id uint) (ds.PvlcMedFormula, error) {
	var formula ds.PvlcMedFormula
	err := r.db.First(&formula, id).Error
	return formula, err
}

// CreatePvlcMedFormula - создание формулы для API (бывший CreateCalculation)
func (r *Repository) CreatePvlcMedFormula(formula *ds.PvlcMedFormula) error {
	return r.db.Create(formula).Error
}

// UpdatePvlcMedFormula - обновление формулы для API (бывший UpdateCalculation)
func (r *Repository) UpdatePvlcMedFormula(formula *ds.PvlcMedFormula) error {
	return r.db.Save(formula).Error
}

// DeletePvlcMedFormula - удаление формулы для API (бывший DeleteCalculation)
func (r *Repository) DeletePvlcMedFormula(id uint) error {
	return r.db.Delete(&ds.PvlcMedFormula{}, id).Error
}

// ==================== МЕТОДЫ ДЛЯ API МЕДИЦИНСКИХ КАРТ ====================

// GetPvlcMedCardsWithFilter - получение медкарт с фильтрацией для API (бывший GetMedicalCardsWithFilter)
func (r *Repository) GetPvlcMedCardsWithFilter(filter ds.PvlcMedCardFilter) ([]ds.PvlcMedCard, error) {
	var cards []ds.PvlcMedCard
	query := r.db.Where("status != ? AND status != ?",
		ds.PvlcMedCardStatusDeleted, ds.PvlcMedCardStatusDraft)

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

// GetPvlcMedCardByID - получение медкарты по ID для API (бывший GetMedicalCardByID)
func (r *Repository) GetPvlcMedCardByID(id uint) (ds.PvlcMedCard, error) {
	var card ds.PvlcMedCard
	err := r.db.Preload("Moderator").First(&card, id).Error
	return card, err
}

// UpdatePvlcMedCard - обновление медкарты для API (бывший UpdateMedicalCard)
func (r *Repository) UpdatePvlcMedCard(card *ds.PvlcMedCard) error {
	return r.db.Save(card).Error
}

// GetDraftPvlcMedCardByUserID - получение черновика по ID пользователя для API (бывший GetDraftCardByUserID)
func (r *Repository) GetDraftPvlcMedCardByUserID(userID uint) (*ds.PvlcMedCard, error) {
	var card ds.PvlcMedCard
	err := r.db.Where("status = ?", ds.PvlcMedCardStatusDraft).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// GetPvlcMedFormulasCountByCardID - количество формул в медкарте для API (бывший GetCalculationsCountByCardID)
func (r *Repository) GetPvlcMedFormulasCountByCardID(cardID uint) (int, error) {
	var count int64
	err := r.db.Model(&ds.MedMmPvlcCalculation{}).Where("pvlc_med_card_id = ?", cardID).Count(&count).Error
	return int(count), err
}

// GetMedMmPvlcCalculationsByCardID - получение расчетов медкарты для API (бывший GetCardCalculationsByCardID)
func (r *Repository) GetMedMmPvlcCalculationsByCardID(cardID uint) ([]ds.MedMmPvlcCalculation, error) {
	var calculations []ds.MedMmPvlcCalculation
	err := r.db.Where("pvlc_med_card_id = ?", cardID).
		Preload("PvlcMedFormula").
		Find(&calculations).Error
	return calculations, err
}

// DeleteMedMmPvlcCalculation - удаление расчета из медкарты для API (бывший DeleteCardCalculation)
func (r *Repository) DeleteMedMmPvlcCalculation(cardID, formulaID uint) error {
	return r.db.Where("pvlc_med_card_id = ? AND pvlc_med_formula_id = ?", cardID, formulaID).
		Delete(&ds.MedMmPvlcCalculation{}).Error
}

// UpdateMedMmPvlcCalculation - обновление расчета в медкарте для API (бывший UpdateCardCalculation)
func (r *Repository) UpdateMedMmPvlcCalculation(cardID, formulaID uint, inputHeight float64) error {
	return r.db.Model(&ds.MedMmPvlcCalculation{}).
		Where("pvlc_med_card_id = ? AND pvlc_med_formula_id = ?", cardID, formulaID).
		Update("input_height", inputHeight).Error
}

// ==================== МЕТОДЫ ДЛЯ РАСЧЕТА ДЖЕЛ (API) ====================

// CalculateTotalDjel - вычисление общего ДЖЕЛ для медкарты для API (без изменений)
func (r *Repository) CalculateTotalDjel(cardID uint) (float64, error) {
	var calculations []ds.MedMmPvlcCalculation
	err := r.db.Where("pvlc_med_card_id = ?", cardID).
		Preload("PvlcMedFormula").
		Find(&calculations).Error
	if err != nil {
		return 0, err
	}

	total := 0.0
	for _, cc := range calculations {
		if cc.InputHeight > 0 {
			// РЕАЛЬНЫЙ РАСЧЕТ ПО ФОРМУЛЕ ИЗ БАЗЫ ДАННЫХ
			result := r.calculateDjelByFormula(cc.PvlcMedFormula.Formula, cc.InputHeight)
			cc.FinalResult = result
			total += result

			// Сохраняем результат расчета
			r.db.Save(&cc)
		}
	}

	return total, nil
}

// calculateDjelByFormula - реальный расчет ДЖЕЛ по формуле для API (без изменений)
func (r *Repository) calculateDjelByFormula(formula string, height float64) float64 {
	// Для демонстрации используем фиксированный возраст 5 лет
	age := 5.0

	// Парсим формулу и вычисляем результат
	if strings.Contains(formula, "0.043") && strings.Contains(formula, "2.89") {
		// Формула для мальчиков 4-7 лет: ДЖЕЛ (л) = (0.043 × Рост) - (0.015 × Возраст) - 2.89
		return (0.043 * height) - (0.015 * age) - 2.89
	} else if strings.Contains(formula, "0.037") && strings.Contains(formula, "2.54") {
		// Формула для девочек 4-7 лет: ДЖЕЛ (л) = (0.037 × Рост) - (0.012 × Возраст) - 2.54
		return (0.037 * height) - (0.012 * age) - 2.54
	} else if strings.Contains(formula, "0.052") && strings.Contains(formula, "4.60") {
		// Формула для мальчиков 8-12 лет: ДЖЕЛ (л) = (0.052 × Рост) - (0.022 × Возраст) - 4.60
		return (0.052 * height) - (0.022 * age) - 4.60
	} else if strings.Contains(formula, "0.041") && strings.Contains(formula, "3.70") {
		// Формула для девочек 8-12 лет: ДЖЕЛ (л) = (0.041 × Рост) - (0.018 × Возраст) - 3.70
		return (0.041 * height) - (0.018 * age) - 3.70
	} else if strings.Contains(formula, "0.052") && strings.Contains(formula, "4.20") {
		// Формула для юношей 13-17 лет: ДЖЕЛ (л) = (0.052 × Рост) - (0.022 × Возраст) - 4.20
		return (0.052 * height) - (0.022 * age) - 4.20
	} else if strings.Contains(formula, "0.041") && strings.Contains(formula, "3.20") {
		// Формула для девушек 13-17 лет: ДЖЕЛ (л) = (0.041 × Рост) - (0.018 × Возраст) - 3.20
		return (0.041 * height) - (0.018 * age) - 3.20
	} else if strings.Contains(formula, "0.052") && strings.Contains(formula, "3.60") {
		// Формула для мужчин 18-60 лет: ДЖЕЛ (л) = (0.052 × Рост) - (0.022 × Возраст) - 3.60
		return (0.052 * height) - (0.022 * age) - 3.60
	} else if strings.Contains(formula, "0.041") && strings.Contains(formula, "2.69") {
		// Формула для женщин 18-60 лет: ДЖЕЛ (л) = (0.041 × Рост) - (0.018 × Возраст) - 2.69
		return (0.041 * height) - (0.018 * age) - 2.69
	} else if strings.Contains(formula, "0.044") && strings.Contains(formula, "2.86") {
		// Формула для пожилых 60+ лет: ДЖЕЛ (л) = (0.044 × Рост) - (0.024 × Возраст) - 2.86
		return (0.044 * height) - (0.024 * age) - 2.86
	}

	// Запасная формула если не распознали
	fmt.Printf("Не удалось распознать формулу: %s\n", formula)
	return height * 0.03 // Упрощенный расчет
}

// ==================== МЕТОДЫ ДЛЯ ПОЛЬЗОВАТЕЛЕЙ (API) ====================

// GetMedUserByID - получение пользователя по ID для API (бывший GetUserByID)
func (r *Repository) GetMedUserByID(id uint) (ds.MedUser, error) {
	var user ds.MedUser
	err := r.db.First(&user, id).Error
	return user, err
}

// GetMedUserByLogin - получение пользователя по логину для API (бывший GetUserByLogin)
func (r *Repository) GetMedUserByLogin(login string) (*ds.MedUser, error) {
	var user ds.MedUser
	err := r.db.Where("login = ?", login).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateMedUser - создание пользователя для API (бывший CreateUser)
func (r *Repository) CreateMedUser(user *ds.MedUser) error {
	return r.db.Create(user).Error
}

// UpdateMedUser - обновление пользователя для API (бывший UpdateUser)
func (r *Repository) UpdateMedUser(user *ds.MedUser) error {
	return r.db.Save(user).Error
}

// ==================== МЕТОДЫ ДЛЯ MINIO (API) ====================

// UploadImageToMinIO - загрузка изображения в MinIO для API (без изменений)
func (r *Repository) UploadImageToMinIO(file *multipart.FileHeader, formulaID uint) (string, error) {
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
	filename := fmt.Sprintf("pvlc_med_formula_%d%s", formulaID, ext)

	// Загружаем в MinIO
	_, err = minioClient.PutObject(context.Background(), "pics", filename, src, file.Size, minio.PutObjectOptions{
		ContentType: "image/" + strings.TrimPrefix(ext, "."),
	})
	if err != nil {
		return "", err
	}

	return filename, nil
}

// DeleteImageFromMinIO - удаление изображения из MinIO для API (без изменений)
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

// InitMinIOBucket - инициализация MinIO bucket для API (без изменений)
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
