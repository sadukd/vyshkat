package domain

import "strings"

type Helmet struct {
	name      string
	id        int
}

func (b Helmet) Name() string {
	return b.name
}

func (b Helmet) Number() int {
	return b.id
}

func NewHelmet(
	name string,
	id int,
) (*Helmet, error) {
	if id < 0 {
		return nil, ErrNegativeItemNumber
	}

	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyName
	}

	return &Helmet{
		name:      name,
		id:        id,
	}, nil
}