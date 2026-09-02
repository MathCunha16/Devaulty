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

type ProblemTools struct {
	problemUseCase *usecase.ProblemUseCase
}

func NewProblemTools(problemUseCase *usecase.ProblemUseCase) *ProblemTools {
	return &ProblemTools{problemUseCase: problemUseCase}
}

func (t *ProblemTools) Register(s *server.MCPServer, opts Options) {

	listProblemTool := mcp.NewTool("list_problems",
		mcp.WithDescription("(paginated) (first page at 0) lists all problems in a Devaulty project. Returns a summary view of the problems"),
		mcp.WithInputSchema[util.ProjectPaginationQuery](),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(listProblemTool, t.list)

	getProblemTool := mcp.NewTool("get_problem",
		mcp.WithDescription("gets a problem by ID (uuid)"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the problem")),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(getProblemTool, t.get)

	if opts.ReadOnly {
		return
	}

	createProblemTool := mcp.NewTool("create_problem",
		mcp.WithDescription("creates a new problem in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project where you want to create a new problem")),
		mcp.WithString("title", mcp.Required(), mcp.Description("The title of the problem")),
		mcp.WithString("errorDescription", mcp.Description("The description and/or the stack trace of the problem"), mcp.Required()),
		mcp.WithString("solution", mcp.Description("The solution to the problem or ways to fix it")),
		mcp.WithString("status", mcp.Description("The status of the problem"), mcp.Enum(util.ToStrings(model.ProblemStatuses)...), mcp.Required()),
		mcp.WithString("severity", mcp.Description("The severity of the problem"), mcp.Enum(util.ToStrings(model.ProblemSeverities)...), mcp.Required()),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(createProblemTool, t.create)

	updateProblemTool := mcp.NewTool("update_problem",
		mcp.WithDescription("updates an existing problem in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the problem")),
		mcp.WithString("title", mcp.Description("The title of the problem")),
		mcp.WithString("errorDescription", mcp.Description("The description and/or the stack trace of the problem")),
		mcp.WithString("solution", mcp.Description("The solution to the problem or ways to fix it")),
		mcp.WithString("severity", mcp.Description("The severity of the problem"), mcp.Enum(util.ToStrings(model.ProblemSeverities)...)),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(updateProblemTool, t.update)

	updateProblemStatusTool := mcp.NewTool("update_problem_status",
		mcp.WithDescription("updates the status of an existing problem in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the problem")),
		mcp.WithString("status", mcp.Description("The status of the problem"), mcp.Enum(util.ToStrings(model.ProblemStatuses)...)),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(updateProblemStatusTool, t.updateStatus)

	if opts.DisableDelete {
		return
	}

	deleteProblemTool := mcp.NewTool("delete_problem",
		mcp.WithDescription("deletes an existing problem in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the problem")),
		mcp.WithToolAnnotation(util.DeleteAnnotations))
	s.AddTool(deleteProblemTool, t.delete)

}

func (t *ProblemTools) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.CreateProblemCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	problem, err := t.problemUseCase.Create(ctx, cmd)
	if err != nil {
		log.Printf("[ProblemTools.create] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(problem)
}

func (t *ProblemTools) list(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var query util.ProjectPaginationQuery
	if err := req.BindArguments(&query); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	util.ValidateProjectQuery(&query)

	projectID, err := uuid.Parse(query.ProjectID)
	if err != nil {
		return mcp.NewToolResultError("Invalid project UUID" + err.Error()), nil
	}

	pagedProblems, err := t.problemUseCase.GetAllByProjectID(ctx, projectID, query.PageNumber, query.PageSize)
	if err != nil {
		log.Printf("[ProblemTools.list] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultJSON(pagedProblems)
}

func (t *ProblemTools) get(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	problemID, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	problem, err := t.problemUseCase.GetByID(ctx, projectID, problemID)
	if err != nil {
		log.Printf("[ProblemTools.get] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(problem)
}

func (t *ProblemTools) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.UpdateProblemCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	problem, err := t.problemUseCase.Update(ctx, cmd)
	if err != nil {
		log.Printf("[ProblemTools.update] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(problem)
}

func (t *ProblemTools) updateStatus(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.UpdateProblemStatusCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	problem, err := t.problemUseCase.UpdateStatus(ctx, cmd)
	if err != nil {
		log.Printf("[ProblemTools.updateStatus] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(problem)
}

func (t *ProblemTools) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	problemID, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = t.problemUseCase.Delete(ctx, projectID, problemID)
	if err != nil {
		log.Printf("[ProblemTools.delete] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Problem Deleted!"), nil
}
