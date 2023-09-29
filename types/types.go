package types

import (
	"gorm.io/gorm"
)

type Widget struct {
	gorm.Model
	Name string
}

type Order struct {
	gorm.Model
	Widget   int
	Quantity int
}
