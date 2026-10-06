package domain

type Scooter struct {
	name string
	id int
	energy float64
	easeOfUse int
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