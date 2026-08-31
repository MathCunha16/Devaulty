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

// This is determined by the Frontend, here we are just providing the same defaults
var tagsColor = []string{"#ef4444", "#f97316", "#eab308", "#22c55e", "#10b981", "#14b8a6", "#06b6d4", "#3b82f6", "#6366f1", "#8b5cf6", "#ec4899", "#f43f5e", "#64748b"}

type TagTools struct {
	tagUseCase *usecase.TagUseCase
}

func NewTagTools(tagUseCase *usecase.TagUseCase) *TagTools {
	return &TagTools{tagUseCase: tagUseCase}
}

func (t *TagTools) Register(s *server.MCPServer, opts Options) {

	listTagTool := mcp.NewTool("list_tags",
		mcp.WithDescription("(non-paged) lists all tags in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(listTagTool, t.list)

	getTagTool := mcp.NewTool("get_tag",
		mcp.WithDescription("gets a tag by ID (uuid)"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the tag")),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(getTagTool, t.get)

	searchTagTool := mcp.NewTool("search_tags_by_name",
		mcp.WithDescription("searches tags in a project by name"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("name", mcp.Required(), mcp.Description("The tag name to search for")),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(searchTagTool, t.searchByName)

	if opts.ReadOnly {
		return
	}

	createTagTool := mcp.NewTool("create_tag",
		mcp.WithDescription("creates a new tag in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("name", mcp.Required(), mcp.Description("The name of the tag"), mcp.MinLength(2), mcp.MaxLength(40)),
		mcp.WithString("color", mcp.Description("The color of the tag"), mcp.Enum(tagsColor...)),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(createTagTool, t.create)

	updateTagTool := mcp.NewTool("update_tag",
		mcp.WithDescription("updates an existing tag in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the tag")),
		mcp.WithString("name", mcp.Description("The name of the tag"), mcp.MinLength(2), mcp.MaxLength(40)),
		mcp.WithString("color", mcp.Description("The color of the tag"), mcp.Enum(tagsColor...)),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(updateTagTool, t.update)

	if opts.DisableDelete {
		return
	}

	deleteTagTool := mcp.NewTool("delete_tag",
		mcp.WithDescription("deletes an existing tag in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the tag")),
		mcp.WithToolAnnotation(util.DeleteAnnotations))
	s.AddTool(deleteTagTool, t.delete)
}

func (t *TagTools) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.CreateTagCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tag, err := t.tagUseCase.Create(ctx, cmd)
	if err != nil {
		log.Printf("[TagTools.create] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(tag)
}

func (t *TagTools) list(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tags, err := t.tagUseCase.GetAllByProjectID(ctx, projectID)
	if err != nil {
		log.Printf("[TagTools.list] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(tags)
}

func (t *TagTools) get(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	id, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tag, err := t.tagUseCase.GetByID(ctx, projectID, id)
	if err != nil {
		log.Printf("[TagTools.get] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(tag)
}

func (t *TagTools) searchByName(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tags, err := t.tagUseCase.SearchByName(ctx, projectID, name)
	if err != nil {
		log.Printf("[TagTools.searchByName] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(tags)
}

func (t *TagTools) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.UpdateTagCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tag, err := t.tagUseCase.Update(ctx, cmd)
	if err != nil {
		log.Printf("[TagTools.update] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(tag)
}

func (t *TagTools) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	id, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = t.tagUseCase.Delete(ctx, projectID, id)
	if err != nil {
		log.Printf("[TagTools.delete] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Tag Deleted!"), nil
}
