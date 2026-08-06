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

type TagRepositoryAdapter struct {
	db *sqlx.DB
}

func NewTagRepository(db *sqlx.DB) port.TagRepository {
	return &TagRepositoryAdapter{db: db}
}

func (r *TagRepositoryAdapter) Save(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	query := `	INSERT INTO tags (id, project_id, name, color, created_at, updated_at) 
				VALUES (:id, :project_id, :name, :color, :created_at, :updated_at)
				ON CONFLICT (id) DO UPDATE SET name = excluded.name, color = excluded.color,
				                               updated_at = excluded.updated_at`
	_, err := r.db.NamedExecContext(ctx, query, tag)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (r *TagRepositoryAdapter) FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Tag, error) {
	query := `SELECT * FROM tags WHERE id = ? AND project_id = ?`
	var tag model.Tag
	err := r.db.GetContext(ctx, &tag, query, id, projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepositoryAdapter) FindAllByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Tag, error) {
	query := `SELECT * FROM tags WHERE project_id = ? ORDER BY created_at DESC`
	var tags []model.Tag
	err := r.db.SelectContext(ctx, &tags, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("error trying to find tags: %w", err)
	}
	return tags, nil
}

func (r *TagRepositoryAdapter) SearchByNameAndProjectID(ctx context.Context, name string, projectID uuid.UUID) ([]model.Tag, error) {
	query := `SELECT * FROM tags WHERE project_id = ? AND (? = '' OR LOWER(name) LIKE LOWER(?)) ORDER BY name ASC`
	var tags []model.Tag
	err := r.db.SelectContext(ctx, &tags, query, projectID, "%"+name+"%", "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("error trying to find tags: %w", err)
	}
	return tags, nil
}

func (r *TagRepositoryAdapter) DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (bool, error) {
	query := `DELETE FROM tags WHERE id = ? AND project_id = ?`
	res, err := r.db.ExecContext(ctx, query, id, projectID)
	if err != nil {
		return false, fmt.Errorf("error trying to delete tag: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error checking deleted rows: %w", err)
	}
	return rows > 0, nil
}

func (r *TagRepositoryAdapter) ExistsByNameAndProjectID(ctx context.Context, name string, projectID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tags WHERE project_id = ? AND LOWER(name) = LOWER(?))`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, projectID, name)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *TagRepositoryAdapter) ExistsByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tags WHERE id = ? AND project_id = ?)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, id, projectID)
	if err != nil {
		return false, err
	}
	return exists, nil
}
