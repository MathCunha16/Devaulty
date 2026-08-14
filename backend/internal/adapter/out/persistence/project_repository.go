package persistence

import (
	"context"
	"database/sql"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProjectRepositoryAdapter struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) port.ProjectRepository {
	return &ProjectRepositoryAdapter{db: db}
}

func (r *ProjectRepositoryAdapter) Save(ctx context.Context, project *model.Project) (*model.Project, error) {
	query := `	INSERT INTO projects (id, name, description, icon, color, archived, created_at, updated_at)
				VALUES (:id, :name, :description, :icon, :color, :archived, :created_at, :updated_at)
				ON CONFLICT (id) DO UPDATE SET name = excluded.name, description = excluded.description,
					icon = excluded.icon, color = excluded.color,
					archived = excluded.archived, updated_at = excluded.updated_at`
	_, err := r.db.NamedExecContext(ctx, query, project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *ProjectRepositoryAdapter) FindAll(ctx context.Context, page int, size int) (model.Page[model.Project], error) {
	countQuery := `SELECT COUNT(*) FROM projects`
	selectQuery := `SELECT * FROM projects ORDER BY created_at DESC`
	return PaginateExec[model.Project](ctx, r.db, countQuery, selectQuery, page, size)
}

func (r *ProjectRepositoryAdapter) FindByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	query := `SELECT * FROM projects WHERE id = ?`
	var project model.Project
	err := r.db.GetContext(ctx, &project, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepositoryAdapter) DeleteByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `DELETE FROM projects WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (r *ProjectRepositoryAdapter) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM projects WHERE id = ?)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, err
	}
	return exists, nil
}
