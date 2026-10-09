package domain

import "strings"

type Bicycle struct {
	name      string
	id        int
	easeOfUse int
	condition int
}

func (b Bicycle) Name() string {
	return b.name
}

func (b Bicycle) Number() int {
	return b.id
}

func (b Bicycle) EaseOfUse() int {
	return b.easeOfUse
}

func (b Bicycle) Condition() int {
	return b.condition
}

func NewBicycle(
	name string,
	id int,
	easeOfUse int,
	condition int,
) (*Bicycle, error) {
	if easeOfUse < 1 || easeOfUse > 10 {
		return nil, ErrInvalidEaseOfUse
	}

	if condition < 0 || condition > 10 {
		return nil, ErrInvalidCondition
	}

	if id < 0 {
		return nil, ErrNegativeItemNumber
	}

	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyName
	}

	return &Bicycle{
		name:      name,
		id:        id,
		easeOfUse: easeOfUse,
		condition: condition,
	}, nil
}
