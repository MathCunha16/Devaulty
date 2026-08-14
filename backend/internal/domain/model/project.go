package model

import "github.com/google/uuid"

type Project struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	Icon        *string   `json:"icon,omitempty" db:"icon"`
	Color       *string   `json:"color,omitempty" db:"color"`
	Archived    bool      `json:"archived" db:"archived"`
	BaseEntity
}
