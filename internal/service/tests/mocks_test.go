package tests

import "vyshkat/internal/domain"

type inspectorMock struct {
	result bool
	calls  []domain.Vehicle
}

func (m *inspectorMock) Inspect(vehicle domain.Vehicle) bool {
	m.calls = append(m.calls, vehicle)
	return m.result
}

type itemMock struct {
	name   string
	number int
}

func (m *itemMock) Name() string { return m.name }

func (m *itemMock) Number() int { return m.number }

type vehicleMock struct {
	itemMock
	easeOfUse int
	condition int
	energy    float64
}

func (m *vehicleMock) EaseOfUse() int { return m.easeOfUse }

func (m *vehicleMock) Condition() int { return m.condition }

func (m *vehicleMock) EnergyConsumption() float64 { return m.energy }

type energyItemMock struct {
	itemMock
	energy float64
}

func (m *energyItemMock) EnergyConsumption() float64 { return m.energy }
