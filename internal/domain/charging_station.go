package domain

import "strings"

type ChargingStation struct {
	name   string
	id     int
	energy float64
}

func (c ChargingStation) Name() string {
	return c.name
}

func (c ChargingStation) Number() int {
	return c.id
}

func (c ChargingStation) EnergyConsumption() float64 {
	return c.energy
}

func NewChargingStation(
	name string,
	id int,
	energy float64,
) (*ChargingStation, error) {
	if id < 0 {
		return nil, ErrNegativeItemNumber
	}

	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyName
	}

	if energy < 0 {
        return nil, ErrInvalidEnergyConsumption
    }

	return &ChargingStation{
		name: name,
		id:   id,
	}, nil
}
