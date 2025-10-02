package ds

import "time"

// Ответы для API

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type PvlcMedFormulaResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Formula     string `json:"formula"`
	ImageURL    string `json:"image_url"`
	Category    string `json:"category"`
	Gender      string `json:"gender"`
	MinAge      int    `json:"min_age"`
	MaxAge      int    `json:"max_age"`
	IsActive    bool   `json:"is_active"`
}

type PvlcMedCardResponse struct {
	ID              uint                           `json:"id"`
	Status          string                         `json:"status"`
	CreatedAt       time.Time                      `json:"created_at"`
	PatientName     string                         `json:"patient_name"`
	DoctorName      string                         `json:"doctor_name"`
	FinalizedAt     *time.Time                     `json:"finalized_at,omitempty"`
	CompletedAt     *time.Time                     `json:"completed_at,omitempty"`
	TotalResult     float64                        `json:"total_result"`
	MedCalculations []MedMmPvlcCalculationResponse `json:"med_calculations"`
}

type MedMmPvlcCalculationResponse struct {
	PvlcMedFormulaID uint    `json:"pvlc_med_formula_id"`
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	Formula          string  `json:"formula"`
	ImageURL         string  `json:"image_url"`
	InputHeight      float64 `json:"input_height"`
	FinalResult      float64 `json:"final_result"`
}

type CartIconResponse struct {
	MedCardID    uint `json:"med_card_id"`
	MedItemCount int  `json:"med_item_count"`
}

type MedUserResponse struct {
	ID          uint   `json:"id"`
	Login       string `json:"login"`
	IsModerator bool   `json:"is_moderator"`
}
