package mcp

import (
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/usecase"

	"github.com/mark3labs/mcp-go/server"
)

type MCPServerAdapter struct {
	opts           Options
	projectUseCase *usecase.ProjectUseCase
	snippetUseCase *usecase.SnippetUseCase
	problemUseCase *usecase.ProblemUseCase
}

func NewMCPServerAdapter(
	opts Options,
	projectUseCase *usecase.ProjectUseCase,
	snippetUseCase *usecase.SnippetUseCase,
	problemUseCase *usecase.ProblemUseCase,
) *MCPServerAdapter {
	return &MCPServerAdapter{
		opts:           opts,
		projectUseCase: projectUseCase,
		snippetUseCase: snippetUseCase,
		problemUseCase: problemUseCase,
	}
}

func (m *MCPServerAdapter) Serve() error {
	mcpServer := server.NewMCPServer(
		"Devaulty MCP Server",
		model.AppVersion,
		server.WithRecovery(),
		server.WithInstructions("Devaulty MCP Server — provides tools to read and write project data"+
			" (boards, cards, notes, snippets, problems, tags) in the user's local Devaulty app."+
			" Does not have access to the Vault/credentials module"),
		server.WithToolCapabilities(false))

	NewProjectTools(m.projectUseCase).Register(mcpServer, m.opts)
	NewSnippetTools(m.snippetUseCase).Register(mcpServer, m.opts)
	NewProblemTools(m.problemUseCase).Register(mcpServer, m.opts)

	return server.ServeStdio(mcpServer)
}
