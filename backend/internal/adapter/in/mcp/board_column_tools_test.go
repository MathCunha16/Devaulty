package mcp

import (
	"context"
	"testing"

	"devaulty-backend/internal/dto"

	"github.com/stretchr/testify/require"
)

func TestBoardColumnTools_CreateListGetUpdateReorderDelete(t *testing.T) {
	app := setupMCPTestApp(t)
	project := app.createProject(t, "Board column project")
	board := app.createBoard(t, project.ID, "Design board")
	tools := NewBoardColumnTools(app.boardColumnUseCase)

	createResult, err := tools.create(context.Background(), newToolRequest("create_board_column", map[string]any{
		"projectID": project.ID.String(),
		"boardID":   board.ID.String(),
		"name":      "Backlog",
		"wipLimit":  5,
	}))
	require.NoError(t, err)
	require.False(t, createResult.IsError)

	created, ok := createResult.StructuredContent.(*dto.BoardColumnView)
	require.True(t, ok)
	require.Equal(t, "Backlog", created.Name)
	require.NotNil(t, created.WipLimit)
	require.EqualValues(t, 5, *created.WipLimit)

	secondResult, err := tools.create(context.Background(), newToolRequest("create_board_column", map[string]any{
		"projectID": project.ID.String(),
		"boardID":   board.ID.String(),
		"name":      "In Progress",
		"wipLimit":  3,
	}))
	require.NoError(t, err)
	require.False(t, secondResult.IsError)

	listResult, err := tools.list(context.Background(), newToolRequest("list_board_columns", map[string]any{
		"projectID": project.ID.String(),
		"boardID":   board.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, listResult.IsError)
	require.NotNil(t, listResult.StructuredContent)

	getResult, err := tools.get(context.Background(), newToolRequest("get_board_column", map[string]any{
		"projectID": project.ID.String(),
		"boardID":   board.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, getResult.IsError)

	fetched, ok := getResult.StructuredContent.(*dto.BoardColumnView)
	require.True(t, ok)
	require.Equal(t, created.ID, fetched.ID)

	updateResult, err := tools.update(context.Background(), newToolRequest("update_board_column", map[string]any{
		"projectID": project.ID.String(),
		"boardID":   board.ID.String(),
		"id":        created.ID.String(),
		"name":      "Planned",
		"wipLimit":  4,
	}))
	require.NoError(t, err)
	require.False(t, updateResult.IsError)

	updated, ok := updateResult.StructuredContent.(*dto.BoardColumnView)
	require.True(t, ok)
	require.Equal(t, "Planned", updated.Name)
	require.EqualValues(t, 4, *updated.WipLimit)

	allColumnsResult, err := tools.list(context.Background(), newToolRequest("list_board_columns", map[string]any{
		"projectID": project.ID.String(),
		"boardID":   board.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, allColumnsResult.IsError)

	columns, ok := allColumnsResult.StructuredContent.([]*dto.BoardColumnView)
	require.True(t, ok)
	require.Len(t, columns, 6)

	orderedIDs := []string{columns[0].ID.String(), columns[2].ID.String(), columns[1].ID.String(), columns[3].ID.String(), columns[4].ID.String(), columns[5].ID.String()}
	reorderResult, err := tools.reorder(context.Background(), newToolRequest("reorder_board_columns", map[string]any{
		"projectID": project.ID.String(),
		"boardID":   board.ID.String(),
		"positions": orderedIDs,
	}))
	require.NoError(t, err)
	require.False(t, reorderResult.IsError)
	require.NotNil(t, reorderResult.StructuredContent)

	deleteResult, err := tools.delete(context.Background(), newToolRequest("delete_board_column", map[string]any{
		"projectID": project.ID.String(),
		"boardID":   board.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, deleteResult.IsError)
	require.Equal(t, "Board Column Deleted!", textFromResult(t, deleteResult))

	invalidResult, err := tools.get(context.Background(), newToolRequest("get_board_column", map[string]any{
		"projectID": project.ID.String(),
		"boardID":   board.ID.String(),
		"id":        "not-a-uuid",
	}))
	require.NoError(t, err)
	require.True(t, invalidResult.IsError)
	require.Contains(t, textFromResult(t, invalidResult), "invalid UUID format")
}
