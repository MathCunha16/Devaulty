package model

import "github.com/google/uuid"

type Drawing struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProjectID uuid.UUID `json:"project_id" db:"project_id"`
	Title     string    `json:"title" db:"title"`
	Content   *string   `json:"content" db:"content"`
	Thumbnail *string   `json:"thumbnail" db:"thumbnail"`
	BaseEntity
}
