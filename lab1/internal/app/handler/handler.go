package handler

import (
	"lab1/internal/app/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetServices(ctx *gin.Context) {
	var services []repository.Service
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		services, err = h.Repository.GetServices()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		services, err = h.Repository.GetServicesByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	calculationsCount, err := h.Repository.GetCalculationsCount()
	if err != nil {
		logrus.Error("Ошибка получения количества расчетов:", err)
		calculationsCount = 0
	}

	ctx.HTML(http.StatusOK, "services.html", gin.H{
		"time":              time.Now().Format("15:04:05"),
		"services":          services,
		"query":             searchQuery,
		"calculationsCount": calculationsCount,
	})
}

func (h *Handler) GetService(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	service, err := h.Repository.GetService(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "service.html", gin.H{
		"service":           service,
		"calculationsCount": 0, // валую если надо будет
	})
}

func (h *Handler) GetCalculation(ctx *gin.Context) {
	calculations, err := h.Repository.GetCalculation()
	if err != nil {
		logrus.Error(err)
		calculations = []repository.Service{} // пустой массив для чека ошибок
	}

	ctx.HTML(http.StatusOK, "calculation.html", gin.H{
		"calculations":      calculations,
		"calculationsCount": len(calculations),
	})
}
