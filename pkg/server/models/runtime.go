package models

import "gorm.io/gorm"

// Runtime is the model of the runtime table.
type Runtime struct {
	gorm.Model
	ID   int    `json:"id"`
	Name string `json:"name"`
}
