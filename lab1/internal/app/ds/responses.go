package ds

import "time"

// Ответы для API

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type CalculationResponse struct {
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

type MedicalCardResponse struct {
	ID           uint                      `json:"id"`
	Status       string                    `json:"status"`
	CreatedAt    time.Time                 `json:"created_at"`
	PatientName  string                    `json:"patient_name"`
	DoctorName   string                    `json:"doctor_name"`
	FinalizedAt  *time.Time                `json:"finalized_at,omitempty"`
	CompletedAt  *time.Time                `json:"completed_at,omitempty"`
	TotalResult  float64                   `json:"total_result"`
	Calculations []CardCalculationResponse `json:"calculations"`
}

type CardCalculationResponse struct {
	CalculationID uint    `json:"calculation_id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Formula       string  `json:"formula"`
	ImageURL      string  `json:"image_url"`
	InputHeight   float64 `json:"input_height"`
	FinalResult   float64 `json:"final_result"`
}

type CartIconResponse struct {
	CardID    uint `json:"card_id"`
	ItemCount int  `json:"item_count"`
}

type UserResponse struct {
	ID          uint   `json:"id"`
	Login       string `json:"login"`
	IsModerator bool   `json:"is_moderator"`
}
