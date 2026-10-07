package domain

import "strings"

type Bicycle struct {
	name string
	id int
	easeOfUse int
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

func NewBycicle(
	name string,
	id int,
	easeOfUse int,
) (*Bicycle, error) {
	if easeOfUse < 1 || easeOfUse > 10 {
		return nil, ErrInvalidEaseOfUse
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
	}, nil
}