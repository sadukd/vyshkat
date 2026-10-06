package domain

type InventoryItem interface {
	ItemName() string
	InventoryNumber() int
}