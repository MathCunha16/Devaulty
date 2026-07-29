package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type ItemTagRepository interface {
	AssociateTagToItem(ctx context.Context, tagID uuid.UUID, itemType model.ItemType, itemID uuid.UUID) error
	DisassembleTagFromItem(ctx context.Context, projectID uuid.UUID, tagID uuid.UUID, itemType model.ItemType, itemID uuid.UUID) error
	RemoveAllTagsFromItem(ctx context.Context, itemType model.ItemType, itemID uuid.UUID) error
	FindTagsForItem(ctx context.Context, itemType model.ItemType, projectID uuid.UUID, itemID uuid.UUID) ([]model.Tag, error)
	FindTagsForItems(ctx context.Context, itemType model.ItemType, projectID uuid.UUID, itemIDs []uuid.UUID) (map[uuid.UUID][]model.Tag, error)
}
