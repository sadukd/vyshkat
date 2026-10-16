package tests

import "vyshkat/internal/domain"

var (
	_ domain.Vehicle        = (*domain.Bicycle)(nil)
	_ domain.Vehicle        = (*domain.Scooter)(nil)
	_ domain.EnergyConsumer = (*domain.Scooter)(nil)
	_ domain.InventoryItem  = (*domain.Helmet)(nil)
	_ domain.InventoryItem  = (*domain.ChargingStation)(nil)
	_ domain.EnergyConsumer = (*domain.ChargingStation)(nil)
)
