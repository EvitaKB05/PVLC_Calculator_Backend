package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Service struct {
	ID          int
	Title       string
	Description string
	Formula     string
	Image       string
	Category    string
	Gender      string
	MinAge      int
	MaxAge      int
	Height      string // Добавлено поле для роста
	Result      string // Добавлено поле для результата ДЖЕЛ
}

func (r *Repository) GetServices() ([]Service, error) {
	services := []Service{
		{
			ID:          1,
			Title:       "Мальчики 4-7 лет",
			Description: "Расчет ДЖЕЛ для мальчиков дошкольного возраста",
			Formula:     "ДЖЕЛ (л) = (0.043 × Рост) - (0.015 × Возраст) - 2.89",
			Image:       "boys_4_7.png",
			Category:    "дети",
			Gender:      "мужской",
			MinAge:      4,
			MaxAge:      7,
			Height:      "127 см", // Добавлено значение роста
			Result:      "2.57 л", // Добавлено значение результата
		},
		{
			ID:          2,
			Title:       "Девочки 4-7 лет",
			Description: "Расчет ДЖЕЛ для девочек дошкольного возраста",
			Formula:     "ДЖЕЛ (л) = (0.037 × Рост) - (0.012 × Возраст) - 2.54",
			Image:       "girls_4_7.png",
			Category:    "дети",
			Gender:      "женский",
			MinAge:      4,
			MaxAge:      7,
			Height:      "120 см", // Добавлено значение роста
			Result:      "1.90 л", // Добавлено значение результата
		},
		{
			ID:          3,
			Title:       "Мальчики 8-12 лет",
			Description: "Расчет ДЖЕЛ для мальчиков младшего школьного возраста",
			Formula:     "ДЖЕЛ (л) = (0.052 × Рост) - (0.022 × Возраст) - 4.60",
			Image:       "boys_8_12.png",
			Category:    "дети",
			Gender:      "мужской",
			MinAge:      8,
			MaxAge:      12,
			Height:      "",
			Result:      "",
		},
		{
			ID:          4,
			Title:       "Девочки 8-12 лет",
			Description: "Расчет ДЖЕЛ для девочек младшего школьного возраста",
			Formula:     "ДЖЕЛ (л) = (0.041 × Рост) - (0.018 × Возраст) - 3.70",
			Image:       "girls_8_12.png",
			Category:    "дети",
			Gender:      "женский",
			MinAge:      8,
			MaxAge:      12,
			Height:      "",
			Result:      "",
		},
		{
			ID:          5,
			Title:       "Юноши 13-17 лет",
			Description: "Расчет ДЖЕЛ для юношей подросткового возраста",
			Formula:     "ДЖЕЛ (л) = (0.052 × Рост) - (0.022 × Возраст) - 4.20",
			Image:       "boys_13_17.png",
			Category:    "подростки",
			Gender:      "мужской",
			MinAge:      13,
			MaxAge:      17,
			Height:      "",
			Result:      "",
		},
		{
			ID:          6,
			Title:       "Девушки 13-17 лет",
			Description: "Расчет ДЖЕЛ для девушек подросткового возраста",
			Formula:     "ДЖЕЛ (л) = (0.041 × Рост) - (0.018 × Возраст) - 3.20",
			Image:       "girls_13_17.png",
			Category:    "подростки",
			Gender:      "женский",
			MinAge:      13,
			MaxAge:      17,
			Height:      "",
			Result:      "",
		},
		{
			ID:          7,
			Title:       "Мужчины 18-60 лет",
			Description: "Расчет ДЖЕЛ для взрослых мужчин",
			Formula:     "ДЖЕЛ (л) = (0.052 × Рост) - (0.022 × Возраст) - 3.60",
			Image:       "men_18_60.png",
			Category:    "взрослые",
			Gender:      "мужской",
			MinAge:      18,
			MaxAge:      60,
			Height:      "",
			Result:      "",
		},
		{
			ID:          8,
			Title:       "Женщины 18-60 лет",
			Description: "Расчет ДЖЕЛ для взрослых женщин",
			Formula:     "ДЖЕЛ (л) = (0.041 × Рост) - (0.018 × Возраст) - 2.69",
			Image:       "women_18_60.png",
			Category:    "взрослые",
			Gender:      "женский",
			MinAge:      18,
			MaxAge:      60,
			Height:      "",
			Result:      "",
		},
		{
			ID:          9,
			Title:       "Пожилые 60+ лет",
			Description: "Расчет ДЖЕЛ для пожилых людей",
			Formula:     "ДЖЕЛ (л) = (0.044 × Рост) - (0.024 × Возраст) - 2.86",
			Image:       "elderly_60plus.png",
			Category:    "пожилые",
			Gender:      "унисекс",
			MinAge:      60,
			MaxAge:      100,
			Height:      "",
			Result:      "",
		},
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("массив услуг пустой")
	}
	return services, nil
}

func (r *Repository) GetService(id int) (Service, error) {
	services, err := r.GetServices()
	if err != nil {
		return Service{}, err
	}

	for _, service := range services {
		if service.ID == id {
			return service, nil
		}
	}
	return Service{}, fmt.Errorf("услуга не найдена")
}

func (r *Repository) GetServicesByTitle(title string) ([]Service, error) {
	services, err := r.GetServices()
	if err != nil {
		return []Service{}, err
	}

	var result []Service
	for _, service := range services {
		if strings.Contains(strings.ToLower(service.Title), strings.ToLower(title)) {
			result = append(result, service)
		}
	}
	return result, nil
}

func (r *Repository) GetCalculation() ([]Service, error) {
	services, err := r.GetServices()
	if err != nil {
		return []Service{}, err
	}

	var result []Service
	for _, service := range services {
		if service.ID == 1 || service.ID == 2 {
			result = append(result, service)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("массив расчетов пустой")
	}

	return result, nil
}

func (r *Repository) GetCalculationsCount() (int, error) {
	calculations, err := r.GetCalculation()
	if err != nil {
		return 0, err
	}
	return len(calculations), nil
}
