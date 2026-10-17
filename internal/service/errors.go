package service

import "errors"

var (
	ErrDuplicateNumber = errors.New(
		"duplicate inventory number",
	)
	ErrInspectionFailed = errors.New(
		"vehicle failed inspection",
	)
	ErrVehicleAsItem = errors.New(
		"vehicles must be added through AddVehicle",
	)
	ErrInventoryItemNotFound = errors.New(
		"inventory item not found",
	)
	ErrNotVehicle = errors.New(
		"inventory item is not a vehicle",
	)
	ErrVehicleRemovalRequired = errors.New(
		"vehicles must be removed through RemoveVehicle",
	)
)
