package domain

type Vehicle interface {
	InventoryItem
	VehicleNumber() int
	EaseOfUse() int8
}