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

type ProblemRepositoryAdapter struct {
	db *sqlx.DB
}

func NewProblemRepository(db *sqlx.DB) port.ProblemRepository {
	return &ProblemRepositoryAdapter{db: db}
}

func (r *ProblemRepositoryAdapter) Save(ctx context.Context, problem *model.Problem) (*model.Problem, error) {
	query := `	INSERT into problems (
	            id, project_id, title, error_description, solution,
	            status, severity, created_at, updated_at
	        	) VALUES (:id , :project_id, :title, :error_description, :solution, :status, :severity, :created_at, :updated_at)
	            ON CONFLICT(id) DO UPDATE SET title = excluded.title, error_description = excluded.error_description,
	            solution = excluded.solution, status = excluded.status,
	            severity = excluded.severity, updated_at = excluded.updated_at`
	_, err := r.db.NamedExecContext(ctx, query, problem)
	if err != nil {
		return nil, fmt.Errorf("error trying to save problem: %w", err)
	}
	return problem, nil
}

func (r *ProblemRepositoryAdapter) FindByID(ctx context.Context, id uuid.UUID) (*model.Problem, error) {
	query := `SELECT * FROM problems WHERE id = ?`
	var problem model.Problem
	err := r.db.GetContext(ctx, &problem, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error trying to find problem: %w", err)
	}
	return &problem, nil
}

func (r *ProblemRepositoryAdapter) FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[model.Problem], error) {
	countQuery := `SELECT COUNT(*) FROM problems WHERE project_id = ?`
	selectQuery := `SELECT * FROM problems WHERE project_id = ? ORDER BY created_at DESC`
	return PaginateExec[model.Problem](ctx, r.db, countQuery, selectQuery, page, size, projectID)
}

func (r *ProblemRepositoryAdapter) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM problems WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error trying to delete problem: %w", err)
	}
	return nil
}

func (r *ProblemRepositoryAdapter) ExistsByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error) {
	return ExistsByIDAndProjectIDExec(ctx, r.db, "problems", id, projectID)
}

func (r *ProblemRepositoryAdapter) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	return FindExistingIDsByProjectIDExec(ctx, r.db, "problems", ids, projectID)
}
