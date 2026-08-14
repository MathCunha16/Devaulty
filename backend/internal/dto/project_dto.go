package dto

import "github.com/google/uuid"

type CreateProjectCommand struct {
	Name        string  `json:"name" binding:"required,min=2,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty,min=1,max=255"`
	Icon        *string `json:"icon,omitempty" binding:"omitempty,max=100"`
	Color       *string `json:"color,omitempty" binding:"omitempty,hexcolor"`
}

type UpdateProjectCommand struct {
	ID          uuid.UUID `json:"-"`
	Name        *string   `json:"name,omitempty" binding:"omitempty,min=2,max=255"`
	Description *string   `json:"description,omitempty" binding:"omitempty,min=1,max=255"`
	Icon        *string   `json:"icon,omitempty" binding:"omitempty,max=100"`
	Color       *string   `json:"color,omitempty" binding:"omitempty,hexcolor"`
}
