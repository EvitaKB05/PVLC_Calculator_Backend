package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ✅ ПЕРЕИМЕНОВАНО: AddToCart -> AddPvlcMedFormulaToCart
func (h *Handler) AddPvlcMedFormulaToCart(ctx *gin.Context) {
	// гет айди
	idStr := ctx.Param("id")
	formulaID, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error("Invalid formula ID:", err)
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Неверный ID формулы"})
		return
	}

	// пока тупо берем по айди
	userID := uint(1)

	// создать черновик
	// ✅ ИСПРАВЛЕНО: GetOrCreatePvlcMedDraftCard -> GetOrCreateDraftPvlcMedCard
	card, err := h.Repository.GetOrCreateDraftPvlcMedCard(userID)
	if err != nil {
		logrus.Error("Error getting draft card:", err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка доступа к корзине"})
		return
	}

	// добавить в папку-корзину
	// ✅ ИСПРАВЛЕНО: AddPvlcMedFormulaToCard -> AddPvlcMedFormulaToCard
	err = h.Repository.AddPvlcMedFormulaToCard(card.ID, uint(formulaID))
	if err != nil {
		logrus.Error("Error adding to cart:", err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка добавления в корзину"})
		return
	}

	// обновленное кол-во уведомление
	// ✅ ИСПРАВЛЕНО: GetPvlcMedFormulasCount -> GetPvlcMedFormulasCount
	count, err := h.Repository.GetPvlcMedFormulasCount()
	if err != nil {
		logrus.Error("Error getting cart count:", err)
		// нон стоп для выполнения
		count = 0
	}

	// получаем ID черновика для ссылки
	// ✅ ИСПРАВЛЕНО: GetPvlcMedDraftCardID -> GetDraftPvlcMedCardID
	draftCardID, err := h.Repository.GetDraftPvlcMedCardID()
	if err != nil {
		logrus.Error("Error getting draft card ID:", err)
		draftCardID = 0
	}

	// список услуг
	// ✅ ИСПРАВЛЕНО: GetPvlcMedFormulas -> GetActivePvlcMedFormulas
	formulas, err := h.Repository.GetActivePvlcMedFormulas()
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка базы данных"})
		return
	}

	// шаблон конверт
	var services []ServiceView
	for _, formula := range formulas {
		// ✅ ПЕРЕИМЕНОВАНО: convertToView -> convertFormulaToView
		services = append(services, convertFormulaToView(formula))
	}

	// обновление данных
	// ✅ ПЕРЕИМЕНОВАНО: djel_patients.html -> pvlc_patients.html
	ctx.HTML(http.StatusOK, "pvlc_patients.html", gin.H{
		"time":              ctx.Query("time"), // сохраняем время
		"services":          services,
		"query":             ctx.Query("query"), // сохраняем поисковый запрос
		"calculationsCount": count,              // обновляем счетчик
		"draftCardID":       draftCardID,        // ДОБАВЛЯЕМ ID ЧЕРНОВИКА
	})
}

// ✅ ПЕРЕИМЕНОВАНО: DeleteCart -> DeletePvlcMedCart
func (h *Handler) DeletePvlcMedCart(ctx *gin.Context) {
	// удалить корзину
	// ✅ ИСПРАВЛЕНО: DeletePvlcMedDraftCard -> DeleteDraftPvlcMedCard
	err := h.Repository.DeleteDraftPvlcMedCard()
	if err != nil {
		logrus.Error("Error deleting cart:", err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка удаления корзины"})
		return
	}

	// редирект на главную
	// ✅ ПЕРЕИМЕНОВАНО: /djel_patients -> /pvlc_patients
	ctx.Redirect(http.StatusFound, "/pvlc_patients")
}
