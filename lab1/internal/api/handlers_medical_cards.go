package api

import (
	"lab1/internal/app/ds"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Домен: Заявки (Medical Cards)

// GET /api/cart/icon - иконка корзины
func (a *API) GetCartIcon(c *gin.Context) {
	// Фиксированный пользователь
	userID := uint(1)

	// Получаем черновик
	card, err := a.repo.GetDraftCardByUserID(userID)
	if err != nil {
		// Если черновика нет - возвращаем пустую корзину
		a.successResponse(c, ds.CartIconResponse{
			CardID:    0,
			ItemCount: 0,
		})
		return
	}

	// Считаем количество услуг в заявке
	count, err := a.repo.GetCalculationsCountByCardID(card.ID)
	if err != nil {
		logrus.Error("Error getting calculations count: ", err)
		count = 0
	}

	a.successResponse(c, ds.CartIconResponse{
		CardID:    card.ID,
		ItemCount: count,
	})
}

// GET /api/medical-cards - список заявок (кроме удаленных и черновика)
func (a *API) GetMedicalCards(c *gin.Context) {
	var filter ds.MedicalCardFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		a.errorResponse(c, http.StatusBadRequest, "Неверные параметры фильтрации")
		return
	}

	cards, err := a.repo.GetMedicalCardsWithFilter(filter)
	if err != nil {
		logrus.Error("Error getting medical cards: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка получения заявок")
		return
	}

	var response []ds.MedicalCardResponse
	for _, card := range cards {
		cardResponse := ds.MedicalCardResponse{
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
		calculations, err := a.repo.GetCardCalculationsByCardID(card.ID)
		if err == nil {
			for _, calc := range calculations {
				cardResponse.Calculations = append(cardResponse.Calculations, ds.CardCalculationResponse{
					CalculationID: calc.CalculationID,
					Title:         calc.Calculation.Title,
					Description:   calc.Calculation.Description,
					Formula:       calc.Calculation.Formula,
					ImageURL:      calc.Calculation.ImageURL,
					InputHeight:   calc.InputHeight,
					FinalResult:   calc.FinalResult,
				})
			}
		}

		response = append(response, cardResponse)
	}

	a.successResponse(c, response)
}

// GET /api/medical-cards/:id - одна заявка
func (a *API) GetMedicalCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID заявки")
		return
	}

	card, err := a.repo.GetMedicalCardByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Заявка не найдена")
		return
	}

	// Не возвращаем удаленные заявки
	if card.Status == ds.MedicalCardStatusDeleted {
		a.errorResponse(c, http.StatusNotFound, "Заявка не найдена")
		return
	}

	response := ds.MedicalCardResponse{
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
	calculations, err := a.repo.GetCardCalculationsByCardID(card.ID)
	if err == nil {
		for _, calc := range calculations {
			response.Calculations = append(response.Calculations, ds.CardCalculationResponse{
				CalculationID: calc.CalculationID,
				Title:         calc.Calculation.Title,
				Description:   calc.Calculation.Description,
				Formula:       calc.Calculation.Formula,
				ImageURL:      calc.Calculation.ImageURL,
				InputHeight:   calc.InputHeight,
				FinalResult:   calc.FinalResult,
			})
		}
	}

	a.successResponse(c, response)
}

// PUT /api/medical-cards/:id - изменение полей заявки
func (a *API) UpdateMedicalCard(c *gin.Context) {
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

	card, err := a.repo.GetMedicalCardByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Заявка не найдена")
		return
	}

	// Можно менять только черновики
	if card.Status != ds.MedicalCardStatusDraft {
		a.errorResponse(c, http.StatusBadRequest, "Можно изменять только черновики")
		return
	}

	if request.PatientName != "" {
		card.PatientName = request.PatientName
	}
	if request.DoctorName != "" {
		card.DoctorName = request.DoctorName
	}

	if err := a.repo.UpdateMedicalCard(&card); err != nil {
		logrus.Error("Error updating medical card: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка обновления заявки")
		return
	}

	a.successResponse(c, gin.H{"message": "Заявка успешно обновлена"})
}

// PUT /api/medical-cards/:id/finalize - сформировать заявку
func (a *API) FinalizeMedicalCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID заявки")
		return
	}

	card, err := a.repo.GetMedicalCardByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Заявка не найдена")
		return
	}

	// Проверяем что заявка в статусе черновика
	if card.Status != ds.MedicalCardStatusDraft {
		a.errorResponse(c, http.StatusBadRequest, "Можно формировать только черновики")
		return
	}

	// Проверяем обязательные поля
	if card.PatientName == "" || card.DoctorName == "" {
		a.errorResponse(c, http.StatusBadRequest, "Заполните все обязательные поля (пациент, врач)")
		return
	}

	// Проверяем что есть расчеты
	count, err := a.repo.GetCalculationsCountByCardID(card.ID)
	if err != nil || count == 0 {
		a.errorResponse(c, http.StatusBadRequest, "Добавьте хотя бы один расчет в заявку")
		return
	}

	// Меняем статус и устанавливаем дату формирования
	now := time.Now()
	card.Status = ds.MedicalCardStatusFormed
	card.FinalizedAt = &now

	if err := a.repo.UpdateMedicalCard(&card); err != nil {
		logrus.Error("Error finalizing medical card: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка формирования заявки")
		return
	}

	a.successResponse(c, gin.H{"message": "Заявка успешно сформирована"})
}

// PUT /api/medical-cards/:id/complete - завершить/отклонить заявку
func (a *API) CompleteMedicalCard(c *gin.Context) {
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

	card, err := a.repo.GetMedicalCardByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Заявка не найдена")
		return
	}

	// Проверяем что заявка в статусе сформирована
	if card.Status != ds.MedicalCardStatusFormed {
		a.errorResponse(c, http.StatusBadRequest, "Можно завершать/отклонять только сформированные заявки")
		return
	}

	// Фиксированный модератор (как требуется)
	moderatorID := uint(2) // admin пользователь

	now := time.Now()
	if request.Action == "complete" {
		card.Status = ds.MedicalCardStatusCompleted

		// ВЫЧИСЛЕНИЕ ДЖЕЛ - реализуем формулу из лабораторной 2
		totalResult, err := a.repo.CalculateTotalDjel(card.ID)
		if err != nil {
			logrus.Error("Error calculating DJEL: ", err)
			a.errorResponse(c, http.StatusInternalServerError, "Ошибка расчета ДЖЕЛ")
			return
		}
		card.TotalResult = totalResult
	} else if request.Action == "reject" {
		card.Status = ds.MedicalCardStatusRejected
	} else {
		a.errorResponse(c, http.StatusBadRequest, "Неверное действие. Используйте 'complete' или 'reject'")
		return
	}

	card.CompletedAt = &now
	card.ModeratorID = &moderatorID

	if err := a.repo.UpdateMedicalCard(&card); err != nil {
		logrus.Error("Error completing medical card: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка завершения заявки")
		return
	}

	a.successResponse(c, gin.H{
		"message":      "Заявка успешно обработана",
		"status":       card.Status,
		"total_result": card.TotalResult,
	})
}

// DELETE /api/medical-cards/:id - удаление заявки
func (a *API) DeleteMedicalCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		a.errorResponse(c, http.StatusBadRequest, "Неверный ID заявки")
		return
	}

	card, err := a.repo.GetMedicalCardByID(uint(id))
	if err != nil {
		a.errorResponse(c, http.StatusNotFound, "Заявка не найдена")
		return
	}

	// Удалять можно только черновики
	if card.Status != ds.MedicalCardStatusDraft {
		a.errorResponse(c, http.StatusBadRequest, "Можно удалять только черновики")
		return
	}

	card.Status = ds.MedicalCardStatusDeleted
	if err := a.repo.UpdateMedicalCard(&card); err != nil {
		logrus.Error("Error deleting medical card: ", err)
		a.errorResponse(c, http.StatusInternalServerError, "Ошибка удаления заявки")
		return
	}

	a.successResponse(c, gin.H{"message": "Заявка успешно удалена"})
}
