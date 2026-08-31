package mcp

import (
	"context"
	"devaulty-backend/internal/adapter/in/mcp/util"
	"devaulty-backend/internal/dto"
	"devaulty-backend/internal/usecase"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type LinkTools struct {
	linkUseCase *usecase.LinkUseCase
}

func NewLinkTools(linkUseCase *usecase.LinkUseCase) *LinkTools {
	return &LinkTools{linkUseCase: linkUseCase}
}

func (t *LinkTools) Register(s *server.MCPServer, opts Options) {

	listLinkTool := mcp.NewTool("list_links",
		mcp.WithDescription("(paginated) lists all links in a Devaulty project"),
		mcp.WithInputSchema[util.ProjectPaginationQuery](),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(listLinkTool, t.list)

	getLinkTool := mcp.NewTool("get_link",
		mcp.WithDescription("gets a link by ID (uuid)"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the link")),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(getLinkTool, t.get)

	if opts.ReadOnly {
		return
	}

	createLinkTool := mcp.NewTool("create_link",
		mcp.WithDescription("creates a new link in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("title", mcp.Required(), mcp.Description("The title of the link")),
		mcp.WithString("url", mcp.Required(), mcp.Description("The URL of the link")),
		mcp.WithString("description", mcp.Description("The description of the link")),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(createLinkTool, t.create)

	updateLinkTool := mcp.NewTool("update_link",
		mcp.WithDescription("updates an existing link in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the link")),
		mcp.WithString("title", mcp.Description("The title of the link")),
		mcp.WithString("url", mcp.Description("The URL of the link")),
		mcp.WithString("description", mcp.Description("The description of the link")),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(updateLinkTool, t.update)

	if opts.DisableDelete {
		return
	}

	deleteLinkTool := mcp.NewTool("delete_link",
		mcp.WithDescription("deletes an existing link in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the link")),
		mcp.WithToolAnnotation(util.DeleteAnnotations))
	s.AddTool(deleteLinkTool, t.delete)
}

func (t *LinkTools) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.CreateLinkCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	link, err := t.linkUseCase.Create(ctx, cmd)
	if err != nil {
		log.Printf("[LinkTools.create] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(link)
}

func (t *LinkTools) list(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var query util.ProjectPaginationQuery
	if err := req.BindArguments(&query); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	util.ValidateProjectQuery(&query)

	pagedLinks, err := t.linkUseCase.GetAllByProjectID(ctx, query.ProjectID, query.PageNumber, query.PageSize)
	if err != nil {
		log.Printf("[LinkTools.list] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(pagedLinks)
}

func (t *LinkTools) get(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	linkID, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	link, err := t.linkUseCase.GetByID(ctx, projectID, linkID)
	if err != nil {
		log.Printf("[LinkTools.get] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(link)
}

func (t *LinkTools) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.UpdateLinkCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	link, err := t.linkUseCase.Update(ctx, cmd)
	if err != nil {
		log.Printf("[LinkTools.update] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(link)
}

func (t *LinkTools) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	linkID, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = t.linkUseCase.Delete(ctx, projectID, linkID)
	if err != nil {
		log.Printf("[LinkTools.delete] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Link Deleted!"), nil
}
