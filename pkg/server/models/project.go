package models

import "gorm.io/gorm"

// Project is the model of the project table.
type Project struct {
	gorm.Model
	ID   int    `json:"id"`
	Name string `json:"name"`
}
