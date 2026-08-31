package mcp

import (
	"context"
	"testing"

	"devaulty-backend/internal/domain/model"

	mcpgo "github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
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

func TestItemTagTools_RejectCredentialAndHideItFromSchema(t *testing.T) {
	app := setupMCPTestApp(t)
	server := mcpserver.NewMCPServer("Devaulty MCP Server", "test-version", mcpserver.WithRecovery())
	NewItemTagTools(app.itemTagUseCase).Register(server, Options{})

	associateTool, ok := server.ListTools()["associate_tag_to_item"]
	require.True(t, ok)
	associateEnums := toolEnumValues(t, associateTool, "itemType")
	require.NotContains(t, associateEnums, string(model.ItemTypeCredential))
	require.NotContains(t, associateEnums, "CREDENTIAL")

	dissociateTool, ok := server.ListTools()["dissociate_tag_from_item"]
	require.True(t, ok)
	dissociateEnums := toolEnumValues(t, dissociateTool, "itemType")
	require.NotContains(t, dissociateEnums, string(model.ItemTypeCredential))
	require.NotContains(t, dissociateEnums, "CREDENTIAL")

	project := app.createProject(t, "Credential validation project")
	tag := app.createTag(t, project.ID, "Must not allow")
	snippet := app.createSnippet(t, project.ID, "Forbidden tag")
	tools := NewItemTagTools(app.itemTagUseCase)

	invalidResult, err := tools.associate(context.Background(), newToolRequest("associate_tag_to_item", map[string]any{
		"projectID": project.ID.String(),
		"itemID":    snippet.ID.String(),
		"tagID":     tag.ID.String(),
		"itemType":  string(model.ItemTypeCredential),
	}))
	require.NoError(t, err)
	require.True(t, invalidResult.IsError)
	require.Contains(t, textFromResult(t, invalidResult), "not supported")
	require.Contains(t, textFromResult(t, invalidResult), "SNIPPET")
}

func toolEnumValues(t *testing.T, tool any, field string) []string {
	t.Helper()

	switch typed := tool.(type) {
	case mcpgo.Tool:
		return enumValuesForTool(t, typed.InputSchema, field)
	case *mcpserver.ServerTool:
		return enumValuesForTool(t, typed.Tool.InputSchema, field)
	default:
		t.Fatalf("unexpected tool type %T", tool)
		return nil
	}
}

func enumValuesForTool(t *testing.T, schema mcpgo.ToolInputSchema, field string) []string {
	t.Helper()

	property, ok := schema.Properties[field].(map[string]any)
	require.True(t, ok)

	raw, ok := property["enum"]
	require.True(t, ok)

	switch values := raw.(type) {
	case []string:
		return values
	case []any:
		result := make([]string, 0, len(values))
		for _, v := range values {
			result = append(result, v.(string))
		}
		return result
	default:
		t.Fatalf("unexpected enum type %T", raw)
		return nil
	}
}
