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

type SnippetRepositoryAdapter struct {
	db *sqlx.DB
}

func NewSnippetRepository(db *sqlx.DB) port.SnippetRepository {
	return &SnippetRepositoryAdapter{db: db}
}

func (r *SnippetRepositoryAdapter) Save(ctx context.Context, snippet *model.Snippet) (*model.Snippet, error) {
	query := `	INSERT INTO snippets (id, project_id, title, description, content, language,
	                      snippet_type, created_at, updated_at)
				VALUES (:id, :project_id, :title, :description, :content, :language,
				        :snippet_type, :created_at, :updated_at)
				ON CONFLICT (id) DO UPDATE SET title = excluded.title, description = excluded.description,
						content = excluded.content, language = excluded.language, 
						snippet_type = excluded.snippet_type, updated_at = excluded.updated_at`
	_, err := r.db.NamedExecContext(ctx, query, snippet)
	if err != nil {
		return nil, fmt.Errorf("error trying to save snippet: %w", err)
	}
	return snippet, nil
}

func (r *SnippetRepositoryAdapter) FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Snippet, error) {
	query := `SELECT * FROM snippets WHERE id = ? AND project_id = ?`
	var snippet model.Snippet
	err := r.db.GetContext(ctx, &snippet, query, id, projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error trying to find snippet: %w", err)
	}
	return &snippet, nil
}

func (r *SnippetRepositoryAdapter) FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[model.Snippet], error) {
	countQuery := `SELECT COUNT(*) FROM snippets WHERE project_id = ?`
	selectQuery := `SELECT * FROM snippets WHERE project_id = ? ORDER BY created_at DESC`
	return PaginateExec[model.Snippet](ctx, r.db, countQuery, selectQuery, page, size, projectID)
}

func (r *SnippetRepositoryAdapter) DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) error {
	query := `DELETE FROM snippets WHERE id = ? AND project_id = ?`
	_, err := r.db.ExecContext(ctx, query, id, projectID)
	if err != nil {
		return fmt.Errorf("error trying to delete snippet: %w", err)
	}
	return nil
}

func (r *SnippetRepositoryAdapter) ExistsByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error) {
	return ExistsByIDAndProjectIDExec(ctx, r.db, "snippets", id, projectID)
}

func (r *SnippetRepositoryAdapter) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	return FindExistingIDsByProjectIDExec(ctx, r.db, "snippets", ids, projectID)
}
