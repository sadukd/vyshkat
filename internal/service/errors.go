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
)
