package model

import "github.com/google/uuid"

type ItemType string

const (
	ItemTypeSnippet    ItemType = "SNIPPET"
	ItemTypeNote       ItemType = "NOTE"
	ItemTypeLink       ItemType = "LINK"
	ItemTypeProblem    ItemType = "PROBLEM"
	ItemTypeCredential ItemType = "CREDENTIAL"
	ItemTypeBoard      ItemType = "BOARD"
	ItemTypeCard       ItemType = "CARD"
)

type Tag struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProjectID uuid.UUID `json:"projectId" db:"project_id"`
	Name      string    `json:"name" db:"name"`
	Color     *string   `json:"color,omitempty" db:"color"`
	BaseEntity
}

type ItemTag struct {
	TagID    uuid.UUID `json:"tagId" db:"tag_id"`
	ItemType ItemType  `json:"itemType" db:"item_type"`
	ItemID   uuid.UUID `json:"itemId" db:"item_id"`
}
