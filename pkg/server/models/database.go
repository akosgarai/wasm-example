package models

import "gorm.io/gorm"

// Dbtype is the model of the dbtype table.
type Dbtype struct {
	gorm.Model
	Id   int    `json:"id"`
	Name string `json:"name"`
}
