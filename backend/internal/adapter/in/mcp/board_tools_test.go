package mcp

import (
	"context"
	"testing"

	"devaulty-backend/internal/dto"

	"github.com/stretchr/testify/require"
)

func TestBoardTools_CreateListGetUpdateDefaultDelete(t *testing.T) {
	app := setupMCPTestApp(t)
	project := app.createProject(t, "Board project")
	tools := NewBoardTools(app.boardUseCase)

	createResult, err := tools.create(context.Background(), newToolRequest("create_board", map[string]any{
		"projectID":   project.ID.String(),
		"name":        "Sprint board",
		"description": "Active sprint board",
		"isDefault":   true,
	}))
	require.NoError(t, err)
	require.False(t, createResult.IsError)

	created, ok := createResult.StructuredContent.(*dto.BoardView)
	require.True(t, ok)
	require.Equal(t, "Sprint board", created.Name)
	require.True(t, created.IsDefault)

	listResult, err := tools.list(context.Background(), newToolRequest("list_boards", map[string]any{
		"projectID":  project.ID.String(),
		"pageNumber": 1,
		"pageSize":   10,
	}))
	require.NoError(t, err)
	require.False(t, listResult.IsError)
	require.NotNil(t, listResult.StructuredContent)

	getResult, err := tools.get(context.Background(), newToolRequest("get_board", map[string]any{
		"projectID": project.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, getResult.IsError)

	fetched, ok := getResult.StructuredContent.(*dto.BoardView)
	require.True(t, ok)
	require.Equal(t, created.ID, fetched.ID)

	defaultResult, err := tools.getDefault(context.Background(), newToolRequest("get_default_board", map[string]any{
		"projectID": project.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, defaultResult.IsError)
	require.NotNil(t, defaultResult.StructuredContent)

	updateResult, err := tools.update(context.Background(), newToolRequest("update_board", map[string]any{
		"projectID":   project.ID.String(),
		"id":          created.ID.String(),
		"name":        "Updated sprint board",
		"description": "Updated board description",
		"isDefault":   false,
	}))
	require.NoError(t, err)
	require.False(t, updateResult.IsError)

	updated, ok := updateResult.StructuredContent.(*dto.BoardView)
	require.True(t, ok)
	require.Equal(t, "Updated sprint board", updated.Name)
	require.False(t, updated.IsDefault)

	deleteResult, err := tools.delete(context.Background(), newToolRequest("delete_board", map[string]any{
		"projectID": project.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, deleteResult.IsError)
	require.Equal(t, "Board Deleted!", textFromResult(t, deleteResult))

	invalidResult, err := tools.get(context.Background(), newToolRequest("get_board", map[string]any{
		"projectID": project.ID.String(),
		"id":        "not-a-uuid",
	}))
	require.NoError(t, err)
	require.True(t, invalidResult.IsError)
	require.Contains(t, textFromResult(t, invalidResult), "invalid UUID format")
}
