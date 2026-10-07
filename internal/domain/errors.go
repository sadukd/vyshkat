package domain

import "errors"

var (
	ErrInvalidEaseOfUse = errors.New(
		"ease of use must be between 1 and 10",
	)
	ErrInvalidEnergyConsumption = errors.New(
		"energy consumption cannot be negative",
	)
	ErrNegativeItemNumber = errors.New(
		"item number cannot be negative",
	)
	ErrEmptyName = errors.New(
		"item name cannot be empty",
	)
)
