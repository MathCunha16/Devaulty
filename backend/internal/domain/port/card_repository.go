package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type CardRepository interface {
	Save(ctx context.Context, card *model.Card) (*model.Card, error)
	FindByIDAndBoardID(ctx context.Context, boardID, id uuid.UUID) (*model.Card, error)
	FindAllByBoardIDAndProjectID(ctx context.Context, projectID, boardID uuid.UUID) ([]model.Card, error)
	GetNextPosition(ctx context.Context, columnID uuid.UUID) (uint16, error)
	MoveCard(ctx context.Context, cardID, sourceColumnID, targetColumnID uuid.UUID, newPosition uint16) error
	DeleteByIDAndBoardID(ctx context.Context, boardID, id uuid.UUID) (bool, error)

	// Support for linked items (card_items)

	SaveLinkedItems(ctx context.Context, cardID uuid.UUID, items []model.CardItem) error
	FindLinkedItemsByCardID(ctx context.Context, cardID uuid.UUID) ([]model.CardItem, error)
	FindLinkedItemsByCardIDs(ctx context.Context, cardIDs []uuid.UUID) (map[uuid.UUID][]model.CardItem, error)

	ProjectScopedRepository
}
