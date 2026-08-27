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

// --- CREATE TESTS ---

func TestBoardColumnUseCase_Create_Success_WithWipLimit(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	wipLimit := uint16(5)

	cmd := dto.CreateBoardColumnCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		Name:      "In Progress",
		WipLimit:  &wipLimit,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("GetNextPosition", ctx, boardID).Return(2, nil)
	mockBoardColumnRepo.On("Save", ctx, mock.MatchedBy(func(c *model.BoardColumn) bool {
		return c.BoardID == boardID && c.Name == "In Progress" && c.Position == 2 && c.WipLimit != nil && *c.WipLimit == 5
	})).Return(&model.BoardColumn{
		ID:         uuid.New(),
		BoardID:    boardID,
		Name:       "In Progress",
		Position:   2,
		WipLimit:   &wipLimit,
		BaseEntity: model.BaseEntity{CreatedAt: time.Now()},
	}, nil)

	created, err := uc.Create(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "In Progress", created.Name)
	assert.Equal(t, uint8(2), created.Position)
	assert.Equal(t, &wipLimit, created.WipLimit)
	assert.Equal(t, projectID, created.ProjectID)
	assert.Equal(t, boardID, created.BoardID)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Create_Success_WithoutWipLimit(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	cmd := dto.CreateBoardColumnCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		Name:      "Backlog",
		WipLimit:  nil,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("GetNextPosition", ctx, boardID).Return(0, nil)
	mockBoardColumnRepo.On("Save", ctx, mock.MatchedBy(func(c *model.BoardColumn) bool {
		return c.BoardID == boardID && c.Name == "Backlog" && c.Position == 0 && c.WipLimit == nil
	})).Return(&model.BoardColumn{
		ID:         uuid.New(),
		BoardID:    boardID,
		Name:       "Backlog",
		Position:   0,
		WipLimit:   nil,
		BaseEntity: model.BaseEntity{CreatedAt: time.Now()},
	}, nil)

	created, err := uc.Create(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Backlog", created.Name)
	assert.Equal(t, uint8(0), created.Position)
	assert.Nil(t, created.WipLimit)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Create_ProjectNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	cmd := dto.CreateBoardColumnCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		Name:      "To Do",
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Create_BoardNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	cmd := dto.CreateBoardColumnCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		Name:      "To Do",
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(false, nil)

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorIs(t, err, usecase.ErrBoardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Create_BoardRepoError(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	cmd := dto.CreateBoardColumnCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		Name:      "To Do",
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(false, errors.New("db error"))

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorContains(t, err, "error checking if board exists")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Create_GetNextPositionError(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	cmd := dto.CreateBoardColumnCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		Name:      "To Do",
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("GetNextPosition", ctx, boardID).Return(0, errors.New("position error"))

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorIs(t, err, usecase.ErrNotPossibleToGetBoardColumnNextPosition)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Create_SaveError(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	cmd := dto.CreateBoardColumnCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		Name:      "To Do",
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("GetNextPosition", ctx, boardID).Return(0, nil)
	mockBoardColumnRepo.On("Save", ctx, mock.Anything).Return(nil, errors.New("save column error"))

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorContains(t, err, "save column error")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

// --- GET BY ID AND BOARD ID TESTS ---

func TestBoardColumnUseCase_GetByIDAndBoardID_Success(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()
	wipLimit := uint16(3)

	expectedColumn := &model.BoardColumn{
		ID:         columnID,
		BoardID:    boardID,
		Name:       "Review",
		Position:   1,
		WipLimit:   &wipLimit,
		BaseEntity: model.BaseEntity{CreatedAt: time.Now()},
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("FindByIDAndBoardID", ctx, boardID, columnID).Return(expectedColumn, nil)

	result, err := uc.GetByIDAndBoardID(ctx, projectID, boardID, columnID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, columnID, result.ID)
	assert.Equal(t, boardID, result.BoardID)
	assert.Equal(t, projectID, result.ProjectID)
	assert.Equal(t, "Review", result.Name)
	assert.Equal(t, uint8(1), result.Position)
	assert.Equal(t, &wipLimit, result.WipLimit)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_GetByIDAndBoardID_ProjectNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	result, err := uc.GetByIDAndBoardID(ctx, projectID, boardID, columnID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_GetByIDAndBoardID_BoardNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(false, nil)

	result, err := uc.GetByIDAndBoardID(ctx, projectID, boardID, columnID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrBoardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_GetByIDAndBoardID_ColumnNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("FindByIDAndBoardID", ctx, boardID, columnID).Return(nil, nil)

	result, err := uc.GetByIDAndBoardID(ctx, projectID, boardID, columnID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrBoardColumnNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_GetByIDAndBoardID_RepoError(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("FindByIDAndBoardID", ctx, boardID, columnID).Return(nil, errors.New("db find error"))

	result, err := uc.GetByIDAndBoardID(ctx, projectID, boardID, columnID)

	assert.Nil(t, result)
	assert.ErrorContains(t, err, "error fetching board column")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

// --- GET ALL BY BOARD ID TESTS ---

func TestBoardColumnUseCase_GetAllByBoardID_Success(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	col1ID := uuid.New()
	col2ID := uuid.New()

	expectedColumns := []model.BoardColumn{
		{ID: col1ID, BoardID: boardID, Name: "To Do", Position: 0},
		{ID: col2ID, BoardID: boardID, Name: "Done", Position: 1},
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("FindAllByBoardIDAndProjectID", ctx, projectID, boardID).Return(expectedColumns, nil)

	views, err := uc.GetAllByBoardID(ctx, projectID, boardID)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(views))
	assert.Equal(t, col1ID, views[0].ID)
	assert.Equal(t, "To Do", views[0].Name)
	assert.Equal(t, col2ID, views[1].ID)
	assert.Equal(t, "Done", views[1].Name)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_GetAllByBoardID_Success_EmptyList(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("FindAllByBoardIDAndProjectID", ctx, projectID, boardID).Return([]model.BoardColumn{}, nil)

	views, err := uc.GetAllByBoardID(ctx, projectID, boardID)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(views))
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_GetAllByBoardID_ProjectNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	views, err := uc.GetAllByBoardID(ctx, projectID, boardID)

	assert.Nil(t, views)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_GetAllByBoardID_BoardNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
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

func TestBoardColumnUseCase_GetAllByBoardID_RepoError(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("FindAllByBoardIDAndProjectID", ctx, projectID, boardID).Return(nil, errors.New("db list error"))

	views, err := uc.GetAllByBoardID(ctx, projectID, boardID)

	assert.Nil(t, views)
	assert.ErrorContains(t, err, "error fetching board columns")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

// --- UPDATE TESTS ---

func TestBoardColumnUseCase_Update_Success_NameAndWipLimit(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()
	newName := "In Progress (Updated)"
	newWip := uint16(10)

	cmd := dto.UpdateBoardColumnCommand{
		ID:        columnID,
		ProjectID: projectID,
		BoardID:   boardID,
		Name:      &newName,
		WipLimit:  &newWip,
	}

	existingColumn := &model.BoardColumn{
		ID:       columnID,
		BoardID:  boardID,
		Name:     "In Progress",
		Position: 1,
		WipLimit: nil,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("FindByIDAndBoardID", ctx, boardID, columnID).Return(existingColumn, nil)
	mockBoardColumnRepo.On("Save", ctx, mock.MatchedBy(func(c *model.BoardColumn) bool {
		return c.ID == columnID && c.Name == "In Progress (Updated)" && c.WipLimit != nil && *c.WipLimit == 10 && c.UpdatedAt != nil
	})).Return(&model.BoardColumn{
		ID:       columnID,
		BoardID:  boardID,
		Name:     "In Progress (Updated)",
		Position: 1,
		WipLimit: &newWip,
	}, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "In Progress (Updated)", updated.Name)
	assert.Equal(t, &newWip, updated.WipLimit)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Update_ProjectNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()

	cmd := dto.UpdateBoardColumnCommand{
		ID:        columnID,
		ProjectID: projectID,
		BoardID:   boardID,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Update_BoardNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()

	cmd := dto.UpdateBoardColumnCommand{
		ID:        columnID,
		ProjectID: projectID,
		BoardID:   boardID,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(false, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrBoardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Update_ColumnNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()

	cmd := dto.UpdateBoardColumnCommand{
		ID:        columnID,
		ProjectID: projectID,
		BoardID:   boardID,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("FindByIDAndBoardID", ctx, boardID, columnID).Return(nil, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrBoardColumnNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Update_FindError(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()

	cmd := dto.UpdateBoardColumnCommand{
		ID:        columnID,
		ProjectID: projectID,
		BoardID:   boardID,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("FindByIDAndBoardID", ctx, boardID, columnID).Return(nil, errors.New("db find error"))

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorContains(t, err, "error fetching board column")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Update_SaveError(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()
	newName := "New Name"

	cmd := dto.UpdateBoardColumnCommand{
		ID:        columnID,
		ProjectID: projectID,
		BoardID:   boardID,
		Name:      &newName,
	}

	existingColumn := &model.BoardColumn{
		ID:      columnID,
		BoardID: boardID,
		Name:    "Old Name",
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("FindByIDAndBoardID", ctx, boardID, columnID).Return(existingColumn, nil)
	mockBoardColumnRepo.On("Save", ctx, mock.Anything).Return(nil, errors.New("db save error"))

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorContains(t, err, "db save error")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

// --- DELETE TESTS ---

func TestBoardColumnUseCase_Delete_Success(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("DeleteByIDAndBoardID", ctx, boardID, columnID).Return(true, nil)

	err := uc.Delete(ctx, projectID, boardID, columnID)

	assert.NoError(t, err)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Delete_ProjectNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	err := uc.Delete(ctx, projectID, boardID, columnID)

	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Delete_BoardNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(false, nil)

	err := uc.Delete(ctx, projectID, boardID, columnID)

	assert.ErrorIs(t, err, usecase.ErrBoardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Delete_ColumnNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("DeleteByIDAndBoardID", ctx, boardID, columnID).Return(false, nil)

	err := uc.Delete(ctx, projectID, boardID, columnID)

	assert.ErrorIs(t, err, usecase.ErrBoardColumnNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Delete_RepoError(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	columnID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("DeleteByIDAndBoardID", ctx, boardID, columnID).Return(false, errors.New("db delete error"))

	err := uc.Delete(ctx, projectID, boardID, columnID)

	assert.ErrorContains(t, err, "error deleting board column")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

// --- REORDER TESTS ---

func TestBoardColumnUseCase_Reorder_Success(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	col1ID := uuid.New()
	col2ID := uuid.New()
	newOrder := []uuid.UUID{col2ID, col1ID}

	cmd := dto.ReorderBoardColumnsCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		Positions: newOrder,
	}

	expectedReorderedColumns := []model.BoardColumn{
		{ID: col2ID, BoardID: boardID, Name: "Done", Position: 0},
		{ID: col1ID, BoardID: boardID, Name: "To Do", Position: 1},
	}

	initialColumns := []model.BoardColumn{
		{ID: col1ID, BoardID: boardID, Name: "To Do", Position: 0},
		{ID: col2ID, BoardID: boardID, Name: "Done", Position: 1},
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("FindAllByBoardIDAndProjectID", ctx, projectID, boardID).Return(initialColumns, nil).Once()
	mockBoardColumnRepo.On("Reorder", ctx, boardID, newOrder).Return(nil)
	mockBoardColumnRepo.On("FindAllByBoardIDAndProjectID", ctx, projectID, boardID).Return(expectedReorderedColumns, nil).Once()

	views, err := uc.Reorder(ctx, cmd)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(views))
	assert.Equal(t, col2ID, views[0].ID)
	assert.Equal(t, uint8(0), views[0].Position)
	assert.Equal(t, col1ID, views[1].ID)
	assert.Equal(t, uint8(1), views[1].Position)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Reorder_ProjectNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	cmd := dto.ReorderBoardColumnsCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		Positions: []uuid.UUID{uuid.New()},
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	views, err := uc.Reorder(ctx, cmd)

	assert.Nil(t, views)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Reorder_BoardNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()

	cmd := dto.ReorderBoardColumnsCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		Positions: []uuid.UUID{uuid.New()},
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(false, nil)

	views, err := uc.Reorder(ctx, cmd)

	assert.Nil(t, views)
	assert.ErrorIs(t, err, usecase.ErrBoardNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Reorder_ColumnNotFound(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	existingColID := uuid.New()
	unknownColID := uuid.New()

	cmd := dto.ReorderBoardColumnsCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		Positions: []uuid.UUID{unknownColID},
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("FindAllByBoardIDAndProjectID", ctx, projectID, boardID).Return([]model.BoardColumn{
		{ID: existingColID, BoardID: boardID, Name: "Backlog"},
	}, nil)

	views, err := uc.Reorder(ctx, cmd)

	assert.Nil(t, views)
	assert.ErrorIs(t, err, usecase.ErrBoardColumnNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Reorder_CountMismatch(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	col1ID := uuid.New()
	col2ID := uuid.New()

	cmd := dto.ReorderBoardColumnsCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		Positions: []uuid.UUID{col1ID}, // only 1 provided when 2 exist
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("FindAllByBoardIDAndProjectID", ctx, projectID, boardID).Return([]model.BoardColumn{
		{ID: col1ID, BoardID: boardID, Name: "Col 1"},
		{ID: col2ID, BoardID: boardID, Name: "Col 2"},
	}, nil)

	views, err := uc.Reorder(ctx, cmd)

	assert.Nil(t, views)
	assert.ErrorContains(t, err, "column count mismatch")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Reorder_DuplicatePositions(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	col1ID := uuid.New()
	col2ID := uuid.New()

	cmd := dto.ReorderBoardColumnsCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		Positions: []uuid.UUID{col1ID, col1ID}, // duplicate [A, A] with same total count
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("FindAllByBoardIDAndProjectID", ctx, projectID, boardID).Return([]model.BoardColumn{
		{ID: col1ID, BoardID: boardID, Name: "Col 1"},
		{ID: col2ID, BoardID: boardID, Name: "Col 2"},
	}, nil)

	views, err := uc.Reorder(ctx, cmd)

	assert.Nil(t, views)
	assert.ErrorIs(t, err, usecase.ErrBoardColumnNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}

func TestBoardColumnUseCase_Reorder_RepoError(t *testing.T) {
	mockBoardColumnRepo := new(MockBoardColumnRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewBoardColumnUseCase(mockBoardColumnRepo, mockBoardRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	boardID := uuid.New()
	colID := uuid.New()
	positions := []uuid.UUID{colID}

	cmd := dto.ReorderBoardColumnsCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		Positions: positions,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
	mockBoardColumnRepo.On("FindAllByBoardIDAndProjectID", ctx, projectID, boardID).Return([]model.BoardColumn{
		{ID: colID, BoardID: boardID, Name: "Col 1"},
	}, nil)
	mockBoardColumnRepo.On("Reorder", ctx, boardID, positions).Return(errors.New("db reorder error"))

	views, err := uc.Reorder(ctx, cmd)

	assert.Nil(t, views)
	assert.ErrorContains(t, err, "error reordering board columns")
	mockProjectRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBoardColumnRepo.AssertExpectations(t)
}
