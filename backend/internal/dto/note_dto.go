package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateNoteCommand struct {
	ProjectID uuid.UUID `json:"projectID"`
	Title     string    `json:"title" binding:"required,min=2,max=255"`
	Content   string    `json:"content" binding:"required"`
}

type UpdateNoteCommand struct {
	ProjectID uuid.UUID `json:"projectID"`
	ID        uuid.UUID `json:"id"`
	Title     *string   `json:"title,omitempty" binding:"omitempty,min=2,max=255"`
	Content   *string   `json:"content,omitempty" binding:"omitempty"`
}

type NoteView struct {
	ID        uuid.UUID    `json:"id"`
	ProjectID uuid.UUID    `json:"projectId"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Archived  bool         `json:"archived"`
	Tags      []TagSummary `json:"tags"`
	CreatedAt *time.Time   `json:"createdAt,omitempty"`
	UpdatedAt *time.Time   `json:"updatedAt,omitempty"`
}

type NoteSummary struct {
	ID        uuid.UUID    `json:"id"`
	ProjectID uuid.UUID    `json:"projectId"`
	Title     string       `json:"title"`
	Archived  bool         `json:"archived"`
	Tags      []TagSummary `json:"tags"`
	CreatedAt *time.Time   `json:"createdAt,omitempty"`
	UpdatedAt *time.Time   `json:"updatedAt,omitempty"`
}
