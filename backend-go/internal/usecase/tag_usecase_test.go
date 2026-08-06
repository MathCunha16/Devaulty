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

type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) Save(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	args := m.Called(ctx, tag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}

func (m *MockTagRepository) FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Tag, error) {
	args := m.Called(ctx, projectID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}

func (m *MockTagRepository) FindAllByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Tag, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Tag), args.Error(1)
}

func (m *MockTagRepository) SearchByNameAndProjectID(ctx context.Context, name string, projectID uuid.UUID) ([]model.Tag, error) {
	args := m.Called(ctx, name, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Tag), args.Error(1)
}

func (m *MockTagRepository) DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, projectID, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockTagRepository) ExistsByNameAndProjectID(ctx context.Context, name string, projectID uuid.UUID) (bool, error) {
	args := m.Called(ctx, name, projectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTagRepository) ExistsByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error) {
	args := m.Called(ctx, id, projectID)
	return args.Bool(0), args.Error(1)
}

// --- UNIT TESTS ---

func TestTagUseCase_Create(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	color := "#FF0000"

	cmd := dto.CreateTagCommand{
		ProjectID: projectID,
		Name:      "  Go Backend  ",
		Color:     color,
	}

	t.Run("Create_Success", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("ExistsByNameAndProjectID", ctx, "Go Backend", projectID).Return(false, nil)
		mockTagRepo.On("Save", ctx, mock.MatchedBy(func(tag *model.Tag) bool {
			return tag.Name == "Go Backend" && tag.ProjectID == projectID && *tag.Color == color
		})).Return(&model.Tag{
			ID:        uuid.New(),
			ProjectID: projectID,
			Name:      "Go Backend",
			Color:     &color,
			BaseEntity: model.BaseEntity{
				CreatedAt: time.Now(),
			},
		}, nil)

		result, err := uc.Create(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Go Backend", result.Name)
		mockProjectRepo.AssertExpectations(t)
		mockTagRepo.AssertExpectations(t)
	})

	t.Run("Create_ProjectNotFound", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		result, err := uc.Create(ctx, cmd)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		assert.Nil(t, result)
	})

	t.Run("Create_TagAlreadyExists", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("ExistsByNameAndProjectID", ctx, "Go Backend", projectID).Return(true, nil)

		result, err := uc.Create(ctx, cmd)

		assert.ErrorIs(t, err, usecase.ErrTagAlreadyExists)
		assert.Nil(t, result)
	})

	t.Run("Create_ProjectRepoError", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		dbErr := errors.New("database failure")
		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, dbErr)

		result, err := uc.Create(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create_ExistsCheckError", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		dbErr := errors.New("query error")
		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("ExistsByNameAndProjectID", ctx, "Go Backend", projectID).Return(false, dbErr)

		result, err := uc.Create(ctx, cmd)

		assert.ErrorIs(t, err, dbErr)
		assert.Nil(t, result)
	})

	t.Run("Create_SaveError", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		dbErr := errors.New("insert error")
		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("ExistsByNameAndProjectID", ctx, "Go Backend", projectID).Return(false, nil)
		mockTagRepo.On("Save", ctx, mock.Anything).Return(nil, dbErr)

		result, err := uc.Create(ctx, cmd)

		assert.ErrorIs(t, err, dbErr)
		assert.Nil(t, result)
	})
}

func TestTagUseCase_GetByID(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	tagID := uuid.New()

	t.Run("GetByID_Success", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		expectedTag := &model.Tag{
			ID:        tagID,
			ProjectID: projectID,
			Name:      "Golang",
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("FindByIDAndProjectID", ctx, projectID, tagID).Return(expectedTag, nil)

		result, err := uc.GetByID(ctx, projectID, tagID)

		assert.NoError(t, err)
		assert.Equal(t, expectedTag, result)
	})

	t.Run("GetByID_ProjectNotFound", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		result, err := uc.GetByID(ctx, projectID, tagID)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		assert.Nil(t, result)
	})

	t.Run("GetByID_TagNotFound", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("FindByIDAndProjectID", ctx, projectID, tagID).Return(nil, nil)

		result, err := uc.GetByID(ctx, projectID, tagID)

		assert.ErrorIs(t, err, usecase.ErrTagNotFound)
		assert.Nil(t, result)
	})

	t.Run("GetByID_RepoError", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		dbErr := errors.New("db failure")
		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("FindByIDAndProjectID", ctx, projectID, tagID).Return(nil, dbErr)

		result, err := uc.GetByID(ctx, projectID, tagID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestTagUseCase_GetAllByProjectID(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()

	t.Run("GetAllByProjectID_Success", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		expectedTags := []model.Tag{
			{ID: uuid.New(), ProjectID: projectID, Name: "Tag 1"},
			{ID: uuid.New(), ProjectID: projectID, Name: "Tag 2"},
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("FindAllByProjectID", ctx, projectID).Return(expectedTags, nil)

		result, err := uc.GetAllByProjectID(ctx, projectID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, expectedTags, result)
	})

	t.Run("GetAllByProjectID_ProjectNotFound", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		result, err := uc.GetAllByProjectID(ctx, projectID)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		assert.Nil(t, result)
	})

	t.Run("GetAllByProjectID_RepoError", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		dbErr := errors.New("connection failed")
		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("FindAllByProjectID", ctx, projectID).Return(nil, dbErr)

		result, err := uc.GetAllByProjectID(ctx, projectID)

		assert.ErrorIs(t, err, dbErr)
		assert.Nil(t, result)
	})
}

func TestTagUseCase_SearchByName(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	searchTerm := "go"

	t.Run("SearchByName_Success", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		expectedTags := []model.Tag{
			{ID: uuid.New(), ProjectID: projectID, Name: "Golang"},
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("SearchByNameAndProjectID", ctx, searchTerm, projectID).Return(expectedTags, nil)

		result, err := uc.SearchByName(ctx, projectID, searchTerm)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Golang", result[0].Name)
	})

	t.Run("SearchByName_ProjectNotFound", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		result, err := uc.SearchByName(ctx, projectID, searchTerm)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		assert.Nil(t, result)
	})

	t.Run("SearchByName_RepoError", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		dbErr := errors.New("search failed")
		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("SearchByNameAndProjectID", ctx, searchTerm, projectID).Return(nil, dbErr)

		result, err := uc.SearchByName(ctx, projectID, searchTerm)

		assert.ErrorIs(t, err, dbErr)
		assert.Nil(t, result)
	})
}

func TestTagUseCase_Update(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	tagID := uuid.New()

	newName := "Updated Tag Name"
	newColor := "#00FF00"

	cmd := dto.UpdateTagCommand{
		ID:        tagID,
		ProjectID: projectID,
		Name:      &newName,
		Color:     &newColor,
	}

	t.Run("Update_Success_ChangeNameAndColor", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		existingTag := &model.Tag{
			ID:        tagID,
			ProjectID: projectID,
			Name:      "Old Name",
			Color:     nil,
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("FindByIDAndProjectID", ctx, projectID, tagID).Return(existingTag, nil)
		mockTagRepo.On("ExistsByNameAndProjectID", ctx, "Updated Tag Name", projectID).Return(false, nil)
		mockTagRepo.On("Save", ctx, mock.MatchedBy(func(tag *model.Tag) bool {
			return tag.Name == "Updated Tag Name" && *tag.Color == "#00FF00" && tag.UpdatedAt != nil
		})).Return(existingTag, nil)

		result, err := uc.Update(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockProjectRepo.AssertExpectations(t)
		mockTagRepo.AssertExpectations(t)
	})

	t.Run("Update_Success_SameNameDifferentCase", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		existingTag := &model.Tag{
			ID:        tagID,
			ProjectID: projectID,
			Name:      "Golang",
		}

		sameNameUpper := "GOLANG"
		updateCmd := dto.UpdateTagCommand{
			ID:        tagID,
			ProjectID: projectID,
			Name:      &sameNameUpper,
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("FindByIDAndProjectID", ctx, projectID, tagID).Return(existingTag, nil)
		// Should NOT call ExistsByNameAndProjectID because name equal (ignoring case)
		mockTagRepo.On("Save", ctx, mock.Anything).Return(existingTag, nil)

		result, err := uc.Update(ctx, updateCmd)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "GOLANG", result.Name)
	})

	t.Run("Update_ProjectNotFound", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		result, err := uc.Update(ctx, cmd)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		assert.Nil(t, result)
	})

	t.Run("Update_TagNotFound", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("FindByIDAndProjectID", ctx, projectID, tagID).Return(nil, nil)

		result, err := uc.Update(ctx, cmd)

		assert.ErrorIs(t, err, usecase.ErrTagNotFound)
		assert.Nil(t, result)
	})

	t.Run("Update_NameAlreadyExists", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		existingTag := &model.Tag{
			ID:        tagID,
			ProjectID: projectID,
			Name:      "Old Name",
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("FindByIDAndProjectID", ctx, projectID, tagID).Return(existingTag, nil)
		mockTagRepo.On("ExistsByNameAndProjectID", ctx, "Updated Tag Name", projectID).Return(true, nil)

		result, err := uc.Update(ctx, cmd)

		assert.ErrorIs(t, err, usecase.ErrTagAlreadyExists)
		assert.Nil(t, result)
	})

	t.Run("Update_FindRepoError", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		dbErr := errors.New("db error")
		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("FindByIDAndProjectID", ctx, projectID, tagID).Return(nil, dbErr)

		result, err := uc.Update(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Update_ExistsCheckError", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		existingTag := &model.Tag{
			ID:        tagID,
			ProjectID: projectID,
			Name:      "Old Name",
		}
		dbErr := errors.New("query failure")

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("FindByIDAndProjectID", ctx, projectID, tagID).Return(existingTag, nil)
		mockTagRepo.On("ExistsByNameAndProjectID", ctx, "Updated Tag Name", projectID).Return(false, dbErr)

		result, err := uc.Update(ctx, cmd)

		assert.ErrorIs(t, err, dbErr)
		assert.Nil(t, result)
	})

	t.Run("Update_SaveError", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		existingTag := &model.Tag{
			ID:        tagID,
			ProjectID: projectID,
			Name:      "Old Name",
		}
		dbErr := errors.New("save failure")

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("FindByIDAndProjectID", ctx, projectID, tagID).Return(existingTag, nil)
		mockTagRepo.On("ExistsByNameAndProjectID", ctx, "Updated Tag Name", projectID).Return(false, nil)
		mockTagRepo.On("Save", ctx, mock.Anything).Return(nil, dbErr)

		result, err := uc.Update(ctx, cmd)

		assert.ErrorIs(t, err, dbErr)
		assert.Nil(t, result)
	})
}

func TestTagUseCase_Delete(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	tagID := uuid.New()

	t.Run("Delete_Success", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("DeleteByIDAndProjectID", ctx, projectID, tagID).Return(true, nil)

		err := uc.Delete(ctx, projectID, tagID)

		assert.NoError(t, err)
		mockProjectRepo.AssertExpectations(t)
		mockTagRepo.AssertExpectations(t)
	})

	t.Run("Delete_ProjectNotFound", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		err := uc.Delete(ctx, projectID, tagID)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	})

	t.Run("Delete_TagNotFound", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("DeleteByIDAndProjectID", ctx, projectID, tagID).Return(false, nil)

		err := uc.Delete(ctx, projectID, tagID)

		assert.ErrorIs(t, err, usecase.ErrTagNotFound)
	})

	t.Run("Delete_RepoError", func(t *testing.T) {
		mockTagRepo := new(MockTagRepository)
		mockProjectRepo := new(MockProjectRepository)
		uc := usecase.NewTagUseCase(mockTagRepo, mockProjectRepo)

		dbErr := errors.New("delete db error")
		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("DeleteByIDAndProjectID", ctx, projectID, tagID).Return(false, dbErr)

		err := uc.Delete(ctx, projectID, tagID)

		assert.Error(t, err)
	})
}
