package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkHandler_Create(t *testing.T) {
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
		linkBody := []byte(`{
			"title": "Go Documentation",
			"url": "https://go.dev/doc",
			"description": "Official Go Documentation"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/links", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, linkBody, true)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("Location"))

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "Go Documentation", result["title"])
		assert.Equal(t, "https://go.dev/doc", result["url"])
		assert.Equal(t, "Official Go Documentation", result["description"])
		assert.Equal(t, projectID, result["projectId"])
		assert.NotEmpty(t, result["id"])
	})

	t.Run("Create failure - project not found", func(t *testing.T) {
		linkBody := []byte(`{
			"title": "Go Documentation",
			"url": "https://go.dev/doc"
		}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects/00000000-0000-0000-0000-000000000000/links", linkBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Create failure - invalid project UUID format", func(t *testing.T) {
		linkBody := []byte(`{
			"title": "Go Documentation",
			"url": "https://go.dev/doc"
		}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects/invalid-project-id/links", linkBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - missing required fields", func(t *testing.T) {
		linkBody := []byte(`{"description": "No title or url"}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/links", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, linkBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - title too short", func(t *testing.T) {
		linkBody := []byte(`{
			"title": "A",
			"url": "https://go.dev/doc"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/links", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, linkBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - invalid URL format", func(t *testing.T) {
		linkBody := []byte(`{
			"title": "Invalid Link",
			"url": "not-a-valid-url"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/links", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, linkBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - unauthorized", func(t *testing.T) {
		linkBody := []byte(`{
			"title": "Go Documentation",
			"url": "https://go.dev/doc"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/links", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, linkBody, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestLinkHandler_Get(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed a project and a link
	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	linkBody := []byte(`{
		"title": "Seeded Link",
		"url": "https://example.com"
	}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/links", projectID)
	respLink := app.DoRequest(t, http.MethodPost, urlCreate, linkBody, true)
	var createdLink map[string]interface{}
	_ = json.NewDecoder(respLink.Body).Decode(&createdLink)
	linkID := createdLink["id"].(string)

	t.Run("Get success", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/links/%s", projectID, linkID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, linkID, result["id"])
		assert.Equal(t, "Seeded Link", result["title"])
		assert.Equal(t, "https://example.com", result["url"])
	})

	t.Run("Get failure - link not found", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/links/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get failure - project not found", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/00000000-0000-0000-0000-000000000000/links/%s", linkID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get failure - invalid project UUID format", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/invalid-project-id/links/%s", linkID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Get failure - invalid link UUID format", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/links/invalid-link-id", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestLinkHandler_GetAll(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Create Project A
	respProjA := app.DoRequest(t, http.MethodPost, "/api/v1/projects", []byte(`{"name":"Project A"}`), true)
	var projA map[string]interface{}
	_ = json.NewDecoder(respProjA.Body).Decode(&projA)
	projectAID := projA["id"].(string)

	// Create Project B (for isolation checks)
	respProjB := app.DoRequest(t, http.MethodPost, "/api/v1/projects", []byte(`{"name":"Project B"}`), true)
	var projB map[string]interface{}
	_ = json.NewDecoder(respProjB.Body).Decode(&projB)
	projectBID := projB["id"].(string)

	// Seed 12 links in Project A
	for i := 1; i <= 12; i++ {
		linkBody := []byte(fmt.Sprintf(`{
			"title": "Project A Link %d",
			"url": "https://example.com/a/%d"
		}`, i, i))
		urlCreate := fmt.Sprintf("/api/v1/projects/%s/links", projectAID)
		_ = app.DoRequest(t, http.MethodPost, urlCreate, linkBody, true)
	}

	// Seed 3 links in Project B
	for i := 1; i <= 3; i++ {
		linkBody := []byte(fmt.Sprintf(`{
			"title": "Project B Link %d",
			"url": "https://example.com/b/%d"
		}`, i, i))
		urlCreate := fmt.Sprintf("/api/v1/projects/%s/links", projectBID)
		_ = app.DoRequest(t, http.MethodPost, urlCreate, linkBody, true)
	}

	t.Run("GetAll success - default pagination", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/links", projectAID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var page map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&page)
		require.NoError(t, err)

		content := page["content"].([]interface{})
		assert.Len(t, content, 10) // default page size is 10
		assert.Equal(t, float64(0), page["number"])
		assert.Equal(t, float64(10), page["size"])
		assert.Equal(t, float64(12), page["totalElements"])
		assert.Equal(t, float64(2), page["totalPages"])

		// Ensure project isolation: all items belong to projectAID
		for _, item := range content {
			linkMap := item.(map[string]interface{})
			assert.Equal(t, projectAID, linkMap["projectId"])
		}
	})

	t.Run("GetAll success - custom page and size", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/links?page=1&size=5", projectAID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var page map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&page)
		require.NoError(t, err)

		content := page["content"].([]interface{})
		assert.Len(t, content, 5) // page 1 with size 5
		assert.Equal(t, float64(1), page["number"])
		assert.Equal(t, float64(5), page["size"])
		assert.Equal(t, float64(12), page["totalElements"])
		assert.Equal(t, float64(3), page["totalPages"])

		for _, item := range content {
			linkMap := item.(map[string]interface{})
			assert.Equal(t, projectAID, linkMap["projectId"])
		}
	})

	t.Run("GetAll failure - project not found", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/00000000-0000-0000-0000-000000000000/links", nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("GetAll failure - invalid project UUID", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/invalid-id/links", nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("GetAll failure - size greater than 100", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/links?page=0&size=101", projectAID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestLinkHandler_Update(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	linkBody := []byte(`{
		"title": "Original Title",
		"url": "https://original.com"
	}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/links", projectID)
	respLink := app.DoRequest(t, http.MethodPost, urlCreate, linkBody, true)
	var createdLink map[string]interface{}
	_ = json.NewDecoder(respLink.Body).Decode(&createdLink)
	linkID := createdLink["id"].(string)

	t.Run("Update success", func(t *testing.T) {
		updateBody := []byte(`{
			"title": "Updated Title",
			"url": "https://updated.com"
		}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/links/%s", projectID, linkID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "Updated Title", result["title"])
		assert.Equal(t, "https://updated.com", result["url"])
	})

	t.Run("Update failure - link not found", func(t *testing.T) {
		updateBody := []byte(`{"title": "Updated Title"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/links/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Update failure - project not found", func(t *testing.T) {
		updateBody := []byte(`{"title": "Updated Title"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/00000000-0000-0000-0000-000000000000/links/%s", linkID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Update failure - invalid link UUID format", func(t *testing.T) {
		updateBody := []byte(`{"title": "Updated Title"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/links/invalid-link-id", projectID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Update failure - invalid project UUID format", func(t *testing.T) {
		updateBody := []byte(`{"title": "Updated Title"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/invalid-project-id/links/%s", linkID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Update failure - invalid URL format", func(t *testing.T) {
		updateBody := []byte(`{"url": "not-a-valid-url"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/links/%s", projectID, linkID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestLinkHandler_Delete(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	linkBody := []byte(`{
		"title": "Link To Delete",
		"url": "https://delete-me.com"
	}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/links", projectID)
	respLink := app.DoRequest(t, http.MethodPost, urlCreate, linkBody, true)
	var createdLink map[string]interface{}
	_ = json.NewDecoder(respLink.Body).Decode(&createdLink)
	linkID := createdLink["id"].(string)

	t.Run("Delete success", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/%s/links/%s", projectID, linkID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Assert persistence: issuing a GET request should now return 404 Not Found
		respGet := app.DoRequest(t, http.MethodGet, urlDelete, nil, true)
		assert.Equal(t, http.StatusNotFound, respGet.StatusCode)
	})

	t.Run("Delete failure - link not found", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/%s/links/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Delete failure - invalid link UUID", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/%s/links/invalid-link-id", projectID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Delete failure - invalid project UUID", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/invalid-project-id/links/%s", linkID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
