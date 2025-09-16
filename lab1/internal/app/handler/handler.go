package handler

import (
	"fmt"
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

	// конверт шаблон
	var services []ServiceView
	for _, calc := range calculations {
		services = append(services, convertToView(calc))
	}

	ctx.HTML(http.StatusOK, "services.html", gin.H{
		"time":              time.Now().Format("15:04:05"),
		"services":          services,
		"query":             searchQuery,
		"calculationsCount": count, // актуально кол-во передача
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

	// конверт шаблон
	service := convertToView(calculation)

	ctx.HTML(http.StatusOK, "service.html", gin.H{
		"service":           service,
		"calculationsCount": count,
	})
}

func (h *Handler) GetCalculation(ctx *gin.Context) {
	// гет расчеты
	calculations, err := h.Repository.GetCalculationWithHeight()
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка загрузки корзины"})
		return
	}

	// гет кол-во
	count, err := h.Repository.GetCalculationsCount()
	if err != nil {
		logrus.Error("Error getting cart count:", err)
		count = 0
	}

	// конверт
	var services []ServiceView
	for _, calc := range calculations {
		service := convertToView(calc.Calculation)
		service.Height = fmt.Sprintf("%.1f", calc.InputHeight) // рост
		services = append(services, service)
	}

	ctx.HTML(http.StatusOK, "calculation.html", gin.H{
		"calculations":      services,
		"calculationsCount": count,
	})
}
