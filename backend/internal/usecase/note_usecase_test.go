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

func TestNoteUseCase_Create(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()

	cmd := dto.CreateNoteCommand{
		ProjectID: projectID,
		Title:     "Architecture Notes",
		Content:   "Hexagonal architecture design notes",
	}

	t.Run("Create_Success", func(t *testing.T) {
		mockNoteRepo := new(MockNoteRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewNoteUseCase(mockNoteRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockNoteRepo.On("Save", ctx, mock.MatchedBy(func(n *model.Note) bool {
			return n.Title == "Architecture Notes" && *n.Content == "Hexagonal architecture design notes" && n.ProjectID == projectID
		})).Return(&model.Note{
			ID:        uuid.New(),
			ProjectID: projectID,
			Title:     "Architecture Notes",
			Content:   &cmd.Content,
			Archived:  false,
			BaseEntity: model.BaseEntity{
				CreatedAt: time.Now(),
			},
		}, nil)

		result, err := uc.Create(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Architecture Notes", result.Title)
		assert.Equal(t, "Hexagonal architecture design notes", result.Content)
		mockProjectRepo.AssertExpectations(t)
		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("Create_ProjectNotFound", func(t *testing.T) {
		mockNoteRepo := new(MockNoteRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewNoteUseCase(mockNoteRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		result, err := uc.Create(ctx, cmd)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		assert.Nil(t, result)
	})
}

func TestNoteUseCase_GetByID(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	noteID := uuid.New()
	content := "Sample content"

	t.Run("GetByID_Success", func(t *testing.T) {
		mockNoteRepo := new(MockNoteRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewNoteUseCase(mockNoteRepo, mockProjectRepo, mockItemTagRepo)

		expectedNote := &model.Note{
			ID:        noteID,
			ProjectID: projectID,
			Title:     "My Note",
			Content:   &content,
			Archived:  false,
			BaseEntity: model.BaseEntity{
				CreatedAt: time.Now(),
			},
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockNoteRepo.On("FindByIDAndProjectID", ctx, projectID, noteID).Return(expectedNote, nil)
		mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeNote, projectID, noteID).Return([]model.Tag{}, nil)

		result, err := uc.GetByID(ctx, projectID, noteID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, noteID, result.ID)
		assert.Equal(t, "My Note", result.Title)
		assert.Equal(t, "Sample content", result.Content)
	})

	t.Run("GetByID_NoteNotFound", func(t *testing.T) {
		mockNoteRepo := new(MockNoteRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewNoteUseCase(mockNoteRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockNoteRepo.On("FindByIDAndProjectID", ctx, projectID, noteID).Return(nil, nil)

		result, err := uc.GetByID(ctx, projectID, noteID)

		assert.ErrorIs(t, err, usecase.ErrNoteNotFound)
		assert.Nil(t, result)
	})

	t.Run("GetByID_RepoError", func(t *testing.T) {
		mockNoteRepo := new(MockNoteRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewNoteUseCase(mockNoteRepo, mockProjectRepo, mockItemTagRepo)

		dbErr := errors.New("db error")
		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockNoteRepo.On("FindByIDAndProjectID", ctx, projectID, noteID).Return(nil, dbErr)

		result, err := uc.GetByID(ctx, projectID, noteID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestNoteUseCase_Delete(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	noteID := uuid.New()

	t.Run("Delete_Success", func(t *testing.T) {
		mockNoteRepo := new(MockNoteRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewNoteUseCase(mockNoteRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockNoteRepo.On("DeleteByIDAndProjectID", ctx, projectID, noteID).Return(true, nil)
		mockItemTagRepo.On("RemoveAllTagsFromItem", ctx, model.ItemTypeNote, noteID).Return(nil)

		err := uc.Delete(ctx, projectID, noteID)

		assert.NoError(t, err)
		mockProjectRepo.AssertExpectations(t)
		mockNoteRepo.AssertExpectations(t)
		mockItemTagRepo.AssertExpectations(t)
	})

	t.Run("Delete_NoteNotFound", func(t *testing.T) {
		mockNoteRepo := new(MockNoteRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewNoteUseCase(mockNoteRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockNoteRepo.On("DeleteByIDAndProjectID", ctx, projectID, noteID).Return(false, nil)

		err := uc.Delete(ctx, projectID, noteID)

		assert.ErrorIs(t, err, usecase.ErrNoteNotFound)
	})
}

func TestNoteUseCase_Archive(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	noteID := uuid.New()

	t.Run("Archive_Success", func(t *testing.T) {
		mockNoteRepo := new(MockNoteRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewNoteUseCase(mockNoteRepo, mockProjectRepo, mockItemTagRepo)

		existingNote := &model.Note{
			ID:        noteID,
			ProjectID: projectID,
			Title:     "Note to archive",
			Archived:  false,
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockNoteRepo.On("FindByIDAndProjectID", ctx, projectID, noteID).Return(existingNote, nil)
		mockNoteRepo.On("Save", ctx, mock.MatchedBy(func(n *model.Note) bool {
			return n.Archived == true && n.UpdatedAt != nil
		})).Return(existingNote, nil)

		err := uc.Archive(ctx, projectID, noteID)

		assert.NoError(t, err)
		mockProjectRepo.AssertExpectations(t)
		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("Archive_AlreadyArchived", func(t *testing.T) {
		mockNoteRepo := new(MockNoteRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewNoteUseCase(mockNoteRepo, mockProjectRepo, mockItemTagRepo)

		existingNote := &model.Note{
			ID:        noteID,
			ProjectID: projectID,
			Title:     "Already archived note",
			Archived:  true,
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockNoteRepo.On("FindByIDAndProjectID", ctx, projectID, noteID).Return(existingNote, nil)

		err := uc.Archive(ctx, projectID, noteID)

		assert.ErrorIs(t, err, usecase.ErrNoteAlreadyArchived)
		mockProjectRepo.AssertExpectations(t)
		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("Archive_NoteNotFound", func(t *testing.T) {
		mockNoteRepo := new(MockNoteRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewNoteUseCase(mockNoteRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockNoteRepo.On("FindByIDAndProjectID", ctx, projectID, noteID).Return(nil, nil)

		err := uc.Archive(ctx, projectID, noteID)

		assert.ErrorIs(t, err, usecase.ErrNoteNotFound)
		mockProjectRepo.AssertExpectations(t)
		mockNoteRepo.AssertExpectations(t)
	})
}

func TestNoteUseCase_Unarchive(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	noteID := uuid.New()

	t.Run("Unarchive_Success", func(t *testing.T) {
		mockNoteRepo := new(MockNoteRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewNoteUseCase(mockNoteRepo, mockProjectRepo, mockItemTagRepo)

		existingNote := &model.Note{
			ID:        noteID,
			ProjectID: projectID,
			Title:     "Archived note",
			Archived:  true,
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockNoteRepo.On("FindByIDAndProjectID", ctx, projectID, noteID).Return(existingNote, nil)
		mockNoteRepo.On("Save", ctx, mock.MatchedBy(func(n *model.Note) bool {
			return n.Archived == false && n.UpdatedAt != nil
		})).Return(existingNote, nil)

		err := uc.Unarchive(ctx, projectID, noteID)

		assert.NoError(t, err)
		mockProjectRepo.AssertExpectations(t)
		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("Unarchive_AlreadyUnarchived", func(t *testing.T) {
		mockNoteRepo := new(MockNoteRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewNoteUseCase(mockNoteRepo, mockProjectRepo, mockItemTagRepo)

		existingNote := &model.Note{
			ID:        noteID,
			ProjectID: projectID,
			Title:     "Open note",
			Archived:  false,
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockNoteRepo.On("FindByIDAndProjectID", ctx, projectID, noteID).Return(existingNote, nil)

		err := uc.Unarchive(ctx, projectID, noteID)

		assert.ErrorIs(t, err, usecase.ErrNoteAlreadyUnarchived)
		mockProjectRepo.AssertExpectations(t)
		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("Unarchive_NoteNotFound", func(t *testing.T) {
		mockNoteRepo := new(MockNoteRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockItemTagRepo := new(MockItemTagRepository)
		uc := usecase.NewNoteUseCase(mockNoteRepo, mockProjectRepo, mockItemTagRepo)

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockNoteRepo.On("FindByIDAndProjectID", ctx, projectID, noteID).Return(nil, nil)

		err := uc.Unarchive(ctx, projectID, noteID)

		assert.ErrorIs(t, err, usecase.ErrNoteNotFound)
		mockProjectRepo.AssertExpectations(t)
		mockNoteRepo.AssertExpectations(t)
	})
}
