package models

import "gorm.io/gorm"

// Client is the model of the client table.
type Client struct {
	gorm.Model
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}
