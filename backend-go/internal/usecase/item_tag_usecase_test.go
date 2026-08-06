package usecase_test

import (
	"context"
	"errors"
	"testing"

	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockItemTagRepository struct {
	mock.Mock
}

func (m *MockItemTagRepository) AssociateTagToItem(ctx context.Context, projectID, tagID uuid.UUID, itemType model.ItemType, itemID uuid.UUID) error {
	args := m.Called(ctx, projectID, tagID, itemType, itemID)
	return args.Error(0)
}

func (m *MockItemTagRepository) DisassembleTagFromItem(ctx context.Context, projectID, tagID uuid.UUID, itemType model.ItemType, itemID uuid.UUID) error {
	args := m.Called(ctx, projectID, tagID, itemType, itemID)
	return args.Error(0)
}

func (m *MockItemTagRepository) RemoveAllTagsFromItem(ctx context.Context, itemType model.ItemType, itemID uuid.UUID) error {
	args := m.Called(ctx, itemType, itemID)
	return args.Error(0)
}

func (m *MockItemTagRepository) FindTagsForItem(ctx context.Context, itemType model.ItemType, projectID, itemID uuid.UUID) ([]model.Tag, error) {
	args := m.Called(ctx, itemType, projectID, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Tag), args.Error(1)
}

func (m *MockItemTagRepository) FindTagsForItems(ctx context.Context, itemType model.ItemType, projectID uuid.UUID, itemIDs []uuid.UUID) (map[uuid.UUID][]model.Tag, error) {
	args := m.Called(ctx, itemType, projectID, itemIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[uuid.UUID][]model.Tag), args.Error(1)
}

// --- UNIT TESTS ---

func SetupItemTagUseCaseTest() (
	*MockItemTagRepository,
	*MockTagRepository,
	*MockProjectRepository,
	*MockSnippetRepository,
	*MockLinkRepository,
	*MockProblemRepository,
	*usecase.ItemTagUseCase,
) {
	mockItemTagRepo := new(MockItemTagRepository)
	mockTagRepo := new(MockTagRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockSnippetRepo := new(MockSnippetRepository)
	mockLinkRepo := new(MockLinkRepository)
	mockProblemRepo := new(MockProblemRepository)

	uc := usecase.NewItemTagUseCase(
		mockItemTagRepo,
		mockTagRepo,
		mockProjectRepo,
		mockSnippetRepo,
		mockLinkRepo,
		mockProblemRepo,
	)

	return mockItemTagRepo, mockTagRepo, mockProjectRepo, mockSnippetRepo, mockLinkRepo, mockProblemRepo, uc
}

func TestItemTagUseCase_AssociateTagToItem(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	tagID := uuid.New()
	snippetID := uuid.New()

	t.Run("Associate_Success", func(t *testing.T) {
		mockItemTagRepo, mockTagRepo, mockProjectRepo, mockSnippetRepo, _, _, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)
		mockSnippetRepo.On("ExistsByIDAndProjectID", ctx, snippetID, projectID).Return(true, nil)
		mockItemTagRepo.On("AssociateTagToItem", ctx, projectID, tagID, model.ItemTypeSnippet, snippetID).Return(nil)

		err := uc.AssociateTagToItem(ctx, projectID, model.ItemTypeSnippet, snippetID, tagID)

		assert.NoError(t, err)
		mockProjectRepo.AssertExpectations(t)
		mockTagRepo.AssertExpectations(t)
		mockSnippetRepo.AssertExpectations(t)
		mockItemTagRepo.AssertExpectations(t)
	})

	t.Run("Associate_ProjectNotFound", func(t *testing.T) {
		_, _, mockProjectRepo, _, _, _, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		err := uc.AssociateTagToItem(ctx, projectID, model.ItemTypeSnippet, snippetID, tagID)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	})

	t.Run("Associate_TagNotFound", func(t *testing.T) {
		_, mockTagRepo, mockProjectRepo, _, _, _, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(false, nil)

		err := uc.AssociateTagToItem(ctx, projectID, model.ItemTypeSnippet, snippetID, tagID)

		assert.ErrorIs(t, err, usecase.ErrTagNotFound)
	})

	t.Run("Associate_UnsupportedItemType", func(t *testing.T) {
		_, mockTagRepo, mockProjectRepo, _, _, _, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)

		err := uc.AssociateTagToItem(ctx, projectID, model.ItemType("UNKNOWN"), snippetID, tagID)

		assert.ErrorIs(t, err, usecase.ErrUnsupportedItemType)
	})

	t.Run("Associate_ItemNotFound", func(t *testing.T) {
		_, mockTagRepo, mockProjectRepo, mockSnippetRepo, _, _, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)
		mockSnippetRepo.On("ExistsByIDAndProjectID", ctx, snippetID, projectID).Return(false, nil)

		err := uc.AssociateTagToItem(ctx, projectID, model.ItemTypeSnippet, snippetID, tagID)

		assert.ErrorIs(t, err, usecase.ErrItemNotFound)
	})

	t.Run("Associate_RepoError", func(t *testing.T) {
		mockItemTagRepo, mockTagRepo, mockProjectRepo, mockSnippetRepo, _, _, uc := SetupItemTagUseCaseTest()

		dbErr := errors.New("db insert failure")
		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)
		mockSnippetRepo.On("ExistsByIDAndProjectID", ctx, snippetID, projectID).Return(true, nil)
		mockItemTagRepo.On("AssociateTagToItem", ctx, projectID, tagID, model.ItemTypeSnippet, snippetID).Return(dbErr)

		err := uc.AssociateTagToItem(ctx, projectID, model.ItemTypeSnippet, snippetID, tagID)

		assert.ErrorIs(t, err, dbErr)
	})
}

func TestItemTagUseCase_DisassociateTagFromItem(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	tagID := uuid.New()
	linkID := uuid.New()

	t.Run("Disassociate_Success", func(t *testing.T) {
		mockItemTagRepo, mockTagRepo, mockProjectRepo, _, mockLinkRepo, _, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)
		mockLinkRepo.On("ExistsByIDAndProjectID", ctx, linkID, projectID).Return(true, nil)
		mockItemTagRepo.On("DisassembleTagFromItem", ctx, projectID, tagID, model.ItemTypeLink, linkID).Return(nil)

		err := uc.DisassociateTagFromItem(ctx, projectID, model.ItemTypeLink, linkID, tagID)

		assert.NoError(t, err)
		mockProjectRepo.AssertExpectations(t)
		mockTagRepo.AssertExpectations(t)
		mockLinkRepo.AssertExpectations(t)
		mockItemTagRepo.AssertExpectations(t)
	})

	t.Run("Disassociate_ProjectNotFound", func(t *testing.T) {
		_, _, mockProjectRepo, _, _, _, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		err := uc.DisassociateTagFromItem(ctx, projectID, model.ItemTypeLink, linkID, tagID)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	})

	t.Run("Disassociate_TagNotFound", func(t *testing.T) {
		_, mockTagRepo, mockProjectRepo, _, _, _, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(false, nil)

		err := uc.DisassociateTagFromItem(ctx, projectID, model.ItemTypeLink, linkID, tagID)

		assert.ErrorIs(t, err, usecase.ErrTagNotFound)
	})

	t.Run("Disassociate_UnsupportedItemType", func(t *testing.T) {
		_, mockTagRepo, mockProjectRepo, _, _, _, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)

		err := uc.DisassociateTagFromItem(ctx, projectID, model.ItemType("INVALID"), linkID, tagID)

		assert.ErrorIs(t, err, usecase.ErrUnsupportedItemType)
	})

	t.Run("Disassociate_ItemNotFound", func(t *testing.T) {
		_, mockTagRepo, mockProjectRepo, _, mockLinkRepo, _, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)
		mockLinkRepo.On("ExistsByIDAndProjectID", ctx, linkID, projectID).Return(false, nil)

		err := uc.DisassociateTagFromItem(ctx, projectID, model.ItemTypeLink, linkID, tagID)

		assert.ErrorIs(t, err, usecase.ErrItemNotFound)
	})

	t.Run("Disassociate_RepoError", func(t *testing.T) {
		mockItemTagRepo, mockTagRepo, mockProjectRepo, _, mockLinkRepo, _, uc := SetupItemTagUseCaseTest()

		dbErr := errors.New("db delete failure")
		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)
		mockLinkRepo.On("ExistsByIDAndProjectID", ctx, linkID, projectID).Return(true, nil)
		mockItemTagRepo.On("DisassembleTagFromItem", ctx, projectID, tagID, model.ItemTypeLink, linkID).Return(dbErr)

		err := uc.DisassociateTagFromItem(ctx, projectID, model.ItemTypeLink, linkID, tagID)

		assert.ErrorIs(t, err, dbErr)
	})
}

func TestItemTagUseCase_GetTagsForItem(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	problemID := uuid.New()

	t.Run("GetTagsForItem_Success", func(t *testing.T) {
		mockItemTagRepo, _, mockProjectRepo, _, _, mockProblemRepo, uc := SetupItemTagUseCaseTest()

		expectedTags := []model.Tag{
			{ID: uuid.New(), ProjectID: projectID, Name: "Bug"},
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockProblemRepo.On("ExistsByIDAndProjectID", ctx, problemID, projectID).Return(true, nil)
		mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeProblem, projectID, problemID).Return(expectedTags, nil)

		result, err := uc.GetTagsForItem(ctx, projectID, model.ItemTypeProblem, problemID)

		assert.NoError(t, err)
		assert.Equal(t, expectedTags, result)
	})

	t.Run("GetTagsForItem_ProjectNotFound", func(t *testing.T) {
		_, _, mockProjectRepo, _, _, _, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		result, err := uc.GetTagsForItem(ctx, projectID, model.ItemTypeProblem, problemID)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		assert.Nil(t, result)
	})

	t.Run("GetTagsForItem_UnsupportedItemType", func(t *testing.T) {
		_, _, mockProjectRepo, _, _, _, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)

		result, err := uc.GetTagsForItem(ctx, projectID, model.ItemType("DUMMY"), problemID)

		assert.ErrorIs(t, err, usecase.ErrUnsupportedItemType)
		assert.Nil(t, result)
	})

	t.Run("GetTagsForItem_ItemNotFound", func(t *testing.T) {
		_, _, mockProjectRepo, _, _, mockProblemRepo, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockProblemRepo.On("ExistsByIDAndProjectID", ctx, problemID, projectID).Return(false, nil)

		result, err := uc.GetTagsForItem(ctx, projectID, model.ItemTypeProblem, problemID)

		assert.ErrorIs(t, err, usecase.ErrItemNotFound)
		assert.Nil(t, result)
	})

	t.Run("GetTagsForItem_RepoError", func(t *testing.T) {
		mockItemTagRepo, _, mockProjectRepo, _, _, mockProblemRepo, uc := SetupItemTagUseCaseTest()

		dbErr := errors.New("find tags error")
		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockProblemRepo.On("ExistsByIDAndProjectID", ctx, problemID, projectID).Return(true, nil)
		mockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeProblem, projectID, problemID).Return(nil, dbErr)

		result, err := uc.GetTagsForItem(ctx, projectID, model.ItemTypeProblem, problemID)

		assert.ErrorIs(t, err, dbErr)
		assert.Nil(t, result)
	})
}

func TestItemTagUseCase_GetTagsForItems(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	id1 := uuid.New()
	id2 := uuid.New()
	itemIDs := []uuid.UUID{id1, id2}

	t.Run("GetTagsForItems_Success", func(t *testing.T) {
		mockItemTagRepo, _, mockProjectRepo, mockSnippetRepo, _, _, uc := SetupItemTagUseCaseTest()

		expectedMap := map[uuid.UUID][]model.Tag{
			id1: {{ID: uuid.New(), ProjectID: projectID, Name: "Tag A"}},
			id2: {{ID: uuid.New(), ProjectID: projectID, Name: "Tag B"}},
		}

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockSnippetRepo.On("FindExistingIDsByProjectID", ctx, itemIDs, projectID).Return(itemIDs, nil)
		mockItemTagRepo.On("FindTagsForItems", ctx, model.ItemTypeSnippet, projectID, itemIDs).Return(expectedMap, nil)

		result, err := uc.GetTagsForItems(ctx, projectID, model.ItemTypeSnippet, itemIDs)

		assert.NoError(t, err)
		assert.Equal(t, expectedMap, result)
	})

	t.Run("GetTagsForItems_Success_EmptyList", func(t *testing.T) {
		_, _, mockProjectRepo, _, _, _, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)

		result, err := uc.GetTagsForItems(ctx, projectID, model.ItemTypeSnippet, []uuid.UUID{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("GetTagsForItems_ProjectNotFound", func(t *testing.T) {
		_, _, mockProjectRepo, _, _, _, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		result, err := uc.GetTagsForItems(ctx, projectID, model.ItemTypeSnippet, itemIDs)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		assert.Nil(t, result)
	})

	t.Run("GetTagsForItems_UnsupportedItemType", func(t *testing.T) {
		_, _, mockProjectRepo, _, _, _, uc := SetupItemTagUseCaseTest()

		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)

		result, err := uc.GetTagsForItems(ctx, projectID, model.ItemType("UNKNOWN"), itemIDs)

		assert.ErrorIs(t, err, usecase.ErrUnsupportedItemType)
		assert.Nil(t, result)
	})

	t.Run("GetTagsForItems_ItemNotFound", func(t *testing.T) {
		_, _, mockProjectRepo, mockSnippetRepo, _, _, uc := SetupItemTagUseCaseTest()

		// Only id1 is returned from db, so length doesn't match itemIDs (2 != 1)
		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockSnippetRepo.On("FindExistingIDsByProjectID", ctx, itemIDs, projectID).Return([]uuid.UUID{id1}, nil)

		result, err := uc.GetTagsForItems(ctx, projectID, model.ItemTypeSnippet, itemIDs)

		assert.ErrorIs(t, err, usecase.ErrItemNotFound)
		assert.Nil(t, result)
	})

	t.Run("GetTagsForItems_RepoError", func(t *testing.T) {
		mockItemTagRepo, _, mockProjectRepo, mockSnippetRepo, _, _, uc := SetupItemTagUseCaseTest()

		dbErr := errors.New("batch query error")
		mockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		mockSnippetRepo.On("FindExistingIDsByProjectID", ctx, itemIDs, projectID).Return(itemIDs, nil)
		mockItemTagRepo.On("FindTagsForItems", ctx, model.ItemTypeSnippet, projectID, itemIDs).Return(nil, dbErr)

		result, err := uc.GetTagsForItems(ctx, projectID, model.ItemTypeSnippet, itemIDs)

		assert.ErrorIs(t, err, dbErr)
		assert.Nil(t, result)
	})
}
