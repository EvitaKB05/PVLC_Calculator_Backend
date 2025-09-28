package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) AddToCart(ctx *gin.Context) {
	// гет айди
	idStr := ctx.Param("id")
	calculationID, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error("Invalid calculation ID:", err)
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Неверный ID расчета"})
		return
	}

	// пока тупо берем по айди
	userID := uint(1)

	// создать черновик
	card, err := h.Repository.GetOrCreateDraftCard(userID)
	if err != nil {
		logrus.Error("Error getting draft card:", err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка доступа к корзине"})
		return
	}

	// добавить в папку-корзину
	err = h.Repository.AddCalculationToCard(card.ID, uint(calculationID))
	if err != nil {
		logrus.Error("Error adding to cart:", err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка добавления в корзину"})
		return
	}

	// обновленное кол-во уведомление
	count, err := h.Repository.GetCalculationsCount()
	if err != nil {
		logrus.Error("Error getting cart count:", err)
		// нон стоп для выполнения
		count = 0
	}

	// получаем ID черновика для ссылки
	draftCardID, err := h.Repository.GetDraftCardID()
	if err != nil {
		logrus.Error("Error getting draft card ID:", err)
		draftCardID = 0
	}

	// список услуг
	calculations, err := h.Repository.GetServices()
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка базы данных"})
		return
	}

	// шаблон конверт
	var services []ServiceView
	for _, calc := range calculations {
		services = append(services, convertToView(calc))
	}

	// обновление данных
	// ПЕРЕИМЕНОВАНО: services.html -> djel_patients.html
	ctx.HTML(http.StatusOK, "djel_patients.html", gin.H{
		"time":              ctx.Query("time"), // сохраняем время
		"services":          services,
		"query":             ctx.Query("query"), // сохраняем поисковый запрос
		"calculationsCount": count,              // обновляем счетчик
		"draftCardID":       draftCardID,        // ДОБАВЛЯЕМ ID ЧЕРНОВИКА
	})
}

func (h *Handler) DeleteCart(ctx *gin.Context) {
	// удалить корзину
	err := h.Repository.DeleteDraftCard()
	if err != nil {
		logrus.Error("Error deleting cart:", err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка удаления корзины"})
		return
	}

	// редирект на главную
	ctx.Redirect(http.StatusFound, "/djel_patients") // ПЕРЕИМЕНОВАНО
}
