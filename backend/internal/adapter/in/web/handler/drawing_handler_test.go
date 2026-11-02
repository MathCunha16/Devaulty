package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDrawingHandler_Create(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed parent project
	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	t.Run("Create success", func(t *testing.T) {
		drawingBody := []byte(`{
			"title": "Architecture Overview"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/excalidraw-drawings", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, drawingBody, true)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("Location"))

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "Architecture Overview", result["title"])
		assert.Equal(t, projectID, result["projectId"])
		assert.NotEmpty(t, result["id"])
		assert.Nil(t, result["content"])
		assert.Nil(t, result["thumbnail"])
	})

	t.Run("Create failure - project not found", func(t *testing.T) {
		drawingBody := []byte(`{"title":"Orphan Drawing"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects/00000000-0000-0000-0000-000000000000/excalidraw-drawings", drawingBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Create failure - missing title", func(t *testing.T) {
		drawingBody := []byte(`{}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/excalidraw-drawings", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, drawingBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestDrawingHandler_GetAll(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project
	projectBody := []byte(`{"name":"Project With Drawings"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	// Create two drawings
	for i := 1; i <= 2; i++ {
		body := []byte(fmt.Sprintf(`{"title": "Drawing %d"}`, i))
		resp := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/excalidraw-drawings", projectID), body, true)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	}

	t.Run("GetAll success with lean payload", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/excalidraw-drawings?page_number=0&page_size=10", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var page map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&page)
		require.NoError(t, err)

		content := page["content"].([]interface{})
		assert.Len(t, content, 2)
		assert.Equal(t, float64(2), page["totalElements"])

		firstItem := content[0].(map[string]interface{})
		assert.NotEmpty(t, firstItem["title"])
		assert.NotEmpty(t, firstItem["id"])
		// Verify lean payload: 'content' must NOT be included in list response
		_, hasContent := firstItem["content"]
		assert.False(t, hasContent, "List endpoint must not return drawing content")
	})

	t.Run("GetAll failure - project not found", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/00000000-0000-0000-0000-000000000000/excalidraw-drawings", nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestDrawingHandler_Get(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", []byte(`{"name":"Project"}`), true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)
	var proj map[string]interface{}
	require.NoError(t, json.NewDecoder(respProject.Body).Decode(&proj))
	projectID := proj["id"].(string)

	// Create drawing
	respCreate := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/excalidraw-drawings", projectID), []byte(`{"title":"Detailed Drawing"}`), true)
	require.Equal(t, http.StatusCreated, respCreate.StatusCode)
	var createdDrawing map[string]interface{}
	require.NoError(t, json.NewDecoder(respCreate.Body).Decode(&createdDrawing))
	drawingID := createdDrawing["id"].(string)

	t.Run("Get success", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/excalidraw-drawings/%s", projectID, drawingID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Equal(t, drawingID, result["id"])
		assert.Equal(t, "Detailed Drawing", result["title"])
		assert.Equal(t, projectID, result["projectId"])
	})

	t.Run("Get failure - drawing not found", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/excalidraw-drawings/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestDrawingHandler_Update(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project & drawing
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", []byte(`{"name":"Project"}`), true)
	var proj map[string]interface{}
	require.NoError(t, json.NewDecoder(respProject.Body).Decode(&proj))
	projectID := proj["id"].(string)

	respCreate := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/excalidraw-drawings", projectID), []byte(`{"title":"Original Title"}`), true)
	var created map[string]interface{}
	require.NoError(t, json.NewDecoder(respCreate.Body).Decode(&created))
	drawingID := created["id"].(string)

	t.Run("Update canvas and title success", func(t *testing.T) {
		updateBody := []byte(`{
			"title": "Updated Architecture",
			"content": "{\"elements\":[{\"type\":\"rectangle\"}]}",
			"thumbnail": "data:image/svg+xml;base64,PHN2Zz4="
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/excalidraw-drawings/%s", projectID, drawingID)
		resp := app.DoRequest(t, http.MethodPatch, urlPath, updateBody, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updated map[string]interface{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&updated))
		assert.Equal(t, "Updated Architecture", updated["title"])
		assert.Equal(t, "{\"elements\":[{\"type\":\"rectangle\"}]}", updated["content"])
		assert.Equal(t, "data:image/svg+xml;base64,PHN2Zz4=", updated["thumbnail"])
	})

	t.Run("Update failure - not found", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/excalidraw-drawings/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodPatch, urlPath, []byte(`{"title":"New"}`), true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestDrawingHandler_Delete(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project & drawing
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", []byte(`{"name":"Project"}`), true)
	var proj map[string]interface{}
	require.NoError(t, json.NewDecoder(respProject.Body).Decode(&proj))
	projectID := proj["id"].(string)

	respCreate := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/excalidraw-drawings", projectID), []byte(`{"title":"To Delete"}`), true)
	var created map[string]interface{}
	require.NoError(t, json.NewDecoder(respCreate.Body).Decode(&created))
	drawingID := created["id"].(string)

	t.Run("Delete success", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/excalidraw-drawings/%s", projectID, drawingID)
		resp := app.DoRequest(t, http.MethodDelete, urlPath, nil, true)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify it no longer exists
		respGet := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusNotFound, respGet.StatusCode)
	})

	t.Run("Delete failure - not found", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/excalidraw-drawings/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodDelete, urlPath, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
