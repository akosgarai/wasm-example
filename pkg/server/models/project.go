package models

import "gorm.io/gorm"

// Project is the model of the project table.
type Project struct {
	gorm.Model
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique"`
}
