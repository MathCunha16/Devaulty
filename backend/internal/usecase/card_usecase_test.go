package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/dto"
	"devaulty-backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MOCK CARD REPOSITORY ---

type MockCardRepository struct {
	mock.Mock
}

func (m *MockCardRepository) Save(ctx context.Context, card *model.Card) (*model.Card, error) {
	args := m.Called(ctx, card)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Card), args.Error(1)
}

func (m *MockCardRepository) FindByIDAndBoardID(ctx context.Context, boardID, id uuid.UUID) (*model.Card, error) {
	args := m.Called(ctx, boardID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Card), args.Error(1)
}

func (m *MockCardRepository) FindAllByBoardIDAndProjectID(ctx context.Context, projectID, boardID uuid.UUID) ([]model.Card, error) {
	args := m.Called(ctx, projectID, boardID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Card), args.Error(1)
}

func (m *MockCardRepository) GetNextPosition(ctx context.Context, columnID uuid.UUID) (uint16, error) {
	args := m.Called(ctx, columnID)
	return uint16(args.Int(0)), args.Error(1)
}

func (m *MockCardRepository) MoveCard(ctx context.Context, cardID, sourceColumnID, targetColumnID uuid.UUID, newPosition uint16) error {
	args := m.Called(ctx, cardID, sourceColumnID, targetColumnID, newPosition)
	return args.Error(0)
}

func (m *MockCardRepository) DeleteByIDAndBoardID(ctx context.Context, boardID, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, boardID, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockCardRepository) SaveLinkedItems(ctx context.Context, cardID uuid.UUID, items []model.CardItem) error {
	args := m.Called(ctx, cardID, items)
	return args.Error(0)
}

func (m *MockCardRepository) FindLinkedItemsByCardID(ctx context.Context, cardID uuid.UUID) ([]model.CardItem, error) {
	args := m.Called(ctx, cardID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.CardItem), args.Error(1)
}

func (m *MockCardRepository) FindLinkedItemsByCardIDs(ctx context.Context, cardIDs []uuid.UUID) (map[uuid.UUID][]model.CardItem, error) {
	args := m.Called(ctx, cardIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[uuid.UUID][]model.CardItem), args.Error(1)
}

func (m *MockCardRepository) ExistsByIDAndProjectID(ctx context.Context, id, projectID uuid.UUID) (bool, error) {
	args := m.Called(ctx, id, projectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockCardRepository) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, ids, projectID)
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

// --- CREATE TESTS ---

func TestCardUseCase_Create_Success_WithLinkedItems(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()
	snippetID := uuid.New()
	priority := model.CardPriorityHigh
	desc := "Fix authentication bug"
	dueDate := time.Now().Add(24 * time.Hour)

	cmd := dto.CreateCardCommand{
		ProjectID:   projectID,
		BoardID:     boardID,
		ColumnID:    columnID,
		Title:       "Fix Bug",
		Description: &desc,
		Priority:    &priority,
		DueDate:     &dueDate,
		LinkedItems: []dto.CreateCardItemCommand{
			{ItemID: snippetID, ItemType: model.ItemTypeSnippet},
		},
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("ExistsByIDAndBoardIDAndProjectID", ctx, columnID, boardID, projectID).Return(true, nil)
	mockCardRepo.On("GetNextPosition", ctx, columnID).Return(0, nil)

	mockCardRepo.On("Save", ctx, mock.MatchedBy(func(c *model.Card) bool {
		return c.ColumnID == columnID && c.Title == "Fix Bug" && c.Position == 0 && *c.Priority == priority
	})).Return(&model.Card{
		ID:          uuid.New(),
		ColumnID:    columnID,
		Title:       "Fix Bug",
		Description: &desc,
		Position:    0,
		Priority:    &priority,
		DueDate:     &dueDate,
		BaseEntity:  model.BaseEntity{CreatedAt: time.Now()},
	}, nil)

	mockCardRepo.On("SaveLinkedItems", ctx, mock.Anything, mock.MatchedBy(func(items []model.CardItem) bool {
		return len(items) == 1 && items[0].ItemID == snippetID && items[0].ItemType == model.ItemTypeSnippet
	})).Return(nil)

	created, err := uc.Create(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Fix Bug", created.Title)
	assert.Equal(t, 1, len(created.LinkedItems))
	assert.Equal(t, snippetID, created.LinkedItems[0].ItemID)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
}

func TestCardUseCase_Create_Success_WithoutLinkedItems(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()

	cmd := dto.CreateCardCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		ColumnID:  columnID,
		Title:     "Simple Card",
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("ExistsByIDAndBoardIDAndProjectID", ctx, columnID, boardID, projectID).Return(true, nil)
	mockCardRepo.On("GetNextPosition", ctx, columnID).Return(3, nil)

	mockCardRepo.On("Save", ctx, mock.MatchedBy(func(c *model.Card) bool {
		return c.ColumnID == columnID && c.Title == "Simple Card" && c.Position == 3
	})).Return(&model.Card{
		ID:         uuid.New(),
		ColumnID:   columnID,
		Title:      "Simple Card",
		Position:   3,
		BaseEntity: model.BaseEntity{CreatedAt: time.Now()},
	}, nil)

	created, err := uc.Create(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Simple Card", created.Title)
	assert.Equal(t, uint16(3), created.Position)
	assert.Nil(t, created.LinkedItems)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
}

func TestCardUseCase_Create_ProjectNotFound(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	cmd := dto.CreateCardCommand{ProjectID: projectID}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestCardUseCase_Create_BoardNotFound(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	cmd := dto.CreateCardCommand{ProjectID: projectID, BoardID: boardID}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(false, nil)

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorIs(t, err, usecase.ErrBoardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestCardUseCase_Create_BoardColumnNotFound(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()
	cmd := dto.CreateCardCommand{ProjectID: projectID, BoardID: boardID, ColumnID: columnID}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("ExistsByIDAndBoardIDAndProjectID", ctx, columnID, boardID, projectID).Return(false, nil)

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorIs(t, err, usecase.ErrBoardColumnNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestCardUseCase_Create_GetNextPositionError(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()
	cmd := dto.CreateCardCommand{ProjectID: projectID, BoardID: boardID, ColumnID: columnID}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("ExistsByIDAndBoardIDAndProjectID", ctx, columnID, boardID, projectID).Return(true, nil)
	mockCardRepo.On("GetNextPosition", ctx, columnID).Return(0, errors.New("position error"))

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorIs(t, err, usecase.ErrNotPossibleToGetCardNextPosition)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
}

func TestCardUseCase_Create_SaveError(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()
	cmd := dto.CreateCardCommand{ProjectID: projectID, BoardID: boardID, ColumnID: columnID, Title: "Card"}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("ExistsByIDAndBoardIDAndProjectID", ctx, columnID, boardID, projectID).Return(true, nil)
	mockCardRepo.On("GetNextPosition", ctx, columnID).Return(0, nil)
	mockCardRepo.On("Save", ctx, mock.Anything).Return(nil, errors.New("save card error"))

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorContains(t, err, "save card error")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
}

func TestCardUseCase_Create_SaveLinkedItemsError(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()
	cmd := dto.CreateCardCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		ColumnID:  columnID,
		Title:     "Card",
		LinkedItems: []dto.CreateCardItemCommand{
			{ItemID: uuid.New(), ItemType: model.ItemTypeNote},
		},
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("ExistsByIDAndBoardIDAndProjectID", ctx, columnID, boardID, projectID).Return(true, nil)
	mockCardRepo.On("GetNextPosition", ctx, columnID).Return(0, nil)
	mockCardRepo.On("Save", ctx, mock.Anything).Return(&model.Card{ID: uuid.New()}, nil)
	mockCardRepo.On("SaveLinkedItems", ctx, mock.Anything, mock.Anything).Return(errors.New("linked items save error"))

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorContains(t, err, "error saving linked items")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
}

// --- GET BY ID TESTS ---

func TestCardUseCase_GetByID_Success(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	cardID := uuid.New()
	columnID := uuid.New()
	expectedCard := &model.Card{
		ID:         cardID,
		ColumnID:   columnID,
		Title:      "Feature Card",
		Position:   1,
		BaseEntity: model.BaseEntity{CreatedAt: time.Now()},
	}
	expectedLinkedItems := []model.CardItem{
		{CardID: cardID, ItemID: uuid.New(), ItemType: model.ItemTypeLink},
	}
	color := "#10B981"
	expectedTags := []model.Tag{
		{ID: uuid.New(), Name: "Backend", Color: &color},
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockCardRepo.On("FindByIDAndBoardID", ctx, boardID, cardID).Return(expectedCard, nil)
	mockCardRepo.On("FindLinkedItemsByCardID", ctx, cardID).Return(expectedLinkedItems, nil)
	mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeCard, projectID, cardID).Return(expectedTags, nil)

	result, err := uc.GetByID(ctx, projectID, boardID, cardID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cardID, result.ID)
	assert.Equal(t, 1, len(result.LinkedItems))
	assert.Equal(t, 1, len(result.Tags))
	assert.Equal(t, "Backend", result.Tags[0].Name)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

func TestCardUseCase_GetByID_ProjectNotFound(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	cardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	result, err := uc.GetByID(ctx, projectID, boardID, cardID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestCardUseCase_GetByID_BoardNotFound(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	cardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(false, nil)

	result, err := uc.GetByID(ctx, projectID, boardID, cardID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrBoardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestCardUseCase_GetByID_CardNotFound(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	cardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockCardRepo.On("FindByIDAndBoardID", ctx, boardID, cardID).Return(nil, nil)

	result, err := uc.GetByID(ctx, projectID, boardID, cardID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrCardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
}

func TestCardUseCase_GetByID_FindError(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	cardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockCardRepo.On("FindByIDAndBoardID", ctx, boardID, cardID).Return(nil, errors.New("db find error"))

	result, err := uc.GetByID(ctx, projectID, boardID, cardID)

	assert.Nil(t, result)
	assert.ErrorContains(t, err, "error fetching card")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
}

func TestCardUseCase_GetByID_LinkedItemsError(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	cardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockCardRepo.On("FindByIDAndBoardID", ctx, boardID, cardID).Return(&model.Card{ID: cardID}, nil)
	mockCardRepo.On("FindLinkedItemsByCardID", ctx, cardID).Return(nil, errors.New("linked items error"))

	result, err := uc.GetByID(ctx, projectID, boardID, cardID)

	assert.Nil(t, result)
	assert.ErrorContains(t, err, "error fetching linked items")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
}

func TestCardUseCase_GetByID_TagsError(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	cardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockCardRepo.On("FindByIDAndBoardID", ctx, boardID, cardID).Return(&model.Card{ID: cardID}, nil)
	mockCardRepo.On("FindLinkedItemsByCardID", ctx, cardID).Return([]model.CardItem{}, nil)
	mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeCard, projectID, cardID).Return(nil, errors.New("tag fetch error"))

	result, err := uc.GetByID(ctx, projectID, boardID, cardID)

	assert.Nil(t, result)
	assert.ErrorContains(t, err, "error fetching tags for card")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

// --- GET ALL BY BOARD ID TESTS ---

func TestCardUseCase_GetAllByBoardID_Success_WithCards(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	card1ID := uuid.New()
	card2ID := uuid.New()

	expectedCards := []model.Card{
		{ID: card1ID, Title: "Card 1", Position: 0},
		{ID: card2ID, Title: "Card 2", Position: 1},
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockCardRepo.On("FindAllByBoardIDAndProjectID", ctx, projectID, boardID).Return(expectedCards, nil)
	mockItemTagRepo.On("FindTagsForItems", ctx, model.ItemTypeCard, projectID, []uuid.UUID{card1ID, card2ID}).Return(map[uuid.UUID][]model.Tag{}, nil)

	views, err := uc.GetAllByBoardID(ctx, projectID, boardID)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(views))
	assert.Equal(t, card1ID, views[0].ID)
	assert.Equal(t, card2ID, views[1].ID)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

func TestCardUseCase_GetAllByBoardID_Success_Empty(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockCardRepo.On("FindAllByBoardIDAndProjectID", ctx, projectID, boardID).Return([]model.Card{}, nil)

	views, err := uc.GetAllByBoardID(ctx, projectID, boardID)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(views))
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
}

func TestCardUseCase_GetAllByBoardID_ProjectNotFound(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	views, err := uc.GetAllByBoardID(ctx, projectID, boardID)

	assert.Nil(t, views)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestCardUseCase_GetAllByBoardID_BoardNotFound(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(false, nil)

	views, err := uc.GetAllByBoardID(ctx, projectID, boardID)

	assert.Nil(t, views)
	assert.ErrorIs(t, err, usecase.ErrBoardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestCardUseCase_GetAllByBoardID_RepoError(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockCardRepo.On("FindAllByBoardIDAndProjectID", ctx, projectID, boardID).Return(nil, errors.New("db find error"))

	views, err := uc.GetAllByBoardID(ctx, projectID, boardID)

	assert.Nil(t, views)
	assert.ErrorContains(t, err, "error fetching cards")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
}

func TestCardUseCase_GetAllByBoardID_TagsError(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	cardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockCardRepo.On("FindAllByBoardIDAndProjectID", ctx, projectID, boardID).Return([]model.Card{{ID: cardID}}, nil)
	mockItemTagRepo.On("FindTagsForItems", ctx, model.ItemTypeCard, projectID, []uuid.UUID{cardID}).Return(nil, errors.New("tags error"))

	views, err := uc.GetAllByBoardID(ctx, projectID, boardID)

	assert.Nil(t, views)
	assert.ErrorContains(t, err, "error fetching tags for cards")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

// --- UPDATE TESTS ---

func TestCardUseCase_Update_Success_WithNewLinkedItems(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()
	cardID := uuid.New()
	newTitle := "Updated Title"
	newPriority := model.CardPriorityLow
	noteID := uuid.New()

	cmd := dto.UpdateCardCommand{
		ID:        cardID,
		ProjectID: projectID,
		BoardID:   boardID,
		Title:     &newTitle,
		Priority:  &newPriority,
		LinkedItems: []dto.CreateCardItemCommand{
			{ItemID: noteID, ItemType: model.ItemTypeNote},
		},
	}

	existingCard := &model.Card{
		ID:       cardID,
		ColumnID: columnID,
		Title:    "Old Title",
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockCardRepo.On("FindByIDAndBoardID", ctx, boardID, cardID).Return(existingCard, nil)
	mockCardRepo.On("Save", ctx, mock.MatchedBy(func(c *model.Card) bool {
		return c.ID == cardID && c.Title == "Updated Name" || c.Title == "Updated Title" && c.UpdatedAt != nil
	})).Return(&model.Card{
		ID:       cardID,
		ColumnID: columnID,
		Title:    "Updated Title",
		Priority: &newPriority,
	}, nil)
	mockCardRepo.On("SaveLinkedItems", ctx, cardID, mock.MatchedBy(func(items []model.CardItem) bool {
		return len(items) == 1 && items[0].ItemID == noteID
	})).Return(nil)
	mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeCard, projectID, cardID).Return([]model.Tag{}, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Title", updated.Title)
	assert.Equal(t, 1, len(updated.LinkedItems))
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

func TestCardUseCase_Update_Success_KeepExistingLinkedItems(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()
	cardID := uuid.New()
	newTitle := "Updated Title"

	cmd := dto.UpdateCardCommand{
		ID:          cardID,
		ProjectID:   projectID,
		BoardID:     boardID,
		Title:       &newTitle,
		LinkedItems: nil,
	}

	existingCard := &model.Card{
		ID:       cardID,
		ColumnID: columnID,
		Title:    "Old Title",
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockCardRepo.On("FindByIDAndBoardID", ctx, boardID, cardID).Return(existingCard, nil)
	mockCardRepo.On("Save", ctx, mock.Anything).Return(&model.Card{ID: cardID, ColumnID: columnID, Title: newTitle}, nil)
	mockCardRepo.On("FindLinkedItemsByCardID", ctx, cardID).Return([]model.CardItem{{CardID: cardID, ItemID: uuid.New()}}, nil)
	mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeCard, projectID, cardID).Return([]model.Tag{}, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, 1, len(updated.LinkedItems))
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

func TestCardUseCase_Update_ProjectNotFound(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	cmd := dto.UpdateCardCommand{ProjectID: projectID}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestCardUseCase_Update_CardNotFound(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	cardID := uuid.New()

	cmd := dto.UpdateCardCommand{
		ID:        cardID,
		ProjectID: projectID,
		BoardID:   boardID,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockCardRepo.On("FindByIDAndBoardID", ctx, boardID, cardID).Return(nil, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrCardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
}

// --- DELETE TESTS ---

func TestCardUseCase_Delete_Success(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	cardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockCardRepo.On("DeleteByIDAndBoardID", ctx, boardID, cardID).Return(true, nil)
	mockItemTagRepo.On("RemoveAllTagsFromItem", ctx, model.ItemTypeCard, cardID).Return(nil)

	err := uc.Delete(ctx, projectID, boardID, cardID)

	assert.NoError(t, err)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

func TestCardUseCase_Delete_NotFound(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	cardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockCardRepo.On("DeleteByIDAndBoardID", ctx, boardID, cardID).Return(false, nil)

	err := uc.Delete(ctx, projectID, boardID, cardID)

	assert.ErrorIs(t, err, usecase.ErrCardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
}

// --- MOVE TESTS ---

func TestCardUseCase_Move_Success(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	sourceColumnID := uuid.New()
	targetColumnID := uuid.New()
	cardID := uuid.New()
	newPosition := uint16(2)

	cmd := dto.MoveCardCommand{
		ID:             cardID,
		ProjectID:      projectID,
		BoardID:        boardID,
		TargetColumnID: targetColumnID,
		Position:       &newPosition,
	}

	existingCard := &model.Card{
		ID:       cardID,
		ColumnID: sourceColumnID,
		Title:    "Move Me",
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("ExistsByIDAndBoardIDAndProjectID", ctx, targetColumnID, boardID, projectID).Return(true, nil)
	mockCardRepo.On("FindByIDAndBoardID", ctx, boardID, cardID).Return(existingCard, nil)
	mockCardRepo.On("MoveCard", ctx, cardID, sourceColumnID, targetColumnID, newPosition).Return(nil)

	err := uc.Move(ctx, cmd)

	assert.NoError(t, err)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
}

func TestCardUseCase_Move_TargetColumnNotFound(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	targetColumnID := uuid.New()
	cardID := uuid.New()
	newPosition := uint16(1)

	cmd := dto.MoveCardCommand{
		ID:             cardID,
		ProjectID:      projectID,
		BoardID:        boardID,
		TargetColumnID: targetColumnID,
		Position:       &newPosition,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("ExistsByIDAndBoardIDAndProjectID", ctx, targetColumnID, boardID, projectID).Return(false, nil)

	err := uc.Move(ctx, cmd)

	assert.ErrorIs(t, err, usecase.ErrBoardColumnNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestCardUseCase_Move_CardNotFound(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	targetColumnID := uuid.New()
	cardID := uuid.New()
	newPosition := uint16(1)

	cmd := dto.MoveCardCommand{
		ID:             cardID,
		ProjectID:      projectID,
		BoardID:        boardID,
		TargetColumnID: targetColumnID,
		Position:       &newPosition,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("ExistsByIDAndBoardIDAndProjectID", ctx, targetColumnID, boardID, projectID).Return(true, nil)
	mockCardRepo.On("FindByIDAndBoardID", ctx, boardID, cardID).Return(nil, nil)

	err := uc.Move(ctx, cmd)

	assert.ErrorIs(t, err, usecase.ErrCardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
}

func TestCardUseCase_Move_RepoError(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewCardUseCase(mockCardRepo, mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	sourceColumnID := uuid.New()
	targetColumnID := uuid.New()
	cardID := uuid.New()
	newPosition := uint16(1)

	cmd := dto.MoveCardCommand{
		ID:             cardID,
		ProjectID:      projectID,
		BoardID:        boardID,
		TargetColumnID: targetColumnID,
		Position:       &newPosition,
	}

	existingCard := &model.Card{
		ID:       cardID,
		ColumnID: sourceColumnID,
		Title:    "Move Me",
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("ExistsByIDAndBoardIDAndProjectID", ctx, targetColumnID, boardID, projectID).Return(true, nil)
	mockCardRepo.On("FindByIDAndBoardID", ctx, boardID, cardID).Return(existingCard, nil)
	mockCardRepo.On("MoveCard", ctx, cardID, sourceColumnID, targetColumnID, newPosition).Return(errors.New("db move error"))

	err := uc.Move(ctx, cmd)

	assert.ErrorContains(t, err, "error moving card")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
}
