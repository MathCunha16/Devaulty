package model

import "github.com/google/uuid"

type BoardColumn struct {
	ID       uuid.UUID `json:"id" db:"id"`
	BoardID  uuid.UUID `json:"boardId" db:"board_id"`
	Name     string    `json:"name" db:"name"`
	Position uint8     `json:"position" db:"position"`
	WipLimit *uint16   `json:"wipLimit,omitempty" db:"wip_limit"`
	BaseEntity
}
