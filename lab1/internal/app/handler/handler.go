package handler

import (
	"database/sql"
	"lab1/internal/app/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
	DB         *sql.DB // для чистых SQL запросов
}

func NewHandler(r *repository.Repository, db *sql.DB) *Handler {
	return &Handler{
		Repository: r,
		DB:         db,
	}
}

func (h *Handler) GetServices(ctx *gin.Context) {
	// по айди замена
	userID := uint(1)

	services, err := h.Repository.GetServices()
	if err != nil {
		logrus.Error("Ошибка получения услуг:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	ordersCount, err := h.Repository.GetOrdersCount(userID)
	if err != nil {
		logrus.Error("Ошибка получения количества заявок:", err)
		ordersCount = 0
	}

	ctx.HTML(http.StatusOK, "services.html", gin.H{
		"time":          time.Now().Format("15:04:05"),
		"services":      services,
		"query":         ctx.Query("query"),
		"ordersCount":   ordersCount,
		"hasDraftOrder": ordersCount > 0,
	})
}

func (h *Handler) GetService(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error("Неверный ID услуги:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	service, err := h.Repository.GetServiceByID(uint(id))
	if err != nil {
		logrus.Error("Ошибка получения услуги:", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Услуга не найдена"})
		return
	}

	ctx.HTML(http.StatusOK, "service.html", gin.H{
		"service": service,
	})
}

func (h *Handler) GetOrder(ctx *gin.Context) {
	// замена на айди
	userID := uint(1)

	order, err := h.Repository.GetDraftOrder(userID)
	if err != nil {
		logrus.Error("Ошибка получения заявки:", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
		return
	}

	orderServices, err := h.Repository.GetOrderServices(order.ID)
	if err != nil {
		logrus.Error("Ошибка получения услуг заявки:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"order":         order,
		"orderServices": orderServices,
	})
}

func (h *Handler) AddServiceToOrder(ctx *gin.Context) {
	// айди айди айди
	userID := uint(1)

	serviceIDStr := ctx.PostForm("service_id")
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		logrus.Error("Неверный ID услуги:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID услуги"})
		return
	}

	// черновик гет или криэйт
	order, err := h.Repository.GetDraftOrder(userID)
	if err != nil {
		// криэйт если нету
		order, err = h.Repository.CreateOrder(userID)
		if err != nil {
			logrus.Error("Ошибка создания заявки:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания заявки"})
			return
		}
	}

	// добавить услугу
	err = h.Repository.AddServiceToOrder(order.ID, uint(serviceID))
	if err != nil {
		logrus.Error("Ошибка добавления услуги в заявку:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления услуги"})
		return
	}

	ctx.Redirect(http.StatusFound, "/services")
}

func (h *Handler) DeleteOrder(ctx *gin.Context) {
	orderIDStr := ctx.Param("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		logrus.Error("Неверный ID заявки:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID заявки"})
		return
	}

	// сырой запрос скл
	query := "UPDATE orders SET status = 'удалён' WHERE id = $1"
	_, err = h.DB.Exec(query, orderID)
	if err != nil {
		logrus.Error("Ошибка удаления заявки:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления заявки"})
		return
	}

	ctx.Redirect(http.StatusFound, "/services")
}
