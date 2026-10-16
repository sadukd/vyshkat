package domain

import "math"

type EnergyConsumer interface {
	EnergyConsumption() float64
}

func validEnergyConsumption(energy float64) bool {
	return energy >= 0 && !math.IsNaN(energy) && !math.IsInf(energy, 0)
}
