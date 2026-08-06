package usecase

import (
	"context"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUnsupportedItemType = errors.New("unsupported item type")
	ErrItemNotFound        = errors.New("item not found")
)

type ItemTagUseCase struct {
	itemTagRepo port.ItemTagRepository
	tagRepo     port.TagRepository
	projectRepo port.ProjectRepository
	itemRepos   map[model.ItemType]port.ProjectScopedRepository
}

func NewItemTagUseCase(
	itemTagRepo port.ItemTagRepository,
	tagRepo port.TagRepository,
	projectRepo port.ProjectRepository,
	snippetRepo port.SnippetRepository,
	linkRepo port.LinkRepository,
	problemRepo port.ProblemRepository,
) *ItemTagUseCase {
	return &ItemTagUseCase{
		itemTagRepo: itemTagRepo,
		tagRepo:     tagRepo,
		projectRepo: projectRepo,
		itemRepos: map[model.ItemType]port.ProjectScopedRepository{
			model.ItemTypeSnippet: snippetRepo,
			model.ItemTypeLink:    linkRepo,
			model.ItemTypeProblem: problemRepo,
		},
	}
}

func (uc *ItemTagUseCase) AssociateTagToItem(ctx context.Context, projectID uuid.UUID, itemType model.ItemType, itemID, tagID uuid.UUID) error {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return err
	}
	tagExists, err := uc.tagRepo.ExistsByIDAndProjectID(ctx, tagID, projectID)
	if err != nil {
		return err
	}
	if !tagExists {
		return ErrTagNotFound
	}
	if err := uc.validateItemOwnership(ctx, projectID, itemType, itemID); err != nil {
		return err
	}

	return uc.itemTagRepo.AssociateTagToItem(ctx, projectID, tagID, itemType, itemID)
}

func (uc *ItemTagUseCase) DisassociateTagFromItem(ctx context.Context, projectID uuid.UUID, itemType model.ItemType, itemID, tagID uuid.UUID) error {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return err
	}
	tagExists, err := uc.tagRepo.ExistsByIDAndProjectID(ctx, tagID, projectID)
	if err != nil {
		return err
	}
	if !tagExists {
		return ErrTagNotFound
	}
	if err := uc.validateItemOwnership(ctx, projectID, itemType, itemID); err != nil {
		return err
	}
	return uc.itemTagRepo.DisassembleTagFromItem(ctx, projectID, tagID, itemType, itemID)

}

// useful methods for listing tags for a project
func (uc *ItemTagUseCase) GetTagsForItem(ctx context.Context, projectID uuid.UUID, itemType model.ItemType, itemID uuid.UUID) ([]model.Tag, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return nil, err
	}
	if err := uc.validateItemOwnership(ctx, projectID, itemType, itemID); err != nil {
		return nil, err
	}

	return uc.itemTagRepo.FindTagsForItem(ctx, itemType, projectID, itemID)
}

func (uc *ItemTagUseCase) GetTagsForItems(ctx context.Context, projectID uuid.UUID, itemType model.ItemType, itemIDs []uuid.UUID) (map[uuid.UUID][]model.Tag, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return nil, err
	}
	if len(itemIDs) == 0 {
		return map[uuid.UUID][]model.Tag{}, nil
	}
	if err := uc.validateItemsOwnership(ctx, projectID, itemType, itemIDs); err != nil {
		return nil, err
	}

	return uc.itemTagRepo.FindTagsForItems(ctx, itemType, projectID, itemIDs)
}

// helper methods
func (uc *ItemTagUseCase) validateItemOwnership(ctx context.Context, projectID uuid.UUID, itemType model.ItemType, itemID uuid.UUID) error {
	repo, ok := uc.itemRepos[itemType]
	if !ok {
		return ErrUnsupportedItemType
	}
	exists, err := repo.ExistsByIDAndProjectID(ctx, itemID, projectID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrItemNotFound
	}
	return nil
}

func (uc *ItemTagUseCase) validateItemsOwnership(ctx context.Context, projectID uuid.UUID, itemType model.ItemType, itemIDs []uuid.UUID) error {
	repo, ok := uc.itemRepos[itemType]
	if !ok {
		return ErrUnsupportedItemType
	}
	existingIDs, err := repo.FindExistingIDsByProjectID(ctx, itemIDs, projectID)
	if err != nil {
		return err
	}
	if len(existingIDs) != len(itemIDs) {
		return ErrItemNotFound
	}
	return nil
}
