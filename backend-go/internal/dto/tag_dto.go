package dto

import "github.com/google/uuid"

type TagSummary struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Color *string   `json:"color,omitempty"`
}

type CreateTagCommand struct {
	ProjectID uuid.UUID `json:"-"`
	Name      string    `json:"name" binding:"required,min=2,max=40"`
	Color     string    `json:"color" binding:"hexcolor"`
}

type UpdateTagCommand struct {
	ID        uuid.UUID `json:"-"`
	ProjectID uuid.UUID `json:"-"`
	Name      *string   `json:"name,omitempty" binding:"omitempty,min=1,max=40"`
	Color     *string   `json:"color,omitempty" binding:"omitempty,hexcolor"`
}
