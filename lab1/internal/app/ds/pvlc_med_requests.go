// internal/app/ds/pvlc_med_requests.go
package ds

// Запросы для API

type CreatePvlcMedFormulaRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Formula     string `json:"formula" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Gender      string `json:"gender" binding:"required"`
	MinAge      int    `json:"min_age" binding:"required"`
	MaxAge      int    `json:"max_age" binding:"required"`
}

type UpdatePvlcMedFormulaRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Formula     string `json:"formula"`
	Category    string `json:"category"`
	Gender      string `json:"gender"`
	MinAge      int    `json:"min_age"`
	MaxAge      int    `json:"max_age"`
	IsActive    *bool  `json:"is_active"`
}

type PvlcMedCardFilter struct {
	Status   string `form:"status"`
	DateFrom string `form:"date_from"`
	DateTo   string `form:"date_to"`
}

type PvlcMedFormulaFilter struct {
	Category string `form:"category"`
	Gender   string `form:"gender"`
	MinAge   int    `form:"min_age"`
	MaxAge   int    `form:"max_age"`
	Active   *bool  `form:"active"`
}

type UpdateMedMmPvlcCalculationRequest struct {
	InputHeight float64 `json:"input_height" binding:"required"`
}

type MedUserRegistrationRequest struct {
	Login       string `json:"login" binding:"required"`
	Password    string `json:"password" binding:"required"`
	IsModerator bool   `json:"is_moderator"`
}

type UpdateMedUserRequest struct {
	Password string `json:"password"`
}

// ДОБАВЛЕНО ДЛЯ ИСПРАВЛЕНИЯ ОШИБОК SWAGGER

// UpdatePvlcMedCardRequest - запрос на обновление заявки
type UpdatePvlcMedCardRequest struct {
	PatientName string `json:"patient_name"`
	DoctorName  string `json:"doctor_name"`
}

// CompletePvlcMedCardRequest - запрос на завершение заявки
type CompletePvlcMedCardRequest struct {
	Action string `json:"action" binding:"required"` // "complete" или "reject"
}
