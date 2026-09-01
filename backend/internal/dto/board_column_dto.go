package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateBoardColumnCommand struct {
	ProjectID uuid.UUID `json:"projectID"`
	BoardID   uuid.UUID `json:"boardID"`
	Name      string    `json:"name" binding:"required,min=2,max=255"`
	WipLimit  *uint16   `json:"wipLimit,omitempty" binding:"omitempty"`
}

type UpdateBoardColumnCommand struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"projectID"`
	BoardID   uuid.UUID `json:"boardID"`
	Name      *string   `json:"name,omitempty" binding:"omitempty,min=2,max=255"`
	WipLimit  *uint16   `json:"wipLimit,omitempty" binding:"omitempty"`
}

type ReorderBoardColumnsCommand struct {
	ProjectID uuid.UUID   `json:"projectID"`
	BoardID   uuid.UUID   `json:"boardID"`
	Positions []uuid.UUID `json:"positions" binding:"required,min=1"`
}

type BoardColumnView struct {
	ID        uuid.UUID  `json:"id"`
	ProjectID uuid.UUID  `json:"projectId"`
	BoardID   uuid.UUID  `json:"boardId"`
	Name      string     `json:"name"`
	Position  uint8      `json:"position"`
	WipLimit  *uint16    `json:"wipLimit"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
