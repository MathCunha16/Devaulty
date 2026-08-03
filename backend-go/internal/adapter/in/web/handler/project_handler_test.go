package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectHandler_Create(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	t.Run("Create success", func(t *testing.T) {
		body := []byte(`{"name":"New Project","description":"Valid desc","icon":"rocket","color":"#4A90E2"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects", body, true)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("Location"))

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "New Project", result["name"])
		assert.NotEmpty(t, result["id"])
	})

	t.Run("Create failure - missing required name", func(t *testing.T) {
		body := []byte(`{"description":"No name"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects", body, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - name too short", func(t *testing.T) {
		body := []byte(`{"name":"A"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects", body, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - invalid hex color", func(t *testing.T) {
		body := []byte(`{"name":"Valid Name","color":"invalid-color"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects", body, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - unauthorized", func(t *testing.T) {
		body := []byte(`{"name":"Valid Name"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects", body, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestProjectHandler_Get(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	body := []byte(`{"name":"Seeded Project"}`)
	respSeed := app.DoRequest(t, http.MethodPost, "/api/v1/projects", body, true)
	var created map[string]interface{}
	_ = json.NewDecoder(respSeed.Body).Decode(&created)
	projectID := created["id"].(string)

	t.Run("Get success", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/"+projectID, nil, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, projectID, result["id"])
	})

	t.Run("Get failure - invalid UUID format", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/not-a-uuid", nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Get failure - not found", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/00000000-0000-0000-0000-000000000000", nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestProjectHandler_GetAll(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	t.Run("GetAll success - default pagination", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects", nil, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GetAll success - custom query page and size", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects?page=0&size=5", nil, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GetAll failure - page less than 0", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects?page=-1&size=10", nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("GetAll failure - size less than 1", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects?page=0&size=0", nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("GetAll failure - size greater than 100", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects?page=0&size=101", nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestProjectHandler_Update(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	body := []byte(`{"name":"Original Name"}`)
	respSeed := app.DoRequest(t, http.MethodPost, "/api/v1/projects", body, true)
	var created map[string]interface{}
	_ = json.NewDecoder(respSeed.Body).Decode(&created)
	projectID := created["id"].(string)

	t.Run("Update success", func(t *testing.T) {
		updateBody := []byte(`{"name":"Updated Name"}`)
		resp := app.DoRequest(t, http.MethodPatch, "/api/v1/projects/"+projectID, updateBody, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "Updated Name", result["name"])
	})

	t.Run("Update failure - invalid UUID format", func(t *testing.T) {
		updateBody := []byte(`{"name":"Updated Name"}`)
		resp := app.DoRequest(t, http.MethodPatch, "/api/v1/projects/invalid-id", updateBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Update failure - project not found", func(t *testing.T) {
		updateBody := []byte(`{"name":"Updated Name"}`)
		resp := app.DoRequest(t, http.MethodPatch, "/api/v1/projects/00000000-0000-0000-0000-000000000000", updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestProjectHandler_ArchiveAndUnarchive(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	body := []byte(`{"name":"Archive Test Project"}`)
	respSeed := app.DoRequest(t, http.MethodPost, "/api/v1/projects", body, true)
	var created map[string]interface{}
	_ = json.NewDecoder(respSeed.Body).Decode(&created)
	projectID := created["id"].(string)

	t.Run("Archive success", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodPatch, "/api/v1/projects/"+projectID+"/archive", nil, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Archive failure - already archived", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodPatch, "/api/v1/projects/"+projectID+"/archive", nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Unarchive success", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodPatch, "/api/v1/projects/"+projectID+"/unarchive", nil, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Unarchive failure - not archived", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodPatch, "/api/v1/projects/"+projectID+"/unarchive", nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Archive failure - invalid UUID", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodPatch, "/api/v1/projects/bad-uuid/archive", nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestProjectHandler_Delete(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	body := []byte(`{"name":"Delete Test Project"}`)
	respSeed := app.DoRequest(t, http.MethodPost, "/api/v1/projects", body, true)
	var created map[string]interface{}
	_ = json.NewDecoder(respSeed.Body).Decode(&created)
	projectID := created["id"].(string)

	t.Run("Delete success", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodDelete, "/api/v1/projects/"+projectID, nil, true)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("Delete failure - invalid UUID", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodDelete, "/api/v1/projects/invalid-id", nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
