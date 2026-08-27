package persistence

import (
	"context"
	"database/sql"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CardRepositoryAdapter struct {
	db *sqlx.DB
}

func NewCardRepository(db *sqlx.DB) port.CardRepository {
	return &CardRepositoryAdapter{db: db}
}

func (r *CardRepositoryAdapter) Save(ctx context.Context, card *model.Card) (*model.Card, error) {
	query := `	INSERT INTO cards (id, column_id, title, description, position, priority, due_date, created_at, updated_at)
				VALUES (:id, :column_id, :title, :description, :position, :priority, :due_date, :created_at, :updated_at)
				ON CONFLICT (id) DO UPDATE SET title = excluded.title, description = excluded.description,
				position = excluded.position, priority = excluded.priority, due_date = excluded.due_date, updated_at = excluded.updated_at`
	_, err := r.db.NamedExecContext(ctx, query, card)
	if err != nil {
		return nil, fmt.Errorf("error trying to save card: %w", err)
	}
	return card, nil
}

func (r *CardRepositoryAdapter) FindByIDAndBoardID(ctx context.Context, boardID, id uuid.UUID) (*model.Card, error) {
	query := `	SELECT c.* FROM cards c                                                                                                        
            	JOIN board_columns col ON c.column_id = col.id                                                                                 
            	WHERE c.id = ? AND col.board_id = ?`
	var card model.Card
	err := r.db.GetContext(ctx, &card, query, id, boardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error trying to find card: %w", err)
	}
	return &card, nil
}

func (r *CardRepositoryAdapter) FindAllByBoardIDAndProjectID(ctx context.Context, projectID, boardID uuid.UUID) ([]model.Card, error) {
	query := `	SELECT c.* FROM cards c
            	JOIN board_columns col ON c.column_id = col.id
            	JOIN boards b ON col.board_id = b.id
            	WHERE col.board_id = ? AND b.project_id = ?
            	ORDER BY col.position ASC, c.position ASC`
	var cards []model.Card
	err := r.db.SelectContext(ctx, &cards, query, boardID, projectID)
	if err != nil {
		return nil, fmt.Errorf("error trying to find cards: %w", err)
	}
	return cards, nil
}

func (r *CardRepositoryAdapter) GetNextPosition(ctx context.Context, columnID uuid.UUID) (uint16, error) {
	query := `SELECT COALESCE(MAX(position) + 1, 0) FROM cards WHERE column_id = ?`
	var nextPos uint16
	err := r.db.GetContext(ctx, &nextPos, query, columnID)
	if err != nil {
		return 0, fmt.Errorf("error trying to get next position for card: %w", err)
	}
	return nextPos, nil
}

func (r *CardRepositoryAdapter) MoveCard(ctx context.Context, cardID, sourceColumnID, targetColumnID uuid.UUID, newPosition uint16) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction to move card: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var currentPos uint16
	err = tx.GetContext(ctx, &currentPos, `SELECT position FROM cards WHERE id = ?`, cardID)
	if err != nil {
		return fmt.Errorf("error fetching current position of card: %w", err)
	}

	if sourceColumnID == targetColumnID {
		if newPosition < currentPos {
			_, err = tx.ExecContext(ctx, `	UPDATE cards 
													SET position = position + 1 
													WHERE column_id = ? AND position >= ? AND position < ? AND id != ?`,
				sourceColumnID, newPosition, currentPos, cardID)
			if err != nil {
				return fmt.Errorf("error shifting cards down: %w", err)
			}
		} else if newPosition > currentPos {
			_, err = tx.ExecContext(ctx, `	UPDATE cards 
													SET position = position - 1 
													WHERE column_id = ? AND position > ? AND position <= ? AND id != ?`,
				sourceColumnID, currentPos, newPosition, cardID)
			if err != nil {
				return fmt.Errorf("error shifting cards up: %w", err)
			}
		}

		_, err = tx.ExecContext(ctx, `UPDATE cards SET position = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, newPosition, cardID)
		if err != nil {
			return fmt.Errorf("error updating card position: %w", err)
		}
	} else {
		// Close gap in source column
		_, err = tx.ExecContext(ctx, `	UPDATE cards 
												SET position = position - 1 
												WHERE column_id = ? AND position > ?`,
			sourceColumnID, currentPos)
		if err != nil {
			return fmt.Errorf("error closing gap in source column: %w", err)
		}

		// Make room in target column
		_, err = tx.ExecContext(ctx, `	UPDATE cards 
												SET position = position + 1 
												WHERE column_id = ? AND position >= ?`,
			targetColumnID, newPosition)
		if err != nil {
			return fmt.Errorf("error making room in target column: %w", err)
		}

		// Move the card
		_, err = tx.ExecContext(ctx, `	UPDATE cards 
												SET column_id = ?, position = ?, updated_at = CURRENT_TIMESTAMP 
												WHERE id = ?`,
			targetColumnID, newPosition, cardID)
		if err != nil {
			return fmt.Errorf("error moving card to target column: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing move card transaction: %w", err)
	}
	return nil
}

func (r *CardRepositoryAdapter) DeleteByIDAndBoardID(ctx context.Context, boardID, id uuid.UUID) (bool, error) {
	query := `	DELETE FROM cards 
				WHERE id = ? AND column_id IN 
				(SELECT id FROM board_columns WHERE board_id = ?)`
	res, err := r.db.ExecContext(ctx, query, id, boardID)
	if err != nil {
		return false, fmt.Errorf("error trying to delete card: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error checking deleted rows: %w", err)
	}
	return rows > 0, nil
}

func (r *CardRepositoryAdapter) SaveLinkedItems(ctx context.Context, cardID uuid.UUID, items []model.CardItem) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction to save linked items: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx, `DELETE FROM card_items WHERE card_id = ?`, cardID)
	if err != nil {
		return fmt.Errorf("error removing existing linked items: %w", err)
	}

	if len(items) > 0 {
		insertQuery := `	INSERT INTO card_items (card_id, item_type, item_id, created_at, updated_at)
							VALUES (:card_id, :item_type, :item_id, :created_at, :updated_at)`
		for _, item := range items {
			item.CardID = cardID
			_, err = tx.NamedExecContext(ctx, insertQuery, item)
			if err != nil {
				return fmt.Errorf("error inserting linked item: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing linked items transaction: %w", err)
	}
	return nil
}

func (r *CardRepositoryAdapter) FindLinkedItemsByCardID(ctx context.Context, cardID uuid.UUID) ([]model.CardItem, error) {
	query := `SELECT card_id, item_type, item_id, created_at, updated_at FROM card_items WHERE card_id = ? ORDER BY created_at ASC`
	var items []model.CardItem
	err := r.db.SelectContext(ctx, &items, query, cardID)
	if err != nil {
		return nil, fmt.Errorf("error trying to find linked items by card id: %w", err)
	}
	return items, nil
}

func (r *CardRepositoryAdapter) FindLinkedItemsByCardIDs(ctx context.Context, cardIDs []uuid.UUID) (map[uuid.UUID][]model.CardItem, error) {
	result := make(map[uuid.UUID][]model.CardItem)
	if len(cardIDs) == 0 {
		return result, nil
	}

	query, args, err := sqlx.In(`SELECT card_id, item_type, item_id, created_at, updated_at FROM card_items WHERE card_id IN (?) ORDER BY created_at ASC`, cardIDs)
	if err != nil {
		return nil, fmt.Errorf("error constructing in query for linked items: %w", err)
	}
	query = r.db.Rebind(query)

	var items []model.CardItem
	err = r.db.SelectContext(ctx, &items, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying linked items by card IDs: %w", err)
	}

	for _, item := range items {
		result[item.CardID] = append(result[item.CardID], item)
	}
	return result, nil
}

func (r *CardRepositoryAdapter) ExistsByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error) {
	query := `	SELECT EXISTS(
				SELECT 1 FROM cards c
				JOIN board_columns col ON c.column_id = col.id
				JOIN boards b ON col.board_id = b.id
				WHERE c.id = ? AND b.project_id = ?)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, id, projectID)
	if err != nil {
		return false, fmt.Errorf("error checking if card exists: %w", err)
	}
	return exists, nil
}

func (r *CardRepositoryAdapter) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	if len(ids) == 0 {
		return []uuid.UUID{}, nil
	}

	query, args, err := sqlx.In(`	SELECT c.id FROM cards c
										JOIN board_columns col ON c.column_id = col.id
										JOIN boards b ON col.board_id = b.id
										WHERE b.project_id = ? AND c.id IN (?)`, projectID, ids)
	if err != nil {
		return nil, fmt.Errorf("error building query for existing card IDs: %w", err)
	}
	query = r.db.Rebind(query)

	var existingIDs []uuid.UUID
	err = r.db.SelectContext(ctx, &existingIDs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error finding existing card IDs: %w", err)
	}
	return existingIDs, nil
}
