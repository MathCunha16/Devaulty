package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoardHandler_Create(t *testing.T) {
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

	t.Run("Create success - auto default board and auto columns", func(t *testing.T) {
		boardBody := []byte(`{
			"name": "Sprint Board",
			"description": "Main sprint tracking"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, boardBody, true)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("Location"))

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "Sprint Board", result["name"])
		assert.Equal(t, projectID, result["projectId"])
		assert.Equal(t, true, result["isDefault"])
		assert.NotEmpty(t, result["id"])

		// Verify 4 default columns were created
		boardID := result["id"].(string)
		colPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
		colResp := app.DoRequest(t, http.MethodGet, colPath, nil, true)
		assert.Equal(t, http.StatusOK, colResp.StatusCode)

		var columns []map[string]interface{}
		err = json.NewDecoder(colResp.Body).Decode(&columns)
		require.NoError(t, err)
		assert.Equal(t, 4, len(columns))
		assert.Equal(t, "Backlog", columns[0]["name"])
		assert.Equal(t, "In Progress", columns[1]["name"])
		assert.Equal(t, "Review", columns[2]["name"])
		assert.Equal(t, "Done", columns[3]["name"])
	})

	t.Run("Create failure - project not found", func(t *testing.T) {
		boardBody := []byte(`{"name": "Valid Name"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects/00000000-0000-0000-0000-000000000000/boards", boardBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Create failure - invalid project UUID format", func(t *testing.T) {
		boardBody := []byte(`{"name": "Valid Name"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/projects/invalid-id/boards", boardBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - missing required name", func(t *testing.T) {
		boardBody := []byte(`{"description": "No name provided"}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, boardBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - name too short", func(t *testing.T) {
		boardBody := []byte(`{"name": "A"}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, boardBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - unauthorized", func(t *testing.T) {
		boardBody := []byte(`{"name": "Unauthorized Board"}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, boardBody, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestBoardHandler_GetAll(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed a project
	projectBody := []byte(`{"name":"Project for Boards"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	// Seed 2 boards
	for i := 1; i <= 2; i++ {
		boardBody := []byte(fmt.Sprintf(`{"name": "Board %d"}`, i))
		resp := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/boards", projectID), boardBody, true)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	}

	t.Run("GetAll success with pagination", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards?page=0&size=10", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, float64(2), result["totalElements"])
		content := result["content"].([]interface{})
		assert.Equal(t, 2, len(content))
	})

	t.Run("GetAll failure - project not found", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/00000000-0000-0000-0000-000000000000/boards", nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("GetAll failure - invalid project UUID", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/not-a-uuid/boards", nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("GetAll failure - unauthorized", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestBoardHandler_Get(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed a project
	projectBody := []byte(`{"name":"Project for Get Board"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	// Seed a board
	boardBody := []byte(`{"name": "Target Board", "description": "Specific board"}`)
	respBoard := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/boards", projectID), boardBody, true)
	require.Equal(t, http.StatusCreated, respBoard.StatusCode)

	var createdBoard map[string]interface{}
	err = json.NewDecoder(respBoard.Body).Decode(&createdBoard)
	require.NoError(t, err)
	boardID := createdBoard["id"].(string)

	t.Run("Get success", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s", projectID, boardID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, boardID, result["id"])
		assert.Equal(t, "Target Board", result["name"])
		assert.Equal(t, "Specific board", result["description"])
	})

	t.Run("Get failure - board not found", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get failure - invalid board UUID", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/invalid-uuid", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Get failure - unauthorized", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s", projectID, boardID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestBoardHandler_GetDefault(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed a project
	projectBody := []byte(`{"name":"Project for Default Board"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	// Create first board (automatically becomes default)
	boardBody := []byte(`{"name": "Initial Default Board"}`)
	respBoard := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/boards", projectID), boardBody, true)
	require.Equal(t, http.StatusCreated, respBoard.StatusCode)

	var createdBoard map[string]interface{}
	err = json.NewDecoder(respBoard.Body).Decode(&createdBoard)
	require.NoError(t, err)
	boardID := createdBoard["id"].(string)

	t.Run("GetDefault success", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/default", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, boardID, result["id"])
		assert.Equal(t, true, result["isDefault"])
		assert.Equal(t, "Initial Default Board", result["name"])
	})

	t.Run("GetDefault failure - project not found", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/projects/00000000-0000-0000-0000-000000000000/boards/default", nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("GetDefault failure - unauthorized", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/default", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestBoardHandler_Update(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed a project
	projectBody := []byte(`{"name":"Project for Update Board"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	// Seed a board
	boardBody := []byte(`{"name": "Original Name"}`)
	respBoard := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/boards", projectID), boardBody, true)
	require.Equal(t, http.StatusCreated, respBoard.StatusCode)

	var createdBoard map[string]interface{}
	err = json.NewDecoder(respBoard.Body).Decode(&createdBoard)
	require.NoError(t, err)
	boardID := createdBoard["id"].(string)

	t.Run("Update success", func(t *testing.T) {
		updateBody := []byte(`{
			"name": "Updated Board Name",
			"description": "Updated Description"
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s", projectID, boardID)
		resp := app.DoRequest(t, http.MethodPatch, urlPath, updateBody, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "Updated Board Name", result["name"])
		assert.Equal(t, "Updated Description", result["description"])
	})

	t.Run("Update failure - board not found", func(t *testing.T) {
		updateBody := []byte(`{"name": "Valid Name"}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodPatch, urlPath, updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Update failure - invalid UUID", func(t *testing.T) {
		updateBody := []byte(`{"name": "Valid Name"}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/invalid-uuid", projectID)
		resp := app.DoRequest(t, http.MethodPatch, urlPath, updateBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Update failure - unauthorized", func(t *testing.T) {
		updateBody := []byte(`{"name": "Unauthorized Update"}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s", projectID, boardID)
		resp := app.DoRequest(t, http.MethodPatch, urlPath, updateBody, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestBoardHandler_Delete(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed a project
	projectBody := []byte(`{"name":"Project for Delete Board"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	// Seed a board
	boardBody := []byte(`{"name": "Board to Delete"}`)
	respBoard := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/boards", projectID), boardBody, true)
	require.Equal(t, http.StatusCreated, respBoard.StatusCode)

	var createdBoard map[string]interface{}
	err = json.NewDecoder(respBoard.Body).Decode(&createdBoard)
	require.NoError(t, err)
	boardID := createdBoard["id"].(string)

	t.Run("Delete success", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s", projectID, boardID)
		resp := app.DoRequest(t, http.MethodDelete, urlPath, nil, true)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify deletion
		getResp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
	})

	t.Run("Delete failure - board not found", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/00000000-0000-0000-0000-000000000000", projectID)
		resp := app.DoRequest(t, http.MethodDelete, urlPath, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Delete failure - invalid UUID", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/invalid-uuid", projectID)
		resp := app.DoRequest(t, http.MethodDelete, urlPath, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Delete failure - unauthorized", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s", projectID, boardID)
		resp := app.DoRequest(t, http.MethodDelete, urlPath, nil, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
