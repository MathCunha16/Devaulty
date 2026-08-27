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

type BoardColumnRepositoryAdapter struct {
	db *sqlx.DB
}

func NewBoardColumnRepository(db *sqlx.DB) port.BoardColumnRepository {
	return &BoardColumnRepositoryAdapter{db: db}
}

func (r *BoardColumnRepositoryAdapter) Save(ctx context.Context, column *model.BoardColumn) (*model.BoardColumn, error) {
	query := `	INSERT INTO board_columns (id, board_id, name, position, wip_limit, created_at, updated_at)
	 			VALUES (:id, :board_id, :name, :position, :wip_limit, :created_at, :updated_at)	
	 			ON CONFLICT (id) DO UPDATE SET name = excluded.name, position = excluded.position, wip_limit = excluded.wip_limit,
				updated_at = excluded.updated_at`
	_, err := r.db.NamedExecContext(ctx, query, column)
	if err != nil {
		return nil, fmt.Errorf("error trying to save board column: %w", err)
	}
	return column, nil
}

func (r *BoardColumnRepositoryAdapter) FindByIDAndBoardID(ctx context.Context, boardID, id uuid.UUID) (*model.BoardColumn, error) {
	query := `SELECT * FROM board_columns WHERE id = ? AND board_id = ?`
	var column model.BoardColumn
	err := r.db.GetContext(ctx, &column, query, id, boardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &column, nil
}

func (r *BoardColumnRepositoryAdapter) FindAllByBoardIDAndProjectID(ctx context.Context, projectID, boardID uuid.UUID) ([]model.BoardColumn, error) {
	query := `	SELECT c.* FROM board_columns c 
				JOIN boards b ON c.board_id = b.id
				WHERE c.board_id = ? AND b.project_id = ?
				ORDER BY c.position ASC`
	var columns []model.BoardColumn
	err := r.db.SelectContext(ctx, &columns, query, boardID, projectID)
	if err != nil {
		return nil, err
	}
	return columns, nil
}

func (r *BoardColumnRepositoryAdapter) GetNextPosition(ctx context.Context, boardID uuid.UUID) (uint8, error) {
	query := `SELECT COALESCE(MAX(position) + 1, 0) FROM board_columns WHERE board_id = ?`
	var nextPos uint8
	err := r.db.GetContext(ctx, &nextPos, query, boardID)
	if err != nil {
		return 0, fmt.Errorf("error trying to get next position: %w", err)
	}
	return nextPos, nil
}

func (r *BoardColumnRepositoryAdapter) Reorder(ctx context.Context, boardID uuid.UUID, columnsIDs []uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error trying to start transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `UPDATE board_columns SET position = ? WHERE board_id = ? AND id = ?`
	for position, columnID := range columnsIDs {
		res, err := tx.ExecContext(ctx, query, position, boardID, columnID)
		if err != nil {
			return fmt.Errorf("error updating position for column %s: %w", columnID, err)
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("error checking rows affected for column %s: %w", columnID, err)
		}
		if rows == 0 {
			return fmt.Errorf("board column %s not found in board %s", columnID, boardID)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing column reorder: %w", err)
	}
	return nil
}

func (r *BoardColumnRepositoryAdapter) DeleteByIDAndBoardID(ctx context.Context, boardID, id uuid.UUID) (bool, error) {
	query := `DELETE FROM board_columns WHERE id = ? AND board_id = ?`
	res, err := r.db.ExecContext(ctx, query, id, boardID)
	if err != nil {
		return false, fmt.Errorf("error trying to delete board column: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error checking deleted rows: %w", err)
	}
	return rows > 0, nil
}

func (r *BoardColumnRepositoryAdapter) ExistsByIDAndBoardIDAndProjectID(ctx context.Context, id, boardID, projectID uuid.UUID) (bool, error) {
	query := `	SELECT EXISTS(SELECT 1 FROM board_columns c 
					JOIN boards b ON c.board_id = b.id
					WHERE c.id = ? AND c.board_id = ? AND b.project_id = ?)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, id, boardID, projectID)
	if err != nil {
		return false, fmt.Errorf("error checking existence in board columns: %w", err)
	}
	return exists, nil
}
