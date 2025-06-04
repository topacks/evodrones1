package models

import (
	"time"
)

type MaintenanceOrderStatus string

const (
	MaintenanceOrderStatusPending    MaintenanceOrderStatus = "pending"
	MaintenanceOrderStatusInProgress MaintenanceOrderStatus = "in_progress"
	MaintenanceOrderStatusCompleted  MaintenanceOrderStatus = "completed"
	MaintenanceOrderStatusCancelled  MaintenanceOrderStatus = "cancelled"
)

type MaintenanceOrder struct {
	ID          string                 `json:"id" db:"id"`
	DroneID     string                 `json:"drone_id" db:"drone_id"`
	UserID      string                 `json:"user_id" db:"user_id"`
	Status      MaintenanceOrderStatus `json:"status" db:"status"`
	Description string                 `json:"description" db:"description"`
	StartDate   time.Time              `json:"start_date" db:"start_date"`
	EndDate     *time.Time             `json:"end_date,omitempty" db:"end_date"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

type CreateMaintenanceOrderInput struct {
	DroneID     string    `json:"drone_id" validate:"required"`
	Description string    `json:"description" validate:"required,min=10,max=1000"`
	StartDate   time.Time `json:"start_date" validate:"required"`
}

type UpdateMaintenanceOrderInput struct {
	Status      *MaintenanceOrderStatus `json:"status,omitempty" validate:"omitempty,oneof=pending in_progress completed cancelled"`
	Description *string                 `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	EndDate     *time.Time              `json:"end_date,omitempty"`
}

type MaintenanceOrderWithDetails struct {
	MaintenanceOrder
	Drone *Drone `json:"drone,omitempty"`
	User  *User  `json:"user,omitempty"`
}
