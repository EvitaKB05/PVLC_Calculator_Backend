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

func convertFormulaToView(formula ds.PvlcMedFormula) ServiceView {
	return ServiceView{
		ID:          formula.ID,
		Title:       formula.Title,
		Description: formula.Description,
		Formula:     formula.Formula,
		Image:       formula.ImageURL,
		Category:    formula.Category,
		Gender:      formula.Gender,
		MinAge:      formula.MinAge,
		MaxAge:      formula.MaxAge,
		Height:      "",
		Result:      "",
	}
}

// GetDjelPatients -> GetPvlcPatients
func (h *Handler) GetPvlcPatients(ctx *gin.Context) {
	var formulas []ds.PvlcMedFormula
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		// GetPvlcMedFormulas -> GetActivePvlcMedFormulas
		formulas, err = h.Repository.GetActivePvlcMedFormulas()
		if err != nil {
			logrus.Error(err)
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка базы данных"})
			return
		}
	} else {
		// GetPvlcMedFormulasByTitle -> GetPvlcMedFormulasByTitle
		formulas, err = h.Repository.GetPvlcMedFormulasByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка поиска"})
			return
		}
	}

	// обработчик обновы количества элементов в корзинке-папке
	// GetPvlcMedFormulasCount -> GetPvlcMedFormulasCount
	count, err := h.Repository.GetPvlcMedFormulasCount()
	if err != nil {
		logrus.Error("Error getting cart count:", err)
		count = 0
	}

	// получаем ID черновика для ссылки
	// GetPvlcMedDraftCardID -> GetDraftPvlcMedCardID
	draftCardID, err := h.Repository.GetDraftPvlcMedCardID()
	if err != nil {
		logrus.Error("Error getting draft card ID:", err)
		draftCardID = 0
	}

	// конверт шаблон
	var services []ServiceView
	for _, formula := range formulas {
		// convertToView -> convertFormulaToView
		services = append(services, convertFormulaToView(formula))
	}

	//  djel_patients.html -> pvlc_patients.html
	ctx.HTML(http.StatusOK, "pvlc_patients.html", gin.H{
		"time":              time.Now().Format("15:04:05"),
		"services":          services,
		"query":             searchQuery,
		"calculationsCount": count,       // актуально кол-во передача
		"draftCardID":       draftCardID, // ДОБАВЛЯЕМ ID ЧЕРНОВИКА
	})
}

// GetDjelPatient -> GetPvlcPatient
func (h *Handler) GetPvlcPatient(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Неверный ID"})
		return
	}

	//  GetPvlcMedFormula -> GetPvlcMedFormulaByIDForHTML
	formula, err := h.Repository.GetPvlcMedFormulaByIDForHTML(id)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Формула не найдена"})
		return
	}

	// гет элементов в корзине
	// GetPvlcMedFormulasCount -> GetPvlcMedFormulasCount
	count, err := h.Repository.GetPvlcMedFormulasCount()
	if err != nil {
		logrus.Error("Error getting cart count:", err)
		count = 0
	}

	// получаем ID черновика для ссылки
	//  GetPvlcMedDraftCardID -> GetDraftPvlcMedCardID
	draftCardID, err := h.Repository.GetDraftPvlcMedCardID()
	if err != nil {
		logrus.Error("Error getting draft card ID:", err)
		draftCardID = 0
	}

	// конверт шаблон
	// convertToView -> convertFormulaToView
	service := convertFormulaToView(formula)

	//  djel_patient.html -> pvlc_patient.html
	ctx.HTML(http.StatusOK, "pvlc_patient.html", gin.H{
		"service":           service,
		"calculationsCount": count,
		"draftCardID":       draftCardID, // ДОБАВЛЯЕМ ID ЧЕРНОВИКА
	})
}

// GetDjelRequest -> GetPvlcMedCalc
func (h *Handler) GetPvlcMedCalc(ctx *gin.Context) {
	// получаем ID из URL
	idStr := ctx.Param("id")
	cardID, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error("Invalid card ID:", err)
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Неверный ID карты"})
		return
	}

	// ПРОВЕРЯЕМ СУЩЕСТВОВАНИЕ КАРТЫ
	//  CheckPvlcMedCardExists -> CheckPvlcMedCardExists
	exists, err := h.Repository.CheckPvlcMedCardExists(uint(cardID))
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

	// корзина по ID карты
	//  GetPvlcMedFormulasByCardID -> GetPvlcMedFormulasByCardIDForHTML
	formulas, err := h.Repository.GetPvlcMedFormulasByCardIDForHTML(uint(cardID))
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка загрузки корзины"})
		return
	}

	// элементы
	//  GetPvlcMedFormulasCount -> GetPvlcMedFormulasCount
	count, err := h.Repository.GetPvlcMedFormulasCount()
	if err != nil {
		logrus.Error("Error getting cart count:", err)
		count = 0
	}

	// фри врачи
	//  GetAvailableMedDoctors -> GetAvailableDoctors
	doctors := h.Repository.GetAvailableDoctors()

	// текущий врач
	//  GetCurrentMedDoctor -> GetCurrentDoctor
	currentDoctor, err := h.Repository.GetCurrentDoctor()
	if err != nil {
		logrus.Error("Error getting current doctor:", err)
		currentDoctor = "Иванов И.И."
	}

	// темплейт конверт
	var services []ServiceView
	for _, formula := range formulas {
		//  convertToView -> convertFormulaToView
		services = append(services, convertFormulaToView(formula))
	}

	//  djel_request.html -> pvlc_med_calc.html
	ctx.HTML(http.StatusOK, "pvlc_med_calc.html", gin.H{
		"calculations":      services,
		"calculationsCount": count,
		"doctors":           doctors,       //
		"currentDoctor":     currentDoctor, //
		"cardID":            cardID,        // ДОБАВЛЯЕМ ID КАРТЫ В ШАБЛОН
	})
}
