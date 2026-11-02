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

type DrawingRepositoryAdapter struct {
	db *sqlx.DB
}

func NewDrawingRepository(db *sqlx.DB) port.DrawingRepository {
	return &DrawingRepositoryAdapter{
		db: db,
	}
}

func (r *DrawingRepositoryAdapter) Save(ctx context.Context, drawing *model.Drawing) (*model.Drawing, error) {
	query := `INSERT INTO drawings (id, project_id, title, content, thumbnail, created_at, updated_at) 
				VALUES (:id, :project_id, :title, :content, :thumbnail, :created_at, :updated_at)
				ON CONFLICT (id) DO UPDATE SET title = excluded.title, content = excluded.content,
				                               thumbnail = excluded.thumbnail, updated_at = excluded.updated_at`
	_, err := r.db.NamedExecContext(ctx, query, drawing)
	if err != nil {
		return nil, fmt.Errorf("error trying to save drawing: %w", err)
	}
	return drawing, nil
}

func (r *DrawingRepositoryAdapter) FindByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*model.Drawing, error) {
	query := `SELECT * FROM drawings WHERE id = ? AND project_id = ?`
	var drawing model.Drawing
	err := r.db.GetContext(ctx, &drawing, query, id, projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error trying to find drawing: %w", err)
	}
	return &drawing, nil
}

func (r *DrawingRepositoryAdapter) FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[model.Drawing], error) {
	countQuery := `SELECT COUNT(*) FROM drawings WHERE project_id = ?`
	selectQuery := `SELECT id, project_id, title, thumbnail, created_at, updated_at FROM drawings WHERE project_id = ? ORDER BY created_at DESC`
	return PaginateExec[model.Drawing](ctx, r.db, countQuery, selectQuery, page, size, projectID)
}

func (r *DrawingRepositoryAdapter) DeleteByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error) {
	query := `DELETE FROM drawings WHERE id = ? AND project_id = ?`
	res, err := r.db.ExecContext(ctx, query, id, projectID)
	if err != nil {
		return false, fmt.Errorf("error trying to delete drawing: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error checking deleted rows: %w", err)
	}
	return rows > 0, nil
}

func (r *DrawingRepositoryAdapter) ExistsByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error) {
	return ExistsByIDAndProjectIDExec(ctx, r.db, "drawings", id, projectID)
}

func (r *DrawingRepositoryAdapter) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	return FindExistingIDsByProjectIDExec(ctx, r.db, "drawings", ids, projectID)
}
