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

type BoardRepositoryAdapter struct {
	db *sqlx.DB
}

func NewBoardRepository(db *sqlx.DB) port.BoardRepository {
	return &BoardRepositoryAdapter{db: db}
}

func (r *BoardRepositoryAdapter) Save(ctx context.Context, board *model.Board) (*model.Board, error) {
	query := ` 	INSERT INTO boards (id, project_id, name, description, is_default, created_at, updated_at)
				VALUES(:id, :project_id, :name, :description, :is_default, :created_at, :updated_at)
 				ON CONFLICT (id) DO UPDATE SET name = excluded.name, description = excluded.description,
 					is_default = excluded.is_default, updated_at = excluded.updated_at`
	_, err := r.db.NamedExecContext(ctx, query, board)
	if err != nil {
		return nil, fmt.Errorf("error trying to save board: %w", err)
	}
	return board, nil
}

func (r *BoardRepositoryAdapter) FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Board, error) {
	query := `SELECT * FROM boards WHERE id = ? AND project_id = ?`
	var board model.Board
	err := r.db.GetContext(ctx, &board, query, id, projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error trying to find board: %w", err)
	}
	return &board, nil
}

func (r *BoardRepositoryAdapter) FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[model.Board], error) {
	countQuery := `SELECT COUNT(*) FROM boards WHERE project_id = ?`
	selectQuery := `SELECT * FROM boards WHERE project_id = ? ORDER BY created_at DESC`
	return PaginateExec[model.Board](ctx, r.db, countQuery, selectQuery, page, size, projectID)
}

func (r *BoardRepositoryAdapter) FindDefaultByProjectID(ctx context.Context, projectID uuid.UUID) (*model.Board, error) {
	query := `SELECT * FROM boards WHERE project_id = ? AND is_default = true`
	var board model.Board
	err := r.db.GetContext(ctx, &board, query, projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error trying to find default board: %w", err)
	}
	return &board, nil
}

func (r *BoardRepositoryAdapter) DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (bool, error) {
	query := `DELETE FROM boards WHERE ID = ? AND project_id = ?`
	res, err := r.db.ExecContext(ctx, query, id, projectID)
	if err != nil {
		return false, fmt.Errorf("error trying to delete board: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error checking deleted rows: %w", err)
	}
	return rows > 0, nil
}

func (r *BoardRepositoryAdapter) UnsetAllDefaultsByProjectID(ctx context.Context, projectID uuid.UUID) (bool, error) {
	query := `UPDATE boards SET is_default = false WHERE project_id = ?`
	_, err := r.db.ExecContext(ctx, query, projectID)
	if err != nil {
		return false, fmt.Errorf("error trying to unset all defaults: %w", err)
	}
	return true, nil
}

func (r *BoardRepositoryAdapter) ExistsByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error) {
	return ExistsByIDAndProjectIDExec(ctx, r.db, "boards", id, projectID)
}

func (r *BoardRepositoryAdapter) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	return FindExistingIDsByProjectIDExec(ctx, r.db, "boards", ids, projectID)
}
