package api

import (
	"lab1/internal/app/ds"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Домен: Медицинские карты (PvlcMedCards)

// GET /api/cart/icon - иконка корзины
func (a *API) GetCartIcon(c *gin.Context) {
	// Фиксированный пользователь
	userID := uint(1)

	// Получаем черновик
	card, err := a.repo.GetDraftPvlcMedCardByUserID(userID)
	if err != nil {
		// Если черновика нет - возвращаем пустую корзину
		a.successResponse(c, ds.CartIconResponse{
			CardID:    0,
			ItemCount: 0,
		})
		return
	}

	// Считаем количество формул в заявке
	count, err := a.repo.GetPvlcMedFormulasCountByCardID(card.ID)
	if err != nil {
		logrus.Error("Error getting pvlc med formulas count: ", err)
		count = 0
	}

	a.successResponse(c, ds.CartIconResponse{
		CardID:    card.ID,
		ItemCount: count,
	})
}

// GET /api/pvlc-med-cards - список заявок (кроме удаленных и черновика)
func (a *API) GetPvlcMedCards(c *gin.Context) {
	var filter ds.PvlcMedCardFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные параметры фильтрации")
		return
	}

	cards, err := a.repo.GetPvlcMedCardsWithFilter(filter)
	if err != nil {
		logrus.Error("Error getting pvlc med cards: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка получения заявок")
		return
	}

	var response []ds.PvlcMedCardResponse
	for _, card := range cards {
		cardResponse := ds.PvlcMedCardResponse{
			ID:          card.ID,
			Status:      card.Status,
			CreatedAt:   card.CreatedAt,
			PatientName: card.PatientName,
			DoctorName:  card.DoctorName,
			TotalResult: card.TotalResult,
		}

		if card.FinalizedAt != nil {
			cardResponse.FinalizedAt = card.FinalizedAt
		}
		if card.CompletedAt != nil {
			cardResponse.CompletedAt = card.CompletedAt
		}

		// Получаем расчеты для этой заявки
		calculations, err := a.repo.GetMedMmPvlcCalculationsByCardID(card.ID)
		if err == nil {
			for _, calc := range calculations {
				cardResponse.Calculations = append(cardResponse.Calculations, ds.MedMmPvlcCalculationResponse{
					PvlcMedFormulaID: calc.PvlcMedFormulaID,
					Title:            calc.PvlcMedFormula.Title,
					Description:      calc.PvlcMedFormula.Description,
					Formula:          calc.PvlcMedFormula.Formula,
					ImageURL:         calc.PvlcMedFormula.ImageURL,
					InputHeight:      calc.InputHeight,
					FinalResult:      calc.FinalResult,
				})
			}
		}

		response = append(response, cardResponse)
	}

	a.successResponse(c, response)
}

// GET /api/pvlc-med-cards/:id - одна заявка
func (a *API) GetPvlcMedCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID заявки")
		return
	}

	card, err := a.repo.GetPvlcMedCardByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Заявка не найдена")
		return
	}

	// Не возвращаем удаленные заявки
	if card.Status == ds.PvlcMedCardStatusDeleted {
		a.errorResponse(c, http.StatusNotFound, "Заявка не найдена")
		return
	}

	response := ds.PvlcMedCardResponse{
		ID:          card.ID,
		Status:      card.Status,
		CreatedAt:   card.CreatedAt,
		PatientName: card.PatientName,
		DoctorName:  card.DoctorName,
		TotalResult: card.TotalResult,
	}

	if card.FinalizedAt != nil {
		response.FinalizedAt = card.FinalizedAt
	}
	if card.CompletedAt != nil {
		response.CompletedAt = card.CompletedAt
	}

	// Получаем расчеты для этой заявки
	calculations, err := a.repo.GetMedMmPvlcCalculationsByCardID(card.ID)
	if err == nil {
		for _, calc := range calculations {
			response.Calculations = append(response.Calculations, ds.MedMmPvlcCalculationResponse{
				PvlcMedFormulaID: calc.PvlcMedFormulaID,
				Title:            calc.PvlcMedFormula.Title,
				Description:      calc.PvlcMedFormula.Description,
				Formula:          calc.PvlcMedFormula.Formula,
				ImageURL:         calc.PvlcMedFormula.ImageURL,
				InputHeight:      calc.InputHeight,
				FinalResult:      calc.FinalResult,
			})
		}
	}

	a.successResponse(c, response)
}

// PUT /api/pvlc-med-cards/:id - изменение полей заявки
func (a *API) UpdatePvlcMedCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID заявки")
		return
	}

	var request struct {
		PatientName string `json:"patient_name"`
		DoctorName  string `json:"doctor_name"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	card, err := a.repo.GetPvlcMedCardByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Заявка не найдена")
		return
	}

	// Можно менять только черновики
	if card.Status != ds.PvlcMedCardStatusDraft {
		a.errorResponse(c, http.StatusBadRequest, "Можно изменять только черновики")
		return
	}

	if request.PatientName != "" {
		card.PatientName = request.PatientName
	}
	if request.DoctorName != "" {
		card.DoctorName = request.DoctorName
	}

	if err := a.repo.UpdatePvlcMedCard(&card); err != nil {
		logrus.Error("Error updating pvlc med card: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка обновления заявки")
		return
	}

	a.successResponse(c, gin.H{"message": "Заявка успешно обновлена"})
}

// PUT /api/pvlc-med-cards/:id/finalize - сформировать заявку
func (a *API) FinalizePvlcMedCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID заявки")
		return
	}

	card, err := a.repo.GetPvlcMedCardByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Заявка не найдена")
		return
	}

	// Проверяем что заявка в статусе черновика
	if card.Status != ds.PvlcMedCardStatusDraft {
		a.errorResponse(c, http.StatusBadRequest, "Можно формировать только черновики")
		return
	}

	// Проверяем обязательные поля
	if card.PatientName == "" || card.DoctorName == "" {
		a.errorResponse(c, http.StatusBadRequest, "Заполните все обязательные поля (пациент, врач)")
		return
	}

	// Проверяем что есть расчеты
	count, err := a.repo.GetPvlcMedFormulasCountByCardID(card.ID)
	if err != nil || count == 0 {
		a.errorResponse(c, http.StatusBadRequest, "Добавьте хотя бы один расчет в заявку")
		return
	}

	// Меняем статус и устанавливаем дату формирования
	now := time.Now()
	card.Status = ds.PvlcMedCardStatusFormed
	card.FinalizedAt = &now

	if err := a.repo.UpdatePvlcMedCard(&card); err != nil {
		logrus.Error("Error finalizing pvlc med card: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка формирования заявки")
		return
	}

	a.successResponse(c, gin.H{"message": "Заявка успешно сформирована"})
}

// PUT /api/pvlc-med-cards/:id/complete - завершить/отклонить заявку
func (a *API) CompletePvlcMedCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID заявки")
		return
	}

	var request struct {
		Action string `json:"action" binding:"required"` // "complete" или "reject"
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	card, err := a.repo.GetPvlcMedCardByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Заявка не найдена")
		return
	}

	// Проверяем что заявка в статусе сформирована
	if card.Status != ds.PvlcMedCardStatusFormed {
		a.errorResponse(c, http.StatusBadRequest, "Можно завершать/отклонять только сформированные заявки")
		return
	}

	// Фиксированный модератор (как требуется)
	moderatorID := uint(2) // admin пользователь

	now := time.Now()
	if request.Action == "complete" {
		card.Status = ds.PvlcMedCardStatusCompleted

		// ВЫЧИСЛЕНИЕ ДЖЕЛ - реализуем формулу из лабораторной 2
		totalResult, err := a.repo.CalculateTotalDjel(card.ID)
		if err != nil {
			logrus.Error("Error calculating DJEL: ", err)
			a.errorResponse(c, http.StatusInternalServerError, "Ошибка расчета ДЖЕЛ")
			return
		}
		card.TotalResult = totalResult
	} else if request.Action == "reject" {
		card.Status = ds.PvlcMedCardStatusRejected
	} else {
		a.errorResponse(c, http.StatusBadRequest, "Неверное действие. Используйте 'complete' или 'reject'")
		return
	}

	card.CompletedAt = &now
	card.ModeratorID = &moderatorID

	if err := a.repo.UpdatePvlcMedCard(&card); err != nil {
		logrus.Error("Error completing pvlc med card: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка завершения заявки")
		return
	}

	a.successResponse(c, gin.H{
		"message":      "Заявка успешно обработана",
		"status":       card.Status,
		"total_result": card.TotalResult,
	})
}

// DELETE /api/pvlc-med-cards/:id - удаление заявки
func (a *API) DeletePvlcMedCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID заявки")
		return
	}

	card, err := a.repo.GetPvlcMedCardByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Заявка не найдена")
		return
	}

	// Удалять можно только черновики
	if card.Status != ds.PvlcMedCardStatusDraft {
		a.errorResponse(c, http.StatusBadRequest, "Можно удалять только черновики")
		return
	}

	card.Status = ds.PvlcMedCardStatusDeleted
	if err := a.repo.UpdatePvlcMedCard(&card); err != nil {
		logrus.Error("Error deleting pvlc med card: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка удаления заявки")
		return
	}

	a.successResponse(c, gin.H{"message": "Заявка успешно удалена"})
}
