package model

import "github.com/google/uuid"

type Board struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ProjectID   uuid.UUID `json:"projectId" db:"project_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	IsDefault   bool      `json:"isDefault" db:"is_default"`
	BaseEntity
}
