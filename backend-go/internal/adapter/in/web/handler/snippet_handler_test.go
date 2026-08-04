package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnippetHandler_Create(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed a project first
	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	t.Run("Create success", func(t *testing.T) {
		snippetBody := []byte(`{
			"title": "Print Hello World",
			"description": "Simple Go print snippet",
			"content": "fmt.Println(\"Hello World\")",
			"language": "GO",
			"snippetType": "CODE"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/snippets", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, snippetBody, true)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("Location"))

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "Print Hello World", result["title"])
		assert.Equal(t, projectID, result["projectId"])
		assert.NotEmpty(t, result["id"])
	})

	t.Run("Create failure - project not found", func(t *testing.T) {
		snippetBody := []byte(`{
			"title": "Valid Title",
			"content": "some content",
			"language": "BASH",
			"snippetType": "COMMAND"
		}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects/00000000-0000-0000-0000-000000000000/snippets", snippetBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Create failure - invalid project UUID format", func(t *testing.T) {
		snippetBody := []byte(`{
			"title": "Valid Title",
			"content": "some content",
			"language": "BASH",
			"snippetType": "COMMAND"
		}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects/invalid-project-id/snippets", snippetBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - missing required fields", func(t *testing.T) {
		snippetBody := []byte(`{"description": "No title or content"}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/snippets", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, snippetBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - title too short", func(t *testing.T) {
		snippetBody := []byte(`{
			"title": "A",
			"content": "some content",
			"language": "BASH",
			"snippetType": "COMMAND"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/snippets", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, snippetBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - unauthorized", func(t *testing.T) {
		snippetBody := []byte(`{
			"title": "Valid Title",
			"content": "some content",
			"language": "BASH",
			"snippetType": "COMMAND"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/snippets", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, snippetBody, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestSnippetHandler_Get(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed a project and a snippet
	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	snippetBody := []byte(`{
		"title": "Seeded Snippet",
		"content": "echo 'hello'",
		"language": "BASH",
		"snippetType": "COMMAND"
	}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/snippets", projectID)
	respSnippet := app.DoRequest(t, http.MethodPost, urlCreate, snippetBody, true)
	var createdSnippet map[string]interface{}
	_ = json.NewDecoder(respSnippet.Body).Decode(&createdSnippet)
	snippetID := createdSnippet["id"].(string)

	t.Run("Get success", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/snippets/%s", projectID, snippetID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, snippetID, result["id"])
		assert.Equal(t, "Seeded Snippet", result["title"])
	})

	t.Run("Get failure - snippet not found", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/snippets/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get failure - invalid project UUID format", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/invalid-project-id/snippets/%s", snippetID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Get failure - invalid snippet UUID format", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/snippets/invalid-snippet-id", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestSnippetHandler_GetAll(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	t.Run("GetAll success - default pagination", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/snippets", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GetAll success - custom page and size", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/snippets?page=0&size=5", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GetAll failure - project not found", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/00000000-0000-0000-0000-000000000000/snippets", nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("GetAll failure - invalid project UUID", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/invalid-id/snippets", nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("GetAll failure - size greater than 100", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/snippets?page=0&size=101", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestSnippetHandler_Update(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	snippetBody := []byte(`{
		"title": "Original Title",
		"content": "Original Content",
		"language": "GO",
		"snippetType": "CODE"
	}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/snippets", projectID)
	respSnippet := app.DoRequest(t, http.MethodPost, urlCreate, snippetBody, true)
	var createdSnippet map[string]interface{}
	_ = json.NewDecoder(respSnippet.Body).Decode(&createdSnippet)
	snippetID := createdSnippet["id"].(string)

	t.Run("Update success", func(t *testing.T) {
		updateBody := []byte(`{"title": "Updated Title"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/snippets/%s", projectID, snippetID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "Updated Title", result["title"])
	})

	t.Run("Update failure - snippet not found", func(t *testing.T) {
		updateBody := []byte(`{"title": "Updated Title"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/snippets/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Update failure - project not found", func(t *testing.T) {
		updateBody := []byte(`{"title": "Updated Title"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/00000000-0000-0000-0000-000000000000/snippets/%s", snippetID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Update failure - invalid snippet UUID format", func(t *testing.T) {
		updateBody := []byte(`{"title": "Updated Title"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/snippets/invalid-snippet-id", projectID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestSnippetHandler_Delete(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	snippetBody := []byte(`{
		"title": "Snippet To Delete",
		"content": "rm -rf /tmp",
		"language": "BASH",
		"snippetType": "COMMAND"
	}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/snippets", projectID)
	respSnippet := app.DoRequest(t, http.MethodPost, urlCreate, snippetBody, true)
	var createdSnippet map[string]interface{}
	_ = json.NewDecoder(respSnippet.Body).Decode(&createdSnippet)
	snippetID := createdSnippet["id"].(string)

	t.Run("Delete success", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/%s/snippets/%s", projectID, snippetID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("Delete failure - invalid snippet UUID", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/%s/snippets/invalid-snippet-id", projectID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Delete failure - invalid project UUID", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/invalid-project-id/snippets/%s", snippetID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
