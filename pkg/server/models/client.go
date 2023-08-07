package models

import "gorm.io/gorm"

// Client is the model of the client table.
type Client struct {
	gorm.Model
	ID   int    `json:"id"`
	Name string `json:"name"`
}
