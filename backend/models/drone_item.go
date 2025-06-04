package models

import (
	"time"
)

type DroneItemType string

const (
	DroneItemTypeBag           DroneItemType = "bag"
	DroneItemTypeRadioControl  DroneItemType = "radio_control"
	DroneItemTypeBatteries     DroneItemType = "batteries"
	DroneItemTypePropellers    DroneItemType = "propellers"
	DroneItemTypeChargingHub   DroneItemType = "charging_hub"
	DroneItemTypeStrap         DroneItemType = "strap"
	DroneItemTypeDrone         DroneItemType = "drone"
)

type DroneItem struct {
	ID        string        `json:"id" db:"id"`
	DroneID   string        `json:"drone_id" db:"drone_id"`
	ItemType  DroneItemType `json:"item_type" db:"item_type"`
	Quantity  int           `json:"quantity" db:"quantity"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}