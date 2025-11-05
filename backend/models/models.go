package models

import (
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	Email    string    `gorm:"uniqueIndex;not null" json:"email"`
	Password string    `gorm:"not null" json:"-"`
	Projects []Project `gorm:"foreignKey:OwnerID" json:"projects,omitempty"`
}

// Project represents a project that contains secrets
type Project struct {
	gorm.Model
	Name    string   `gorm:"not null" json:"name"`
	OwnerID uint     `gorm:"not null" json:"owner_id"`
	Owner   User     `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Secrets []Secret `gorm:"foreignKey:ProjectID" json:"secrets,omitempty"`
}

// Secret represents an encrypted secret key-value pair
type Secret struct {
	gorm.Model
	Key       string  `gorm:"not null" json:"key"`
	Value     string  `gorm:"not null" json:"value"` // This will be encrypted
	ProjectID uint    `gorm:"not null" json:"project_id"`
	Project   Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}
