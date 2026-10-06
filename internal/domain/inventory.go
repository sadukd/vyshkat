package domain

type InventoryItem interface {
	Name() string
	Number() int
}