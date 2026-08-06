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

func (r *ProblemRepositoryAdapter) FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Problem, error) {
	query := `SELECT * FROM problems WHERE id = ? AND project_id = ?`
	var problem model.Problem
	err := r.db.GetContext(ctx, &problem, query, id, projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error trying to find problem: %w", err)
	}
	return &problem, nil
}

func (r *ProblemRepositoryAdapter) FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[port.ProblemSummary], error) {
	countQuery := `SELECT COUNT(*) FROM problems WHERE project_id = ?`
	selectQuery := `SELECT id, project_id, title, status, severity, created_at, updated_at FROM problems
                                                                       WHERE project_id = ? ORDER BY created_at DESC`
	return PaginateExec[port.ProblemSummary](ctx, r.db, countQuery, selectQuery, page, size, projectID)
}

func (r *ProblemRepositoryAdapter) DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (bool, error) {
	query := `DELETE FROM problems WHERE id = ? AND project_id = ?`
	res, err := r.db.ExecContext(ctx, query, id, projectID)
	if err != nil {
		return false, fmt.Errorf("error trying to delete problem: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error checking deleted rows: %w", err)
	}

	return rows > 0, nil
}

func (r *ProblemRepositoryAdapter) ExistsByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error) {
	return ExistsByIDAndProjectIDExec(ctx, r.db, "problems", id, projectID)
}

func (r *ProblemRepositoryAdapter) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	return FindExistingIDsByProjectIDExec(ctx, r.db, "problems", ids, projectID)
}
