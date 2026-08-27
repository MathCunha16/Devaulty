package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoardColumnHandler_Create(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project and board
	projectBody := []byte(`{"name":"Project for Columns"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	boardBody := []byte(`{"name": "Sprint Board"}`)
	respBoard := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/boards", projectID), boardBody, true)
	require.Equal(t, http.StatusCreated, respBoard.StatusCode)

	var createdBoard map[string]interface{}
	err = json.NewDecoder(respBoard.Body).Decode(&createdBoard)
	require.NoError(t, err)
	boardID := createdBoard["id"].(string)

	t.Run("Create success - auto position after default columns", func(t *testing.T) {
		colBody := []byte(`{
			"name": "QA Testing",
			"wipLimit": 5
		}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, colBody, true)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("Location"))

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "QA Testing", result["name"])
		assert.Equal(t, projectID, result["projectId"])
		assert.Equal(t, boardID, result["boardId"])
		assert.Equal(t, float64(4), result["position"]) // 4 default columns (0,1,2,3), so next is 4
		assert.Equal(t, float64(5), result["wipLimit"])
	})

	t.Run("Create failure - board not found", func(t *testing.T) {
		colBody := []byte(`{"name": "QA"}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/00000000-0000-0000-0000-000000000000/columns", projectID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, colBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Create failure - project not found", func(t *testing.T) {
		colBody := []byte(`{"name": "QA"}`)
		urlPath := fmt.Sprintf("/api/v1/projects/00000000-0000-0000-0000-000000000000/boards/%s/columns", boardID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, colBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Create failure - missing name", func(t *testing.T) {
		colBody := []byte(`{"wipLimit": 3}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, colBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - unauthorized", func(t *testing.T) {
		colBody := []byte(`{"name": "QA"}`)
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
		resp := app.DoRequest(t, http.MethodPost, urlPath, colBody, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestBoardColumnHandler_GetAll(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project and board
	projectBody := []byte(`{"name":"Project for Column List"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	boardBody := []byte(`{"name": "Sprint Board"}`)
	respBoard := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/boards", projectID), boardBody, true)
	require.Equal(t, http.StatusCreated, respBoard.StatusCode)

	var createdBoard map[string]interface{}
	err = json.NewDecoder(respBoard.Body).Decode(&createdBoard)
	require.NoError(t, err)
	boardID := createdBoard["id"].(string)

	t.Run("GetAll success - default columns returned in order", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var columns []map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&columns)
		require.NoError(t, err)
		assert.Equal(t, 4, len(columns))
		assert.Equal(t, "Backlog", columns[0]["name"])
		assert.Equal(t, float64(0), columns[0]["position"])
		assert.Equal(t, "In Progress", columns[1]["name"])
		assert.Equal(t, float64(1), columns[1]["position"])
	})

	t.Run("GetAll failure - board not found", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/00000000-0000-0000-0000-000000000000/columns", projectID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("GetAll failure - unauthorized", func(t *testing.T) {
		urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
		resp := app.DoRequest(t, http.MethodGet, urlPath, nil, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestBoardColumnHandler_Get(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project and board
	projectBody := []byte(`{"name":"Project for Column Get"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	boardBody := []byte(`{"name": "Sprint Board"}`)
	respBoard := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/boards", projectID), boardBody, true)
	require.Equal(t, http.StatusCreated, respBoard.StatusCode)

	var createdBoard map[string]interface{}
	err = json.NewDecoder(respBoard.Body).Decode(&createdBoard)
	require.NoError(t, err)
	boardID := createdBoard["id"].(string)

	// Fetch columns to get first column ID
	urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
	respCols := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
	require.Equal(t, http.StatusOK, respCols.StatusCode)

	var columns []map[string]interface{}
	err = json.NewDecoder(respCols.Body).Decode(&columns)
	require.NoError(t, err)
	columnID := columns[0]["id"].(string)

	t.Run("Get success", func(t *testing.T) {
		colUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s", projectID, boardID, columnID)
		resp := app.DoRequest(t, http.MethodGet, colUrl, nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, columnID, result["id"])
		assert.Equal(t, "Backlog", result["name"])
	})

	t.Run("Get failure - column not found", func(t *testing.T) {
		colUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/00000000-0000-0000-0000-000000000000", projectID, boardID)
		resp := app.DoRequest(t, http.MethodGet, colUrl, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get failure - unauthorized", func(t *testing.T) {
		colUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s", projectID, boardID, columnID)
		resp := app.DoRequest(t, http.MethodGet, colUrl, nil, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestBoardColumnHandler_Update(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project and board
	projectBody := []byte(`{"name":"Project for Column Update"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	boardBody := []byte(`{"name": "Sprint Board"}`)
	respBoard := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/boards", projectID), boardBody, true)
	require.Equal(t, http.StatusCreated, respBoard.StatusCode)

	var createdBoard map[string]interface{}
	err = json.NewDecoder(respBoard.Body).Decode(&createdBoard)
	require.NoError(t, err)
	boardID := createdBoard["id"].(string)

	// Fetch columns to get first column ID
	urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
	respCols := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
	require.Equal(t, http.StatusOK, respCols.StatusCode)

	var columns []map[string]interface{}
	err = json.NewDecoder(respCols.Body).Decode(&columns)
	require.NoError(t, err)
	columnID := columns[0]["id"].(string)

	t.Run("Update success", func(t *testing.T) {
		updateBody := []byte(`{
			"name": "To Do (Updated)",
			"wipLimit": 10
		}`)
		colUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s", projectID, boardID, columnID)
		resp := app.DoRequest(t, http.MethodPatch, colUrl, updateBody, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "To Do (Updated)", result["name"])
		assert.Equal(t, float64(10), result["wipLimit"])
	})

	t.Run("Update failure - column not found", func(t *testing.T) {
		updateBody := []byte(`{"name": "New Name"}`)
		colUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/00000000-0000-0000-0000-000000000000", projectID, boardID)
		resp := app.DoRequest(t, http.MethodPatch, colUrl, updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Update failure - unauthorized", func(t *testing.T) {
		updateBody := []byte(`{"name": "New Name"}`)
		colUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s", projectID, boardID, columnID)
		resp := app.DoRequest(t, http.MethodPatch, colUrl, updateBody, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestBoardColumnHandler_Delete(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project and board
	projectBody := []byte(`{"name":"Project for Column Delete"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	boardBody := []byte(`{"name": "Sprint Board"}`)
	respBoard := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/boards", projectID), boardBody, true)
	require.Equal(t, http.StatusCreated, respBoard.StatusCode)

	var createdBoard map[string]interface{}
	err = json.NewDecoder(respBoard.Body).Decode(&createdBoard)
	require.NoError(t, err)
	boardID := createdBoard["id"].(string)

	// Fetch columns to get last column ID
	urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
	respCols := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
	require.Equal(t, http.StatusOK, respCols.StatusCode)

	var columns []map[string]interface{}
	err = json.NewDecoder(respCols.Body).Decode(&columns)
	require.NoError(t, err)
	columnID := columns[3]["id"].(string)

	t.Run("Delete success", func(t *testing.T) {
		colUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s", projectID, boardID, columnID)
		resp := app.DoRequest(t, http.MethodDelete, colUrl, nil, true)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify deleted
		getResp := app.DoRequest(t, http.MethodGet, colUrl, nil, true)
		assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
	})

	t.Run("Delete failure - column not found", func(t *testing.T) {
		colUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/00000000-0000-0000-0000-000000000000", projectID, boardID)
		resp := app.DoRequest(t, http.MethodDelete, colUrl, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Delete failure - unauthorized", func(t *testing.T) {
		colUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s", projectID, boardID, columnID)
		resp := app.DoRequest(t, http.MethodDelete, colUrl, nil, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestBoardColumnHandler_Reorder(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project and board
	projectBody := []byte(`{"name":"Project for Column Reorder"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	boardBody := []byte(`{"name": "Sprint Board"}`)
	respBoard := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/boards", projectID), boardBody, true)
	require.Equal(t, http.StatusCreated, respBoard.StatusCode)

	var createdBoard map[string]interface{}
	err = json.NewDecoder(respBoard.Body).Decode(&createdBoard)
	require.NoError(t, err)
	boardID := createdBoard["id"].(string)

	// Fetch columns
	urlPath := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
	respCols := app.DoRequest(t, http.MethodGet, urlPath, nil, true)
	require.Equal(t, http.StatusOK, respCols.StatusCode)

	var columns []map[string]interface{}
	err = json.NewDecoder(respCols.Body).Decode(&columns)
	require.NoError(t, err)
	require.Equal(t, 4, len(columns))

	id0 := columns[0]["id"].(string)
	id1 := columns[1]["id"].(string)
	id2 := columns[2]["id"].(string)
	id3 := columns[3]["id"].(string)

	t.Run("Reorder success - invert order", func(t *testing.T) {
		reorderBody := []byte(fmt.Sprintf(`{
			"positions": ["%s", "%s", "%s", "%s"]
		}`, id3, id2, id1, id0))

		reorderUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/reorder", projectID, boardID)
		resp := app.DoRequest(t, http.MethodPatch, reorderUrl, reorderBody, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var reorderedCols []map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&reorderedCols)
		require.NoError(t, err)
		assert.Equal(t, 4, len(reorderedCols))
		assert.Equal(t, id3, reorderedCols[0]["id"])
		assert.Equal(t, float64(0), reorderedCols[0]["position"])
		assert.Equal(t, id0, reorderedCols[3]["id"])
		assert.Equal(t, float64(3), reorderedCols[3]["position"])
	})

	t.Run("Reorder failure - empty positions list", func(t *testing.T) {
		reorderBody := []byte(`{"positions": []}`)
		reorderUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/reorder", projectID, boardID)
		resp := app.DoRequest(t, http.MethodPatch, reorderUrl, reorderBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Reorder failure - board not found", func(t *testing.T) {
		reorderBody := []byte(fmt.Sprintf(`{"positions": ["%s"]}`, id0))
		reorderUrl := fmt.Sprintf("/api/v1/projects/%s/boards/00000000-0000-0000-0000-000000000000/columns/reorder", projectID)
		resp := app.DoRequest(t, http.MethodPatch, reorderUrl, reorderBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Reorder failure - unauthorized", func(t *testing.T) {
		reorderBody := []byte(fmt.Sprintf(`{"positions": ["%s"]}`, id0))
		reorderUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/reorder", projectID, boardID)
		resp := app.DoRequest(t, http.MethodPatch, reorderUrl, reorderBody, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
