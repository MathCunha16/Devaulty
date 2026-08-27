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

type MockNoteRepository struct {
	mock.Mock
}

func (m *MockNoteRepository) Save(ctx context.Context, note *model.Note) (*model.Note, error) {
	args := m.Called(ctx, note)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Note), args.Error(1)
}

func (m *MockNoteRepository) FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Note, error) {
	args := m.Called(ctx, projectID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Note), args.Error(1)
}

func (m *MockNoteRepository) FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[model.Note], error) {
	args := m.Called(ctx, projectID, page, size)
	return args.Get(0).(model.Page[model.Note]), args.Error(1)
}

func (m *MockNoteRepository) DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, projectID, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockNoteRepository) ExistsByIDAndProjectID(ctx context.Context, id, projectID uuid.UUID) (bool, error) {
	args := m.Called(ctx, id, projectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockNoteRepository) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, ids, projectID)
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

// --- UNIT TESTS ---

type ItemTagTestSetup struct {
	MockItemTagRepo *MockItemTagRepository
	MockTagRepo     *MockTagRepository
	MockProjectRepo *MockProjectRepository
	MockSnippetRepo *MockSnippetRepository
	MockCredRepo    *MockCredentialRepository
	MockLinkRepo    *MockLinkRepository
	MockProblemRepo *MockProblemRepository
	MockNoteRepo    *MockNoteRepository
	MockBoardRepo   *MockBoardRepository
	MockCardRepo    *MockCardRepository
	UC              *usecase.ItemTagUseCase
}

func SetupItemTagUseCaseTest() ItemTagTestSetup {
	mockItemTagRepo := new(MockItemTagRepository)
	mockTagRepo := new(MockTagRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockSnippetRepo := new(MockSnippetRepository)
	mockCredRepo := new(MockCredentialRepository)
	mockLinkRepo := new(MockLinkRepository)
	mockProblemRepo := new(MockProblemRepository)
	mockNoteRepo := new(MockNoteRepository)
	mockBoardRepo := new(MockBoardRepository)
	mockCardRepo := new(MockCardRepository)

	uc := usecase.NewItemTagUseCase(
		mockItemTagRepo,
		mockTagRepo,
		mockProjectRepo,
		mockSnippetRepo,
		mockCredRepo,
		mockLinkRepo,
		mockProblemRepo,
		mockNoteRepo,
		mockBoardRepo,
		mockCardRepo,
	)

	return ItemTagTestSetup{
		MockItemTagRepo: mockItemTagRepo,
		MockTagRepo:     mockTagRepo,
		MockProjectRepo: mockProjectRepo,
		MockSnippetRepo: mockSnippetRepo,
		MockCredRepo:    mockCredRepo,
		MockLinkRepo:    mockLinkRepo,
		MockProblemRepo: mockProblemRepo,
		MockNoteRepo:    mockNoteRepo,
		MockBoardRepo:   mockBoardRepo,
		MockCardRepo:    mockCardRepo,
		UC:              uc,
	}
}

func TestItemTagUseCase_AssociateTagToItem(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	tagID := uuid.New()
	snippetID := uuid.New()

	t.Run("Associate_Success_Snippet", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)
		s.MockSnippetRepo.On("ExistsByIDAndProjectID", ctx, snippetID, projectID).Return(true, nil)
		s.MockItemTagRepo.On("AssociateTagToItem", ctx, projectID, tagID, model.ItemTypeSnippet, snippetID).Return(nil)

		err := s.UC.AssociateTagToItem(ctx, projectID, model.ItemTypeSnippet, snippetID, tagID)

		assert.NoError(t, err)
		s.MockProjectRepo.AssertExpectations(t)
		s.MockTagRepo.AssertExpectations(t)
		s.MockSnippetRepo.AssertExpectations(t)
		s.MockItemTagRepo.AssertExpectations(t)
	})

	t.Run("Associate_Success_Board", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()
		boardID := uuid.New()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)
		s.MockBoardRepo.On("ExistsByIDAndProjectID", ctx, boardID, projectID).Return(true, nil)
		s.MockItemTagRepo.On("AssociateTagToItem", ctx, projectID, tagID, model.ItemTypeBoard, boardID).Return(nil)

		err := s.UC.AssociateTagToItem(ctx, projectID, model.ItemTypeBoard, boardID, tagID)

		assert.NoError(t, err)
		s.MockProjectRepo.AssertExpectations(t)
		s.MockTagRepo.AssertExpectations(t)
		s.MockBoardRepo.AssertExpectations(t)
		s.MockItemTagRepo.AssertExpectations(t)
	})

	t.Run("Associate_Success_Card", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()
		cardID := uuid.New()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)
		s.MockCardRepo.On("ExistsByIDAndProjectID", ctx, cardID, projectID).Return(true, nil)
		s.MockItemTagRepo.On("AssociateTagToItem", ctx, projectID, tagID, model.ItemTypeCard, cardID).Return(nil)

		err := s.UC.AssociateTagToItem(ctx, projectID, model.ItemTypeCard, cardID, tagID)

		assert.NoError(t, err)
		s.MockProjectRepo.AssertExpectations(t)
		s.MockTagRepo.AssertExpectations(t)
		s.MockCardRepo.AssertExpectations(t)
		s.MockItemTagRepo.AssertExpectations(t)
	})

	t.Run("Associate_ProjectNotFound", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		err := s.UC.AssociateTagToItem(ctx, projectID, model.ItemTypeSnippet, snippetID, tagID)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	})

	t.Run("Associate_TagNotFound", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(false, nil)

		err := s.UC.AssociateTagToItem(ctx, projectID, model.ItemTypeSnippet, snippetID, tagID)

		assert.ErrorIs(t, err, usecase.ErrTagNotFound)
	})

	t.Run("Associate_UnsupportedItemType", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)

		err := s.UC.AssociateTagToItem(ctx, projectID, model.ItemType("UNKNOWN"), snippetID, tagID)

		assert.ErrorIs(t, err, usecase.ErrUnsupportedItemType)
	})

	t.Run("Associate_ItemNotFound", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)
		s.MockSnippetRepo.On("ExistsByIDAndProjectID", ctx, snippetID, projectID).Return(false, nil)

		err := s.UC.AssociateTagToItem(ctx, projectID, model.ItemTypeSnippet, snippetID, tagID)

		assert.ErrorIs(t, err, usecase.ErrItemNotFound)
	})

	t.Run("Associate_RepoError", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		dbErr := errors.New("db insert failure")
		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)
		s.MockSnippetRepo.On("ExistsByIDAndProjectID", ctx, snippetID, projectID).Return(true, nil)
		s.MockItemTagRepo.On("AssociateTagToItem", ctx, projectID, tagID, model.ItemTypeSnippet, snippetID).Return(dbErr)

		err := s.UC.AssociateTagToItem(ctx, projectID, model.ItemTypeSnippet, snippetID, tagID)

		assert.ErrorIs(t, err, dbErr)
	})
}

func TestItemTagUseCase_DisassociateTagFromItem(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	tagID := uuid.New()
	linkID := uuid.New()

	t.Run("Disassociate_Success", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)
		s.MockLinkRepo.On("ExistsByIDAndProjectID", ctx, linkID, projectID).Return(true, nil)
		s.MockItemTagRepo.On("DisassembleTagFromItem", ctx, projectID, tagID, model.ItemTypeLink, linkID).Return(nil)

		err := s.UC.DisassociateTagFromItem(ctx, projectID, model.ItemTypeLink, linkID, tagID)

		assert.NoError(t, err)
		s.MockProjectRepo.AssertExpectations(t)
		s.MockTagRepo.AssertExpectations(t)
		s.MockLinkRepo.AssertExpectations(t)
		s.MockItemTagRepo.AssertExpectations(t)
	})

	t.Run("Disassociate_ProjectNotFound", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		err := s.UC.DisassociateTagFromItem(ctx, projectID, model.ItemTypeLink, linkID, tagID)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
	})

	t.Run("Disassociate_TagNotFound", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(false, nil)

		err := s.UC.DisassociateTagFromItem(ctx, projectID, model.ItemTypeLink, linkID, tagID)

		assert.ErrorIs(t, err, usecase.ErrTagNotFound)
	})

	t.Run("Disassociate_UnsupportedItemType", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)

		err := s.UC.DisassociateTagFromItem(ctx, projectID, model.ItemType("INVALID"), linkID, tagID)

		assert.ErrorIs(t, err, usecase.ErrUnsupportedItemType)
	})

	t.Run("Disassociate_ItemNotFound", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)
		s.MockLinkRepo.On("ExistsByIDAndProjectID", ctx, linkID, projectID).Return(false, nil)

		err := s.UC.DisassociateTagFromItem(ctx, projectID, model.ItemTypeLink, linkID, tagID)

		assert.ErrorIs(t, err, usecase.ErrItemNotFound)
	})

	t.Run("Disassociate_RepoError", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		dbErr := errors.New("db delete failure")
		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockTagRepo.On("ExistsByIDAndProjectID", ctx, tagID, projectID).Return(true, nil)
		s.MockLinkRepo.On("ExistsByIDAndProjectID", ctx, linkID, projectID).Return(true, nil)
		s.MockItemTagRepo.On("DisassembleTagFromItem", ctx, projectID, tagID, model.ItemTypeLink, linkID).Return(dbErr)

		err := s.UC.DisassociateTagFromItem(ctx, projectID, model.ItemTypeLink, linkID, tagID)

		assert.ErrorIs(t, err, dbErr)
	})
}

func TestItemTagUseCase_GetTagsForItem(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	problemID := uuid.New()

	t.Run("GetTagsForItem_Success", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		expectedTags := []model.Tag{
			{ID: uuid.New(), ProjectID: projectID, Name: "Bug"},
		}

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockProblemRepo.On("ExistsByIDAndProjectID", ctx, problemID, projectID).Return(true, nil)
		s.MockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeProblem, projectID, problemID).Return(expectedTags, nil)

		result, err := s.UC.GetTagsForItem(ctx, projectID, model.ItemTypeProblem, problemID)

		assert.NoError(t, err)
		assert.Equal(t, expectedTags, result)
	})

	t.Run("GetTagsForItem_ProjectNotFound", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		result, err := s.UC.GetTagsForItem(ctx, projectID, model.ItemTypeProblem, problemID)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		assert.Nil(t, result)
	})

	t.Run("GetTagsForItem_UnsupportedItemType", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)

		result, err := s.UC.GetTagsForItem(ctx, projectID, model.ItemType("DUMMY"), problemID)

		assert.ErrorIs(t, err, usecase.ErrUnsupportedItemType)
		assert.Nil(t, result)
	})

	t.Run("GetTagsForItem_ItemNotFound", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockProblemRepo.On("ExistsByIDAndProjectID", ctx, problemID, projectID).Return(false, nil)

		result, err := s.UC.GetTagsForItem(ctx, projectID, model.ItemTypeProblem, problemID)

		assert.ErrorIs(t, err, usecase.ErrItemNotFound)
		assert.Nil(t, result)
	})

	t.Run("GetTagsForItem_RepoError", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		dbErr := errors.New("find tags error")
		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockProblemRepo.On("ExistsByIDAndProjectID", ctx, problemID, projectID).Return(true, nil)
		s.MockItemTagRepo.On("FindTagsForItem", ctx, model.ItemTypeProblem, projectID, problemID).Return(nil, dbErr)

		result, err := s.UC.GetTagsForItem(ctx, projectID, model.ItemTypeProblem, problemID)

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
		s := SetupItemTagUseCaseTest()

		expectedMap := map[uuid.UUID][]model.Tag{
			id1: {{ID: uuid.New(), ProjectID: projectID, Name: "Tag A"}},
			id2: {{ID: uuid.New(), ProjectID: projectID, Name: "Tag B"}},
		}

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockSnippetRepo.On("FindExistingIDsByProjectID", ctx, itemIDs, projectID).Return(itemIDs, nil)
		s.MockItemTagRepo.On("FindTagsForItems", ctx, model.ItemTypeSnippet, projectID, itemIDs).Return(expectedMap, nil)

		result, err := s.UC.GetTagsForItems(ctx, projectID, model.ItemTypeSnippet, itemIDs)

		assert.NoError(t, err)
		assert.Equal(t, expectedMap, result)
	})

	t.Run("GetTagsForItems_Success_EmptyList", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)

		result, err := s.UC.GetTagsForItems(ctx, projectID, model.ItemTypeSnippet, []uuid.UUID{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("GetTagsForItems_ProjectNotFound", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		result, err := s.UC.GetTagsForItems(ctx, projectID, model.ItemTypeSnippet, itemIDs)

		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		assert.Nil(t, result)
	})

	t.Run("GetTagsForItems_UnsupportedItemType", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)

		result, err := s.UC.GetTagsForItems(ctx, projectID, model.ItemType("UNKNOWN"), itemIDs)

		assert.ErrorIs(t, err, usecase.ErrUnsupportedItemType)
		assert.Nil(t, result)
	})

	t.Run("GetTagsForItems_ItemNotFound", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		// Only id1 is returned from db, so length doesn't match itemIDs (2 != 1)
		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockSnippetRepo.On("FindExistingIDsByProjectID", ctx, itemIDs, projectID).Return([]uuid.UUID{id1}, nil)

		result, err := s.UC.GetTagsForItems(ctx, projectID, model.ItemTypeSnippet, itemIDs)

		assert.ErrorIs(t, err, usecase.ErrItemNotFound)
		assert.Nil(t, result)
	})

	t.Run("GetTagsForItems_RepoError", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		dbErr := errors.New("batch query error")
		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockSnippetRepo.On("FindExistingIDsByProjectID", ctx, itemIDs, projectID).Return(itemIDs, nil)
		s.MockItemTagRepo.On("FindTagsForItems", ctx, model.ItemTypeSnippet, projectID, itemIDs).Return(nil, dbErr)

		result, err := s.UC.GetTagsForItems(ctx, projectID, model.ItemTypeSnippet, itemIDs)

		assert.ErrorIs(t, err, dbErr)
		assert.Nil(t, result)
	})

	t.Run("GetTagsForItems_DuplicateIDsSuccess", func(t *testing.T) {
		s := SetupItemTagUseCaseTest()

		duplicateIDs := []uuid.UUID{id1, id1}
		existingIDs := []uuid.UUID{id1} // DB returns unique existing ID

		expectedTags := map[uuid.UUID][]model.Tag{
			id1: {{ID: uuid.New(), Name: "Go"}},
		}

		s.MockProjectRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		s.MockSnippetRepo.On("FindExistingIDsByProjectID", ctx, duplicateIDs, projectID).Return(existingIDs, nil)
		s.MockItemTagRepo.On("FindTagsForItems", ctx, model.ItemTypeSnippet, projectID, duplicateIDs).Return(expectedTags, nil)

		result, err := s.UC.GetTagsForItems(ctx, projectID, model.ItemTypeSnippet, duplicateIDs)

		assert.NoError(t, err)
		assert.Equal(t, expectedTags, result)
	})
}
