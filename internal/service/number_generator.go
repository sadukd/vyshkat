package service

import (
	"errors"
	"math"
)

var ErrNumbersExhausted = errors.New("inventory numbers exhausted")

type NumberGenerator struct {
	next int
}

func NewNumberGenerator() *NumberGenerator {
	return &NumberGenerator{next: 1}
}

func (g *NumberGenerator) Next() (int, error) {
	if g.next <= 0 {
		return 0, ErrNumbersExhausted
	}

	number := g.next

	if g.next == math.MaxInt {
		g.next = 0
	} else {
		g.next++
	}

	return number, nil
}
