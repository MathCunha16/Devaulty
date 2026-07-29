package model

import "github.com/google/uuid"

type Note struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProjectID uuid.UUID `json:"projectId" db:"project_id"`
	Title     string    `json:"title" db:"title"`
	Content   *string   `json:"content,omitempty" db:"content"`
	Archived  bool      `json:"archived" db:"archived"`
	BaseEntity
}
