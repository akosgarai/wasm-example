package models

import "gorm.io/gorm"

// Dbtype is the model of the dbtype table.
type Dbtype struct {
	gorm.Model
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}
