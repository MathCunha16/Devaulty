package mcp

import (
	"context"
	"testing"

	"devaulty-backend/internal/domain/model"

	"github.com/stretchr/testify/require"
)

func TestItemTagTools_AssociateAndDissociate(t *testing.T) {
	app := setupMCPTestApp(t)
	project := app.createProject(t, "Item tag project")
	tag := app.createTag(t, project.ID, "Bug")
	snippet := app.createSnippet(t, project.ID, "Auth bootstrap")
	tools := NewItemTagTools(app.itemTagUseCase)

	associateResult, err := tools.associate(context.Background(), newToolRequest("associate_tag_to_item", map[string]any{
		"projectID": project.ID.String(),
		"itemID":    snippet.ID.String(),
		"tagID":     tag.ID.String(),
		"itemType":  string(model.ItemTypeSnippet),
	}))
	require.NoError(t, err)
	require.False(t, associateResult.IsError)
	require.Equal(t, "Item Tag Associated!", textFromResult(t, associateResult))

	dissociateResult, err := tools.disassociate(context.Background(), newToolRequest("dissociate_tag_from_item", map[string]any{
		"projectID": project.ID.String(),
		"itemID":    snippet.ID.String(),
		"tagID":     tag.ID.String(),
		"itemType":  "snippet",
	}))
	require.NoError(t, err)
	require.False(t, dissociateResult.IsError)
	require.Equal(t, "Item Tag Dissociated!", textFromResult(t, dissociateResult))

	invalidResult, err := tools.associate(context.Background(), newToolRequest("associate_tag_to_item", map[string]any{
		"projectID": project.ID.String(),
		"itemID":    "not-a-uuid",
		"tagID":     tag.ID.String(),
		"itemType":  string(model.ItemTypeSnippet),
	}))
	require.NoError(t, err)
	require.True(t, invalidResult.IsError)
	require.Contains(t, textFromResult(t, invalidResult), "invalid UUID format")
}
