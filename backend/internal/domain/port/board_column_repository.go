package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type BoardColumnRepository interface {
	Save(ctx context.Context, column *model.BoardColumn) (*model.BoardColumn, error)
	FindByIDAndBoardID(ctx context.Context, boardID, id uuid.UUID) (*model.BoardColumn, error)
	FindAllByBoardIDAndProjectID(ctx context.Context, projectID, boardID uuid.UUID) ([]model.BoardColumn, error)
	GetNextPosition(ctx context.Context, boardID uuid.UUID) (uint8, error)
	Reorder(ctx context.Context, boardID uuid.UUID, columnsIDs []uuid.UUID) error
	DeleteByIDAndBoardID(ctx context.Context, boardID, id uuid.UUID) (bool, error)
	ExistsByIDAndBoardIDAndProjectID(ctx context.Context, id, boardID, projectID uuid.UUID) (bool, error)
}
