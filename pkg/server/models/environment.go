package models

import "gorm.io/gorm"

// Environment is the model of the environment table.
type Environment struct {
	gorm.Model
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique"`
	Host []Host `gorm:"foreignKey:EnvironmentID"`
}
