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

type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Save(ctx context.Context, project *model.Project) (*model.Project, error) {
	args := m.Called(ctx, project)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockProjectRepository) FindAll(ctx context.Context, page, size int) (model.Page[model.Project], error) {
	args := m.Called(ctx, page, size)
	return args.Get(0).(model.Page[model.Project]), args.Error(1)
}

func (m *MockProjectRepository) DeleteByID(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// --- UNIT TESTS ---

func TestProjectUseCase_Create_Success(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	uc := usecase.NewProjectUseCase(mockRepo)
	ctx := context.Background()

	desc := "Project description"
	cmd := usecase.CreateProjectCommand{
		Name:        "Devaulty Go",
		Description: &desc,
	}

	mockRepo.On("Save", ctx, mock.MatchedBy(func(p *model.Project) bool {
		return p.Name == "Devaulty Go" && p.Archived == false
	})).Return(&model.Project{
		ID:          uuid.New(),
		Name:        "Devaulty Go",
		Description: &desc,
		Archived:    false,
	}, nil)

	created, err := uc.Create(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Devaulty Go", created.Name)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_GetByID_Success(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	uc := usecase.NewProjectUseCase(mockRepo)
	ctx := context.Background()
	projectID := uuid.New()

	expectedProject := &model.Project{
		ID:   projectID,
		Name: "Existing Project",
	}

	mockRepo.On("FindByID", ctx, projectID).Return(expectedProject, nil)

	result, err := uc.GetByID(ctx, projectID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, projectID, result.ID)
	assert.Equal(t, "Existing Project", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	uc := usecase.NewProjectUseCase(mockRepo)
	ctx := context.Background()
	projectID := uuid.New()

	mockRepo.On("FindByID", ctx, projectID).Return(nil, nil)

	result, err := uc.GetByID(ctx, projectID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_GetAll_Success(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	uc := usecase.NewProjectUseCase(mockRepo)
	ctx := context.Background()

	expectedPage := model.NewPage([]model.Project{
		{ID: uuid.New(), Name: "Project 1"},
		{ID: uuid.New(), Name: "Project 2"},
	}, 0, 10, 2)

	mockRepo.On("FindAll", ctx, 0, 10).Return(expectedPage, nil)

	result, err := uc.GetAll(ctx, 0, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result.Content))
	assert.Equal(t, int64(2), result.TotalElements)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_Update_Success(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	uc := usecase.NewProjectUseCase(mockRepo)
	ctx := context.Background()
	projectID := uuid.New()

	existingProject := &model.Project{
		ID:   projectID,
		Name: "Old Name",
	}

	newName := "Updated Name"
	cmd := usecase.UpdateProjectCommand{
		ID:   projectID,
		Name: &newName,
	}

	mockRepo.On("FindByID", ctx, projectID).Return(existingProject, nil)
	mockRepo.On("Save", ctx, mock.MatchedBy(func(p *model.Project) bool {
		return p.Name == "Updated Name" && p.UpdatedAt != nil
	})).Return(&model.Project{
		ID:   projectID,
		Name: "Updated Name",
	}, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Name", updated.Name)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_Update_NotFound(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	uc := usecase.NewProjectUseCase(mockRepo)
	ctx := context.Background()
	projectID := uuid.New()

	newName := "Updated Name"
	cmd := usecase.UpdateProjectCommand{
		ID:   projectID,
		Name: &newName,
	}

	mockRepo.On("FindByID", ctx, projectID).Return(nil, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_Delete_Success(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	uc := usecase.NewProjectUseCase(mockRepo)
	ctx := context.Background()
	projectID := uuid.New()

	mockRepo.On("DeleteByID", ctx, projectID).Return(true, nil)

	err := uc.Delete(ctx, projectID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_Delete_NotFound(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	uc := usecase.NewProjectUseCase(mockRepo)
	ctx := context.Background()
	projectID := uuid.New()

	mockRepo.On("DeleteByID", ctx, projectID).Return(false, nil)

	err := uc.Delete(ctx, projectID)

	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_Archive_Success(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	uc := usecase.NewProjectUseCase(mockRepo)
	ctx := context.Background()
	projectID := uuid.New()

	existingProject := &model.Project{
		ID:       projectID,
		Name:     "Project To Archive",
		Archived: false,
	}

	mockRepo.On("FindByID", ctx, projectID).Return(existingProject, nil)
	mockRepo.On("Save", ctx, mock.MatchedBy(func(p *model.Project) bool {
		return p.Archived == true
	})).Return(existingProject, nil)

	err := uc.Archive(ctx, projectID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_Archive_AlreadyArchived(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	uc := usecase.NewProjectUseCase(mockRepo)
	ctx := context.Background()
	projectID := uuid.New()

	existingProject := &model.Project{
		ID:       projectID,
		Name:     "Already Archived Project",
		Archived: true,
	}

	mockRepo.On("FindByID", ctx, projectID).Return(existingProject, nil)

	err := uc.Archive(ctx, projectID)

	assert.ErrorIs(t, err, usecase.ErrProjectAlreadyArchived)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_Unarchive_Success(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	uc := usecase.NewProjectUseCase(mockRepo)
	ctx := context.Background()
	projectID := uuid.New()

	existingProject := &model.Project{
		ID:       projectID,
		Name:     "Archived Project",
		Archived: true,
	}

	mockRepo.On("FindByID", ctx, projectID).Return(existingProject, nil)
	mockRepo.On("Save", ctx, mock.MatchedBy(func(p *model.Project) bool {
		return p.Archived == false
	})).Return(existingProject, nil)

	err := uc.Unarchive(ctx, projectID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProjectUseCase_Unarchive_NotArchived(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	uc := usecase.NewProjectUseCase(mockRepo)
	ctx := context.Background()
	projectID := uuid.New()

	existingProject := &model.Project{
		ID:       projectID,
		Name:     "Active Project",
		Archived: false,
	}

	mockRepo.On("FindByID", ctx, projectID).Return(existingProject, nil)

	err := uc.Unarchive(ctx, projectID)

	assert.ErrorIs(t, err, usecase.ErrProjectAlreadyUnarchived)
	mockRepo.AssertExpectations(t)
}
