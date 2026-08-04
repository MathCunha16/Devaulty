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
	ErrLinkNotFound = errors.New("link not found")
)

type LinkUseCase struct {
	linkRepo    port.LinkRepository
	projectRepo port.ProjectRepository
}

type CreateLinkCommand struct {
	ProjectID   uuid.UUID
	Title       string  `json:"title" binding:"required,min=2,max=255"`
	URL         string  `json:"url" binding:"required,url"`
	Description *string `json:"description,omitempty" binding:"omitempty"`
}

type UpdateLinkCommand struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Title       *string `json:"title,omitempty" binding:"omitempty,min=2,max=255"`
	URL         *string `json:"url,omitempty" binding:"omitempty,url"`
	Description *string `json:"description,omitempty" binding:"omitempty"`
}

func NewLinkUseCase(linkRepo port.LinkRepository, projectRepo port.ProjectRepository) *LinkUseCase {
	return &LinkUseCase{
		linkRepo:    linkRepo,
		projectRepo: projectRepo,
	}
}

func (uc *LinkUseCase) Create(ctx context.Context, cmd CreateLinkCommand) (*model.Link, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}

	link := model.Link{
		ID:          uuid.New(),
		ProjectID:   cmd.ProjectID,
		Title:       cmd.Title,
		Url:         cmd.URL,
		Description: cmd.Description,
		BaseEntity: model.BaseEntity{
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		},
	}

	return uc.linkRepo.Save(ctx, &link)
}

func (uc *LinkUseCase) GetByID(ctx context.Context, projectID, id uuid.UUID) (*model.Link, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return nil, err
	}
	link, err := uc.linkRepo.FindByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error trying to find link: %w", err)
	}

	if link == nil {
		return nil, ErrLinkNotFound
	}

	return link, nil
}

func (uc *LinkUseCase) GetAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[model.Link], error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return model.Page[model.Link]{}, err
	}
	return uc.linkRepo.FindAllByProjectID(ctx, projectID, page, size)
}

func (uc *LinkUseCase) Update(ctx context.Context, cmd UpdateLinkCommand) (*model.Link, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}

	link, err := uc.linkRepo.FindByIDAndProjectID(ctx, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error trying to find link: %w", err)
	}
	if link == nil {
		return nil, ErrLinkNotFound
	}

	if cmd.Title != nil {
		link.Title = *cmd.Title
	}
	if cmd.URL != nil {
		link.Url = *cmd.URL
	}
	if cmd.Description != nil {
		link.Description = cmd.Description
	}
	now := time.Now()
	link.UpdatedAt = &now
	return uc.linkRepo.Save(ctx, link)
}

func (uc *LinkUseCase) Delete(ctx context.Context, projectID, id uuid.UUID) error {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return err
	}
	deleted, err := uc.linkRepo.DeleteByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return fmt.Errorf("error deleting link: %w", err)
	}
	if !deleted {
		return ErrLinkNotFound
	}
	return nil
}
