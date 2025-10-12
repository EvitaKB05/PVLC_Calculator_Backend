// internal/app/repository/auth_repo.go
package repository

import (
	"lab1/internal/app/ds"
)

// GetPvlcMedCardsByUserID возвращает заявки конкретного пользователя
// ДОБАВЛЕНО ДЛЯ ЛАБОРАТОРНОЙ РАБОТЫ 4
func (r *Repository) GetPvlcMedCardsByUserID(userID uint, filter ds.PvlcMedCardFilter) ([]ds.PvlcMedCard, error) {
	var cards []ds.PvlcMedCard
	query := r.db.Where("user_id = ? AND status != ? AND status != ?",
		userID, ds.PvlcMedCardStatusDeleted, ds.PvlcMedCardStatusDraft)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.DateFrom != "" {
		// Здесь должна быть логика парсинга даты (упрощенно)
		query = query.Where("created_at >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("created_at <= ?", filter.DateTo)
	}

	err := query.Preload("Moderator").Find(&cards).Error
	return cards, err
}

// GetPvlcMedCardsForModerator возвращает все заявки для модератора
func (r *Repository) GetPvlcMedCardsForModerator(filter ds.PvlcMedCardFilter) ([]ds.PvlcMedCard, error) {
	var cards []ds.PvlcMedCard
	query := r.db.Where("status != ? AND status != ?",
		ds.PvlcMedCardStatusDeleted, ds.PvlcMedCardStatusDraft)

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
	return cards, err
}

// UpdatePvlcMedCardUserID обновляет владельца заявки
func (r *Repository) UpdatePvlcMedCardUserID(cardID uint, userID uint) error {
	return r.db.Model(&ds.PvlcMedCard{}).Where("id = ?", cardID).Update("user_id", userID).Error
}
