package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateBoardCommand struct {
	ProjectID   uuid.UUID `json:"projectID"`
	Name        string    `json:"name" binding:"required,min=2,max=255"`
	Description *string   `json:"description,omitempty" binding:"omitempty,min=1"`
	IsDefault   bool      `json:"isDefault"`
}

type UpdateBoardCommand struct {
	ID          uuid.UUID `json:"id"`
	ProjectID   uuid.UUID `json:"projectID"`
	Name        *string   `json:"name,omitempty" binding:"omitempty,min=2,max=255"`
	Description *string   `json:"description,omitempty" binding:"omitempty,min=1,max=255"`
	IsDefault   *bool     `json:"isDefault,omitempty" binding:"omitempty"`
}

type BoardView struct {
	ID          uuid.UUID    `json:"id"`
	ProjectID   uuid.UUID    `json:"projectId"`
	Name        string       `json:"name"`
	Description *string      `json:"description"`
	IsDefault   bool         `json:"isDefault"`
	Tags        []TagSummary `json:"tags"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   *time.Time   `json:"updatedAt,omitempty"`
}
