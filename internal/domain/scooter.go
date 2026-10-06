package domain

type Scooter struct {
	name string
	id int
	energy float64
}

func (s Scooter) VehicleName() string {
	return s.name
}

func (s Scooter) VehicleNumber() int {
	return s.id
}

func (s Scooter) EnergyConsumption() float64 {
	return s.energy
}