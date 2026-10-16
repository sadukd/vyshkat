package domain

import "errors"

var (
	ErrInvalidEaseOfUse = errors.New(
		"ease of use must be between 1 and 10",
	)
	ErrInvalidCondition = errors.New(
		"condition must be between 0 and 10",
	)
	ErrInvalidEnergyConsumption = errors.New(
		"energy consumption must be a non-negative finite number",
	)
	ErrNegativeItemNumber = errors.New(
		"item number cannot be negative",
	)
	ErrEmptyName = errors.New(
		"item name cannot be empty",
	)
)
