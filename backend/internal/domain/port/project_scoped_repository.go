package port

import (
	"context"

	"github.com/google/uuid"
)

type ProjectScopedRepository interface {
	ExistsByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error)
	FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error)
}
