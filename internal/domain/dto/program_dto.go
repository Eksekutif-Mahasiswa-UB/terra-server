package dto

import "time"

// CreateProgramRequest represents the request to create a new program
type CreateProgramRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	ImageURL     string  `json:"image_url"`
	TargetAmount float64 `json:"target_amount" binding:"required,gt=0"`
}

// UpdateProgramRequest represents the request to update a program
type UpdateProgramRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	ImageURL     string  `json:"image_url"`
	TargetAmount float64 `json:"target_amount" binding:"required,gt=0"`
}

// ProgramResponse represents the program data transfer object for API responses
type ProgramResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ImageURL     string    `json:"image_url"`
	TargetAmount float64   `json:"target_amount"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
