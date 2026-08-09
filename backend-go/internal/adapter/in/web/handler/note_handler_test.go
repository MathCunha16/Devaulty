package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoteHandler_Create(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project
	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	t.Run("Create success", func(t *testing.T) {
		noteBody := []byte(`{
			"title": "Meeting Notes",
			"content": "Discussed system architecture and hexagonal pattern."
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/notes", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, noteBody, true)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("Location"))

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "Meeting Notes", result["title"])
		assert.Equal(t, "Discussed system architecture and hexagonal pattern.", result["content"])
		assert.Equal(t, projectID, result["projectId"])
		assert.NotEmpty(t, result["id"])
		assert.Equal(t, false, result["archived"])
	})

	t.Run("Create failure - project not found", func(t *testing.T) {
		noteBody := []byte(`{
			"title": "Orphan Note",
			"content": "Content for non-existing project"
		}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects/00000000-0000-0000-0000-000000000000/notes", noteBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Create failure - invalid project UUID format", func(t *testing.T) {
		noteBody := []byte(`{
			"title": "Invalid UUID Note",
			"content": "Content"
		}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects/invalid-uuid/notes", noteBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - title too short", func(t *testing.T) {
		noteBody := []byte(`{
			"title": "A",
			"content": "Content"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/notes", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, noteBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - missing content", func(t *testing.T) {
		noteBody := []byte(`{
			"title": "Note Without Content"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/notes", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, noteBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - unauthorized", func(t *testing.T) {
		noteBody := []byte(`{
			"title": "Unauthorized Note",
			"content": "Content"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/notes", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, noteBody, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestNoteHandler_Get(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project and note
	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	noteBody := []byte(`{
		"title": "Seeded Note",
		"content": "Seeded Content"
	}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/notes", projectID)
	respNote := app.DoRequest(t, http.MethodPost, urlCreate, noteBody, true)
	var createdNote map[string]interface{}
	_ = json.NewDecoder(respNote.Body).Decode(&createdNote)
	noteID := createdNote["id"].(string)

	t.Run("Get success", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/notes/%s", projectID, noteID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, noteID, result["id"])
		assert.Equal(t, "Seeded Note", result["title"])
		assert.Equal(t, "Seeded Content", result["content"])
	})

	t.Run("Get failure - note not found", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/notes/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get failure - project not found", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/00000000-0000-0000-0000-000000000000/notes/%s", noteID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get failure - invalid project UUID", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/invalid-uuid/notes/%s", noteID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Get failure - invalid note UUID", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/notes/invalid-uuid", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestNoteHandler_GetAll(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	respProjA := app.DoRequest(t, http.MethodPost, "/api/v1/projects", []byte(`{"name":"Project A"}`), true)
	var projA map[string]interface{}
	_ = json.NewDecoder(respProjA.Body).Decode(&projA)
	projectAID := projA["id"].(string)

	respProjB := app.DoRequest(t, http.MethodPost, "/api/v1/projects", []byte(`{"name":"Project B"}`), true)
	var projB map[string]interface{}
	_ = json.NewDecoder(respProjB.Body).Decode(&projB)
	projectBID := projB["id"].(string)

	// Seed 2 notes in Project A
	for i := 1; i <= 2; i++ {
		body := []byte(fmt.Sprintf(`{"title":"Note A%d","content":"Content A%d"}`, i, i))
		urlPath := fmt.Sprintf("/api/v1/projects/%s/notes", projectAID)
		_ = app.DoRequest(t, http.MethodPost, urlPath, body, true)
	}

	t.Run("GetAll success - returns paginated content", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/notes?page=0&size=10", projectAID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var page map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&page)
		require.NoError(t, err)

		content := page["content"].([]interface{})
		assert.Len(t, content, 2)
	})

	t.Run("GetAll success - empty page for project with no notes", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/notes?page=0&size=10", projectBID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var page map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&page)
		require.NoError(t, err)

		content := page["content"].([]interface{})
		assert.Len(t, content, 0)
	})

	t.Run("GetAll failure - project not found", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/00000000-0000-0000-0000-000000000000/notes", nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("GetAll failure - invalid project UUID", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/invalid-uuid/notes", nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestNoteHandler_Update(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	noteBody := []byte(`{"title":"Original Title","content":"Original Content"}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/notes", projectID)
	respNote := app.DoRequest(t, http.MethodPost, urlCreate, noteBody, true)
	var createdNote map[string]interface{}
	_ = json.NewDecoder(respNote.Body).Decode(&createdNote)
	noteID := createdNote["id"].(string)

	t.Run("Update success", func(t *testing.T) {
		updateBody := []byte(`{"title":"Updated Title","content":"Updated Content"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/notes/%s", projectID, noteID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "Updated Title", result["title"])
		assert.Equal(t, "Updated Content", result["content"])
	})

	t.Run("Update failure - note not found", func(t *testing.T) {
		updateBody := []byte(`{"title":"Updated Title"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/notes/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Update failure - project not found", func(t *testing.T) {
		updateBody := []byte(`{"title":"Updated Title"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/00000000-0000-0000-0000-000000000000/notes/%s", noteID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Update failure - title too short", func(t *testing.T) {
		updateBody := []byte(`{"title":"A"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/notes/%s", projectID, noteID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Update failure - invalid note UUID", func(t *testing.T) {
		updateBody := []byte(`{"title":"Updated Title"}`)
		urlUpdate := fmt.Sprintf("/api/v1/projects/%s/notes/invalid-uuid", projectID)
		resp := app.DoRequest(t, http.MethodPatch, urlUpdate, updateBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestNoteHandler_Delete(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectBody := []byte(`{"name":"Parent Project"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	var createdProject map[string]interface{}
	_ = json.NewDecoder(respProject.Body).Decode(&createdProject)
	projectID := createdProject["id"].(string)

	noteBody := []byte(`{"title":"Note To Delete","content":"Content"}`)
	urlCreate := fmt.Sprintf("/api/v1/projects/%s/notes", projectID)
	respNote := app.DoRequest(t, http.MethodPost, urlCreate, noteBody, true)
	var createdNote map[string]interface{}
	_ = json.NewDecoder(respNote.Body).Decode(&createdNote)
	noteID := createdNote["id"].(string)

	t.Run("Delete success", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/%s/notes/%s", projectID, noteID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Assert persistence: GET request should return 404
		respGet := app.DoRequest(t, http.MethodGet, urlDelete, nil, true)
		assert.Equal(t, http.StatusNotFound, respGet.StatusCode)
	})

	t.Run("Delete failure - note not found", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/%s/notes/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Delete failure - project not found", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/00000000-0000-0000-0000-000000000000/notes/%s", noteID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Delete failure - invalid note UUID", func(t *testing.T) {
		urlDelete := fmt.Sprintf("/api/v1/projects/%s/notes/invalid-uuid", projectID)
		resp := app.DoRequest(t, http.MethodDelete, urlDelete, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
