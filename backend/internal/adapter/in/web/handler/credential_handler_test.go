package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUnlockedVaultWithProject(t *testing.T, app *TestApp) string {
	// Setup master password & unlock vault
	setupBody := []byte(`{"masterPassword":"mySuperSecretMasterPassword123"}`)
	respSetup := app.DoRequest(t, http.MethodPost, "/api/v1/security/master-password", setupBody, true)
	require.Equal(t, http.StatusNoContent, respSetup.StatusCode)

	unlockBody := []byte(`{"masterPassword":"mySuperSecretMasterPassword123"}`)
	respUnlock := app.DoRequest(t, http.MethodPost, "/api/v1/security/vault/unlock", unlockBody, true)
	require.Equal(t, http.StatusOK, respUnlock.StatusCode)

	// Create project
	projBody := []byte(`{"name":"Integration Credential Project"}`)
	respProj := app.DoRequest(t, http.MethodPost, "/api/v1/projects", projBody, true)
	require.Equal(t, http.StatusCreated, respProj.StatusCode)

	var proj map[string]interface{}
	err := json.NewDecoder(respProj.Body).Decode(&proj)
	require.NoError(t, err)

	return proj["id"].(string)
}

func TestCredentialHandler_Create(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectID := setupUnlockedVaultWithProject(t, app)

	t.Run("Create LOGIN credential success", func(t *testing.T) {
		body := []byte(`{
			"title": "My Login DB",
			"secretType": "LOGIN",
			"username": "dbuser",
			"password": "dbpassword",
			"notes": "some notes",
			"relatedUrl": "http://db-server"
		}`)
		url := fmt.Sprintf("/api/v1/projects/%s/credentials", projectID)
		resp := app.DoRequest(t, http.MethodPost, url, body, true)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("Location"))

		var view map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&view)
		require.NoError(t, err)

		assert.Equal(t, projectID, view["projectId"])
		assert.Equal(t, "My Login DB", view["title"])
		assert.Equal(t, "LOGIN", view["secretType"])
		payload := view["decryptedPayload"].(map[string]interface{})
		assert.Equal(t, "dbuser", payload["username"])
		assert.Equal(t, "dbpassword", payload["password"])
	})

	t.Run("Create API_KEY credential success", func(t *testing.T) {
		body := []byte(`{
			"title": "Stripe API Key",
			"secretType": "API_KEY",
			"apiKey": "sk_test_123456789"
		}`)
		url := fmt.Sprintf("/api/v1/projects/%s/credentials", projectID)
		resp := app.DoRequest(t, http.MethodPost, url, body, true)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var view map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&view)
		require.NoError(t, err)

		assert.Equal(t, "API_KEY", view["secretType"])
		payload := view["decryptedPayload"].(map[string]interface{})
		assert.Equal(t, "sk_test_123456789", payload["apiKey"])
	})

	t.Run("Create RAW_TEXT credential success", func(t *testing.T) {
		body := []byte(`{
			"title": "Private Certificate",
			"secretType": "RAW_TEXT",
			"rawTextContent": "certificate content"
		}`)
		url := fmt.Sprintf("/api/v1/projects/%s/credentials", projectID)
		resp := app.DoRequest(t, http.MethodPost, url, body, true)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var view map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&view)
		require.NoError(t, err)

		assert.Equal(t, "RAW_TEXT", view["secretType"])
		payload := view["decryptedPayload"].(map[string]interface{})
		assert.Equal(t, "certificate content", payload["rawText"])
	})

	t.Run("Create failure - missing password for LOGIN", func(t *testing.T) {
		body := []byte(`{
			"title": "Invalid Login",
			"secretType": "LOGIN",
			"username": "admin"
		}`)
		url := fmt.Sprintf("/api/v1/projects/%s/credentials", projectID)
		resp := app.DoRequest(t, http.MethodPost, url, body, true)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - title too short", func(t *testing.T) {
		body := []byte(`{
			"title": "A",
			"secretType": "API_KEY",
			"apiKey": "key"
		}`)
		url := fmt.Sprintf("/api/v1/projects/%s/credentials", projectID)
		resp := app.DoRequest(t, http.MethodPost, url, body, true)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failure - project not found", func(t *testing.T) {
		randomProjID := uuid.New().String()
		body := []byte(`{
			"title": "Valid Title",
			"secretType": "API_KEY",
			"apiKey": "key"
		}`)
		url := fmt.Sprintf("/api/v1/projects/%s/credentials", randomProjID)
		resp := app.DoRequest(t, http.MethodPost, url, body, true)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestCredentialHandler_GetAll(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectID := setupUnlockedVaultWithProject(t, app)

	// Create 1 credential
	body1 := []byte(`{"title":"Cred 1","secretType":"API_KEY","apiKey":"key1"}`)
	url := fmt.Sprintf("/api/v1/projects/%s/credentials", projectID)
	resp1 := app.DoRequest(t, http.MethodPost, url, body1, true)
	require.Equal(t, http.StatusCreated, resp1.StatusCode)

	t.Run("GetAll credentials success", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, url+"?page=0&size=10", nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var page map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&page)
		require.NoError(t, err)

		content := page["content"].([]interface{})
		assert.Len(t, content, 1)

		summary := content[0].(map[string]interface{})
		assert.Equal(t, "Cred 1", summary["title"])
		assert.Nil(t, summary["decryptedPayload"]) // Summary must NOT leak payload
	})

	t.Run("GetAll failure - project not found", func(t *testing.T) {
		randomProjID := uuid.New().String()
		url := fmt.Sprintf("/api/v1/projects/%s/credentials?page=0&size=10", randomProjID)
		resp := app.DoRequest(t, http.MethodGet, url, nil, true)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestCredentialHandler_Get(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectID := setupUnlockedVaultWithProject(t, app)

	body := []byte(`{"title":"Get Credential","secretType":"LOGIN","username":"user","password":"pass"}`)
	urlPost := fmt.Sprintf("/api/v1/projects/%s/credentials", projectID)
	respPost := app.DoRequest(t, http.MethodPost, urlPost, body, true)
	require.Equal(t, http.StatusCreated, respPost.StatusCode)

	var createdView map[string]interface{}
	_ = json.NewDecoder(respPost.Body).Decode(&createdView)
	credID := createdView["id"].(string)

	t.Run("Get credential by ID success", func(t *testing.T) {
		urlGet := fmt.Sprintf("/api/v1/projects/%s/credentials/%s", projectID, credID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var view map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&view)
		require.NoError(t, err)

		assert.Equal(t, credID, view["id"])
		payload := view["decryptedPayload"].(map[string]interface{})
		assert.Equal(t, "user", payload["username"])
	})

	t.Run("Get credential by ID failure - credential not found", func(t *testing.T) {
		randomID := uuid.New().String()
		urlGet := fmt.Sprintf("/api/v1/projects/%s/credentials/%s", projectID, randomID)
		resp := app.DoRequest(t, http.MethodGet, urlGet, nil, true)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestCredentialHandler_Update(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectID := setupUnlockedVaultWithProject(t, app)

	body := []byte(`{"title":"Original Title","secretType":"API_KEY","apiKey":"orig_key"}`)
	urlPost := fmt.Sprintf("/api/v1/projects/%s/credentials", projectID)
	respPost := app.DoRequest(t, http.MethodPost, urlPost, body, true)
	require.Equal(t, http.StatusCreated, respPost.StatusCode)

	var createdView map[string]interface{}
	_ = json.NewDecoder(respPost.Body).Decode(&createdView)
	credID := createdView["id"].(string)

	t.Run("Update credential title only success", func(t *testing.T) {
		updateBody := []byte(`{"title":"Updated Title"}`)
		urlPatch := fmt.Sprintf("/api/v1/projects/%s/credentials/%s", projectID, credID)
		resp := app.DoRequest(t, http.MethodPatch, urlPatch, updateBody, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var view map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&view)
		require.NoError(t, err)

		assert.Equal(t, "Updated Title", view["title"])
		payload := view["decryptedPayload"].(map[string]interface{})
		assert.Equal(t, "orig_key", payload["apiKey"])
	})
}

func TestCredentialHandler_Delete(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectID := setupUnlockedVaultWithProject(t, app)

	body := []byte(`{"title":"Cred To Delete","secretType":"API_KEY","apiKey":"key"}`)
	urlPost := fmt.Sprintf("/api/v1/projects/%s/credentials", projectID)
	respPost := app.DoRequest(t, http.MethodPost, urlPost, body, true)
	require.Equal(t, http.StatusCreated, respPost.StatusCode)

	var createdView map[string]interface{}
	_ = json.NewDecoder(respPost.Body).Decode(&createdView)
	credID := createdView["id"].(string)

	t.Run("Delete credential success", func(t *testing.T) {
		urlDel := fmt.Sprintf("/api/v1/projects/%s/credentials/%s", projectID, credID)
		resp := app.DoRequest(t, http.MethodDelete, urlDel, nil, true)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify it's gone
		respGet := app.DoRequest(t, http.MethodGet, urlDel, nil, true)
		assert.Equal(t, http.StatusNotFound, respGet.StatusCode)
	})

	t.Run("Delete failure - credential not found", func(t *testing.T) {
		randomID := uuid.New().String()
		urlDel := fmt.Sprintf("/api/v1/projects/%s/credentials/%s", projectID, randomID)
		resp := app.DoRequest(t, http.MethodDelete, urlDel, nil, true)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestCredentialHandler_VaultLockedSecurity(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	projectID := setupUnlockedVaultWithProject(t, app)

	// Lock the vault
	respLock := app.DoRequest(t, http.MethodPost, "/api/v1/security/vault/lock", nil, true)
	require.Equal(t, http.StatusNoContent, respLock.StatusCode)

	t.Run("Create while vault locked returns 401 Unauthorized", func(t *testing.T) {
		body := []byte(`{"title":"Locked","secretType":"API_KEY","apiKey":"key"}`)
		url := fmt.Sprintf("/api/v1/projects/%s/credentials", projectID)
		resp := app.DoRequest(t, http.MethodPost, url, body, true)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
