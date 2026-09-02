package mcp

import (
	"context"
	"devaulty-backend/internal/adapter/in/mcp/util"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/dto"
	"devaulty-backend/internal/usecase"
	"log"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type SnippetTools struct {
	snippetUseCase *usecase.SnippetUseCase
}

func NewSnippetTools(snippetUseCase *usecase.SnippetUseCase) *SnippetTools {
	return &SnippetTools{snippetUseCase: snippetUseCase}
}

func (t *SnippetTools) Register(s *server.MCPServer, opts Options) {

	listSnippetTool := mcp.NewTool("list_snippets",
		mcp.WithDescription("(paginated) (first page at 0) lists all snippets in a Devaulty project"),
		mcp.WithInputSchema[util.ProjectPaginationQuery](),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(listSnippetTool, t.list)

	getSnippetTool := mcp.NewTool("get_snippet",
		mcp.WithDescription("gets a snippet by ID (uuid)"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the snippet")),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(getSnippetTool, t.get)

	if opts.ReadOnly {
		return
	}

	createSnippetTool := mcp.NewTool("create_snippet",
		mcp.WithDescription("creates a new snippet in a Devaulty project"),
		mcp.WithString(
			"projectID",
			mcp.Required(),
			mcp.Description("The projectID (UUID) of the project where you want to create a new snippet"),
		),
		mcp.WithString("title", mcp.Required(), mcp.Description("The title of the snippet"),
			mcp.MinLength(2), mcp.MaxLength(255)),
		mcp.WithString("description", mcp.Description("The description of the snippet"),
			mcp.MinLength(1)),
		mcp.WithString("content", mcp.Required(), mcp.Description("The content of the snippet")),
		mcp.WithString("language", mcp.Description("The language of the snippet"), mcp.Enum(util.ToStrings(model.SnippetLanguages)...)),
		mcp.WithString("snippetType", mcp.Description("The type of snippet"), mcp.Enum(util.ToStrings(model.SnippetTypes)...)),
		mcp.WithToolAnnotation(util.WriteAnnotations),
	)
	s.AddTool(createSnippetTool, t.create)

	updateSnippetTool := mcp.NewTool("update_snippet",
		mcp.WithDescription("updates an existing snippet in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project where you want to update a snippet")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the snippet")),
		mcp.WithString("title", mcp.Description("The title of the snippet"),
			mcp.MinLength(2), mcp.MaxLength(255)),
		mcp.WithString("description", mcp.Description("The description of the snippet"),
			mcp.MinLength(1)),
		mcp.WithString("content", mcp.Description("The content of the snippet")),
		mcp.WithString("language", mcp.Description("The language of the snippet"), mcp.Enum(util.ToStrings(model.SnippetLanguages)...)),
		mcp.WithString("snippetType", mcp.Description("The type of snippet"), mcp.Enum(util.ToStrings(model.SnippetTypes)...)),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(updateSnippetTool, t.update)

	if opts.DisableDelete {
		return
	}

	deleteSnippetTool := mcp.NewTool("delete_snippet",
		mcp.WithDescription("deletes an existing snippet in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project where you want to delete a snippet")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the snippet")),
		mcp.WithToolAnnotation(util.DeleteAnnotations))
	s.AddTool(deleteSnippetTool, t.delete)
}

func (t *SnippetTools) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.CreateSnippetCommand
	err := req.BindArguments(&cmd)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	snippet, err := t.snippetUseCase.Create(ctx, cmd)
	if err != nil {
		log.Printf("[SnippetTools.create] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultJSON(snippet)
}

func (t *SnippetTools) list(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var query util.ProjectPaginationQuery
	if err := req.BindArguments(&query); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	util.ValidateProjectQuery(&query)

	projectID, err := uuid.Parse(query.ProjectID)
	if err != nil {
		return mcp.NewToolResultError("Invalid project UUID" + err.Error()), nil
	}

	pagedSnippets, err := t.snippetUseCase.GetAllByProjectID(ctx, projectID, query.PageNumber, query.PageSize)
	if err != nil {
		log.Printf("[SnippetTools.list] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultJSON(pagedSnippets)
}

func (t *SnippetTools) get(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	snippetID, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	snippet, err := t.snippetUseCase.GetByID(ctx, projectID, snippetID)
	if err != nil {
		log.Printf("[SnippetTools.get] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultJSON(snippet)
}

func (t *SnippetTools) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.UpdateSnippetCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	snippet, err := t.snippetUseCase.Update(ctx, cmd)
	if err != nil {
		log.Printf("[SnippetTools.update] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(snippet)
}

func (t *SnippetTools) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	snippetID, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = t.snippetUseCase.Delete(ctx, projectID, snippetID)
	if err != nil {
		log.Printf("[SnippetTools.delete] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Snippet Deleted!"), nil
}
