package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemTagHandler_Associate(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed Project
	respProj := app.DoRequest(t, http.MethodPost, "/api/v1/projects", []byte(`{"name":"Parent Project"}`), true)
	var proj map[string]interface{}
	_ = json.NewDecoder(respProj.Body).Decode(&proj)
	projectID := proj["id"].(string)

	// Seed Snippet
	urlSnippet := fmt.Sprintf("/api/v1/projects/%s/snippets", projectID)
	respSnippet := app.DoRequest(t, http.MethodPost, urlSnippet, []byte(`{"title":"Sample Snippet","content":"code","language":"GO","snippetType":"CODE"}`), true)
	var snippet map[string]interface{}
	_ = json.NewDecoder(respSnippet.Body).Decode(&snippet)
	snippetID := snippet["id"].(string)

	// Seed Tag
	urlTag := fmt.Sprintf("/api/v1/projects/%s/tags", projectID)
	respTag := app.DoRequest(t, http.MethodPost, urlTag, []byte(`{"name":"Backend","color":"#FF0000"}`), true)
	var tag map[string]interface{}
	_ = json.NewDecoder(respTag.Body).Decode(&tag)
	tagID := tag["id"].(string)

	t.Run("Associate success", func(t *testing.T) {
		urlAssoc := fmt.Sprintf("/api/v1/projects/%s/items/snippet/%s/tags/%s", projectID, snippetID, tagID)
		resp := app.DoRequest(t, http.MethodPut, urlAssoc, nil, true)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("Associate failure - project not found", func(t *testing.T) {
		urlAssoc := fmt.Sprintf("/api/v1/projects/00000000-0000-0000-0000-000000000000/items/snippet/%s/tags/%s", snippetID, tagID)
		resp := app.DoRequest(t, http.MethodPut, urlAssoc, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Associate failure - tag not found", func(t *testing.T) {
		urlAssoc := fmt.Sprintf("/api/v1/projects/%s/items/snippet/%s/tags/00000000-0000-0000-0000-000000000000", projectID, snippetID)
		resp := app.DoRequest(t, http.MethodPut, urlAssoc, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Associate failure - item not found", func(t *testing.T) {
		urlAssoc := fmt.Sprintf("/api/v1/projects/%s/items/snippet/00000000-0000-0000-0000-000000000000/tags/%s", projectID, tagID)
		resp := app.DoRequest(t, http.MethodPut, urlAssoc, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Associate failure - unsupported item type", func(t *testing.T) {
		urlAssoc := fmt.Sprintf("/api/v1/projects/%s/items/invalid_type/%s/tags/%s", projectID, snippetID, tagID)
		resp := app.DoRequest(t, http.MethodPut, urlAssoc, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Associate failure - invalid project UUID", func(t *testing.T) {
		urlAssoc := fmt.Sprintf("/api/v1/projects/invalid-id/items/snippet/%s/tags/%s", snippetID, tagID)
		resp := app.DoRequest(t, http.MethodPut, urlAssoc, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Associate failure - invalid item UUID", func(t *testing.T) {
		urlAssoc := fmt.Sprintf("/api/v1/projects/%s/items/snippet/invalid-id/tags/%s", projectID, tagID)
		resp := app.DoRequest(t, http.MethodPut, urlAssoc, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Associate failure - invalid tag UUID", func(t *testing.T) {
		urlAssoc := fmt.Sprintf("/api/v1/projects/%s/items/snippet/%s/tags/invalid-id", projectID, snippetID)
		resp := app.DoRequest(t, http.MethodPut, urlAssoc, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestItemTagHandler_Disassociate(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	// Seed Project
	respProj := app.DoRequest(t, http.MethodPost, "/api/v1/projects", []byte(`{"name":"Parent Project"}`), true)
	var proj map[string]interface{}
	_ = json.NewDecoder(respProj.Body).Decode(&proj)
	projectID := proj["id"].(string)

	// Seed Snippet
	urlSnippet := fmt.Sprintf("/api/v1/projects/%s/snippets", projectID)
	respSnippet := app.DoRequest(t, http.MethodPost, urlSnippet, []byte(`{"title":"Sample Snippet","content":"code","language":"GO","snippetType":"CODE"}`), true)
	var snippet map[string]interface{}
	_ = json.NewDecoder(respSnippet.Body).Decode(&snippet)
	snippetID := snippet["id"].(string)

	// Seed Tag
	urlTag := fmt.Sprintf("/api/v1/projects/%s/tags", projectID)
	respTag := app.DoRequest(t, http.MethodPost, urlTag, []byte(`{"name":"Backend","color":"#FF0000"}`), true)
	var tag map[string]interface{}
	_ = json.NewDecoder(respTag.Body).Decode(&tag)
	tagID := tag["id"].(string)

	// Associate first
	urlAssoc := fmt.Sprintf("/api/v1/projects/%s/items/snippet/%s/tags/%s", projectID, snippetID, tagID)
	respAssoc := app.DoRequest(t, http.MethodPut, urlAssoc, nil, true)
	require.Equal(t, http.StatusNoContent, respAssoc.StatusCode)

	t.Run("Disassociate success", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodDelete, urlAssoc, nil, true)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("Disassociate failure - project not found", func(t *testing.T) {
		urlDisassoc := fmt.Sprintf("/api/v1/projects/00000000-0000-0000-0000-000000000000/items/snippet/%s/tags/%s", snippetID, tagID)
		resp := app.DoRequest(t, http.MethodDelete, urlDisassoc, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Disassociate failure - tag not found", func(t *testing.T) {
		urlDisassoc := fmt.Sprintf("/api/v1/projects/%s/items/snippet/%s/tags/00000000-0000-0000-0000-000000000000", projectID, snippetID)
		resp := app.DoRequest(t, http.MethodDelete, urlDisassoc, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Disassociate failure - item not found", func(t *testing.T) {
		urlDisassoc := fmt.Sprintf("/api/v1/projects/%s/items/snippet/00000000-0000-0000-0000-000000000000/tags/%s", projectID, tagID)
		resp := app.DoRequest(t, http.MethodDelete, urlDisassoc, nil, true)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Disassociate failure - unsupported item type", func(t *testing.T) {
		urlDisassoc := fmt.Sprintf("/api/v1/projects/%s/items/invalid_type/%s/tags/%s", projectID, snippetID, tagID)
		resp := app.DoRequest(t, http.MethodDelete, urlDisassoc, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Disassociate failure - invalid project UUID", func(t *testing.T) {
		urlDisassoc := fmt.Sprintf("/api/v1/projects/invalid-id/items/snippet/%s/tags/%s", snippetID, tagID)
		resp := app.DoRequest(t, http.MethodDelete, urlDisassoc, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Disassociate failure - invalid item UUID", func(t *testing.T) {
		urlDisassoc := fmt.Sprintf("/api/v1/projects/%s/items/snippet/invalid-id/tags/%s", projectID, tagID)
		resp := app.DoRequest(t, http.MethodDelete, urlDisassoc, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Disassociate failure - invalid tag UUID", func(t *testing.T) {
		urlDisassoc := fmt.Sprintf("/api/v1/projects/%s/items/snippet/%s/tags/invalid-id", projectID, snippetID)
		resp := app.DoRequest(t, http.MethodDelete, urlDisassoc, nil, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
