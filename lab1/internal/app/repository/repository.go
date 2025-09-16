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
	_ = godotenv.Load()
	dsnString := dsn.FromEnv()
	db, err := gorm.Open(postgres.Open(dsnString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

// гет заявки по статусу активности
func (r *Repository) GetServices() ([]ds.Calculation, error) {
	var calculations []ds.Calculation
	err := r.db.Where("is_active = ?", true).Find(&calculations).Error
	if err != nil {
		return nil, err
	}
	return calculations, nil
}

// поиск по названию
func (r *Repository) GetServicesByTitle(title string) ([]ds.Calculation, error) {
	var calculations []ds.Calculation
	err := r.db.Where("is_active = ? AND title ILIKE ?", true, "%"+title+"%").Find(&calculations).Error
	if err != nil {
		return nil, err
	}
	return calculations, nil
}

// гет расчет по айди
func (r *Repository) GetService(id int) (ds.Calculation, error) {
	var calculation ds.Calculation
	err := r.db.Where("is_active = ? AND id = ?", true, id).First(&calculation).Error
	if err != nil {
		return ds.Calculation{}, err
	}
	return calculation, nil
}

// найти создать черновик
func (r *Repository) GetOrCreateDraftCard(userID uint) (*ds.MedicalCard, error) {
	var card ds.MedicalCard

	// найти
	err := r.db.Where("status = ?", ds.MedicalCardStatusDraft).First(&card).Error

	if err != nil {
		// создать
		card = ds.MedicalCard{
			Status:      ds.MedicalCardStatusDraft,
			PatientName: "Врач", // пока так хз
		}

		err = r.db.Create(&card).Error
		if err != nil {
			return nil, err
		}
	}

	return &card, nil
}

// добавить расчёт в медкарту
func (r *Repository) AddCalculationToCard(cardID uint, calculationID uint) error {
	// проверяем дубликаты
	var existingLink ds.CardCalculation
	err := r.db.Where("medical_card_id = ? AND calculation_id = ?", cardID, calculationID).First(&existingLink).Error

	// не добавляем если есть дубликат
	if err == nil {
		return nil
	}

	// новая связь
	link := ds.CardCalculation{
		MedicalCardID: cardID,
		CalculationID: calculationID,
		InputHeight:   0, // пока 0
		FinalResult:   0, // пока 0
	}

	return r.db.Create(&link).Error
}

// гет кол-во в корзинке
func (r *Repository) GetCalculationsCount() (int, error) {
	var count int64

	// найти черновик
	var card ds.MedicalCard
	err := r.db.Where("status = ?", ds.MedicalCardStatusDraft).First(&card).Error
	if err != nil {
		// нет черновика - 0
		return 0, nil
	}

	// считаем связанные расчёты
	err = r.db.Model(&ds.CardCalculation{}).Where("medical_card_id = ?", card.ID).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// гет расчёты
func (r *Repository) GetCalculation() ([]ds.Calculation, error) {
	// найти (майонез) черновик
	var card ds.MedicalCard
	err := r.db.Where("status = ?", ds.MedicalCardStatusDraft).First(&card).Error
	if err != nil {
		// пустой массив если нет черновика
		return []ds.Calculation{}, nil
	}

	// получаем все расчеты, связанные с этой картой
	var calculations []ds.Calculation
	err = r.db.Table("calculations").
		Joins("JOIN card_calculations ON calculations.id = card_calculations.calculation_id").
		Where("card_calculations.medical_card_id = ?", card.ID).
		Find(&calculations).Error

	if err != nil {
		return nil, err
	}

	return calculations, nil
}

// с параметром роста
func (r *Repository) GetCalculationWithHeight() ([]ds.CalculationWithHeight, error) {
	// найти черновик
	var card ds.MedicalCard
	err := r.db.Where("status = ?", ds.MedicalCardStatusDraft).First(&card).Error
	if err != nil {
		// нет черновика - возвращаем пустой массив
		return []ds.CalculationWithHeight{}, nil
	}

	// гет расчеты
	var calculations []ds.CalculationWithHeight
	err = r.db.Table("calculations c").
		Select("c.*, cc.input_height").
		Joins("JOIN card_calculations cc ON c.id = cc.calculation_id").
		Where("cc.medical_card_id = ?", card.ID).
		Scan(&calculations).Error

	if err != nil {
		return nil, err
	}

	return calculations, nil
}

// удаляем апдейтом ненапрямую
func (r *Repository) DeleteDraftCard() error {
	// !!!!!
	return r.db.Exec("UPDATE medical_cards SET status = ? WHERE status = ?",
		ds.MedicalCardStatusDeleted, ds.MedicalCardStatusDraft).Error
}
