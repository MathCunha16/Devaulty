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

type MockDrawingRepository struct {
	mock.Mock
}

func (m *MockDrawingRepository) Save(ctx context.Context, drawing *model.Drawing) (*model.Drawing, error) {
	args := m.Called(ctx, drawing)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Drawing), args.Error(1)
}

func (m *MockDrawingRepository) FindByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*model.Drawing, error) {
	args := m.Called(ctx, id, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Drawing), args.Error(1)
}

func (m *MockDrawingRepository) FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[model.Drawing], error) {
	args := m.Called(ctx, projectID, page, size)
	if args.Get(0) == nil {
		return model.Page[model.Drawing]{}, args.Error(1)
	}
	return args.Get(0).(model.Page[model.Drawing]), args.Error(1)
}

func (m *MockDrawingRepository) DeleteByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error) {
	args := m.Called(ctx, id, projectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDrawingRepository) ExistsByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error) {
	args := m.Called(ctx, id, projectID)
	if args.Get(0) == nil {
		return false, args.Error(1)
	}
	return args.Bool(0), args.Error(1)
}

func (m *MockDrawingRepository) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, ids, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

// --- UNIT TESTS ---

func TestDrawingUseCase_Create(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()

	cmd := dto.CreateDrawingCommand{
		ProjectID: projectID,
		Title:     "System Architecture",
	}

	t.Run("Create_Success", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("Save", ctx, mock.MatchedBy(func(d *model.Drawing) bool {
			return d.Title == "System Architecture" && d.ProjectID == projectID && d.Content == nil && d.Thumbnail == nil
		})).Return(&model.Drawing{
			ID:        uuid.New(),
			ProjectID: projectID,
			Title:     "System Architecture",
			Content:   nil,
			Thumbnail: nil,
			BaseEntity: model.BaseEntity{
				CreatedAt: time.Now(),
			},
		}, nil)

		result, err := uc.Create(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "System Architecture", result.Title)
		assert.Nil(t, result.Content)
		assert.Nil(t, result.Thumbnail)
		assert.Empty(t, result.Tags)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
	})

	t.Run("Create_ProjectNotFound", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		result, err := uc.Create(ctx, cmd)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		assert.Nil(t, result)
		mockProjectRepo.AssertExpectations(t)
	})

	t.Run("Create_RepoError", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("Save", ctx, mock.Anything).Return(nil, errors.New("db save error"))

		result, err := uc.Create(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
	})
}

func TestDrawingUseCase_GetByID(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	drawingID := uuid.New()
	content := `{"elements":[]}`
	thumbnail := "data:image/svg+xml;base64,PHN2Zz4="

	t.Run("GetByID_Success", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		expectedDrawing := &model.Drawing{
			ID:        drawingID,
			ProjectID: projectID,
			Title:     "Architecture Diagram",
			Content:   &content,
			Thumbnail: &thumbnail,
			BaseEntity: model.BaseEntity{
				CreatedAt: time.Now(),
			},
		}

		tagID := uuid.New()
		tags := []model.Tag{
			{ID: tagID, Name: "Architecture"},
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("FindByIDAndProjectID", ctx, drawingID, projectID).Return(expectedDrawing, nil)
		mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeDrawing, projectID, drawingID).Return(tags, nil)

		result, err := uc.GetByID(ctx, projectID, drawingID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, drawingID, result.ID)
		assert.Equal(t, "Architecture Diagram", result.Title)
		assert.Equal(t, &content, result.Content)
		assert.Equal(t, &thumbnail, result.Thumbnail)
		assert.Len(t, result.Tags, 1)
		assert.Equal(t, "Architecture", result.Tags[0].Name)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
		mockItemTagRepo.AssertExpectations(t)
	})

	t.Run("GetByID_ProjectNotFound", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		result, err := uc.GetByID(ctx, projectID, drawingID)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		assert.Nil(t, result)
		mockProjectRepo.AssertExpectations(t)
	})

	t.Run("GetByID_DrawingNotFound", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("FindByIDAndProjectID", ctx, drawingID, projectID).Return(nil, nil)

		result, err := uc.GetByID(ctx, projectID, drawingID)

		assert.ErrorIs(t, err, usecase.ErrDrawingNotFound)
		assert.Nil(t, result)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
	})

	t.Run("GetByID_RepoError", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("FindByIDAndProjectID", ctx, drawingID, projectID).Return(nil, errors.New("db error"))

		result, err := uc.GetByID(ctx, projectID, drawingID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
	})

	t.Run("GetByID_TagsRepoError", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		expectedDrawing := &model.Drawing{
			ID:        drawingID,
			ProjectID: projectID,
			Title:     "Architecture Diagram",
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("FindByIDAndProjectID", ctx, drawingID, projectID).Return(expectedDrawing, nil)
		mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeDrawing, projectID, drawingID).Return(nil, errors.New("tags error"))

		result, err := uc.GetByID(ctx, projectID, drawingID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
		mockItemTagRepo.AssertExpectations(t)
	})
}

func TestDrawingUseCase_GetAllByProjectID(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	d1ID := uuid.New()
	d2ID := uuid.New()
	thumb := "data:image/svg+xml;base64,preview"

	t.Run("GetAllByProjectID_Success", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		expectedPage := model.NewPage([]model.Drawing{
			{ID: d1ID, ProjectID: projectID, Title: "Drawing 1", Thumbnail: &thumb},
			{ID: d2ID, ProjectID: projectID, Title: "Drawing 2"},
		}, 0, 10, 2)

		tagsMap := map[uuid.UUID][]model.Tag{
			d1ID: {{ID: uuid.New(), Name: "Draft"}},
			d2ID: {},
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("FindAllByProjectID", ctx, projectID, 0, 10).Return(expectedPage, nil)
		mockItemTagRepo.On("FindTagsForItems", ctx, model.ItemTypeDrawing, projectID, []uuid.UUID{d1ID, d2ID}).Return(tagsMap, nil)

		result, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(result.Content))
		assert.Equal(t, int64(2), result.TotalElements)
		assert.Equal(t, "Drawing 1", result.Content[0].Title)
		assert.Equal(t, &thumb, result.Content[0].Thumbnail)
		assert.Len(t, result.Content[0].Tags, 1)
		assert.Equal(t, "Drawing 2", result.Content[1].Title)
		assert.Nil(t, result.Content[1].Thumbnail)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
		mockItemTagRepo.AssertExpectations(t)
	})

	t.Run("GetAllByProjectID_EmptyList", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		emptyPage := model.NewPage([]model.Drawing{}, 0, 10, 0)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("FindAllByProjectID", ctx, projectID, 0, 10).Return(emptyPage, nil)

		result, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

		assert.NoError(t, err)
		assert.Empty(t, result.Content)
		assert.Equal(t, int64(0), result.TotalElements)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
		mockItemTagRepo.AssertNotCalled(t, "FindTagsForItems", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("GetAllByProjectID_ProjectNotFound", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		result, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		assert.Empty(t, result.Content)
		mockProjectRepo.AssertExpectations(t)
	})

	t.Run("GetAllByProjectID_RepoError", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("FindAllByProjectID", ctx, projectID, 0, 10).Return(nil, errors.New("db error"))

		result, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

		assert.Error(t, err)
		assert.Empty(t, result.Content)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
	})

	t.Run("GetAllByProjectID_TagsRepoError", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		expectedPage := model.NewPage([]model.Drawing{
			{ID: d1ID, ProjectID: projectID, Title: "Drawing 1"},
		}, 0, 10, 1)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("FindAllByProjectID", ctx, projectID, 0, 10).Return(expectedPage, nil)
		mockItemTagRepo.On("FindTagsForItems", ctx, model.ItemTypeDrawing, projectID, []uuid.UUID{d1ID}).Return(nil, errors.New("tags error"))

		result, err := uc.GetAllByProjectID(ctx, projectID, 0, 10)

		assert.Error(t, err)
		assert.Empty(t, result.Content)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
		mockItemTagRepo.AssertExpectations(t)
	})
}

func TestDrawingUseCase_Update(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	drawingID := uuid.New()

	existingDrawing := &model.Drawing{
		ID:        drawingID,
		ProjectID: projectID,
		Title:     "Old Title",
		Content:   nil,
		Thumbnail: nil,
		BaseEntity: model.BaseEntity{
			CreatedAt: time.Now(),
		},
	}

	newTitle := "Updated Title"
	newContent := `{"elements":[{"type":"rectangle"}]}`
	newThumb := "data:image/svg+xml;base64,newpreview"

	t.Run("Update_Success_AllFields", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		cmd := dto.UpdateDrawingCommand{
			ProjectID: projectID,
			ID:        drawingID,
			Title:     &newTitle,
			Content:   &newContent,
			Thumbnail: &newThumb,
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("FindByIDAndProjectID", ctx, drawingID, projectID).Return(existingDrawing, nil)
		mockDrawingRepo.On("Save", ctx, mock.MatchedBy(func(d *model.Drawing) bool {
			return d.Title == newTitle &&
				d.Content != nil && *d.Content == newContent &&
				d.Thumbnail != nil && *d.Thumbnail == newThumb &&
				d.UpdatedAt != nil
		})).Return(&model.Drawing{
			ID:        drawingID,
			ProjectID: projectID,
			Title:     newTitle,
			Content:   &newContent,
			Thumbnail: &newThumb,
		}, nil)
		mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeDrawing, projectID, drawingID).Return([]model.Tag{}, nil)

		updated, err := uc.Update(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, newTitle, updated.Title)
		assert.Equal(t, &newContent, updated.Content)
		assert.Equal(t, &newThumb, updated.Thumbnail)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
		mockItemTagRepo.AssertExpectations(t)
	})

	t.Run("Update_Success_PartialContentOnly", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		cmd := dto.UpdateDrawingCommand{
			ProjectID: projectID,
			ID:        drawingID,
			Content:   &newContent,
		}

		currentDrawing := &model.Drawing{
			ID:        drawingID,
			ProjectID: projectID,
			Title:     "Old Title",
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("FindByIDAndProjectID", ctx, drawingID, projectID).Return(currentDrawing, nil)
		mockDrawingRepo.On("Save", ctx, mock.MatchedBy(func(d *model.Drawing) bool {
			return d.Title == "Old Title" && d.Content == &newContent && d.UpdatedAt != nil
		})).Return(&model.Drawing{
			ID:        drawingID,
			ProjectID: projectID,
			Title:     "Old Title",
			Content:   &newContent,
		}, nil)
		mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeDrawing, projectID, drawingID).Return([]model.Tag{}, nil)

		updated, err := uc.Update(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, "Old Title", updated.Title)
		assert.Equal(t, &newContent, updated.Content)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
	})

	t.Run("Update_ProjectNotFound", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		cmd := dto.UpdateDrawingCommand{
			ProjectID: projectID,
			ID:        drawingID,
			Title:     &newTitle,
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		updated, err := uc.Update(ctx, cmd)

		assert.Nil(t, updated)
		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		mockProjectRepo.AssertExpectations(t)
	})

	t.Run("Update_DrawingNotFound", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		cmd := dto.UpdateDrawingCommand{
			ProjectID: projectID,
			ID:        drawingID,
			Title:     &newTitle,
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("FindByIDAndProjectID", ctx, drawingID, projectID).Return(nil, nil)

		updated, err := uc.Update(ctx, cmd)

		assert.Nil(t, updated)
		assert.ErrorIs(t, err, usecase.ErrDrawingNotFound)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
	})

	t.Run("Update_FindRepoError", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		cmd := dto.UpdateDrawingCommand{
			ProjectID: projectID,
			ID:        drawingID,
			Title:     &newTitle,
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("FindByIDAndProjectID", ctx, drawingID, projectID).Return(nil, errors.New("db find error"))

		updated, err := uc.Update(ctx, cmd)

		assert.Nil(t, updated)
		assert.Error(t, err)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
	})

	t.Run("Update_SaveRepoError", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		cmd := dto.UpdateDrawingCommand{
			ProjectID: projectID,
			ID:        drawingID,
			Title:     &newTitle,
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("FindByIDAndProjectID", ctx, drawingID, projectID).Return(existingDrawing, nil)
		mockDrawingRepo.On("Save", ctx, mock.Anything).Return(nil, errors.New("db save error"))

		updated, err := uc.Update(ctx, cmd)

		assert.Nil(t, updated)
		assert.Error(t, err)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
	})

	t.Run("Update_TagsRepoError", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		cmd := dto.UpdateDrawingCommand{
			ProjectID: projectID,
			ID:        drawingID,
			Title:     &newTitle,
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("FindByIDAndProjectID", ctx, drawingID, projectID).Return(existingDrawing, nil)
		mockDrawingRepo.On("Save", ctx, mock.Anything).Return(existingDrawing, nil)
		mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeDrawing, projectID, drawingID).Return(nil, errors.New("tags error"))

		updated, err := uc.Update(ctx, cmd)

		assert.Nil(t, updated)
		assert.Error(t, err)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
		mockItemTagRepo.AssertExpectations(t)
	})
}

func TestDrawingUseCase_Delete(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	drawingID := uuid.New()

	t.Run("Delete_Success", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("DeleteByIDAndProjectID", ctx, drawingID, projectID).Return(true, nil)
		mockItemTagRepo.On("RemoveAllTagsFromItem", ctx, model.ItemTypeDrawing, drawingID).Return(nil)

		err := uc.Delete(ctx, projectID, drawingID)

		assert.NoError(t, err)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
		mockItemTagRepo.AssertExpectations(t)
	})

	t.Run("Delete_ProjectNotFound", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		err := uc.Delete(ctx, projectID, drawingID)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		mockProjectRepo.AssertExpectations(t)
	})

	t.Run("Delete_DrawingNotFound", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("DeleteByIDAndProjectID", ctx, drawingID, projectID).Return(false, nil)

		err := uc.Delete(ctx, projectID, drawingID)

		assert.ErrorIs(t, err, usecase.ErrDrawingNotFound)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
	})

	t.Run("Delete_RepoError", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("DeleteByIDAndProjectID", ctx, drawingID, projectID).Return(false, errors.New("delete error"))

		err := uc.Delete(ctx, projectID, drawingID)

		assert.Error(t, err)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
	})

	t.Run("Delete_TagsRepoError", func(t *testing.T) {
		mockDrawingRepo := new(MockDrawingRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewDrawingUseCase(mockDrawingRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockDrawingRepo.On("DeleteByIDAndProjectID", ctx, drawingID, projectID).Return(true, nil)
		mockItemTagRepo.On("RemoveAllTagsFromItem", ctx, model.ItemTypeDrawing, drawingID).Return(errors.New("tags error"))

		err := uc.Delete(ctx, projectID, drawingID)

		assert.Error(t, err)
		mockProjectRepo.AssertExpectations(t)
		mockDrawingRepo.AssertExpectations(t)
		mockItemTagRepo.AssertExpectations(t)
	})
}
