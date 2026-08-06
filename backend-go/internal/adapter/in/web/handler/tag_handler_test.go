package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagHandler_Create(t *testing.T) {
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
		tagBody := []byte(`{
			"name": "Backend",
			"color": "#FF0000"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/tags", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, tagBody, true)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("Location"))

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "Backend", result["name"])
		assert.Equal(t, "#FF0000", result["color"])
		assert.Equal(t, projectID, result["projectId"])
		assert.NotEmpty(t, result["id"])
	})

	t.Run("Create failure - tag name already exists", func(t *testing.T) {
		tagBody := []byte(`{
			"name": "Backend",
			"color": "#00FF00"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/tags", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, tagBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - project not found", func(t *testing.T) {
		tagBody := []byte(`{
			"name": "Frontend",
			"color": "#0000FF"
		}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects/00000000-0000-0000-0000-000000000000/tags", tagBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Create failure - invalid project UUID format", func(t *testing.T) {
		tagBody := []byte(`{
			"name": "Frontend",
			"color": "#0000FF"
		}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects/invalid-project-id/tags", tagBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - name too short", func(t *testing.T) {
		tagBody := []byte(`{
			"name": "A",
			"color": "#0000FF"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/tags", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, tagBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - invalid hex color", func(t *testing.T) {
		tagBody := []byte(`{
			"name": "Database",
			"color": "invalid-color"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/tags", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, tagBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - unauthorized", func(t *testing.T) {
		tagBody := []byte(`{
			"name": "DevOps",
			"color": "#123456"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/tags", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, tagBody, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestTagHandler_Get(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed a project and a tag
	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	tagBody := []byte(`{
		"name": "Seeded Tag",
		"color": "#112233"
	}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/tags", projectID)
	respTag := app.DoRequest(t, http.MethodPost, urlCreate, tagBody, true)
	var createdTag map[string]interface{}
	_ = json.NewDecoder(respTag.Body).Decode(&createdTag)
	tagID := createdTag["id"].(string)

	t.Run("Get success", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/tags/%s", projectID, tagID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, tagID, result["id"])
		assert.Equal(t, "Seeded Tag", result["name"])
		assert.Equal(t, "#112233", result["color"])
	})

	t.Run("Get failure - tag not found", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/tags/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get failure - project not found", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/00000000-0000-0000-0000-000000000000/tags/%s", tagID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get failure - invalid project UUID format", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/invalid-project-id/tags/%s", tagID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Get failure - invalid tag UUID format", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/tags/invalid-tag-id", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestTagHandler_GetAll(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Create Project A
	respProjA := app.DoRequest(t, http.MethodPost, "/api/v1/projects", []byte(`{"name":"Project A"}`), true)
	var projA map[string]interface{}
	_ = json.NewDecoder(respProjA.Body).Decode(&projA)
	projectAID := projA["id"].(string)

	// Create Project B (isolation check)
	respProjB := app.DoRequest(t, http.MethodPost, "/api/v1/projects", []byte(`{"name":"Project B"}`), true)
	var projB map[string]interface{}
	_ = json.NewDecoder(respProjB.Body).Decode(&projB)
	projectBID := projB["id"].(string)

	// Seed 3 tags in Project A
	for i := 1; i <= 3; i++ {
		tagBody := []byte(fmt.Sprintf(`{"name": "Tag A%d", "color": "#00000%d"}`, i, i))
		urlCreate := fmt.Sprintf("/api/v1/projects/%s/tags", projectAID)
		_ = app.DoRequest(t, http.MethodPost, urlCreate, tagBody, true)
	}

	// Seed 2 tags in Project B
	for i := 1; i <= 2; i++ {
		tagBody := []byte(fmt.Sprintf(`{"name": "Tag B%d", "color": "#00000%d"}`, i, i))
		urlCreate := fmt.Sprintf("/api/v1/projects/%s/tags", projectBID)
		_ = app.DoRequest(t, http.MethodPost, urlCreate, tagBody, true)
	}

	t.Run("GetAll success", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/tags", projectAID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var tags []map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&tags)
		require.NoError(t, err)

		assert.Len(t, tags, 3)
		for _, tagMap := range tags {
			assert.Equal(t, projectAID, tagMap["projectId"])
		}
	})

	t.Run("GetAll failure - project not found", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/00000000-0000-0000-0000-000000000000/tags", nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("GetAll failure - invalid project UUID", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/invalid-id/tags", nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestTagHandler_SearchByName(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	// Seed tags
	tagNames := []string{"Golang", "Go-Fiber", "Python", "Rust"}
	for _, name := range tagNames {
		tagBody := []byte(fmt.Sprintf(`{"name": "%s", "color": "#123456"}`, name))
		urlCreate := fmt.Sprintf("/api/v1/projects/%s/tags", projectID)
		_ = app.DoRequest(t, http.MethodPost, urlCreate, tagBody, true)
	}

	t.Run("SearchByName success - matching term", func(t *testing.T) {
		urlSearch := fmt.Sprintf("/api/v1/projects/%s/tags/search?tag_name=go", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlSearch, nil, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var tags []map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&tags)
		require.NoError(t, err)

		assert.Len(t, tags, 2) // Golang, Go-Fiber
	})

	t.Run("SearchByName failure - project not found", func(t *testing.T) {
		urlSearch := "/api/v1/projects/00000000-0000-0000-0000-000000000000/tags/search?tag_name=go"
		resp := app.DoRequest(t, http.MethodGet, urlSearch, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("SearchByName failure - invalid project UUID", func(t *testing.T) {
		urlSearch := "/api/v1/projects/invalid-id/tags/search?tag_name=go"
		resp := app.DoRequest(t, http.MethodGet, urlSearch, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestTagHandler_Update(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	tagBody := []byte(`{"name": "Original Name", "color": "#111111"}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/tags", projectID)
	respTag := app.DoRequest(t, http.MethodPost, urlCreate, tagBody, true)
	var createdTag map[string]interface{}
	_ = json.NewDecoder(respTag.Body).Decode(&createdTag)
	tagID := createdTag["id"].(string)

	// Create another tag to test name collision
	_ = app.DoRequest(t, http.MethodPost, urlCreate, []byte(`{"name": "Existing Tag", "color": "#222222"}`), true)

	t.Run("Update success", func(t *testing.T) {
		updateBody := []byte(`{"name": "Updated Name", "color": "#333333"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/tags/%s", projectID, tagID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "Updated Name", result["name"])
		assert.Equal(t, "#333333", result["color"])
	})

	t.Run("Update failure - name collision", func(t *testing.T) {
		updateBody := []byte(`{"name": "Existing Tag"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/tags/%s", projectID, tagID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Update failure - tag not found", func(t *testing.T) {
		updateBody := []byte(`{"name": "New Name"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/tags/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Update failure - project not found", func(t *testing.T) {
		updateBody := []byte(`{"name": "New Name"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/00000000-0000-0000-0000-000000000000/tags/%s", tagID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Update failure - invalid tag UUID format", func(t *testing.T) {
		updateBody := []byte(`{"name": "New Name"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/tags/invalid-tag-id", projectID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Update failure - invalid project UUID format", func(t *testing.T) {
		updateBody := []byte(`{"name": "New Name"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/invalid-project-id/tags/%s", tagID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestTagHandler_Delete(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	tagBody := []byte(`{"name": "Tag To Delete", "color": "#000000"}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/tags", projectID)
	respTag := app.DoRequest(t, http.MethodPost, urlCreate, tagBody, true)
	var createdTag map[string]interface{}
	_ = json.NewDecoder(respTag.Body).Decode(&createdTag)
	tagID := createdTag["id"].(string)

	t.Run("Delete success", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/%s/tags/%s", projectID, tagID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Assert persistence: GET request should return 404
		respGet := app.DoRequest(t, http.MethodGet, urlDelete, nil, true)
		assert.Equal(t, http.StatusNotFound, respGet.StatusCode)
	})

	t.Run("Delete failure - tag not found", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/%s/tags/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Delete failure - invalid tag UUID", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/%s/tags/invalid-tag-id", projectID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Delete failure - invalid project UUID", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/invalid-project-id/tags/%s", tagID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
