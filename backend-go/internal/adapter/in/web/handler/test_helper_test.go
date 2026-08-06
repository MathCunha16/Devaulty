package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"devaulty-backend/internal/adapter/in/web"
	"devaulty-backend/internal/adapter/in/web/handler"
	"devaulty-backend/internal/adapter/out/persistence"
	"devaulty-backend/internal/usecase"

	"github.com/stretchr/testify/require"
)

type TestApp struct {
	Server *httptest.Server
	Token  string
}

func SetupTestApp(t *testing.T) *TestApp {
	db, err := persistence.InitDB(":memory:", "../../../../../migrations")
	require.NoError(t, err)

	itemTagRepo := persistence.NewItemTagRepository(db)

	projectRepo := persistence.NewProjectRepository(db)
	projectUseCase := usecase.NewProjectUseCase(projectRepo)
	projectHandler := handler.NewProjectHandler(projectUseCase)

	snippetRepo := persistence.NewSnippetRepository(db)
	snippetUseCase := usecase.NewSnippetUseCase(snippetRepo, projectRepo, itemTagRepo)
	snippetHandler := handler.NewSnippetHandler(snippetUseCase)

	linkRepo := persistence.NewLinkRepository(db)
	linkUseCase := usecase.NewLinkUseCase(linkRepo, projectRepo, itemTagRepo)
	linkHandler := handler.NewLinkHandler(linkUseCase)

	problemRepo := persistence.NewProblemRepository(db)
	problemUseCase := usecase.NewProblemUseCase(problemRepo, projectRepo, itemTagRepo)
	problemHandler := handler.NewProblemHandler(problemUseCase)

	handlers := &web.Handlers{
		Project: projectHandler,
		Snippet: snippetHandler,
		Link:    linkHandler,
		Problem: problemHandler,
	}

	token := "test-internal-token-12345"
	router := web.SetupRouter(handlers, token)

	ts := httptest.NewServer(router)

	return &TestApp{
		Server: ts,
		Token:  token,
	}
}

func (app *TestApp) DoRequest(t *testing.T, method, path string, body []byte, sendToken bool) *http.Response {
	var bodyReader *bytes.Buffer
	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	} else {
		bodyReader = bytes.NewBuffer([]byte{})
	}

	req, err := http.NewRequest(method, app.Server.URL+path, bodyReader)
	require.NoError(t, err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if sendToken {
		req.Header.Set("DEVAULTY_INTERNAL_TOKEN", app.Token)
	}

	resp, err := app.Server.Client().Do(req)
	require.NoError(t, err)
	return resp
}
