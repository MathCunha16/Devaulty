package mcp

import (
	"context"
	"testing"

	"devaulty-backend/internal/dto"

	"github.com/stretchr/testify/require"
)

func TestTagTools_CreateListGetSearchUpdateDelete(t *testing.T) {
	app := setupMCPTestApp(t)
	project := app.createProject(t, "Tag project")
	tools := NewTagTools(app.tagUseCase)

	createResult, err := tools.create(context.Background(), newToolRequest("create_tag", map[string]any{
		"projectID": project.ID.String(),
		"name":      "Frontend",
		"color":     "#8b5cf6",
	}))
	require.NoError(t, err)
	require.False(t, createResult.IsError)

	created, ok := createResult.StructuredContent.(*dto.TagView)
	require.True(t, ok)
	require.Equal(t, "Frontend", created.Name)

	listResult, err := tools.list(context.Background(), newToolRequest("list_tags", map[string]any{
		"projectID": project.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, listResult.IsError)
	require.NotNil(t, listResult.StructuredContent)

	getResult, err := tools.get(context.Background(), newToolRequest("get_tag", map[string]any{
		"projectID": project.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, getResult.IsError)

	fetched, ok := getResult.StructuredContent.(*dto.TagView)
	require.True(t, ok)
	require.Equal(t, created.ID, fetched.ID)

	searchResult, err := tools.searchByName(context.Background(), newToolRequest("search_tags_by_name", map[string]any{
		"projectID": project.ID.String(),
		"name":      "front",
	}))
	require.NoError(t, err)
	require.False(t, searchResult.IsError)
	require.NotNil(t, searchResult.StructuredContent)

	updateResult, err := tools.update(context.Background(), newToolRequest("update_tag", map[string]any{
		"projectID": project.ID.String(),
		"id":        created.ID.String(),
		"name":      "Frontend UI",
		"color":     "#ec4899",
	}))
	require.NoError(t, err)
	require.False(t, updateResult.IsError)

	updated, ok := updateResult.StructuredContent.(*dto.TagView)
	require.True(t, ok)
	require.Equal(t, "Frontend UI", updated.Name)

	deleteResult, err := tools.delete(context.Background(), newToolRequest("delete_tag", map[string]any{
		"projectID": project.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, deleteResult.IsError)
	require.Equal(t, "Tag Deleted!", textFromResult(t, deleteResult))

	invalidResult, err := tools.get(context.Background(), newToolRequest("get_tag", map[string]any{
		"projectID": project.ID.String(),
		"id":        "not-a-uuid",
	}))
	require.NoError(t, err)
	require.True(t, invalidResult.IsError)
	require.Contains(t, textFromResult(t, invalidResult), "invalid UUID format")
}
