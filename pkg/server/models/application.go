package models

import "gorm.io/gorm"

// Application is the model of the application table.
type Application struct {
	gorm.Model
	ID            uint         `json:"id" gorm:"primaryKey"`
	ProjectID     uint         `json:"project_id"`
	Project       *Project     `json:"project"`
	ClientID      uint         `json:"client_id"`
	Client        *Client      `json:"client"`
	OwnerEmail    string       `json:"owner_email"`
	RuntimeID     int          `json:"runtime_id"`
	Runtime       *Runtime     `json:"runtime"`
	DatabaseID    int          `json:"database_id"`
	Database      *Dbtype      `json:"database"`
	EnvironmentID int          `json:"environment_id"`
	Environment   *Environment `json:"environment"`
}
