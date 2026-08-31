package model

import (
	"time"

	"github.com/google/uuid"
)

type CardPriority string

const (
	CardPriorityLow           CardPriority = "LOW"
	CardPriorityMedium        CardPriority = "MEDIUM"
	CardPriorityHigh          CardPriority = "HIGH"
	CardPriorityExtremelyHigh CardPriority = "EXTREMELY_HIGH"
)

var CardPriorities = []CardPriority{CardPriorityLow, CardPriorityMedium, CardPriorityHigh, CardPriorityExtremelyHigh}

type CardItem struct {
	CardID   uuid.UUID `json:"cardId" db:"card_id"`
	ItemID   uuid.UUID `json:"itemId" db:"item_id"`
	ItemType ItemType  `json:"itemType" db:"item_type"`
	BaseEntity
}

type Card struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	ColumnID    uuid.UUID     `json:"columnId" db:"column_id"`
	Title       string        `json:"title" db:"title"`
	Description *string       `json:"description,omitempty" db:"description"`
	Position    uint16        `json:"position" db:"position"`
	Priority    *CardPriority `json:"priority,omitempty" db:"priority"`
	DueDate     *time.Time    `json:"dueDate,omitempty" db:"due_date"`
	LinkedItems []CardItem    `json:"linkedItems,omitempty" db:"-"`
	BaseEntity
}
