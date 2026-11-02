package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDrawingCommand struct {
	ProjectID uuid.UUID `json:"projectID"`
	Title     string    `json:"title" binding:"required,min=1,max=255"`
}

type UpdateDrawingCommand struct {
	ProjectID uuid.UUID `json:"projectID"`
	ID        uuid.UUID `json:"id"`
	Title     *string   `json:"title,omitempty" binding:"omitempty,min=1,max=255"`
	Content   *string   `json:"content,omitempty" binding:"omitempty"`
	Thumbnail *string   `json:"thumbnail,omitempty" binding:"omitempty"`
}

type DrawingView struct {
	ID        uuid.UUID    `json:"id"`
	ProjectID uuid.UUID    `json:"projectId"`
	Title     string       `json:"title"`
	Content   *string      `json:"content"`
	Thumbnail *string      `json:"thumbnail,omitempty"`
	Tags      []TagSummary `json:"tags"`
	CreatedAt *time.Time   `json:"createdAt,omitempty"`
	UpdatedAt *time.Time   `json:"updatedAt,omitempty"`
}

type DrawingSummary struct {
	ID        uuid.UUID    `json:"id"`
	ProjectID uuid.UUID    `json:"projectId"`
	Title     string       `json:"title"`
	Thumbnail *string      `json:"thumbnail,omitempty"`
	Tags      []TagSummary `json:"tags"`
	CreatedAt *time.Time   `json:"createdAt,omitempty"`
	UpdatedAt *time.Time   `json:"updatedAt,omitempty"`
}
