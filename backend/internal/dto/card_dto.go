package dto

import (
	"devaulty-backend/internal/domain/model"
	"time"

	"github.com/google/uuid"
)

type CreateCardItemCommand struct {
	ItemType model.ItemType `json:"itemType" binding:"required"`
	ItemID   uuid.UUID      `json:"itemId" binding:"required"`
}

type CreateCardCommand struct {
	ProjectID   uuid.UUID               `json:"projectID"`
	BoardID     uuid.UUID               `json:"boardID"`
	ColumnID    uuid.UUID               `json:"columnID"`
	Title       string                  `json:"title" binding:"required,min=2,max=255"`
	Description *string                 `json:"description,omitempty" binding:"omitempty,min=1"`
	Priority    *model.CardPriority     `json:"priority,omitempty"`
	DueDate     *time.Time              `json:"dueDate,omitempty"`
	LinkedItems []CreateCardItemCommand `json:"linkedItems,omitempty"`
}

type UpdateCardCommand struct {
	ID          uuid.UUID               `json:"id"`
	ProjectID   uuid.UUID               `json:"projectID"`
	BoardID     uuid.UUID               `json:"boardID"`
	Title       *string                 `json:"title,omitempty" binding:"omitempty,min=2,max=255"`
	Description *string                 `json:"description,omitempty" binding:"omitempty,min=1"`
	Priority    *model.CardPriority     `json:"priority,omitempty"`
	DueDate     *time.Time              `json:"dueDate,omitempty"`
	LinkedItems []CreateCardItemCommand `json:"linkedItems,omitempty"`
}

type MoveCardCommand struct {
	ID             uuid.UUID `json:"id"`
	ProjectID      uuid.UUID `json:"projectID"`
	BoardID        uuid.UUID `json:"boardID"`
	TargetColumnID uuid.UUID `json:"targetColumnId" binding:"required"`
	Position       *uint16   `json:"position" binding:"required"`
}

type CardSummaryView struct {
	ID        uuid.UUID           `json:"id"`
	ProjectID uuid.UUID           `json:"projectId"`
	BoardID   uuid.UUID           `json:"boardId"`
	ColumnID  uuid.UUID           `json:"columnId"`
	Title     string              `json:"title"`
	Position  uint16              `json:"position"`
	Priority  *model.CardPriority `json:"priority,omitempty"`
	DueDate   *time.Time          `json:"dueDate,omitempty"`
	Tags      []TagSummary        `json:"tags"`
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt *time.Time          `json:"updatedAt,omitempty"`
}

type CardView struct {
	ID          uuid.UUID           `json:"id"`
	ProjectID   uuid.UUID           `json:"projectId"`
	BoardID     uuid.UUID           `json:"boardId"`
	ColumnID    uuid.UUID           `json:"columnId"`
	Title       string              `json:"title"`
	Description *string             `json:"description,omitempty"`
	Position    uint16              `json:"position"`
	Priority    *model.CardPriority `json:"priority,omitempty"`
	DueDate     *time.Time          `json:"dueDate,omitempty"`
	LinkedItems []model.CardItem    `json:"linkedItems,omitempty"`
	Tags        []TagSummary        `json:"tags"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   *time.Time          `json:"updatedAt,omitempty"`
}
