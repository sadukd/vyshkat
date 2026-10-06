package domain

type Scooter struct {
	name string
	id int
	energy float64
	dif_level int8
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

func (s Scooter) EaseOfUse() int8 {
	return s.dif_level
}