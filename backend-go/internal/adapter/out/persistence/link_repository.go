package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LinkRepositoryAdapter struct {
	db *sqlx.DB
}

func NewLinkRepository(db *sqlx.DB) port.LinkRepository {
	return &LinkRepositoryAdapter{db: db}
}

func (r *LinkRepositoryAdapter) Save(ctx context.Context, link *model.Link) (*model.Link, error) {
	query := `	INSERT INTO links (id, project_id, title, url, description, created_at, updated_at) 
				VALUES (:id, :project_id, :title, :url, :description, :created_at, :updated_at)
				ON CONFLICT(id) DO UPDATE SET title = excluded.title, url = excluded.url, description = excluded.description, updated_at = excluded.updated_at`
	_, err := r.db.NamedExecContext(ctx, query, link)
	if err != nil {
		return nil, fmt.Errorf("error trying to save link: %w", err)
	}

	return link, nil
}

func (r *LinkRepositoryAdapter) FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Link, error) {
	query := `SELECT * FROM links WHERE id = ? AND project_id = ?`
	var link model.Link
	err := r.db.GetContext(ctx, &link, query, id, projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error trying to find link: %w", err)
	}

	return &link, nil
}

func (r *LinkRepositoryAdapter) FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[model.Link], error) {
	countQuery := `SELECT COUNT(*) FROM links WHERE project_id = ?`
	selectQuery := `SELECT * FROM links WHERE project_id = ? ORDER BY created_at DESC`
	return PaginateExec[model.Link](ctx, r.db, countQuery, selectQuery, page, size, projectID)
}

func (r *LinkRepositoryAdapter) DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (bool, error) {
	query := `DELETE FROM links WHERE id = ? AND project_id = ?`
	res, err := r.db.ExecContext(ctx, query, id, projectID)
	if err != nil {
		return false, fmt.Errorf("error trying to delete link: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error checking deleted rows: %w", err)
	}
	return rows > 0, nil
}

func (r *LinkRepositoryAdapter) ExistsByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error) {
	return ExistsByIDAndProjectIDExec(ctx, r.db, "links", id, projectID)
}

func (r *LinkRepositoryAdapter) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	return FindExistingIDsByProjectIDExec(ctx, r.db, "links", ids, projectID)
}
