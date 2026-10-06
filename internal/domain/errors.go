package domain

import "errors"

var (
	ErrInvalidEaseOfUse = errors.New(
		"Error: Ease of use should be between 1 and 10",
	)
	ErrInvalidEnergyConsumption = errors.New(
		"Error: Energy consumption cannot be negative",
	)
	ErrNegativeItemNumber = errors.New(
		"Error: Item number cannot be negative",
	)
	ErrItemNumberExists = errors.New(
		"Error: Item number cannot be duplicate",
	)
)