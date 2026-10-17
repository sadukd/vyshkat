package service

import "vyshkat/internal/domain"

type VehicleInspector interface {
	Inspect(vehicle domain.Vehicle) bool
}

type Service struct {
	inspector VehicleInspector
	inventory []domain.InventoryItem
}

func NewService(inspector VehicleInspector) *Service {
	return &Service{
		inventory: make([]domain.InventoryItem, 0),
		inspector: inspector,
	}
}

func (s *Service) hasNumber(number int) bool {
	_, _, found := s.findInventoryItem(number)
	return found
}

func (s *Service) findInventoryItem(number int) (int, domain.InventoryItem, bool) {
	for index, item := range s.inventory {
		if item.Number() == number {
			return index, item, true
		}
	}

	return -1, nil, false
}

func (s *Service) AddVehicle(vehicle domain.Vehicle) error {
	if s.hasNumber(vehicle.Number()) {
		return ErrDuplicateNumber
	}

	if !s.inspector.Inspect(vehicle) {
		return ErrInspectionFailed
	}

	s.inventory = append(s.inventory, vehicle)
	return nil
}

func (s *Service) AddItem(item domain.InventoryItem) error {
	if _, ok := item.(domain.Vehicle); ok {
		return ErrVehicleAsItem
	}

	if s.hasNumber(item.Number()) {
		return ErrDuplicateNumber
	}

	s.inventory = append(s.inventory, item)
	return nil
}

func (s *Service) Inventory() []domain.InventoryItem {
	return append([]domain.InventoryItem(nil), s.inventory...)
}

func (s *Service) FindInventoryItem(number int) (domain.InventoryItem, error) {
	_, item, found := s.findInventoryItem(number)
	if !found {
		return nil, ErrInventoryItemNotFound
	}
	return item, nil
}

func (s *Service) RemoveVehicle(number int) error {
	index, item, found := s.findInventoryItem(number)
	if !found {
		return ErrInventoryItemNotFound
	}
	if _, ok := item.(domain.Vehicle); !ok {
		return ErrNotVehicle
	}
	s.removeAt(index)
	return nil
}

func (s *Service) RemoveItem(number int) error {
	index, item, found := s.findInventoryItem(number)
	if !found {
		return ErrInventoryItemNotFound
	}
	if _, ok := item.(domain.Vehicle); ok {
		return ErrVehicleRemovalRequired
	}
	s.removeAt(index)
	return nil
}

func (s *Service) removeAt(index int) {
	copy(s.inventory[index:], s.inventory[index+1:])
	s.inventory[len(s.inventory)-1] = nil
	s.inventory = s.inventory[:len(s.inventory)-1]
}

func (s *Service) VehicleCount() int {
	count := 0

	for _, item := range s.inventory {
		if _, ok := item.(domain.Vehicle); ok {
			count++
		}
	}

	return count
}

func (s *Service) InventoryCount() int {
	return len(s.inventory)
}

func (s *Service) TotalEnergyConsumption() float64 {
	total := 0.0

	for _, item := range s.inventory {
		if consumer, ok := item.(domain.EnergyConsumer); ok {
			total += consumer.EnergyConsumption()
		}
	}

	return total
}

func (s *Service) BeginnerVehicles() []domain.Vehicle {
	var vehicles []domain.Vehicle

	for _, item := range s.inventory {
		vehicle, ok := item.(domain.Vehicle)

		if ok && vehicle.EaseOfUse() >= 6 {
			vehicles = append(vehicles, vehicle)
		}
	}

	return vehicles
}
