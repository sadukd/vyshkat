package inspection

import "vyshkat/internal/domain"

const MinAcceptedCondition = 7

type Inspector struct{}

func NewInspector() *Inspector {
    return &Inspector{}
}

func (i *Inspector) Inspect(vehicle domain.Vehicle) bool {
    return vehicle.Condition() >= MinAcceptedCondition
}