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

type MockLinkRepository struct {
	mock.Mock
}

func (m *MockLinkRepository) Save(ctx context.Context, link *model.Link) (*model.Link, error) {
	args := m.Called(ctx, link)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Link), args.Error(1)
}

func (m *MockLinkRepository) FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Link, error) {
	args := m.Called(ctx, projectID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Link), args.Error(1)
}

func (m *MockLinkRepository) FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[model.Link], error) {
	args := m.Called(ctx, projectID, page, size)
	return args.Get(0).(model.Page[model.Link]), args.Error(1)
}

func (m *MockLinkRepository) DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, projectID, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockLinkRepository) ExistsByIDAndProjectID(ctx context.Context, id, projectID uuid.UUID) (bool, error) {
	args := m.Called(ctx, id, projectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockLinkRepository) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, ids, projectID)
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

// --- UNIT TESTS ---

func TestLinkUseCase_Create_Success(t *testing.T) {
	mockLinkRepo := new(MockLinkRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewLinkUseCase(mockLinkRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	desc := "Link description"

	cmd := usecase.CreateLinkCommand{
		ProjectID:   projectID,
		Title:       "Go Documentation",
		URL:         "https://go.dev/doc",
		Description: &desc,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockLinkRepo.On("Save", ctx, mock.MatchedBy(func(l *model.Link) bool {
		return l.ProjectID == projectID && l.Title == "Go Documentation" && l.Url == "https://go.dev/doc"
	})).Return(&model.Link{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Title:       "Go Documentation",
		Url:         "https://go.dev/doc",
		Description: &desc,
	}, nil)

	created, err := uc.Create(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Go Documentation", created.Title)
	mockProjectRepo.AssertExpectations(t)
	mockLinkRepo.AssertExpectations(t)
}

func TestLinkUseCase_Create_ProjectNotFound(t *testing.T) {
	mockLinkRepo := new(MockLinkRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewLinkUseCase(mockLinkRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	cmd := usecase.CreateLinkCommand{
		ProjectID: projectID,
		Title:     "Go Documentation",
		URL:       "https://go.dev/doc",
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	created, err := uc.Create(ctx, cmd)

	assert.Nil(t, created)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestLinkUseCase_GetByID_Success(t *testing.T) {
	mockLinkRepo := new(MockLinkRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewLinkUseCase(mockLinkRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	linkID := uuid.New()
	expectedLink := &model.Link{
		ID:        linkID,
		ProjectID: projectID,
		Title:     "Go Documentation",
		Url:       "https://go.dev/doc",
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockLinkRepo.On("FindByIDAndProjectID", ctx, projectID, linkID).Return(expectedLink, nil)

	result, err := uc.GetByID(ctx, projectID, linkID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, linkID, result.ID)
	assert.Equal(t, projectID, result.ProjectID)
	mockProjectRepo.AssertExpectations(t)
	mockLinkRepo.AssertExpectations(t)
}

func TestLinkUseCase_GetByID_ProjectNotFound(t *testing.T) {
	mockLinkRepo := new(MockLinkRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewLinkUseCase(mockLinkRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	linkID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	result, err := uc.GetByID(ctx, projectID, linkID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestLinkUseCase_GetByID_NotFound(t *testing.T) {
	mockLinkRepo := new(MockLinkRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewLinkUseCase(mockLinkRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	linkID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockLinkRepo.On("FindByIDAndProjectID", ctx, projectID, linkID).Return(nil, nil)

	result, err := uc.GetByID(ctx, projectID, linkID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, usecase.ErrLinkNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockLinkRepo.AssertExpectations(t)
}

func TestLinkUseCase_GetAllByProjectID_Success(t *testing.T) {
	mockLinkRepo := new(MockLinkRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewLinkUseCase(mockLinkRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	expectedPage := model.NewPage([]model.Link{
		{ID: uuid.New(), ProjectID: projectID, Title: "Link 1", Url: "https://link1.com"},
		{ID: uuid.New(), ProjectID: projectID, Title: "Link 2", Url: "https://link2.com"},
	}, 0, 10, 2)

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockLinkRepo.On("FindAllByProjectID", ctx, projectID, 0, 10).Return(expectedPage, nil)

	result, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result.Content))
	assert.Equal(t, int64(2), result.TotalElements)
	mockProjectRepo.AssertExpectations(t)
	mockLinkRepo.AssertExpectations(t)
}

func TestLinkUseCase_GetAllByProjectID_ProjectNotFound(t *testing.T) {
	mockLinkRepo := new(MockLinkRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewLinkUseCase(mockLinkRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	_, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestLinkUseCase_Update_Success(t *testing.T) {
	mockLinkRepo := new(MockLinkRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewLinkUseCase(mockLinkRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	linkID := uuid.New()

	existingLink := &model.Link{
		ID:        linkID,
		ProjectID: projectID,
		Title:     "Old Title",
		Url:       "https://old-url.com",
	}

	newTitle := "Updated Title"
	cmd := usecase.UpdateLinkCommand{
		ProjectID: projectID,
		ID:        linkID,
		Title:     &newTitle,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockLinkRepo.On("FindByIDAndProjectID", ctx, projectID, linkID).Return(existingLink, nil)
	mockLinkRepo.On("Save", ctx, mock.MatchedBy(func(l *model.Link) bool {
		return l.Title == "Updated Title" && l.UpdatedAt != nil
	})).Return(&model.Link{
		ID:        linkID,
		ProjectID: projectID,
		Title:     "Updated Title",
		Url:       "https://old-url.com",
	}, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Title", updated.Title)
	mockProjectRepo.AssertExpectations(t)
	mockLinkRepo.AssertExpectations(t)
}

func TestLinkUseCase_Update_ProjectNotFound(t *testing.T) {
	mockLinkRepo := new(MockLinkRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewLinkUseCase(mockLinkRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	linkID := uuid.New()
	newTitle := "Updated Title"
	cmd := usecase.UpdateLinkCommand{
		ProjectID: projectID,
		ID:        linkID,
		Title:     &newTitle,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestLinkUseCase_Update_LinkNotFound(t *testing.T) {
	mockLinkRepo := new(MockLinkRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewLinkUseCase(mockLinkRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	linkID := uuid.New()
	newTitle := "Updated Title"
	cmd := usecase.UpdateLinkCommand{
		ProjectID: projectID,
		ID:        linkID,
		Title:     &newTitle,
	}

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockLinkRepo.On("FindByIDAndProjectID", ctx, projectID, linkID).Return(nil, nil)

	updated, err := uc.Update(ctx, cmd)

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, usecase.ErrLinkNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockLinkRepo.AssertExpectations(t)
}

func TestLinkUseCase_Delete_Success(t *testing.T) {
	mockLinkRepo := new(MockLinkRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewLinkUseCase(mockLinkRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	linkID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockLinkRepo.On("DeleteByIDAndProjectID", ctx, projectID, linkID).Return(true, nil)
	mockItemTagRepo.On("RemoveAllTagsFromItem", ctx, model.ItemTypeLink, linkID).Return(nil)

	err := uc.Delete(ctx, projectID, linkID)

	assert.NoError(t, err)
	mockProjectRepo.AssertExpectations(t)
	mockLinkRepo.AssertExpectations(t)
	mockItemTagRepo.AssertExpectations(t)
}

func TestLinkUseCase_Delete_ProjectNotFound(t *testing.T) {
	mockLinkRepo := new(MockLinkRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewLinkUseCase(mockLinkRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	linkID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

	err := uc.Delete(ctx, projectID, linkID)

	assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	mockProjectRepo.AssertExpectations(t)
}

func TestLinkUseCase_Delete_LinkNotFound(t *testing.T) {
	mockLinkRepo := new(MockLinkRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockItemTagRepo := new(MockItemTagRepository)
	uc := usecase.NewLinkUseCase(mockLinkRepo, mockProjectRepo, mockItemTagRepo)
	ctx := context.Background()

	projectID := uuid.New()
	linkID := uuid.New()

	mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
	mockLinkRepo.On("DeleteByIDAndProjectID", ctx, projectID, linkID).Return(false, nil)

	err := uc.Delete(ctx, projectID, linkID)

	assert.ErrorIs(t, err, usecase.ErrLinkNotFound)
	mockProjectRepo.AssertExpectations(t)
	mockLinkRepo.AssertExpectations(t)
}
