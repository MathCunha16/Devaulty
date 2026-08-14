package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityHandler_SetupMasterPassword(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	t.Run("Setup success - first time", func(t *testing.T) {
		body := []byte(`{"masterPassword":"mySuperSecretPassword123"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/security/master-password", body, true)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify that setup is no longer required
		respReq := app.DoRequest(t, http.MethodGet, "/api/v1/security/master-password/required-status", nil, true)
		assert.Equal(t, http.StatusOK, respReq.StatusCode)
		var result map[string]interface{}
		err := json.NewDecoder(respReq.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, false, result["isRequired"])
	})

	t.Run("Setup failure - password too short", func(t *testing.T) {
		body := []byte(`{"masterPassword":"short"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/security/master-password", body, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Setup failure - already configured", func(t *testing.T) {
		// Second setup attempt on already initialized vault
		body := []byte(`{"masterPassword":"anotherPassword123"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/security/master-password", body, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Setup failure - unauthorized", func(t *testing.T) {
		body := []byte(`{"masterPassword":"mySuperSecretPassword123"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/security/master-password", body, false)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestSecurityHandler_CheckMasterPasswordSetup(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	t.Run("Returns true when not configured", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/security/master-password/required-status", nil, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, true, result["isRequired"])
	})
}

func TestSecurityHandler_SessionStatusAndLockUnlock(t *testing.T) {
	app := SetupTestApp(t)
	defer app.Server.Close()

	t.Run("Unlock failure - not configured", func(t *testing.T) {
		body := []byte(`{"masterPassword":"mySuperSecretPassword123"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/security/vault/unlock", body, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Session status - locked by default", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodGet, "/api/v1/security/vault/status", nil, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var status map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&status)
		require.NoError(t, err)
		assert.Equal(t, false, status["active"])
		assert.Equal(t, float64(0), status["secondsLeft"])
	})

	// Setup master password for unlock/lock tests
	setupBody := []byte(`{"masterPassword":"mySuperSecretPassword123"}`)
	respSetup := app.DoRequest(t, http.MethodPost, "/api/v1/security/master-password", setupBody, true)
	require.Equal(t, http.StatusNoContent, respSetup.StatusCode)

	t.Run("Unlock failure - wrong password", func(t *testing.T) {
		body := []byte(`{"masterPassword":"wrongPassword123"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/security/vault/unlock", body, true)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Unlock failure - password too short", func(t *testing.T) {
		body := []byte(`{"masterPassword":"short"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/security/vault/unlock", body, true)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Unlock success - correct password", func(t *testing.T) {
		// First, lock vault so we test unlocking
		_ = app.DoRequest(t, http.MethodPost, "/api/v1/security/vault/lock", nil, true)

		body := []byte(`{"masterPassword":"mySuperSecretPassword123"}`)
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/security/vault/unlock", body, true)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result bool
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.True(t, result)

		// Verify session status is active
		respStatus := app.DoRequest(t, http.MethodGet, "/api/v1/security/vault/status", nil, true)
		assert.Equal(t, http.StatusOK, respStatus.StatusCode)
		var status map[string]interface{}
		_ = json.NewDecoder(respStatus.Body).Decode(&status)
		assert.Equal(t, true, status["active"])
		assert.Greater(t, status["secondsLeft"].(float64), float64(0))
	})

	t.Run("Lock success - clears active session", func(t *testing.T) {
		resp := app.DoRequest(t, http.MethodPost, "/api/v1/security/vault/lock", nil, true)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify session status is inactive
		respStatus := app.DoRequest(t, http.MethodGet, "/api/v1/security/vault/status", nil, true)
		assert.Equal(t, http.StatusOK, respStatus.StatusCode)
		var status map[string]interface{}
		_ = json.NewDecoder(respStatus.Body).Decode(&status)
		assert.Equal(t, false, status["active"])
		assert.Equal(t, float64(0), status["secondsLeft"])
	})
}
