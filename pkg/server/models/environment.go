package models

import "gorm.io/gorm"

// Environment is the model of the environment table.
type Environment struct {
	gorm.Model
	ID   int    `json:"id"`
	Name string `json:"name"`
	IP   string `json:"ip"`
}
