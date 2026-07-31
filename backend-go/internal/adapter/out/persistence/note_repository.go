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

type NoteRepositoryAdapter struct {
	db *sqlx.DB
}

func NewNoteRepository(db *sqlx.DB) port.NoteRepository {
	return &NoteRepositoryAdapter{db: db}
}

func (r *NoteRepositoryAdapter) Save(ctx context.Context, note *model.Note) (*model.Note, error) {
	query := `	INSERT INTO notes (id, project_id, title, content, archived, created_at, updated_at) 
				VALUES (:id, :project_id, :title, :content, :archived, :created_at, :updated_at)
				ON CONFLICT (id) DO UPDATE SET title = excluded.title, content = excluded.content,
				                               archived = excluded.archived, updated_at = excluded.updated_at`
	_, err := r.db.NamedExecContext(ctx, query, note)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (r *NoteRepositoryAdapter) FindByID(ctx context.Context, id uuid.UUID) (*model.Note, error) {
	query := `SELECT * FROM notes WHERE id = ?`
	var note model.Note
	err := r.db.GetContext(ctx, &note, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error trying to find note: %w", err)
	}
	return &note, nil
}

func (r *NoteRepositoryAdapter) FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[model.Note], error) {
	countQuery := `SELECT COUNT(*) FROM notes WHERE project_id = ?`
	selectQuery := `SELECT * FROM notes WHERE project_id = ? ORDER BY created_at DESC`
	return PaginateExec[model.Note](ctx, r.db, countQuery, selectQuery, page, size, projectID)
}

func (r *NoteRepositoryAdapter) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notes WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *NoteRepositoryAdapter) ExistsByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error) {
	return ExistsByIDAndProjectIDExec(ctx, r.db, "notes", id, projectID)
}

func (r *NoteRepositoryAdapter) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	return FindExistingIDsByProjectIDExec(ctx, r.db, "notes", ids, projectID)
}
