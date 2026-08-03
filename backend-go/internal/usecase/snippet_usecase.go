package usecase

import (
	"context"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSnippetNotFound = errors.New("snippet not found")
)

type SnippetUseCase struct {
	snippetRepo port.SnippetRepository
	projectRepo port.ProjectRepository
}

type CreateSnippetCommand struct {
	ProjectID   uuid.UUID
	Title       string                `json:"title" binding:"required,min=2,max=255"`
	Description *string               `json:"description,omitempty" binding:"omitempty,min=1,max=255"`
	Content     string                `json:"content" binding:"required,min=1"`
	Language    model.SnippetLanguage `json:"language" binding:"required"`
	SnippetType model.SnippetType     `json:"snippetType" binding:"required"`
}

type UpdateSnippetCommand struct {
	ProjectID   uuid.UUID
	ID          uuid.UUID
	Title       *string                `json:"title,omitempty" binding:"omitempty,min=2,max=255"`
	Description *string                `json:"description,omitempty" binding:"omitempty,min=1,max=255"`
	Content     *string                `json:"content,omitempty" binding:"omitempty,min=1"`
	Language    *model.SnippetLanguage `json:"language,omitempty" binding:"omitempty"`
	SnippetType *model.SnippetType     `json:"snippetType" binding:"omitempty"`
}

func NewSnippetUseCase(snippetRepo port.SnippetRepository, projectRepo port.ProjectRepository) *SnippetUseCase {
	return &SnippetUseCase{
		snippetRepo: snippetRepo,
		projectRepo: projectRepo,
	}
}

func (uc *SnippetUseCase) Create(ctx context.Context, cmd CreateSnippetCommand) (*model.Snippet, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}

	snippet := model.Snippet{
		ID:          uuid.New(),
		ProjectID:   cmd.ProjectID,
		Title:       cmd.Title,
		Description: cmd.Description,
		Content:     cmd.Content,
		Language:    &cmd.Language,
		SnippetType: cmd.SnippetType,
		BaseEntity: model.BaseEntity{
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		},
	}
	return uc.snippetRepo.Save(ctx, &snippet)
}

func (uc *SnippetUseCase) GetByID(ctx context.Context, projectID, id uuid.UUID) (*model.Snippet, error) {
	snippet, err := uc.snippetRepo.FindByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error trying to find snippet: %w", err)
	}
	if snippet == nil {
		return nil, ErrSnippetNotFound
	}
	return snippet, nil
}

func (uc *SnippetUseCase) GetAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[model.Snippet], error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return model.Page[model.Snippet]{}, err
	}
	return uc.snippetRepo.FindAllByProjectID(ctx, projectID, page, size)
}

func (uc *SnippetUseCase) Update(ctx context.Context, cmd UpdateSnippetCommand) (*model.Snippet, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}

	snippet, err := uc.snippetRepo.FindByIDAndProjectID(ctx, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error trying to find snippet: %w", err)
	}
	if snippet == nil {
		return nil, ErrSnippetNotFound
	}
	if cmd.Title != nil {
		snippet.Title = *cmd.Title
	}
	if cmd.Description != nil {
		snippet.Description = cmd.Description
	}
	if cmd.Content != nil {
		snippet.Content = *cmd.Content
	}
	if cmd.Language != nil {
		snippet.Language = cmd.Language
	}
	if cmd.SnippetType != nil {
		snippet.SnippetType = *cmd.SnippetType
	}
	now := time.Now()
	snippet.UpdatedAt = &now
	return uc.snippetRepo.Save(ctx, snippet)
}

func (uc *SnippetUseCase) Delete(ctx context.Context, projectID, id uuid.UUID) error {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return err
	}
	return uc.snippetRepo.DeleteByIDAndProjectID(ctx, projectID, id)
}
