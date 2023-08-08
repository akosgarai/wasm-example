package models

import "gorm.io/gorm"

// Runtime is the model of the runtime table.
type Runtime struct {
	gorm.Model
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique"`
}
