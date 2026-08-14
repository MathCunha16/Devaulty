package model

import "github.com/google/uuid"

type Link struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ProjectID   uuid.UUID `json:"projectId" db:"project_id"`
	Title       string    `json:"title" db:"title"`
	Url         string    `json:"url" db:"url"`
	Description *string   `json:"description,omitempty" db:"description"`
	BaseEntity
}
