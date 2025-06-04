package models

import (
	"time"
)

type Drone struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Model     string    `db:"model" json:"model"`
	Serial    string    `db:"serial" json:"serial"`
	Status    string    `db:"status" json:"status"`
	LastCheck time.Time `db:"last_check" json:"last_check"`
	NextCheck time.Time `db:"next_check" json:"next_check"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateDroneInput struct {
	Model  string `json:"model" validate:"required"`
	Serial string `json:"serial" validate:"required"`
}

type UpdateDroneInput struct {
	Model  string `json:"model" validate:"required"`
	Serial string `json:"serial" validate:"required"`
	Status string `json:"status" validate:"required,oneof=active maintenance inactive"`
}

func NewDrone(userID string, input CreateDroneInput) *Drone {
	now := time.Now()
	return &Drone{
		UserID:    userID,
		Model:     input.Model,
		Serial:    input.Serial,
		Status:    "active",
		LastCheck: now,
		NextCheck: now.AddDate(0, 1, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
}
