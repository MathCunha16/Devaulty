package persistence

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func ExistsByIDAndProjectIDExec(ctx context.Context, db *sqlx.DB, tableName string, id, projectID uuid.UUID) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = ? AND project_id = ?)", tableName)
	var exists bool
	err := db.GetContext(ctx, &exists, query, id, projectID)
	if err != nil {
		return false, fmt.Errorf("error checking existence in %s: %w", tableName, err)
	}
	return exists, nil
}

func FindExistingIDsByProjectIDExec(ctx context.Context, db *sqlx.DB, tableName string, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	if len(ids) == 0 {
		return []uuid.UUID{}, nil
	}

	query := fmt.Sprintf("SELECT id FROM %s WHERE id IN (?) AND project_id = ?", tableName)
	query, args, err := sqlx.In(query, ids, projectID)
	if err != nil {
		return nil, fmt.Errorf("error trying to build query: %w", err)
	}
	query = db.Rebind(query)
	var existingIDs []uuid.UUID
	err = db.SelectContext(ctx, &existingIDs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error trying to find existing ids in %s: %w", tableName, err)
	}

	return existingIDs, nil
}
