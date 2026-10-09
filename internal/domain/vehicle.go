package domain

type Vehicle interface {
	InventoryItem
	EaseOfUse() int
	Condition() int
}
