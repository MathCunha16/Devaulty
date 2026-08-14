package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProblemHandler_Create(t *testing.T) {
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
		problemBody := []byte(`{
			"title": "Database Connection Timeout",
			"errorDescription": "Failed to connect to database at localhost:5432 after 30s",
			"solution": "Restart PostgreSQL container",
			"status": "OPEN",
			"severity": "HIGH"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/problems", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, problemBody, true)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("Location"))

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "Database Connection Timeout", result["title"])
		assert.Equal(t, "Failed to connect to database at localhost:5432 after 30s", result["errorDescription"])
		assert.Equal(t, "Restart PostgreSQL container", result["solution"])
		assert.Equal(t, "OPEN", result["status"])
		assert.Equal(t, "HIGH", result["severity"])
		assert.Equal(t, projectID, result["projectId"])
		assert.NotEmpty(t, result["id"])
	})

	t.Run("Create failure - project not found", func(t *testing.T) {
		problemBody := []byte(`{
			"title": "Valid Title",
			"errorDescription": "Valid Error Description",
			"status": "OPEN",
			"severity": "LOW"
		}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects/00000000-0000-0000-0000-000000000000/problems", problemBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Create failure - invalid project UUID format", func(t *testing.T) {
		problemBody := []byte(`{
			"title": "Valid Title",
			"errorDescription": "Valid Error Description",
			"status": "OPEN",
			"severity": "LOW"
		}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects/invalid-project-id/problems", problemBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - missing required fields", func(t *testing.T) {
		problemBody := []byte(`{"title": "Only Title"}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/problems", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, problemBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - title too short", func(t *testing.T) {
		problemBody := []byte(`{
			"title": "A",
			"errorDescription": "Valid Error Description",
			"status": "OPEN",
			"severity": "LOW"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/problems", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, problemBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - unauthorized", func(t *testing.T) {
		problemBody := []byte(`{
			"title": "Valid Title",
			"errorDescription": "Valid Error Description",
			"status": "OPEN",
			"severity": "LOW"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/problems", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, problemBody, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestProblemHandler_Get(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed a project and a problem
	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	problemBody := []byte(`{
		"title": "Seeded Problem",
		"errorDescription": "Seeded Error Description",
		"status": "OPEN",
		"severity": "MEDIUM"
	}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/problems", projectID)
	respProblem := app.DoRequest(t, http.MethodPost, urlCreate, problemBody, true)
	var createdProblem map[string]interface{}
	_ = json.NewDecoder(respProblem.Body).Decode(&createdProblem)
	problemID := createdProblem["id"].(string)

	t.Run("Get success", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/problems/%s", projectID, problemID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, problemID, result["id"])
		assert.Equal(t, "Seeded Problem", result["title"])
		assert.Equal(t, "OPEN", result["status"])
		assert.Equal(t, "MEDIUM", result["severity"])
	})

	t.Run("Get failure - problem not found", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/problems/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get failure - project not found", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/00000000-0000-0000-0000-000000000000/problems/%s", problemID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get failure - invalid project UUID format", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/invalid-project-id/problems/%s", problemID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Get failure - invalid problem UUID format", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/problems/invalid-problem-id", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestProblemHandler_GetAll(t *testing.T) {
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

	// Seed 12 problems in Project A
	for i := 1; i <= 12; i++ {
		problemBody := []byte(fmt.Sprintf(`{
			"title": "Project A Problem %d",
			"errorDescription": "Error description %d",
			"status": "OPEN",
			"severity": "LOW"
		}`, i, i))
		urlCreate := fmt.Sprintf("/api/v1/projects/%s/problems", projectAID)
		_ = app.DoRequest(t, http.MethodPost, urlCreate, problemBody, true)
	}

	// Seed 3 problems in Project B
	for i := 1; i <= 3; i++ {
		problemBody := []byte(fmt.Sprintf(`{
			"title": "Project B Problem %d",
			"errorDescription": "Error description %d",
			"status": "OPEN",
			"severity": "HIGH"
		}`, i, i))
		urlCreate := fmt.Sprintf("/api/v1/projects/%s/problems", projectBID)
		_ = app.DoRequest(t, http.MethodPost, urlCreate, problemBody, true)
	}

	t.Run("GetAll success - default pagination", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/problems", projectAID)
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
			problemMap := item.(map[string]interface{})
			assert.Equal(t, projectAID, problemMap["projectId"])
		}
	})

	t.Run("GetAll success - custom page and size", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/problems?page=1&size=5", projectAID)
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
			problemMap := item.(map[string]interface{})
			assert.Equal(t, projectAID, problemMap["projectId"])
		}
	})

	t.Run("GetAll failure - project not found", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/00000000-0000-0000-0000-000000000000/problems", nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("GetAll failure - invalid project UUID", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/invalid-id/problems", nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("GetAll failure - size greater than 100", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/problems?page=0&size=101", projectAID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestProblemHandler_Update(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	problemBody := []byte(`{
		"title": "Original Title",
		"errorDescription": "Original Error",
		"status": "OPEN",
		"severity": "LOW"
	}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/problems", projectID)
	respProblem := app.DoRequest(t, http.MethodPost, urlCreate, problemBody, true)
	var createdProblem map[string]interface{}
	_ = json.NewDecoder(respProblem.Body).Decode(&createdProblem)
	problemID := createdProblem["id"].(string)

	t.Run("Update success", func(t *testing.T) {
		updateBody := []byte(`{
			"title": "Updated Title",
			"severity": "CRITICAL"
		}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/problems/%s", projectID, problemID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "Updated Title", result["title"])
		assert.Equal(t, "CRITICAL", result["severity"])
	})

	t.Run("Update failure - problem not found", func(t *testing.T) {
		updateBody := []byte(`{"title": "Updated Title"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/problems/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Update failure - project not found", func(t *testing.T) {
		updateBody := []byte(`{"title": "Updated Title"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/00000000-0000-0000-0000-000000000000/problems/%s", problemID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Update failure - invalid problem UUID format", func(t *testing.T) {
		updateBody := []byte(`{"title": "Updated Title"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/problems/invalid-problem-id", projectID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Update failure - invalid project UUID format", func(t *testing.T) {
		updateBody := []byte(`{"title": "Updated Title"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/invalid-project-id/problems/%s", problemID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestProblemHandler_UpdateStatus(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	problemBody := []byte(`{
		"title": "Problem To Update Status",
		"errorDescription": "Error description",
		"status": "OPEN",
		"severity": "MEDIUM"
	}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/problems", projectID)
	respProblem := app.DoRequest(t, http.MethodPost, urlCreate, problemBody, true)
	var createdProblem map[string]interface{}
	_ = json.NewDecoder(respProblem.Body).Decode(&createdProblem)
	problemID := createdProblem["id"].(string)

	t.Run("UpdateStatus success", func(t *testing.T) {
		updateBody := []byte(`{"status": "RESOLVED"}`)
		urlUpdateStatus := fmt.Sprintf("/api/v1/projects/%s/problems/%s/status", projectID, problemID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdateStatus, updateBody, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "RESOLVED", result["status"])
	})

	t.Run("UpdateStatus failure - problem not found", func(t *testing.T) {
		updateBody := []byte(`{"status": "RESOLVED"}`)
		urlUpdateStatus := fmt.Sprintf("/api/v1/projects/%s/problems/00000000-0000-0000-0000-000000000000/status", projectID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdateStatus, updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("UpdateStatus failure - project not found", func(t *testing.T) {
		updateBody := []byte(`{"status": "RESOLVED"}`)
		urlUpdateStatus := fmt.Sprintf("/api/v1/projects/00000000-0000-0000-0000-000000000000/problems/%s/status", problemID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdateStatus, updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("UpdateStatus failure - missing status field", func(t *testing.T) {
		updateBody := []byte(`{}`)
		urlUpdateStatus := fmt.Sprintf("/api/v1/projects/%s/problems/%s/status", projectID, problemID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdateStatus, updateBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("UpdateStatus failure - invalid problem UUID format", func(t *testing.T) {
		updateBody := []byte(`{"status": "RESOLVED"}`)
		urlUpdateStatus := fmt.Sprintf("/api/v1/projects/%s/problems/invalid-problem-id/status", projectID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdateStatus, updateBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestProblemHandler_Delete(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	problemBody := []byte(`{
		"title": "Problem To Delete",
		"errorDescription": "Error description",
		"status": "OPEN",
		"severity": "LOW"
	}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/problems", projectID)
	respProblem := app.DoRequest(t, http.MethodPost, urlCreate, problemBody, true)
	var createdProblem map[string]interface{}
	_ = json.NewDecoder(respProblem.Body).Decode(&createdProblem)
	problemID := createdProblem["id"].(string)

	t.Run("Delete success", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/%s/problems/%s", projectID, problemID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Assert persistence: issuing a GET request should now return 404 Not Found
		respGet := app.DoRequest(t, http.MethodGet, urlDelete, nil, true)
		assert.Equal(t, http.StatusNotFound, respGet.StatusCode)
	})

	t.Run("Delete failure - problem not found", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/%s/problems/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Delete failure - invalid problem UUID", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/%s/problems/invalid-problem-id", projectID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Delete failure - invalid project UUID", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/invalid-project-id/problems/%s", problemID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
