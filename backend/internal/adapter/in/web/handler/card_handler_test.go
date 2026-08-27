package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCardHandler_Create(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project, note (for linking), and board
	projectBody := []byte(`{"name":"Project for Card Create"}`)
	respProject := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projectBody, true)
	require.Equal(t, http.StatusCreated, respProject.StatusCode)

	var createdProject map[string]interface{}
	err := json.NewDecoder(respProject.Body).Decode(&createdProject)
	require.NoError(t, err)
	projectID := createdProject["id"].(string)

	// Seed note to link
	noteBody := []byte(`{"title":"API Specs","content":"Details about endpoints"}`)
	respNote := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/notes", projectID), noteBody, true)
	require.Equal(t, http.StatusCreated, respNote.StatusCode)

	var createdNote map[string]interface{}
	err = json.NewDecoder(respNote.Body).Decode(&createdNote)
	require.NoError(t, err)
	noteID := createdNote["id"].(string)

	// Seed board
	boardBody := []byte(`{"name": "Sprint Board"}`)
	respBoard := app.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/boards", projectID), boardBody, true)
	require.Equal(t, http.StatusCreated, respBoard.StatusCode)

	var createdBoard map[string]interface{}
	err = json.NewDecoder(respBoard.Body).Decode(&createdBoard)
	require.NoError(t, err)
	boardID := createdBoard["id"].(string)

	// Get Backlog column ID
	colUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
	respCols := app.DoRequest(t, http.MethodGet, colUrl, nil, true)
	require.Equal(t, http.StatusOK, respCols.StatusCode)

	var columns []map[string]interface{}
	err = json.NewDecoder(respCols.Body).Decode(&columns)
	require.NoError(t, err)
	columnID := columns[0]["id"].(string)

	t.Run("Create success - with linked note and priority", func(t *testing.T) {
		cardBody := []byte(fmt.Sprintf(`{
			"title": "Build Auth Flow",
			"description": "Implement OAuth2",
			"priority": "HIGH",
			"dueDate": "2026-10-15T23:59:59Z",
			"linkedItems": [
				{
					"itemType": "NOTE",
					"itemId": "%s"
				}
			]
		}`, noteID))

		cardUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s/cards", projectID, boardID, columnID)
		resp := app.DoRequest(t, http.MethodPost, cardUrl, cardBody, true)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("Location"))

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.NotEmpty(t, result["id"])
		assert.Equal(t, "Build Auth Flow", result["title"])
		assert.Equal(t, "Implement OAuth2", result["description"])
		assert.Equal(t, "HIGH", result["priority"])
		assert.Equal(t, float64(0), result["position"])

		linkedItems := result["linkedItems"].([]interface{})
		assert.Equal(t, 1, len(linkedItems))
		firstItem := linkedItems[0].(map[string]interface{})
		assert.Equal(t, noteID, firstItem["itemId"])
		assert.Equal(t, "NOTE", firstItem["itemType"])
	})

	t.Run("Create success - simple card without linked items", func(t *testing.T) {
		cardBody := []byte(`{"title": "Simple Task"}`)
		cardUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s/cards", projectID, boardID, columnID)
		resp := app.DoRequest(t, http.MethodPost, cardUrl, cardBody, true)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "Simple Task", result["title"])
		assert.Equal(t, float64(1), result["position"]) // Next position in Backlog is 1
	})

	t.Run("Create failure - column not found", func(t *testing.T) {
		cardBody := []byte(`{"title": "Valid Title"}`)
		cardUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/00000000-0000-0000-0000-000000000000/cards", projectID, boardID)
		resp := app.DoRequest(t, http.MethodPost, cardUrl, cardBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Create failure - board not found", func(t *testing.T) {
		cardBody := []byte(`{"title": "Valid Title"}`)
		cardUrl := fmt.Sprintf("/api/v1/projects/%s/boards/00000000-0000-0000-0000-000000000000/columns/%s/cards", projectID, columnID)
		resp := app.DoRequest(t, http.MethodPost, cardUrl, cardBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Create failure - missing title", func(t *testing.T) {
		cardBody := []byte(`{"description": "No title"}`)
		cardUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s/cards", projectID, boardID, columnID)
		resp := app.DoRequest(t, http.MethodPost, cardUrl, cardBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - title too short", func(t *testing.T) {
		cardBody := []byte(`{"title": "X"}`)
		cardUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s/cards", projectID, boardID, columnID)
		resp := app.DoRequest(t, http.MethodPost, cardUrl, cardBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - unauthorized", func(t *testing.T) {
		cardBody := []byte(`{"title": "Unauthorized Card"}`)
		cardUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s/cards", projectID, boardID, columnID)
		resp := app.DoRequest(t, http.MethodPost, cardUrl, cardBody, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestCardHandler_GetAll(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project, board, and cards
	projectBody := []byte(`{"name":"Project for Card List"}`)
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

	// Get Backlog column
	colUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
	respCols := app.DoRequest(t, http.MethodGet, colUrl, nil, true)
	require.Equal(t, http.StatusOK, respCols.StatusCode)

	var columns []map[string]interface{}
	err = json.NewDecoder(respCols.Body).Decode(&columns)
	require.NoError(t, err)
	columnID := columns[0]["id"].(string)

	// Create 2 cards
	for i := 1; i <= 2; i++ {
		cardBody := []byte(fmt.Sprintf(`{"title": "Card %d"}`, i))
		cardUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s/cards", projectID, boardID, columnID)
		resp := app.DoRequest(t, http.MethodPost, cardUrl, cardBody, true)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	}

	t.Run("GetAll success", func(t *testing.T) {
		getAllUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards", projectID, boardID)
		resp := app.DoRequest(t, http.MethodGet, getAllUrl, nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var cards []map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&cards)
		require.NoError(t, err)
		assert.Equal(t, 2, len(cards))
		assert.Equal(t, "Card 1", cards[0]["title"])
		assert.Equal(t, "Card 2", cards[1]["title"])
	})

	t.Run("GetAll failure - board not found", func(t *testing.T) {
		getAllUrl := fmt.Sprintf("/api/v1/projects/%s/boards/00000000-0000-0000-0000-000000000000/cards", projectID)
		resp := app.DoRequest(t, http.MethodGet, getAllUrl, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("GetAll failure - unauthorized", func(t *testing.T) {
		getAllUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards", projectID, boardID)
		resp := app.DoRequest(t, http.MethodGet, getAllUrl, nil, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestCardHandler_Get(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project, board, and card
	projectBody := []byte(`{"name":"Project for Card Get"}`)
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

	// Get Backlog column
	colUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
	respCols := app.DoRequest(t, http.MethodGet, colUrl, nil, true)
	require.Equal(t, http.StatusOK, respCols.StatusCode)

	var columns []map[string]interface{}
	err = json.NewDecoder(respCols.Body).Decode(&columns)
	require.NoError(t, err)
	columnID := columns[0]["id"].(string)

	// Create card
	cardBody := []byte(`{"title": "Target Card", "description": "Details here"}`)
	cardUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s/cards", projectID, boardID, columnID)
	respCard := app.DoRequest(t, http.MethodPost, cardUrl, cardBody, true)
	require.Equal(t, http.StatusCreated, respCard.StatusCode)

	var createdCard map[string]interface{}
	err = json.NewDecoder(respCard.Body).Decode(&createdCard)
	require.NoError(t, err)
	cardID := createdCard["id"].(string)

	t.Run("Get success", func(t *testing.T) {
		getUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards/%s", projectID, boardID, cardID)
		resp := app.DoRequest(t, http.MethodGet, getUrl, nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, cardID, result["id"])
		assert.Equal(t, "Target Card", result["title"])
		assert.Equal(t, "Details here", result["description"])
	})

	t.Run("Get failure - card not found", func(t *testing.T) {
		getUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards/00000000-0000-0000-0000-000000000000", projectID, boardID)
		resp := app.DoRequest(t, http.MethodGet, getUrl, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get failure - unauthorized", func(t *testing.T) {
		getUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards/%s", projectID, boardID, cardID)
		resp := app.DoRequest(t, http.MethodGet, getUrl, nil, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestCardHandler_Update(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project, board, and card
	projectBody := []byte(`{"name":"Project for Card Update"}`)
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

	// Get Backlog column
	colUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
	respCols := app.DoRequest(t, http.MethodGet, colUrl, nil, true)
	require.Equal(t, http.StatusOK, respCols.StatusCode)

	var columns []map[string]interface{}
	err = json.NewDecoder(respCols.Body).Decode(&columns)
	require.NoError(t, err)
	columnID := columns[0]["id"].(string)

	// Create card
	cardBody := []byte(`{"title": "Initial Card Title"}`)
	cardUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s/cards", projectID, boardID, columnID)
	respCard := app.DoRequest(t, http.MethodPost, cardUrl, cardBody, true)
	require.Equal(t, http.StatusCreated, respCard.StatusCode)

	var createdCard map[string]interface{}
	err = json.NewDecoder(respCard.Body).Decode(&createdCard)
	require.NoError(t, err)
	cardID := createdCard["id"].(string)

	t.Run("Update success", func(t *testing.T) {
		updateBody := []byte(`{
			"title": "Updated Card Title",
			"description": "Added new description",
			"priority": "EXTREMELY_HIGH"
		}`)
		updateUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards/%s", projectID, boardID, cardID)
		resp := app.DoRequest(t, http.MethodPatch, updateUrl, updateBody, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "Updated Card Title", result["title"])
		assert.Equal(t, "Added new description", result["description"])
		assert.Equal(t, "EXTREMELY_HIGH", result["priority"])
	})

	t.Run("Update failure - card not found", func(t *testing.T) {
		updateBody := []byte(`{"title": "New Title"}`)
		updateUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards/00000000-0000-0000-0000-000000000000", projectID, boardID)
		resp := app.DoRequest(t, http.MethodPatch, updateUrl, updateBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Update failure - unauthorized", func(t *testing.T) {
		updateBody := []byte(`{"title": "New Title"}`)
		updateUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards/%s", projectID, boardID, cardID)
		resp := app.DoRequest(t, http.MethodPatch, updateUrl, updateBody, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestCardHandler_Move(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project and board
	projectBody := []byte(`{"name":"Project for Card Move"}`)
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

	// Fetch columns: 0=Backlog, 1=In Progress
	colUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
	respCols := app.DoRequest(t, http.MethodGet, colUrl, nil, true)
	require.Equal(t, http.StatusOK, respCols.StatusCode)

	var columns []map[string]interface{}
	err = json.NewDecoder(respCols.Body).Decode(&columns)
	require.NoError(t, err)
	backlogColID := columns[0]["id"].(string)
	inProgressColID := columns[1]["id"].(string)

	// Create card in Backlog
	cardBody := []byte(`{"title": "Card to Move"}`)
	cardUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s/cards", projectID, boardID, backlogColID)
	respCard := app.DoRequest(t, http.MethodPost, cardUrl, cardBody, true)
	require.Equal(t, http.StatusCreated, respCard.StatusCode)

	var createdCard map[string]interface{}
	err = json.NewDecoder(respCard.Body).Decode(&createdCard)
	require.NoError(t, err)
	cardID := createdCard["id"].(string)

	t.Run("Move success - move from Backlog to In Progress", func(t *testing.T) {
		moveBody := []byte(fmt.Sprintf(`{
			"targetColumnId": "%s",
			"position": 0
		}`, inProgressColID))

		moveUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards/%s/move", projectID, boardID, cardID)
		resp := app.DoRequest(t, http.MethodPatch, moveUrl, moveBody, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify card now belongs to In Progress column
		getUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards/%s", projectID, boardID, cardID)
		getResp := app.DoRequest(t, http.MethodGet, getUrl, nil, true)
		assert.Equal(t, http.StatusOK, getResp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(getResp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, inProgressColID, result["columnId"])
		assert.Equal(t, float64(0), result["position"])
	})

	t.Run("Move failure - target column not found", func(t *testing.T) {
		moveBody := []byte(`{
			"targetColumnId": "11111111-1111-1111-1111-111111111111",
			"position": 0
		}`)
		moveUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards/%s/move", projectID, boardID, cardID)
		resp := app.DoRequest(t, http.MethodPatch, moveUrl, moveBody, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Move failure - missing required fields", func(t *testing.T) {
		moveBody := []byte(`{"position": 0}`)
		moveUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards/%s/move", projectID, boardID, cardID)
		resp := app.DoRequest(t, http.MethodPatch, moveUrl, moveBody, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Move failure - unauthorized", func(t *testing.T) {
		moveBody := []byte(fmt.Sprintf(`{"targetColumnId": "%s", "position": 0}`, inProgressColID))
		moveUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards/%s/move", projectID, boardID, cardID)
		resp := app.DoRequest(t, http.MethodPatch, moveUrl, moveBody, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestCardHandler_Delete(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed project, board, and card
	projectBody := []byte(`{"name":"Project for Card Delete"}`)
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

	// Get Backlog column
	colUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns", projectID, boardID)
	respCols := app.DoRequest(t, http.MethodGet, colUrl, nil, true)
	require.Equal(t, http.StatusOK, respCols.StatusCode)

	var columns []map[string]interface{}
	err = json.NewDecoder(respCols.Body).Decode(&columns)
	require.NoError(t, err)
	columnID := columns[0]["id"].(string)

	// Create card
	cardBody := []byte(`{"title": "Card to Delete"}`)
	cardUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/columns/%s/cards", projectID, boardID, columnID)
	respCard := app.DoRequest(t, http.MethodPost, cardUrl, cardBody, true)
	require.Equal(t, http.StatusCreated, respCard.StatusCode)

	var createdCard map[string]interface{}
	err = json.NewDecoder(respCard.Body).Decode(&createdCard)
	require.NoError(t, err)
	cardID := createdCard["id"].(string)

	t.Run("Delete success", func(t *testing.T) {
		deleteUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards/%s", projectID, boardID, cardID)
		resp := app.DoRequest(t, http.MethodDelete, deleteUrl, nil, true)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify deletion
		getResp := app.DoRequest(t, http.MethodGet, deleteUrl, nil, true)
		assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
	})

	t.Run("Delete failure - card not found", func(t *testing.T) {
		deleteUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards/00000000-0000-0000-0000-000000000000", projectID, boardID)
		resp := app.DoRequest(t, http.MethodDelete, deleteUrl, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Delete failure - unauthorized", func(t *testing.T) {
		deleteUrl := fmt.Sprintf("/api/v1/projects/%s/boards/%s/cards/%s", projectID, boardID, cardID)
		resp := app.DoRequest(t, http.MethodDelete, deleteUrl, nil, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
