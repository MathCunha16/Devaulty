package usecase

import (
	"context"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"
	"devaulty-backend/internal/dto"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSnippetNotFound = errors.New("snippet not found")
)

type SnippetUseCase struct {
	snippetRepo port.SnippetRepository
	projectRepo port.ProjectRepository
	itemTagRepo port.ItemTagRepository
}

func NewSnippetUseCase(snippetRepo port.SnippetRepository, projectRepo port.ProjectRepository, itemTagRepo port.ItemTagRepository) *SnippetUseCase {
	return &SnippetUseCase{
		snippetRepo: snippetRepo,
		projectRepo: projectRepo,
		itemTagRepo: itemTagRepo,
	}
}

func (uc *SnippetUseCase) Create(ctx context.Context, cmd dto.CreateSnippetCommand) (*dto.SnippetView, error) {
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
	saved, err := uc.snippetRepo.Save(ctx, &snippet)
	if err != nil {
		return nil, err
	}
	return mapSnippetToView(saved, nil), nil
}

func (uc *SnippetUseCase) GetByID(ctx context.Context, projectID, id uuid.UUID) (*dto.SnippetView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return nil, err
	}
	snippet, err := uc.snippetRepo.FindByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error trying to find snippet: %w", err)
	}
	if snippet == nil {
		return nil, ErrSnippetNotFound
	}
	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeSnippet, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for snippet: %w", err)
	}
	return mapSnippetToView(snippet, tags), nil
}

func (uc *SnippetUseCase) GetAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[dto.SnippetView], error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return model.Page[dto.SnippetView]{}, err
	}
	snippetPage, err := uc.snippetRepo.FindAllByProjectID(ctx, projectID, page, size)
	if err != nil {
		return model.Page[dto.SnippetView]{}, err
	}
	if len(snippetPage.Content) == 0 {
		return model.NewPage([]dto.SnippetView{}, snippetPage.Number, snippetPage.Size, snippetPage.TotalElements), nil
	}

	snippetIDs := make([]uuid.UUID, len(snippetPage.Content))
	for i, s := range snippetPage.Content {
		snippetIDs[i] = s.ID
	}

	tagsMap, err := uc.itemTagRepo.FindTagsForItems(ctx, model.ItemTypeSnippet, projectID, snippetIDs)
	if err != nil {
		return model.Page[dto.SnippetView]{}, fmt.Errorf("error fetching tags for snippets: %w", err)
	}

	views := make([]dto.SnippetView, len(snippetPage.Content))
	for i, s := range snippetPage.Content {
		tags := tagsMap[s.ID]
		views[i] = *mapSnippetToView(&s, tags)
	}

	return model.NewPage(views, snippetPage.Number, snippetPage.Size, snippetPage.TotalElements), nil
}

func (uc *SnippetUseCase) Update(ctx context.Context, cmd dto.UpdateSnippetCommand) (*dto.SnippetView, error) {
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
	saved, err := uc.snippetRepo.Save(ctx, snippet)
	if err != nil {
		return nil, err
	}
	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeSnippet, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for snippet: %w", err)
	}
	return mapSnippetToView(saved, tags), nil
}

func (uc *SnippetUseCase) Delete(ctx context.Context, projectID, id uuid.UUID) error {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return err
	}
	deleted, err := uc.snippetRepo.DeleteByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return fmt.Errorf("error deleting snippet: %w", err)
	}
	if !deleted {
		return ErrSnippetNotFound
	}

	if err := uc.itemTagRepo.RemoveAllTagsFromItem(ctx, model.ItemTypeSnippet, id); err != nil {
		log.Printf("warning: failed to remove tags from snippet %s: %v", id, err)
	}
	return nil
}

func mapSnippetToView(snippet *model.Snippet, tags []model.Tag) *dto.SnippetView {
	tagSummaries := make([]dto.TagSummary, len(tags))
	for i, t := range tags {
		tagSummaries[i] = dto.TagSummary{
			ID:    t.ID,
			Name:  t.Name,
			Color: t.Color,
		}
	}
	return &dto.SnippetView{
		ID:          snippet.ID,
		ProjectID:   snippet.ProjectID,
		Title:       snippet.Title,
		Description: snippet.Description,
		Content:     snippet.Content,
		Language:    snippet.Language,
		SnippetType: snippet.SnippetType,
		Tags:        tagSummaries,
		CreatedAt:   &snippet.CreatedAt,
		UpdatedAt:   snippet.UpdatedAt,
	}
}
