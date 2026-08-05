package usecase_test

import (
	"context"
	"testing"

	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProblemRepository struct {
	mock.Mock
}

func (m *MockProblemRepository) Save(ctx context.Context, problem *model.Problem) (*model.Problem, error) {
	args := m.Called(ctx, problem)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Problem), args.Error(1)
}

func (m *MockProblemRepository) FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Problem, error) {
	args := m.Called(ctx, projectID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Problem), args.Error(1)
}

func (m *MockProblemRepository) FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[model.Problem], error) {
	args := m.Called(ctx, projectID, page, size)
	return args.Get(0).(model.Page[model.Problem]), args.Error(1)
}

func (m *MockProblemRepository) DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, projectID, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockProblemRepository) ExistsByIDAndProjectID(ctx context.Context, id, projectID uuid.UUID) (bool, error) {
	args := m.Called(ctx, id, projectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockProblemRepository) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, ids, projectID)
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

// --- UNIT TESTS ---

func TestProblemUseCase_Create_Success(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	solution := "Restart database service"

	cmd := usecase.CreateProblemCommand{
		ProjectID:        projectID,
		Title:            "Connection Refused",
		ErrorDescription: "Failed to connect to database at localhost:5432",
		Solution:         &solution,
		Status:           model.ProblemStatusOpen,
		Severity:         model.ProblemSeverityHigh,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockProblemRepo.On("Save", ctx, mock.MatchedBy(func(p *model.Problem) bool {
		return p.ProjectID == projectID && p.Title == "Connection Refused" && p.Severity == model.ProblemSeverityHigh
	})).Return(&model.Problem{
		ID:               uuid.New(),
		ProjectID:        projectID,
		Title:            "Connection Refused",
		ErrorDescription: "Failed to connect to database at localhost:5432",
		Solution:         &solution,
		Status:           model.ProblemStatusOpen,
		Severity:         model.ProblemSeverityHigh,
	}, nil)

	created, err := uc.Create(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Connection Refused", created.Title)
	assert.Equal(t, model.ProblemSeverityHigh, created.Severity)
	mockProjectRepo.AssertExpectations(t)
	mockProblemRepo.AssertExpectations(t)
}

func TestProblemUseCase_Create_ProjectNotFound(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	cmd := usecase.CreateProblemCommand{
		ProjectID:        projectID,
		Title:            "Connection Refused",
		ErrorDescription: "Failed to connect",
		Status:           model.ProblemStatusOpen,
		Severity:         model.ProblemSeverityHigh,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestProblemUseCase_GetByID_Success(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	problemID := uuid.New()
	expectedProblem := &model.Problem{
		ID:               problemID,
		ProjectID:        projectID,
		Title:            "Memory Leak",
		ErrorDescription: "Process OOM killed",
		Status:           model.ProblemStatusWorkingOnIt,
		Severity:         model.ProblemSeverityCritical,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockProblemRepo.On("FindByIDAndProjectID", ctx, projectID, problemID).Return(expectedProblem, nil)

	result, err := uc.GetByID(ctx, projectID, problemID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, problemID, result.ID)
	assert.Equal(t, projectID, result.ProjectID)
	assert.Equal(t, model.ProblemSeverityCritical, result.Severity)
	mockProjectRepo.AssertExpectations(t)
	mockProblemRepo.AssertExpectations(t)
}

func TestProblemUseCase_GetByID_ProjectNotFound(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	problemID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	result, err := uc.GetByID(ctx, projectID, problemID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestProblemUseCase_GetByID_NotFound(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	problemID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockProblemRepo.On("FindByIDAndProjectID", ctx, projectID, problemID).Return(nil, nil)

	result, err := uc.GetByID(ctx, projectID, problemID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrProblemNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockProblemRepo.AssertExpectations(t)
}

func TestProblemUseCase_GetAllByProjectID_Success(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	expectedPage := model.NewPage([]model.Problem{
		{ID: uuid.New(), ProjectID: projectID, Title: "Problem 1", Severity: model.ProblemSeverityLow},
		{ID: uuid.New(), ProjectID: projectID, Title: "Problem 2", Severity: model.ProblemSeverityMedium},
	}, 0, 10, 2)

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockProblemRepo.On("FindAllByProjectID", ctx, projectID, 0, 10).Return(expectedPage, nil)

	result, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result.Content))
	assert.Equal(t, int64(2), result.TotalElements)
	mockProjectRepo.AssertExpectations(t)
	mockProblemRepo.AssertExpectations(t)
}

func TestProblemUseCase_GetAllByProjectID_ProjectNotFound(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	_, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestProblemUseCase_Update_Success(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	problemID := uuid.New()

	existingProblem := &model.Problem{
		ID:               problemID,
		ProjectID:        projectID,
		Title:            "Old Title",
		ErrorDescription: "Old Desc",
		Status:           model.ProblemStatusOpen,
		Severity:         model.ProblemSeverityLow,
	}

	newTitle := "Updated Title"
	newSeverity := model.ProblemSeverityCritical
	cmd := usecase.UpdateProblemCommand{
		ProjectID: projectID,
		ID:        problemID,
		Title:     &newTitle,
		Severity:  &newSeverity,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockProblemRepo.On("FindByIDAndProjectID", ctx, projectID, problemID).Return(existingProblem, nil)
	mockProblemRepo.On("Save", ctx, mock.MatchedBy(func(p *model.Problem) bool {
		return p.Title == "Updated Title" && p.Severity == model.ProblemSeverityCritical && p.UpdatedAt != nil
	})).Return(&model.Problem{
		ID:        problemID,
		ProjectID: projectID,
		Title:     "Updated Title",
		Severity:  model.ProblemSeverityCritical,
	}, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Title", updated.Title)
	assert.Equal(t, model.ProblemSeverityCritical, updated.Severity)
	mockProjectRepo.AssertExpectations(t)
	mockProblemRepo.AssertExpectations(t)
}

func TestProblemUseCase_Update_ProjectNotFound(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	problemID := uuid.New()
	newTitle := "Updated Title"
	cmd := usecase.UpdateProblemCommand{
		ProjectID: projectID,
		ID:        problemID,
		Title:     &newTitle,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestProblemUseCase_Update_ProblemNotFound(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	problemID := uuid.New()
	newTitle := "Updated Title"
	cmd := usecase.UpdateProblemCommand{
		ProjectID: projectID,
		ID:        problemID,
		Title:     &newTitle,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockProblemRepo.On("FindByIDAndProjectID", ctx, projectID, problemID).Return(nil, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrProblemNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockProblemRepo.AssertExpectations(t)
}

func TestProblemUseCase_UpdateStatus_Success(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	problemID := uuid.New()

	existingProblem := &model.Problem{
		ID:        problemID,
		ProjectID: projectID,
		Title:     "Sample Problem",
		Status:    model.ProblemStatusOpen,
	}

	cmd := usecase.UpdateProblemStatusCommand{
		ProjectID: projectID,
		ID:        problemID,
		Status:    model.ProblemStatusResolved,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockProblemRepo.On("FindByIDAndProjectID", ctx, projectID, problemID).Return(existingProblem, nil)
	mockProblemRepo.On("Save", ctx, mock.MatchedBy(func(p *model.Problem) bool {
		return p.Status == model.ProblemStatusResolved && p.UpdatedAt != nil
	})).Return(&model.Problem{
		ID:        problemID,
		ProjectID: projectID,
		Title:     "Sample Problem",
		Status:    model.ProblemStatusResolved,
	}, nil)

	updated, err := uc.UpdateStatus(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, model.ProblemStatusResolved, updated.Status)
	mockProjectRepo.AssertExpectations(t)
	mockProblemRepo.AssertExpectations(t)
}

func TestProblemUseCase_UpdateStatus_ProjectNotFound(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	problemID := uuid.New()

	cmd := usecase.UpdateProblemStatusCommand{
		ProjectID: projectID,
		ID:        problemID,
		Status:    model.ProblemStatusResolved,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	updated, err := uc.UpdateStatus(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestProblemUseCase_UpdateStatus_ProblemNotFound(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	problemID := uuid.New()

	cmd := usecase.UpdateProblemStatusCommand{
		ProjectID: projectID,
		ID:        problemID,
		Status:    model.ProblemStatusResolved,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockProblemRepo.On("FindByIDAndProjectID", ctx, projectID, problemID).Return(nil, nil)

	updated, err := uc.UpdateStatus(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrProblemNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockProblemRepo.AssertExpectations(t)
}

func TestProblemUseCase_Delete_Success(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	problemID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockProblemRepo.On("DeleteByIDAndProjectID", ctx, projectID, problemID).Return(true, nil)

	err := uc.Delete(ctx, projectID, problemID)

	assert.NoError(t, err)
	mockProjectRepo.AssertExpectations(t)
	mockProblemRepo.AssertExpectations(t)
}

func TestProblemUseCase_Delete_ProjectNotFound(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	problemID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	err := uc.Delete(ctx, projectID, problemID)

	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestProblemUseCase_Delete_ProblemNotFound(t *testing.T) {
	mockProblemRepo := new(MockProblemRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewProblemUseCase(mockProblemRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	problemID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockProblemRepo.On("DeleteByIDAndProjectID", ctx, projectID, problemID).Return(false, nil)

	err := uc.Delete(ctx, projectID, problemID)

	assert.ErrorIs(t, err, usecase.ErrProblemNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockProblemRepo.AssertExpectations(t)
}
