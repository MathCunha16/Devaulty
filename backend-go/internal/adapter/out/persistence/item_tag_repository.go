package persistence

import (
	"context"
	"fmt"

	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ItemTagRepositoryAdapter struct {
	db *sqlx.DB
}

func NewItemTagRepository(db *sqlx.DB) port.ItemTagRepository {
	return &ItemTagRepositoryAdapter{db: db}
}

func (r *ItemTagRepositoryAdapter) AssociateTagToItem(ctx context.Context, tagID uuid.UUID, itemType model.ItemType, itemID uuid.UUID) error {
	query := `INSERT INTO item_tags (tag_id, item_type, item_id) VALUES (?, ?, ?) ON CONFLICT(tag_id, item_type, item_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, tagID, itemType, itemID)
	if err != nil {
		return fmt.Errorf("error trying to associate tag to item: %w", err)
	}
	return nil
}

func (r *ItemTagRepositoryAdapter) DisassembleTagFromItem(ctx context.Context, projectID uuid.UUID, tagID uuid.UUID, itemType model.ItemType, itemID uuid.UUID) error {
	query := `DELETE FROM item_tags WHERE tag_id = ? AND item_type = ? AND item_id = ? AND tag_id IN (SELECT id FROM tags WHERE project_id = ?)`
	_, err := r.db.ExecContext(ctx, query, tagID, itemType, itemID, projectID)
	if err != nil {
		return fmt.Errorf("error trying to disassemble tag from item: %w", err)
	}
	return nil
}

func (r *ItemTagRepositoryAdapter) RemoveAllTagsFromItem(ctx context.Context, itemType model.ItemType, itemID uuid.UUID) error {
	query := `DELETE FROM item_tags WHERE item_type = ? AND item_id = ?`
	_, err := r.db.ExecContext(ctx, query, itemType, itemID)
	if err != nil {
		return fmt.Errorf("error trying to remove all tags from item: %w", err)
	}
	return nil
}

func (r *ItemTagRepositoryAdapter) FindTagsForItem(ctx context.Context, itemType model.ItemType, projectID uuid.UUID, itemID uuid.UUID) ([]model.Tag, error) {
	query := `
		SELECT t.id, t.project_id, t.name, t.color, t.created_at, t.updated_at 
		FROM tags t
		INNER JOIN item_tags it ON t.id = it.tag_id
		WHERE it.item_type = ? AND it.item_id = ? AND t.project_id = ?
		ORDER BY t.name ASC
	`
	var tags []model.Tag
	err := r.db.SelectContext(ctx, &tags, query, itemType, itemID, projectID)
	if err != nil {
		return nil, fmt.Errorf("error trying to find tags for item: %w", err)
	}
	if tags == nil {
		tags = []model.Tag{}
	}
	return tags, nil
}

type itemTagJoinRow struct {
	ItemID uuid.UUID `db:"item_id"`
	model.Tag
}

func (r *ItemTagRepositoryAdapter) FindTagsForItems(ctx context.Context, itemType model.ItemType, projectID uuid.UUID, itemIDs []uuid.UUID) (map[uuid.UUID][]model.Tag, error) {
	result := make(map[uuid.UUID][]model.Tag)
	if len(itemIDs) == 0 {
		return result, nil
	}

	query := `
		SELECT it.item_id, t.id, t.project_id, t.name, t.color, t.created_at, t.updated_at 
		FROM tags t
		INNER JOIN item_tags it ON t.id = it.tag_id
		WHERE it.item_type = ? AND it.item_id IN (?) AND t.project_id = ?
		ORDER BY t.name ASC
	`
	query, args, err := sqlx.In(query, itemType, itemIDs, projectID)
	if err != nil {
		return nil, fmt.Errorf("error trying to build query for item tags: %w", err)
	}
	query = r.db.Rebind(query)

	var rows []itemTagJoinRow
	err = r.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error trying to find tags for items: %w", err)
	}

	for _, row := range rows {
		result[row.ItemID] = append(result[row.ItemID], row.Tag)
	}

	return result, nil
}
