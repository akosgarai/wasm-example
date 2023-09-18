package models

import "gorm.io/gorm"

// Host is the model of the host table.
type Host struct {
	gorm.Model
	ID            uint         `json:"id" gorm:"primaryKey"`
	Name          string       `json:"name"`                                                     // host name
	IP            string       `json:"ip" gorm:"uniqueIndex:host_environment_index;size:256"`    // host ip or domain
	EnvironmentID uint         `json:"environment_id" gorm:"uniqueIndex:host_environment_index"` // environment id
	Environment   *Environment `gorm:"foreignKey:EnvironmentID"`
	SSHUser       string       `json:"ssh_user"` // ssh user
	SSHPort       int          `json:"ssh_port"` // ssh port
	SSHKey        string       `json:"ssh_key"`  // ssh key
}
