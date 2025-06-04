package models

import (
	"time"
)

type TechnicalReport struct {
	ID              string    `json:"id" db:"id"`
	MaintenanceID   string    `json:"maintenance_id" db:"maintenance_id"`
	UserID          string    `json:"user_id" db:"user_id"`
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description" db:"description"`
	Findings        string    `json:"findings" db:"findings"`
	Recommendations string    `json:"recommendations" db:"recommendations"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type CreateTechnicalReportInput struct {
	MaintenanceID   string `json:"maintenance_id" validate:"required"`
	Title           string `json:"title" validate:"required,min=5,max=100"`
	Description     string `json:"description" validate:"required,min=10,max=1000"`
	Findings        string `json:"findings" validate:"required,min=10,max=1000"`
	Recommendations string `json:"recommendations" validate:"required,min=10,max=1000"`
}

type UpdateTechnicalReportInput struct {
	Title           *string `json:"title,omitempty" validate:"omitempty,min=5,max=100"`
	Description     *string `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	Findings        *string `json:"findings,omitempty" validate:"omitempty,min=10,max=1000"`
	Recommendations *string `json:"recommendations,omitempty" validate:"omitempty,min=10,max=1000"`
}

type TechnicalReportWithDetails struct {
	TechnicalReport
	MaintenanceOrder *MaintenanceOrderWithDetails `json:"maintenance_order,omitempty"`
	User             *User                        `json:"user,omitempty"`
}
