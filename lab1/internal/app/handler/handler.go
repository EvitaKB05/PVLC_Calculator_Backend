package handler

import (
	"lab1/internal/app/ds"
	"lab1/internal/app/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// вью - структура для отображения в шаблоне
type ServiceView struct {
	ID          uint
	Title       string
	Description string
	Formula     string
	Image       string
	Category    string
	Gender      string
	MinAge      int
	MaxAge      int
	Height      string
	Result      string
	Comment     string
}

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// конвертим вью в шаблон
func convertToView(calc ds.Calculation) ServiceView {
	return ServiceView{
		ID:          calc.ID,
		Title:       calc.Title,
		Description: calc.Description,
		Formula:     calc.Formula,
		Image:       calc.ImageURL,
		Category:    calc.Category,
		Gender:      calc.Gender,
		MinAge:      calc.MinAge,
		MaxAge:      calc.MaxAge,
		Height:      "",
		Result:      "",
	}
}

func (h *Handler) GetServices(ctx *gin.Context) {
	var calculations []ds.Calculation
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		calculations, err = h.Repository.GetServices()
		if err != nil {
			logrus.Error(err)
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка базы данных"})
			return
		}
	} else {
		calculations, err = h.Repository.GetServicesByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка поиска"})
			return
		}
	}

	// обработчик обновы количества элементов в корзинке-папке
	count, err := h.Repository.GetCalculationsCount()
	if err != nil {
		logrus.Error("Error getting cart count:", err)
		count = 0
	}

	// получаем ID черновика для ссылки
	draftCardID, err := h.Repository.GetDraftCardID()
	if err != nil {
		logrus.Error("Error getting draft card ID:", err)
		draftCardID = 0
	}

	// конверт шаблон
	var services []ServiceView
	for _, calc := range calculations {
		services = append(services, convertToView(calc))
	}

	ctx.HTML(http.StatusOK, "services.html", gin.H{
		"time":              time.Now().Format("15:04:05"),
		"services":          services,
		"query":             searchQuery,
		"calculationsCount": count,       // актуально кол-во передача
		"draftCardID":       draftCardID, // ДОБАВЛЯЕМ ID ЧЕРНОВИКА
	})
}

func (h *Handler) GetService(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Неверный ID"})
		return
	}

	calculation, err := h.Repository.GetService(id)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Услуга не найдена"})
		return
	}

	// гет элементов в корзине
	count, err := h.Repository.GetCalculationsCount()
	if err != nil {
		logrus.Error("Error getting cart count:", err)
		count = 0
	}

	// получаем ID черновика для ссылки
	draftCardID, err := h.Repository.GetDraftCardID()
	if err != nil {
		logrus.Error("Error getting draft card ID:", err)
		draftCardID = 0
	}

	// конверт шаблон
	service := convertToView(calculation)

	ctx.HTML(http.StatusOK, "service.html", gin.H{
		"service":           service,
		"calculationsCount": count,
		"draftCardID":       draftCardID, // ДОБАВЛЯЕМ ID ЧЕРНОВИКА
	})
}

// fix 3 start!
// ИЗМЕНЯЕМ МЕТОД - ДОБАВЛЯЕМ ПАРАМЕТР ID
func (h *Handler) GetCalculation(ctx *gin.Context) {
	// получаем ID из URL
	idStr := ctx.Param("id")
	cardID, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error("Invalid card ID:", err)
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Неверный ID карты"})
		return
	}

	// ПРОВЕРЯЕМ СУЩЕСТВОВАНИЕ КАРТЫ
	exists, err := h.Repository.CheckCardExists(uint(cardID))
	if err != nil {
		logrus.Error("Error checking card existence:", err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка проверки карты"})
		return
	}

	// ЕСЛИ КАРТА НЕ СУЩЕСТВУЕТ ИЛИ УДАЛЕНА - 404
	if !exists {
		logrus.Warnf("Card with ID %d not found or deleted", cardID)
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Карта не найдена"})
		return
	}
	// fix 3 end
	// корзина по ID карты
	calculations, err := h.Repository.GetCalculationByCardID(uint(cardID))
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка загрузки корзины"})
		return
	}

	// элементы
	count, err := h.Repository.GetCalculationsCount()
	if err != nil {
		logrus.Error("Error getting cart count:", err)
		count = 0
	}

	// фри врачи
	doctors := h.Repository.GetAvailableDoctors()

	// текущий врач
	currentDoctor, err := h.Repository.GetCurrentDoctor()
	if err != nil {
		logrus.Error("Error getting current doctor:", err)
		currentDoctor = "Иванов И.И."
	}

	// темплейт конверт
	var services []ServiceView
	for _, calc := range calculations {
		services = append(services, convertToView(calc))
	}

	ctx.HTML(http.StatusOK, "calculation.html", gin.H{
		"calculations":      services,
		"calculationsCount": count,
		"doctors":           doctors,       //
		"currentDoctor":     currentDoctor, //
		"cardID":            cardID,        //
	})
}
