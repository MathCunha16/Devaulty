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

// --- MOCKS ---

type MockBoardRepository struct {
	mock.Mock
}

func (m *MockBoardRepository) Save(ctx context.Context, board *model.Board) (*model.Board, error) {
	args := m.Called(ctx, board)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Board), args.Error(1)
}

func (m *MockBoardRepository) FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Board, error) {
	args := m.Called(ctx, projectID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Board), args.Error(1)
}

func (m *MockBoardRepository) FindDefaultByProjectID(ctx context.Context, projectID uuid.UUID) (*model.Board, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Board), args.Error(1)
}

func (m *MockBoardRepository) FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[model.Board], error) {
	args := m.Called(ctx, projectID, page, size)
	return args.Get(0).(model.Page[model.Board]), args.Error(1)
}

func (m *MockBoardRepository) UnsetAllDefaultsByProjectID(ctx context.Context, projectID uuid.UUID) (bool, error) {
	args := m.Called(ctx, projectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockBoardRepository) DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, projectID, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockBoardRepository) ExistsByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error) {
	args := m.Called(ctx, id, projectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockBoardRepository) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, ids, projectID)
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

type MockBoardColumnRepository struct {
	mock.Mock
}

func (m *MockBoardColumnRepository) Save(ctx context.Context, column *model.BoardColumn) (*model.BoardColumn, error) {
	args := m.Called(ctx, column)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.BoardColumn), args.Error(1)
}

func (m *MockBoardColumnRepository) FindByIDAndBoardID(ctx context.Context, boardID, id uuid.UUID) (*model.BoardColumn, error) {
	args := m.Called(ctx, boardID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.BoardColumn), args.Error(1)
}

func (m *MockBoardColumnRepository) FindAllByBoardIDAndProjectID(ctx context.Context, projectID, boardID uuid.UUID) ([]model.BoardColumn, error) {
	args := m.Called(ctx, projectID, boardID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.BoardColumn), args.Error(1)
}

func (m *MockBoardColumnRepository) GetNextPosition(ctx context.Context, boardID uuid.UUID) (uint8, error) {
	args := m.Called(ctx, boardID)
	return uint8(args.Int(0)), args.Error(1)
}

func (m *MockBoardColumnRepository) Reorder(ctx context.Context, boardID uuid.UUID, columnsIDs []uuid.UUID) error {
	args := m.Called(ctx, boardID, columnsIDs)
	return args.Error(0)
}

func (m *MockBoardColumnRepository) DeleteByIDAndBoardID(ctx context.Context, boardID, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, boardID, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockBoardColumnRepository) ExistsByIDAndBoardIDAndProjectID(ctx context.Context, id, boardID, projectID uuid.UUID) (bool, error) {
	args := m.Called(ctx, id, boardID, projectID)
	return args.Bool(0), args.Error(1)
}

// --- CREATE TESTS ---

func TestBoardUseCase_Create_Success_ExplicitDefault(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	desc := "Main sprint board"
	cmd := dto.CreateBoardCommand{
		ProjectID:   projectID,
		Name:        "Sprint 1",
		Description: &desc,
		IsDefault:   true,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("UnsetAllDefaultsByProjectID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("Save", ctx, mock.MatchedBy(func(b *model.Board) bool {
		return b.ProjectID == projectID && b.Name == "Sprint 1" && b.IsDefault == true
	})).Return(&model.Board{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Name:        "Sprint 1",
		Description: &desc,
		IsDefault:   true,
		BaseEntity:  model.BaseEntity{CreatedAt: time.Now()},
	}, nil)

	// Expect 4 default columns to be saved: Backlog, In Progress (WIP 3), Review (WIP 3), Done
	mockBoardColumnRepo.On("Save", ctx, mock.MatchedBy(func(c *model.BoardColumn) bool {
		return c.Name == "Backlog" && c.Position == 0 && c.WipLimit == nil
	})).Return(&model.BoardColumn{}, nil)
	mockBoardColumnRepo.On("Save", ctx, mock.MatchedBy(func(c *model.BoardColumn) bool {
		return c.Name == "In Progress" && c.Position == 1 && c.WipLimit != nil && *c.WipLimit == 3
	})).Return(&model.BoardColumn{}, nil)
	mockBoardColumnRepo.On("Save", ctx, mock.MatchedBy(func(c *model.BoardColumn) bool {
		return c.Name == "Review" && c.Position == 2 && c.WipLimit != nil && *c.WipLimit == 3
	})).Return(&model.BoardColumn{}, nil)
	mockBoardColumnRepo.On("Save", ctx, mock.MatchedBy(func(c *model.BoardColumn) bool {
		return c.Name == "Done" && c.Position == 3 && c.WipLimit == nil
	})).Return(&model.BoardColumn{}, nil)

	created, err := uc.Create(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Sprint 1", created.Name)
	assert.True(t, created.IsDefault)
	assert.Equal(t, projectID, created.ProjectID)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardUseCase_Create_Success_AutoDefault_WhenNoDefaultExists(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	cmd := dto.CreateBoardCommand{
		ProjectID: projectID,
		Name:      "First Board",
		IsDefault: false,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindDefaultByProjectID", ctx, projectID).Return(nil, nil)
	mockBoardRepo.On("Save", ctx, mock.MatchedBy(func(b *model.Board) bool {
		return b.ProjectID == projectID && b.Name == "First Board" && b.IsDefault == true
	})).Return(&model.Board{
		ID:         uuid.New(),
		ProjectID:  projectID,
		Name:       "First Board",
		IsDefault:  true,
		BaseEntity: model.BaseEntity{CreatedAt: time.Now()},
	}, nil)

	mockBoardColumnRepo.On("Save", ctx, mock.Anything).Return(&model.BoardColumn{}, nil).Times(4)

	created, err := uc.Create(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.True(t, created.IsDefault)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardUseCase_Create_Success_KeepNonDefault_WhenDefaultAlreadyExists(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	cmd := dto.CreateBoardCommand{
		ProjectID: projectID,
		Name:      "Second Board",
		IsDefault: false,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindDefaultByProjectID", ctx, projectID).Return(&model.Board{ID: uuid.New(), IsDefault: true}, nil)
	mockBoardRepo.On("Save", ctx, mock.MatchedBy(func(b *model.Board) bool {
		return b.ProjectID == projectID && b.Name == "Second Board" && b.IsDefault == false
	})).Return(&model.Board{
		ID:         uuid.New(),
		ProjectID:  projectID,
		Name:       "Second Board",
		IsDefault:  false,
		BaseEntity: model.BaseEntity{CreatedAt: time.Now()},
	}, nil)

	mockBoardColumnRepo.On("Save", ctx, mock.Anything).Return(&model.BoardColumn{}, nil).Times(4)

	created, err := uc.Create(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.False(t, created.IsDefault)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardUseCase_Create_ProjectNotFound(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	cmd := dto.CreateBoardCommand{
		ProjectID: projectID,
		Name:      "New Board",
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestBoardUseCase_Create_ProjectRepoError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	cmd := dto.CreateBoardCommand{
		ProjectID: projectID,
		Name:      "New Board",
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, errors.New("db connection failed"))

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorContains(t, err, "db connection failed")
	mockProjectRepo.AssertExpectations(t)
}

func TestBoardUseCase_Create_UnsetDefaultError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	cmd := dto.CreateBoardCommand{
		ProjectID: projectID,
		Name:      "New Board",
		IsDefault: true,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("UnsetAllDefaultsByProjectID", ctx, projectID).Return(false, errors.New("failed to unset"))

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorContains(t, err, "failed to unset")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_Create_UnsetDefaultFailed(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	cmd := dto.CreateBoardCommand{
		ProjectID: projectID,
		Name:      "New Board",
		IsDefault: true,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("UnsetAllDefaultsByProjectID", ctx, projectID).Return(false, nil)

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorIs(t, err, usecase.ErrNotPossibleToUnsetDefaultBoard)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_Create_FindDefaultError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	cmd := dto.CreateBoardCommand{
		ProjectID: projectID,
		Name:      "New Board",
		IsDefault: false,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindDefaultByProjectID", ctx, projectID).Return(nil, errors.New("find default error"))

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorContains(t, err, "find default error")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_Create_SaveError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	cmd := dto.CreateBoardCommand{
		ProjectID: projectID,
		Name:      "New Board",
		IsDefault: true,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("UnsetAllDefaultsByProjectID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("Save", ctx, mock.Anything).Return(nil, errors.New("save board failed"))

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorContains(t, err, "save board failed")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_Create_CreateDefaultColumnsError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	cmd := dto.CreateBoardCommand{
		ProjectID: projectID,
		Name:      "New Board",
		IsDefault: true,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("UnsetAllDefaultsByProjectID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("Save", ctx, mock.Anything).Return(&model.Board{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      "New Board",
		IsDefault: true,
	}, nil)
	mockBoardColumnRepo.On("Save", ctx, mock.Anything).Return(nil, errors.New("column insert failed"))

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorContains(t, err, "error creating default board columns")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

// --- GET BY ID TESTS ---

func TestBoardUseCase_GetByID_Success(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	expectedBoard := &model.Board{
		ID:        boardID,
		ProjectID: projectID,
		Name:      "Kanban Board",
		IsDefault: true,
	}
	color := "#3B82F6"
	expectedTags := []model.Tag{
		{ID: uuid.New(), Name: "Sprint", Color: &color},
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindByIDAndProjectID", ctx, projectID, boardID).Return(expectedBoard, nil)
	mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeBoard, projectID, boardID).Return(expectedTags, nil)

	result, err := uc.GetByID(ctx, projectID, boardID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, boardID, result.ID)
	assert.Equal(t, projectID, result.ProjectID)
	assert.Equal(t, "Kanban Board", result.Name)
	assert.Equal(t, 1, len(result.Tags))
	assert.Equal(t, "Sprint", result.Tags[0].Name)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

func TestBoardUseCase_GetByID_ProjectNotFound(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	result, err := uc.GetByID(ctx, projectID, boardID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestBoardUseCase_GetByID_NotFound(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindByIDAndProjectID", ctx, projectID, boardID).Return(nil, nil)

	result, err := uc.GetByID(ctx, projectID, boardID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrBoardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_GetByID_RepoError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindByIDAndProjectID", ctx, projectID, boardID).Return(nil, errors.New("db error"))

	result, err := uc.GetByID(ctx, projectID, boardID)

	assert.Nil(t, result)
	assert.ErrorContains(t, err, "error fetching board")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_GetByID_TagsError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	expectedBoard := &model.Board{ID: boardID, ProjectID: projectID, Name: "Kanban"}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindByIDAndProjectID", ctx, projectID, boardID).Return(expectedBoard, nil)
	mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeBoard, projectID, boardID).Return(nil, errors.New("tag fetch error"))

	result, err := uc.GetByID(ctx, projectID, boardID)

	assert.Nil(t, result)
	assert.ErrorContains(t, err, "error fetching tags for board")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

// --- GET DEFAULT BY PROJECT ID TESTS ---

func TestBoardUseCase_GetDefaultByProjectID_Success(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	expectedBoard := &model.Board{
		ID:        boardID,
		ProjectID: projectID,
		Name:      "Default Board",
		IsDefault: true,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindDefaultByProjectID", ctx, projectID).Return(expectedBoard, nil)
	mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeBoard, projectID, boardID).Return([]model.Tag{}, nil)

	result, err := uc.GetDefaultByProjectID(ctx, projectID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, boardID, result.ID)
	assert.True(t, result.IsDefault)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

func TestBoardUseCase_GetDefaultByProjectID_ProjectNotFound(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	result, err := uc.GetDefaultByProjectID(ctx, projectID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestBoardUseCase_GetDefaultByProjectID_NotFound(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindDefaultByProjectID", ctx, projectID).Return(nil, nil)

	result, err := uc.GetDefaultByProjectID(ctx, projectID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrBoardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_GetDefaultByProjectID_RepoError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindDefaultByProjectID", ctx, projectID).Return(nil, errors.New("db error"))

	result, err := uc.GetDefaultByProjectID(ctx, projectID)

	assert.Nil(t, result)
	assert.ErrorContains(t, err, "db error")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_GetDefaultByProjectID_TagsError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	expectedBoard := &model.Board{ID: boardID, ProjectID: projectID, Name: "Default"}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindDefaultByProjectID", ctx, projectID).Return(expectedBoard, nil)
	mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeBoard, projectID, boardID).Return(nil, errors.New("tags error"))

	result, err := uc.GetDefaultByProjectID(ctx, projectID)

	assert.Nil(t, result)
	assert.ErrorContains(t, err, "error fetching tags for board")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

// --- GET ALL BY PROJECT ID TESTS ---

func TestBoardUseCase_GetAllByProjectID_Success_WithContent(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	board1ID := uuid.New()
	board2ID := uuid.New()

	expectedPage := model.NewPage([]model.Board{
		{ID: board1ID, ProjectID: projectID, Name: "Board 1"},
		{ID: board2ID, ProjectID: projectID, Name: "Board 2"},
	}, 0, 10, 2)

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindAllByProjectID", ctx, projectID, 0, 10).Return(expectedPage, nil)
	mockItemTagRepo.On("FindTagsForItems", ctx, model.ItemTypeBoard, projectID, []uuid.UUID{board1ID, board2ID}).Return(map[uuid.UUID][]model.Tag{}, nil)

	result, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result.Content))
	assert.Equal(t, int64(2), result.TotalElements)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

func TestBoardUseCase_GetAllByProjectID_Success_EmptyPage(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	emptyPage := model.NewPage([]model.Board{}, 0, 10, 0)

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindAllByProjectID", ctx, projectID, 0, 10).Return(emptyPage, nil)

	result, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(result.Content))
	assert.Equal(t, int64(0), result.TotalElements)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_GetAllByProjectID_ProjectNotFound(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	_, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestBoardUseCase_GetAllByProjectID_RepoError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindAllByProjectID", ctx, projectID, 0, 10).Return(model.Page[model.Board]{}, errors.New("db error"))

	_, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

	assert.ErrorContains(t, err, "db error")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_GetAllByProjectID_TagsError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	expectedPage := model.NewPage([]model.Board{
		{ID: boardID, ProjectID: projectID, Name: "Board 1"},
	}, 0, 10, 1)

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindAllByProjectID", ctx, projectID, 0, 10).Return(expectedPage, nil)
	mockItemTagRepo.On("FindTagsForItems", ctx, model.ItemTypeBoard, projectID, []uuid.UUID{boardID}).Return(nil, errors.New("tags error"))

	_, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

	assert.ErrorContains(t, err, "error fetching tags for boards")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

// --- UPDATE TESTS ---

func TestBoardUseCase_Update_Success_NameAndDescription(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	newName := "Updated Name"
	newDesc := "Updated Description"

	cmd := dto.UpdateBoardCommand{
		ID:          boardID,
		ProjectID:   projectID,
		Name:        &newName,
		Description: &newDesc,
	}

	existingBoard := &model.Board{
		ID:          boardID,
		ProjectID:   projectID,
		Name:        "Old Name",
		Description: nil,
		IsDefault:   false,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindByIDAndProjectID", ctx, projectID, boardID).Return(existingBoard, nil)
	mockBoardRepo.On("Save", ctx, mock.MatchedBy(func(b *model.Board) bool {
		return b.ID == boardID && b.Name == "Updated Name" && *b.Description == "Updated Description" && b.UpdatedAt != nil
	})).Return(&model.Board{
		ID:          boardID,
		ProjectID:   projectID,
		Name:        "Updated Name",
		Description: &newDesc,
		IsDefault:   false,
	}, nil)
	mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeBoard, projectID, boardID).Return([]model.Tag{}, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, &newDesc, updated.Description)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

func TestBoardUseCase_Update_Success_SetDefault_True(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	isDefault := true

	cmd := dto.UpdateBoardCommand{
		ID:        boardID,
		ProjectID: projectID,
		IsDefault: &isDefault,
	}

	existingBoard := &model.Board{
		ID:        boardID,
		ProjectID: projectID,
		Name:      "Board",
		IsDefault: false,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindByIDAndProjectID", ctx, projectID, boardID).Return(existingBoard, nil)
	mockBoardRepo.On("UnsetAllDefaultsByProjectID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("Save", ctx, mock.MatchedBy(func(b *model.Board) bool {
		return b.ID == boardID && b.IsDefault == true
	})).Return(&model.Board{
		ID:        boardID,
		ProjectID: projectID,
		Name:      "Board",
		IsDefault: true,
	}, nil)
	mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeBoard, projectID, boardID).Return([]model.Tag{}, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.True(t, updated.IsDefault)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

func TestBoardUseCase_Update_Success_SetDefault_False(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	isDefault := false

	cmd := dto.UpdateBoardCommand{
		ID:        boardID,
		ProjectID: projectID,
		IsDefault: &isDefault,
	}

	existingBoard := &model.Board{
		ID:        boardID,
		ProjectID: projectID,
		Name:      "Board",
		IsDefault: true,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindByIDAndProjectID", ctx, projectID, boardID).Return(existingBoard, nil)
	mockBoardRepo.On("Save", ctx, mock.MatchedBy(func(b *model.Board) bool {
		return b.ID == boardID && b.IsDefault == false
	})).Return(&model.Board{
		ID:        boardID,
		ProjectID: projectID,
		Name:      "Board",
		IsDefault: false,
	}, nil)
	mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeBoard, projectID, boardID).Return([]model.Tag{}, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.False(t, updated.IsDefault)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

func TestBoardUseCase_Update_ProjectNotFound(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	cmd := dto.UpdateBoardCommand{
		ID:        uuid.New(),
		ProjectID: projectID,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestBoardUseCase_Update_NotFound(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	cmd := dto.UpdateBoardCommand{
		ID:        boardID,
		ProjectID: projectID,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindByIDAndProjectID", ctx, projectID, boardID).Return(nil, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrBoardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_Update_FindError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	cmd := dto.UpdateBoardCommand{
		ID:        boardID,
		ProjectID: projectID,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindByIDAndProjectID", ctx, projectID, boardID).Return(nil, errors.New("db find error"))

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorContains(t, err, "error fetching board")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_Update_UnsetDefaultError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	isDefault := true
	cmd := dto.UpdateBoardCommand{
		ID:        boardID,
		ProjectID: projectID,
		IsDefault: &isDefault,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindByIDAndProjectID", ctx, projectID, boardID).Return(&model.Board{ID: boardID, ProjectID: projectID}, nil)
	mockBoardRepo.On("UnsetAllDefaultsByProjectID", ctx, projectID).Return(false, errors.New("unset error"))

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorContains(t, err, "unset error")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_Update_UnsetDefaultFailed(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	isDefault := true
	cmd := dto.UpdateBoardCommand{
		ID:        boardID,
		ProjectID: projectID,
		IsDefault: &isDefault,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindByIDAndProjectID", ctx, projectID, boardID).Return(&model.Board{ID: boardID, ProjectID: projectID}, nil)
	mockBoardRepo.On("UnsetAllDefaultsByProjectID", ctx, projectID).Return(false, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrNotPossibleToUnsetDefaultBoard)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_Update_SaveError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	newName := "Name"
	cmd := dto.UpdateBoardCommand{
		ID:        boardID,
		ProjectID: projectID,
		Name:      &newName,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindByIDAndProjectID", ctx, projectID, boardID).Return(&model.Board{ID: boardID, ProjectID: projectID}, nil)
	mockBoardRepo.On("Save", ctx, mock.Anything).Return(nil, errors.New("save error"))

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorContains(t, err, "save error")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_Update_TagsError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	newName := "Name"
	cmd := dto.UpdateBoardCommand{
		ID:        boardID,
		ProjectID: projectID,
		Name:      &newName,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("FindByIDAndProjectID", ctx, projectID, boardID).Return(&model.Board{ID: boardID, ProjectID: projectID}, nil)
	mockBoardRepo.On("Save", ctx, mock.Anything).Return(&model.Board{ID: boardID, ProjectID: projectID, Name: newName}, nil)
	mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeBoard, projectID, boardID).Return(nil, errors.New("tags error"))

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorContains(t, err, "error fetching tags for board")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

// --- DELETE TESTS ---

func TestBoardUseCase_Delete_Success(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("DeleteByIDAndProjectID", ctx, projectID, boardID).Return(true, nil)
	mockItemTagRepo.On("RemoveAllTagsFromItem", ctx, model.ItemTypeBoard, boardID).Return(nil)

	err := uc.Delete(ctx, projectID, boardID)

	assert.NoError(t, err)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

func TestBoardUseCase_Delete_ProjectNotFound(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	err := uc.Delete(ctx, projectID, boardID)

	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestBoardUseCase_Delete_NotFound(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("DeleteByIDAndProjectID", ctx, projectID, boardID).Return(false, nil)

	err := uc.Delete(ctx, projectID, boardID)

	assert.ErrorIs(t, err, usecase.ErrBoardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_Delete_RepoError(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("DeleteByIDAndProjectID", ctx, projectID, boardID).Return(false, errors.New("db delete error"))

	err := uc.Delete(ctx, projectID, boardID)

	assert.ErrorContains(t, err, "error deleting board")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardUseCase_Delete_TagsError_NonBlocking(t *testing.T) {
	mockBoardRepo := new(MockBoardRepository)
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewBoardUseCase(mockBoardRepo, mockBoardColumnRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("DeleteByIDAndProjectID", ctx, projectID, boardID).Return(true, nil)
	mockItemTagRepo.On("RemoveAllTagsFromItem", ctx, model.ItemTypeBoard, boardID).Return(errors.New("tag cleanup failed"))

	err := uc.Delete(ctx, projectID, boardID)

	assert.NoError(t, err) // Tag cleanup failure should not fail the deletion
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}
