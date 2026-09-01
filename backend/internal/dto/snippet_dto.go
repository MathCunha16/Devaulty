package dto

import (
	"time"

	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type CreateSnippetCommand struct {
	ProjectID   uuid.UUID             `json:"projectID"`
	Title       string                `json:"title" binding:"required,min=2,max=255"`
	Description *string               `json:"description,omitempty" binding:"omitempty,min=1"`
	Content     string                `json:"content" binding:"required,min=1"`
	Language    model.SnippetLanguage `json:"language" binding:"required"`
	SnippetType model.SnippetType     `json:"snippetType" binding:"required"`
}

type UpdateSnippetCommand struct {
	ProjectID   uuid.UUID              `json:"projectID"`
	ID          uuid.UUID              `json:"id"`
	Title       *string                `json:"title,omitempty" binding:"omitempty,min=2,max=255"`
	Description *string                `json:"description,omitempty" binding:"omitempty,min=1"`
	Content     *string                `json:"content,omitempty" binding:"omitempty,min=1"`
	Language    *model.SnippetLanguage `json:"language,omitempty" binding:"omitempty"`
	SnippetType *model.SnippetType     `json:"snippetType,omitempty" binding:"omitempty"`
}

type SnippetView struct {
	ID          uuid.UUID              `json:"id"`
	ProjectID   uuid.UUID              `json:"projectId"`
	Title       string                 `json:"title"`
	Description *string                `json:"description,omitempty"`
	Content     string                 `json:"content"`
	Language    *model.SnippetLanguage `json:"language,omitempty"`
	SnippetType model.SnippetType      `json:"snippetType"`
	Tags        []TagSummary           `json:"tags"`
	CreatedAt   *time.Time             `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time             `json:"updatedAt,omitempty"`
}
