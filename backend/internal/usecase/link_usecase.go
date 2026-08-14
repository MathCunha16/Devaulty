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
	ErrLinkNotFound = errors.New("link not found")
)

type LinkUseCase struct {
	linkRepo    port.LinkRepository
	projectRepo port.ProjectRepository
	itemTagRepo port.ItemTagRepository
}

func NewLinkUseCase(linkRepo port.LinkRepository, projectRepo port.ProjectRepository, itemTagRepo port.ItemTagRepository) *LinkUseCase {
	return &LinkUseCase{
		linkRepo:    linkRepo,
		projectRepo: projectRepo,
		itemTagRepo: itemTagRepo,
	}
}

func (uc *LinkUseCase) Create(ctx context.Context, cmd dto.CreateLinkCommand) (*dto.LinkView, error) {
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

	saved, err := uc.linkRepo.Save(ctx, &link)
	if err != nil {
		return nil, err
	}
	return mapLinkToView(saved, nil), nil
}

func (uc *LinkUseCase) GetByID(ctx context.Context, projectID, id uuid.UUID) (*dto.LinkView, error) {
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

	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeLink, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for link: %w", err)
	}

	return mapLinkToView(link, tags), nil
}

func (uc *LinkUseCase) GetAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[dto.LinkView], error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return model.Page[dto.LinkView]{}, err
	}
	linkPage, err := uc.linkRepo.FindAllByProjectID(ctx, projectID, page, size)
	if err != nil {
		return model.Page[dto.LinkView]{}, err
	}
	if len(linkPage.Content) == 0 {
		return model.NewPage([]dto.LinkView{}, linkPage.Number, linkPage.Size, linkPage.TotalElements), nil
	}

	linkIDs := make([]uuid.UUID, len(linkPage.Content))
	for i, l := range linkPage.Content {
		linkIDs[i] = l.ID
	}

	tagsMap, err := uc.itemTagRepo.FindTagsForItems(ctx, model.ItemTypeLink, projectID, linkIDs)
	if err != nil {
		return model.Page[dto.LinkView]{}, fmt.Errorf("error fetching tags for links: %w", err)
	}

	views := make([]dto.LinkView, len(linkPage.Content))
	for i, l := range linkPage.Content {
		tags := tagsMap[l.ID]
		views[i] = *mapLinkToView(&l, tags)
	}

	return model.NewPage(views, linkPage.Number, linkPage.Size, linkPage.TotalElements), nil
}

func (uc *LinkUseCase) Update(ctx context.Context, cmd dto.UpdateLinkCommand) (*dto.LinkView, error) {
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
	saved, err := uc.linkRepo.Save(ctx, link)
	if err != nil {
		return nil, err
	}
	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeLink, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for link: %w", err)
	}
	return mapLinkToView(saved, tags), nil
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
	if err := uc.itemTagRepo.RemoveAllTagsFromItem(ctx, model.ItemTypeLink, id); err != nil {
		log.Printf("warning: failed to remove tags from link %s: %v", id, err)
	}
	return nil
}

func mapLinkToView(link *model.Link, tags []model.Tag) *dto.LinkView {
	tagSummaries := make([]dto.TagSummary, len(tags))
	for i, t := range tags {
		tagSummaries[i] = dto.TagSummary{
			ID:    t.ID,
			Name:  t.Name,
			Color: t.Color,
		}
	}
	return &dto.LinkView{
		ID:          link.ID,
		ProjectID:   link.ProjectID,
		Title:       link.Title,
		URL:         link.Url,
		Description: link.Description,
		Tags:        tagSummaries,
		CreatedAt:   &link.CreatedAt,
		UpdatedAt:   link.UpdatedAt,
	}
}
