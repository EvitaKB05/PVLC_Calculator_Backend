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

// ServiceView - структура для отображения в шаблоне
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
	Height      string // временно оставляем как было
	Result      string // временно оставляем как было
}

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// convertToView преобразует ds.Calculation в ServiceView для шаблона
func convertToView(calc ds.Calculation) ServiceView {
	return ServiceView{
		ID:          calc.ID,
		Title:       calc.Title,
		Description: calc.Description,
		Formula:     calc.Formula,
		Image:       calc.ImageURL, // Маппим ImageURL -> Image
		Category:    calc.Category,
		Gender:      calc.Gender,
		MinAge:      calc.MinAge,
		MaxAge:      calc.MaxAge,
		Height:      "", // Пока оставляем пустым
		Result:      "", // Пока оставляем пустым
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

	// Конвертируем для шаблона
	var services []ServiceView
	for _, calc := range calculations {
		services = append(services, convertToView(calc))
	}

	// Временная заглушка - потом реализуем настоящий счетчик
	calculationsCount := 0

	ctx.HTML(http.StatusOK, "services.html", gin.H{
		"time":              time.Now().Format("15:04:05"),
		"services":          services, // Теперь передаем преобразованные данные
		"query":             searchQuery,
		"calculationsCount": calculationsCount,
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

	// Конвертируем для шаблона
	service := convertToView(calculation)

	ctx.HTML(http.StatusOK, "service.html", gin.H{
		"service":           service,
		"calculationsCount": 0,
	})
}

// Временная заглушка - потом реализуем
func (h *Handler) GetCalculation(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "calculation.html", gin.H{
		"calculations":      []ServiceView{},
		"calculationsCount": 0,
	})
}
