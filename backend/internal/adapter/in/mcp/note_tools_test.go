package mcp

import (
	"context"
	"testing"

	"devaulty-backend/internal/dto"

	"github.com/stretchr/testify/require"
)

func TestNoteTools_CreateListGetUpdateDelete(t *testing.T) {
	app := setupMCPTestApp(t)
	project := app.createProject(t, "Note project")
	tools := NewNoteTools(app.noteUseCase)

	createResult, err := tools.create(context.Background(), newToolRequest("create_note", map[string]any{
		"projectID": project.ID.String(),
		"title":     "Daily notes",
		"content":   "# Notes\n- Conference follow-up",
	}))
	require.NoError(t, err)
	require.False(t, createResult.IsError)

	created, ok := createResult.StructuredContent.(*dto.NoteView)
	require.True(t, ok)
	require.Equal(t, "Daily notes", created.Title)

	listResult, err := tools.list(context.Background(), newToolRequest("list_notes", map[string]any{
		"projectID":  project.ID.String(),
		"pageNumber": 1,
		"pageSize":   10,
	}))
	require.NoError(t, err)
	require.False(t, listResult.IsError)
	require.NotNil(t, listResult.StructuredContent)

	getResult, err := tools.get(context.Background(), newToolRequest("get_note", map[string]any{
		"projectID": project.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, getResult.IsError)

	fetched, ok := getResult.StructuredContent.(*dto.NoteView)
	require.True(t, ok)
	require.Equal(t, created.ID, fetched.ID)

	updateResult, err := tools.update(context.Background(), newToolRequest("update_note", map[string]any{
		"projectID": project.ID.String(),
		"id":        created.ID.String(),
		"title":     "Updated daily notes",
		"content":   "# Updated notes\n- Follow-up executed",
	}))
	require.NoError(t, err)
	require.False(t, updateResult.IsError)

	updated, ok := updateResult.StructuredContent.(*dto.NoteView)
	require.True(t, ok)
	require.Equal(t, "Updated daily notes", updated.Title)

	deleteResult, err := tools.delete(context.Background(), newToolRequest("delete_note", map[string]any{
		"projectID": project.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, deleteResult.IsError)
	require.Equal(t, "Note Deleted!", textFromResult(t, deleteResult))

	invalidResult, err := tools.get(context.Background(), newToolRequest("get_note", map[string]any{
		"projectID": project.ID.String(),
		"id":        "not-a-uuid",
	}))
	require.NoError(t, err)
	require.True(t, invalidResult.IsError)
	require.Contains(t, textFromResult(t, invalidResult), "invalid UUID format")
}
