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

type MockSnippetRepository struct {
	mock.Mock
}

func (m *MockSnippetRepository) Save(ctx context.Context, snippet *model.Snippet) (*model.Snippet, error) {
	args := m.Called(ctx, snippet)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Snippet), args.Error(1)
}

func (m *MockSnippetRepository) FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Snippet, error) {
	args := m.Called(ctx, projectID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Snippet), args.Error(1)
}

func (m *MockSnippetRepository) FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[model.Snippet], error) {
	args := m.Called(ctx, projectID, page, size)
	return args.Get(0).(model.Page[model.Snippet]), args.Error(1)
}

func (m *MockSnippetRepository) DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) error {
	args := m.Called(ctx, projectID, id)
	return args.Error(0)
}

func (m *MockSnippetRepository) ExistsByIDAndProjectID(ctx context.Context, id, projectID uuid.UUID) (bool, error) {
	args := m.Called(ctx, id, projectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSnippetRepository) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, ids, projectID)
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

// --- UNIT TESTS ---

func TestSnippetUseCase_Create_Success(t *testing.T) {
	mockSnippetRepo := new(MockSnippetRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewSnippetUseCase(mockSnippetRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	desc := "Snippet description"
	lang := model.SnippetLangGo

	cmd := usecase.CreateSnippetCommand{
		ProjectID:   projectID,
		Title:       "Print Hello World",
		Description: &desc,
		Content:     `fmt.Println("Hello World")`,
		Language:    lang,
		SnippetType: model.SnippetTypeCode,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockSnippetRepo.On("Save", ctx, mock.MatchedBy(func(s *model.Snippet) bool {
		return s.ProjectID == projectID && s.Title == "Print Hello World" && s.Content == `fmt.Println("Hello World")`
	})).Return(&model.Snippet{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Title:       "Print Hello World",
		Description: &desc,
		Content:     `fmt.Println("Hello World")`,
		Language:    &lang,
		SnippetType: model.SnippetTypeCode,
	}, nil)

	created, err := uc.Create(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Print Hello World", created.Title)
	mockProjectRepo.AssertExpectations(t)
	mockSnippetRepo.AssertExpectations(t)
}

func TestSnippetUseCase_Create_ProjectNotFound(t *testing.T) {
	mockSnippetRepo := new(MockSnippetRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewSnippetUseCase(mockSnippetRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	cmd := usecase.CreateSnippetCommand{
		ProjectID:   projectID,
		Title:       "Print Hello World",
		Content:     `fmt.Println("Hello World")`,
		Language:    model.SnippetLangGo,
		SnippetType: model.SnippetTypeCode,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestSnippetUseCase_GetByID_Success(t *testing.T) {
	mockSnippetRepo := new(MockSnippetRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewSnippetUseCase(mockSnippetRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	snippetID := uuid.New()
	expectedSnippet := &model.Snippet{
		ID:        snippetID,
		ProjectID: projectID,
		Title:     "Sample Snippet",
	}

	mockSnippetRepo.On("FindByIDAndProjectID", ctx, projectID, snippetID).Return(expectedSnippet, nil)

	result, err := uc.GetByID(ctx, projectID, snippetID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, snippetID, result.ID)
	assert.Equal(t, projectID, result.ProjectID)
	mockSnippetRepo.AssertExpectations(t)
}

func TestSnippetUseCase_GetByID_NotFound(t *testing.T) {
	mockSnippetRepo := new(MockSnippetRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewSnippetUseCase(mockSnippetRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	snippetID := uuid.New()

	mockSnippetRepo.On("FindByIDAndProjectID", ctx, projectID, snippetID).Return(nil, nil)

	result, err := uc.GetByID(ctx, projectID, snippetID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrSnippetNotFound)
	mockSnippetRepo.AssertExpectations(t)
}

func TestSnippetUseCase_GetAllByProjectID_Success(t *testing.T) {
	mockSnippetRepo := new(MockSnippetRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewSnippetUseCase(mockSnippetRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	expectedPage := model.NewPage([]model.Snippet{
		{ID: uuid.New(), ProjectID: projectID, Title: "Snippet 1"},
		{ID: uuid.New(), ProjectID: projectID, Title: "Snippet 2"},
	}, 0, 10, 2)

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockSnippetRepo.On("FindAllByProjectID", ctx, projectID, 0, 10).Return(expectedPage, nil)

	result, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result.Content))
	assert.Equal(t, int64(2), result.TotalElements)
	mockProjectRepo.AssertExpectations(t)
	mockSnippetRepo.AssertExpectations(t)
}

func TestSnippetUseCase_GetAllByProjectID_ProjectNotFound(t *testing.T) {
	mockSnippetRepo := new(MockSnippetRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewSnippetUseCase(mockSnippetRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	_, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestSnippetUseCase_Update_Success(t *testing.T) {
	mockSnippetRepo := new(MockSnippetRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewSnippetUseCase(mockSnippetRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	snippetID := uuid.New()

	existingSnippet := &model.Snippet{
		ID:        snippetID,
		ProjectID: projectID,
		Title:     "Old Title",
		Content:   "Old Content",
	}

	newTitle := "Updated Title"
	cmd := usecase.UpdateSnippetCommand{
		ProjectID: projectID,
		ID:        snippetID,
		Title:     &newTitle,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockSnippetRepo.On("FindByIDAndProjectID", ctx, projectID, snippetID).Return(existingSnippet, nil)
	mockSnippetRepo.On("Save", ctx, mock.MatchedBy(func(s *model.Snippet) bool {
		return s.Title == "Updated Title" && s.UpdatedAt != nil
	})).Return(&model.Snippet{
		ID:        snippetID,
		ProjectID: projectID,
		Title:     "Updated Title",
	}, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Title", updated.Title)
	mockProjectRepo.AssertExpectations(t)
	mockSnippetRepo.AssertExpectations(t)
}

func TestSnippetUseCase_Update_ProjectNotFound(t *testing.T) {
	mockSnippetRepo := new(MockSnippetRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewSnippetUseCase(mockSnippetRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	snippetID := uuid.New()
	newTitle := "Updated Title"
	cmd := usecase.UpdateSnippetCommand{
		ProjectID: projectID,
		ID:        snippetID,
		Title:     &newTitle,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestSnippetUseCase_Update_SnippetNotFound(t *testing.T) {
	mockSnippetRepo := new(MockSnippetRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewSnippetUseCase(mockSnippetRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	snippetID := uuid.New()
	newTitle := "Updated Title"
	cmd := usecase.UpdateSnippetCommand{
		ProjectID: projectID,
		ID:        snippetID,
		Title:     &newTitle,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockSnippetRepo.On("FindByIDAndProjectID", ctx, projectID, snippetID).Return(nil, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrSnippetNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockSnippetRepo.AssertExpectations(t)
}

func TestSnippetUseCase_Delete_Success(t *testing.T) {
	mockSnippetRepo := new(MockSnippetRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewSnippetUseCase(mockSnippetRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	snippetID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockSnippetRepo.On("DeleteByIDAndProjectID", ctx, projectID, snippetID).Return(nil)

	err := uc.Delete(ctx, projectID, snippetID)

	assert.NoError(t, err)
	mockProjectRepo.AssertExpectations(t)
	mockSnippetRepo.AssertExpectations(t)
}

func TestSnippetUseCase_Delete_ProjectNotFound(t *testing.T) {
	mockSnippetRepo := new(MockSnippetRepository)
	mockProjectRepo := new(MockProjectRepository)
	uc := usecase.NewSnippetUseCase(mockSnippetRepo, mockProjectRepo)
	ctx := context.Background()

	projectID := uuid.New()
	snippetID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	err := uc.Delete(ctx, projectID, snippetID)

	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}
