package domain

import "strings"

type Scooter struct {
	name      string
	id        int
	energy    float64
	easeOfUse int
	condition int
}

func (s Scooter) Name() string {
	return s.name
}

func (s Scooter) Number() int {
	return s.id
}

func (s Scooter) EnergyConsumption() float64 {
	return s.energy
}

func (s Scooter) EaseOfUse() int {
	return s.easeOfUse
}

func (s Scooter) Condition() int {
	return s.condition
}

func NewScooter(
	name string,
	id int,
	energy float64,
	easeOfUse int,
	condition int,
) (*Scooter, error) {
	if easeOfUse < 1 || easeOfUse > 10 {
		return nil, ErrInvalidEaseOfUse
	}

	if condition < 0 || condition > 10 {
		return nil, ErrInvalidCondition
	}

	if id < 0 {
		return nil, ErrNegativeItemNumber
	}

	if energy < 0 {
		return nil, ErrInvalidEnergyConsumption
	}

	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyName
	}

	return &Scooter{
		name:      name,
		id:        id,
		energy:    energy,
		easeOfUse: easeOfUse,
		condition: condition,
	}, nil
}
