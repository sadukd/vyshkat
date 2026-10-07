package domain

import "strings"

type Helmet struct {
	name string
	id   int
}

func (h Helmet) Name() string {
	return h.name
}

func (h Helmet) Number() int {
	return h.id
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
		name: name,
		id:   id,
	}, nil
}
