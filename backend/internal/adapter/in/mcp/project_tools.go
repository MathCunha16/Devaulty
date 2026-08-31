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

var ( // This is determined by the Frontend, here we are just providing the same defaults
	projectIcons  = []string{"Folder", "Terminal", "Database", "Globe", "Cpu", "Activity", "BookOpen", "Code"}
	projectColors = []string{"#ef4444", "#f97316", "#facc15", "#22c55e", "#06b6d4", "#3b82f6", "#8b5cf6", "#ec4899"}
)

type ProjectTools struct {
	projectUseCase *usecase.ProjectUseCase
}

func NewProjectTools(projectUseCase *usecase.ProjectUseCase) *ProjectTools {
	return &ProjectTools{projectUseCase: projectUseCase}
}

func (t *ProjectTools) Register(s *server.MCPServer, opts Options) {

	listProjectTool := mcp.NewTool("list_projects",
		mcp.WithDescription("(paginated) lists all projects in Devaulty"),
		mcp.WithInputSchema[util.MCPPaginationQuery](),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations),
	)
	s.AddTool(listProjectTool, t.list)

	getProjectTool := mcp.NewTool("get_project",
		mcp.WithDescription("gets a project by ID (uuid)"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("UUID of the project")),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations),
	)
	s.AddTool(getProjectTool, t.get)

	if opts.ReadOnly {
		return
	}

	createProjectTool := mcp.NewTool("create_project",
		mcp.WithDescription("creates a new project in Devaulty"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name of the project"),
			mcp.MinLength(2), mcp.MaxLength(255)),
		mcp.WithString("description", mcp.Description("Description of the project"),
			mcp.MinLength(1), mcp.MaxLength(255)),
		mcp.WithString("icon", mcp.Description("Icon of the project"), mcp.Enum(projectIcons...)),
		mcp.WithString("color", mcp.Description("Color of the project"), mcp.Enum(projectColors...)),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(createProjectTool, t.create)

	updateProjectTool := mcp.NewTool("update_project",
		mcp.WithDescription("updates an existing project in Devaulty"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("UUID of the project")),
		mcp.WithString("name", mcp.Description("Name of the project"),
			mcp.MinLength(2), mcp.MaxLength(255)),
		mcp.WithString("description", mcp.Description("Description of the project"),
			mcp.MinLength(1), mcp.MaxLength(255)),
		mcp.WithString("icon", mcp.Description("Icon of the project"), mcp.Enum(projectIcons...)),
		mcp.WithString("color", mcp.Description("Color of the project"), mcp.Enum(projectColors...)),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(updateProjectTool, t.update)

	archiveProjectTool := mcp.NewTool("archive_project",
		mcp.WithDescription("archives an existing project in Devaulty"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("UUID of the project")),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(archiveProjectTool, t.archive)

	unarchiveProjectTool := mcp.NewTool("unarchive_project",
		mcp.WithDescription("unarchives an existing project in Devaulty"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("UUID of the project")),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(unarchiveProjectTool, t.unarchive)

	if opts.DisableDelete {
		return
	}

	deleteProjectTool := mcp.NewTool("delete_project",
		mcp.WithDescription("deletes an existing project in Devaulty"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("UUID of the project")),
		mcp.WithToolAnnotation(util.DeleteAnnotations))
	s.AddTool(deleteProjectTool, t.delete)
}

func (t *ProjectTools) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.CreateProjectCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	project, err := t.projectUseCase.Create(ctx, cmd)
	if err != nil {
		log.Printf("[ProjectTools.create] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(project)
}

func (t *ProjectTools) list(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var query util.MCPPaginationQuery
	if err := req.BindArguments(&query); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	util.ValidateQuery(&query)

	pagedProjects, err := t.projectUseCase.GetAll(ctx, query.PageNumber, query.PageSize)
	if err != nil {
		log.Printf("[ProjectTools.list] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultJSON(pagedProjects)
}

func (t *ProjectTools) get(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	project, err := t.projectUseCase.GetByID(ctx, projectID)
	if err != nil {
		log.Printf("[ProjectTools.get] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultJSON(project)
}

func (t *ProjectTools) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.UpdateProjectCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	project, err := t.projectUseCase.Update(ctx, cmd)
	if err != nil {
		log.Printf("[ProjectTools.update] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(project)
}

func (t *ProjectTools) archive(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = t.projectUseCase.Archive(ctx, projectID)
	if err != nil {
		log.Printf("[ProjectTools.archive] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Project Archived!"), nil
}

func (t *ProjectTools) unarchive(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = t.projectUseCase.Unarchive(ctx, projectID)
	if err != nil {
		log.Printf("[ProjectTools.unarchive] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Project Unarchived!"), nil
}

func (t *ProjectTools) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = t.projectUseCase.Delete(ctx, projectID)
	if err != nil {
		log.Printf("[ProjectTools.delete] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Project Deleted!"), nil
}
