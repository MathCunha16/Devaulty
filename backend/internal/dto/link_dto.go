package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateLinkCommand struct {
	ProjectID   uuid.UUID `json:"projectID"`
	Title       string    `json:"title" binding:"required,min=2,max=255"`
	URL         string    `json:"url" binding:"required,url"`
	Description *string   `json:"description,omitempty" binding:"omitempty"`
}

type UpdateLinkCommand struct {
	ID          uuid.UUID `json:"id"`
	ProjectID   uuid.UUID `json:"projectID"`
	Title       *string   `json:"title,omitempty" binding:"omitempty,min=2,max=255"`
	URL         *string   `json:"url,omitempty" binding:"omitempty,url"`
	Description *string   `json:"description,omitempty" binding:"omitempty"`
}

type LinkView struct {
	ID          uuid.UUID    `json:"id"`
	ProjectID   uuid.UUID    `json:"projectId"`
	Title       string       `json:"title"`
	URL         string       `json:"url"`
	Description *string      `json:"description,omitempty"`
	Tags        []TagSummary `json:"tags"`
	CreatedAt   *time.Time   `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time   `json:"updatedAt,omitempty"`
}
