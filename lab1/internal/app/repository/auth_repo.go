// internal/app/repository/auth_repo.go
package repository

import (
	"lab1/internal/app/ds"
)

// GetPvlcMedCardsByUserID возвращает заявки конкретного пользователя
// ИСПРАВЛЕНО: включаем черновики в список
func (r *Repository) GetPvlcMedCardsByUserID(userID uint, filter ds.PvlcMedCardFilter) ([]ds.PvlcMedCard, error) {
	var cards []ds.PvlcMedCard
	// ИСПРАВЛЕНО: убираем исключение черновиков
	query := r.db.Where("user_id = ? AND status != ?", userID, ds.PvlcMedCardStatusDeleted)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.DateFrom != "" {
		query = query.Where("created_at >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("created_at <= ?", filter.DateTo)
	}

	err := query.Preload("Moderator").Find(&cards).Error

	// ИСПРАВЛЕНО: гарантируем возврат пустого массива вместо nil
	if cards == nil {
		cards = make([]ds.PvlcMedCard, 0)
	}

	return cards, err
}

// GetPvlcMedCardsForModerator возвращает все заявки для модератора
// ИСПРАВЛЕНО: включаем черновики в список
func (r *Repository) GetPvlcMedCardsForModerator(filter ds.PvlcMedCardFilter) ([]ds.PvlcMedCard, error) {
	var cards []ds.PvlcMedCard
	// ИСПРАВЛЕНО: убираем исключение черновиков
	query := r.db.Where("status != ?", ds.PvlcMedCardStatusDeleted)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.DateFrom != "" {
		query = query.Where("created_at >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("created_at <= ?", filter.DateTo)
	}

	err := query.Preload("Moderator").Preload("User").Find(&cards).Error

	// ИСПРАВЛЕНО: гарантируем возврат пустого массива вместо nil
	if cards == nil {
		cards = make([]ds.PvlcMedCard, 0)
	}

	return cards, err
}

// UpdatePvlcMedCardUserID обновляет владельца заявки
func (r *Repository) UpdatePvlcMedCardUserID(cardID uint, userID uint) error {
	return r.db.Model(&ds.PvlcMedCard{}).Where("id = ?", cardID).Update("user_id", userID).Error
}
