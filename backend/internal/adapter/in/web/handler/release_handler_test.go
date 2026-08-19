package handler_test

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"devaulty-backend/internal/adapter/in/web"
	"devaulty-backend/internal/adapter/in/web/handler"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"
	"devaulty-backend/internal/dto"
	"devaulty-backend/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- MOCKS FOR INTEGRATION TESTS ---

type MockReleasePort struct {
	mock.Mock
}

func (m *MockReleasePort) GetLatestRelease(ctx context.Context) (*port.LatestReleaseInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*port.LatestReleaseInfo), args.Error(1)
}

func (m *MockReleasePort) DownloadAsset(ctx context.Context, downloadUrl string, destinationPath string, progressCb port.ReleaseProgressCallback) error {
	args := m.Called(ctx, downloadUrl, destinationPath, progressCb)
	return args.Error(0)
}

type MockAppUpdater struct {
	mock.Mock
}

func (m *MockAppUpdater) InstallUpdate(ctx context.Context, tempDownloadFilePath string) error {
	args := m.Called(ctx, tempDownloadFilePath)
	return args.Error(0)
}

func (m *MockAppUpdater) RestartApp() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockAppUpdater) CleanupResidualFiles() error {
	args := m.Called()
	return args.Error(0)
}

// --- TEST APP SETUP WITH MOCKS ---

type ReleaseTestApp struct {
	Server          *httptest.Server
	Token           string
	MockReleasePort *MockReleasePort
	MockUpdater     *MockAppUpdater
}

func SetupReleaseTestApp(t *testing.T) *ReleaseTestApp {
	mockPort := new(MockReleasePort)
	mockUpdater := new(MockAppUpdater)

	uc := usecase.NewReleaseUseCase(mockPort, mockUpdater)
	releaseHandler := handler.NewReleaseHandler(uc)

	handlers := &web.Handlers{
		Release: releaseHandler,
	}

	token := "test-release-token-999"
	router := web.SetupRouter(handlers, token)
	ts := httptest.NewServer(router)

	t.Cleanup(func() {
		ts.Close()
	})

	return &ReleaseTestApp{
		Server:          ts,
		Token:           token,
		MockReleasePort: mockPort,
		MockUpdater:     mockUpdater,
	}
}

// --- INTEGRATION TESTS ---

func TestReleaseHandler_GetCurrentVersion(t *testing.T) {
	t.Run("Get current version - Success", func(t *testing.T) {
		app := SetupTestApp(t)
		defer app.Server.Close()

		resp := app.DoRequest(t, http.MethodGet, "/api/v1/releases/current-app-version", nil, true)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result dto.CurrentVersionView
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, model.AppVersion, result.CurrentVersion)
	})

	t.Run("Get current version - Unauthorized missing token", func(t *testing.T) {
		app := SetupTestApp(t)
		defer app.Server.Close()

		resp := app.DoRequest(t, http.MethodGet, "/api/v1/releases/current-app-version", nil, false)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestReleaseHandler_CheckUpdates(t *testing.T) {
	t.Run("Check updates - Success update available", func(t *testing.T) {
		app := SetupReleaseTestApp(t)

		sampleRelease := &port.LatestReleaseInfo{
			TagName:     "v9.9.9",
			Name:        "Release v9.9.9",
			Body:        "New features and bug fixes",
			HtmlUrl:     "https://github.com/owner/repo/releases/tag/v9.9.9",
			PublishedAt: "2026-03-10T12:00:00Z",
			Assets: []port.ReleaseAssetInfo{
				{
					FileName:    "Devaulty_9.9.9_amd64.deb",
					DownloadUrl: "https://github.com/owner/repo/releases/download/v9.9.9/Devaulty_9.9.9_amd64.deb",
					SizeInBytes: 15000000,
					ContentType: "application/octet-stream",
				},
			},
		}

		app.MockReleasePort.On("GetLatestRelease", mock.Anything).Return(sampleRelease, nil)

		req, err := http.NewRequest(http.MethodGet, app.Server.URL+"/api/v1/releases/check", nil)
		require.NoError(t, err)
		req.Header.Set("DEVAULTY_INTERNAL_TOKEN", app.Token)

		resp, err := app.Server.Client().Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result dto.AppUpdateInfoResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.True(t, result.UpdateAvailable)
		assert.Equal(t, "v9.9.9", result.LatestVersion)
		assert.Equal(t, "Release v9.9.9", result.ReleaseTitle)
		assert.Equal(t, "New features and bug fixes", result.ReleaseNotes)
		assert.Equal(t, int64(15000000), result.DownloadSizeInBytes)
		app.MockReleasePort.AssertExpectations(t)
	})

	t.Run("Check updates - Internal server error when upstream fails", func(t *testing.T) {
		app := SetupReleaseTestApp(t)

		app.MockReleasePort.On("GetLatestRelease", mock.Anything).Return(nil, errors.New("github rate limit reached"))

		req, err := http.NewRequest(http.MethodGet, app.Server.URL+"/api/v1/releases/check", nil)
		require.NoError(t, err)
		req.Header.Set("DEVAULTY_INTERNAL_TOKEN", app.Token)

		resp, err := app.Server.Client().Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var errResult map[string]string
		err = json.NewDecoder(resp.Body).Decode(&errResult)
		require.NoError(t, err)
		assert.Contains(t, errResult["error"], "github rate limit reached")
		app.MockReleasePort.AssertExpectations(t)
	})

	t.Run("Check updates - Unauthorized missing token", func(t *testing.T) {
		app := SetupReleaseTestApp(t)

		req, err := http.NewRequest(http.MethodGet, app.Server.URL+"/api/v1/releases/check", nil)
		require.NoError(t, err)

		resp, err := app.Server.Client().Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestReleaseHandler_DownloadAndInstall(t *testing.T) {
	t.Run("Download and install - SSE streaming success", func(t *testing.T) {
		app := SetupReleaseTestApp(t)

		sampleRelease := &port.LatestReleaseInfo{
			TagName: "v9.9.9",
			Assets: []port.ReleaseAssetInfo{
				{
					FileName:    "Devaulty_9.9.9_amd64.deb",
					DownloadUrl: "https://github.com/owner/repo/releases/download/v9.9.9/Devaulty_9.9.9_amd64.deb",
					SizeInBytes: 20000000,
				},
			},
		}

		app.MockReleasePort.On("GetLatestRelease", mock.Anything).Return(sampleRelease, nil)

		expectedUrl := sampleRelease.Assets[0].DownloadUrl
		app.MockReleasePort.On("DownloadAsset", mock.Anything, expectedUrl, mock.AnythingOfType("string"), mock.Anything).
			Run(func(args mock.Arguments) {
				destPath := args.String(2)
				cb := args.Get(3).(port.ReleaseProgressCallback)

				// Teardown guarantee: remove temp file created during test
				t.Cleanup(func() {
					_ = os.Remove(destPath)
				})

				cb(10000000, 20000000) // 50%
				cb(20000000, 20000000) // 100%
			}).
			Return(nil)

		app.MockUpdater.On("InstallUpdate", mock.Anything, mock.AnythingOfType("string")).Return(nil)
		app.MockUpdater.On("RestartApp").Return(nil)

		req, err := http.NewRequest(http.MethodPost, app.Server.URL+"/api/v1/releases/download-and-install", nil)
		require.NoError(t, err)
		req.Header.Set("DEVAULTY_INTERNAL_TOKEN", app.Token)
		req.Header.Set("Accept", "text/event-stream")

		resp, err := app.Server.Client().Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("Content-Type"), "text/event-stream")

		// Read SSE stream output
		reader := bufio.NewReader(resp.Body)
		var receivedProgressEvents []dto.UpdateDownloadProgressView

		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			require.NoError(t, err)

			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "data:") {
				jsonPayload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
				var progressView dto.UpdateDownloadProgressView
				if err := json.Unmarshal([]byte(jsonPayload), &progressView); err == nil {
					receivedProgressEvents = append(receivedProgressEvents, progressView)
				}
			}
		}

		assert.NotEmpty(t, receivedProgressEvents)

		// Verify event sequence contains status values
		var statuses []dto.UpdateDownloadStatus
		for _, event := range receivedProgressEvents {
			statuses = append(statuses, event.Status)
		}
		assert.Contains(t, statuses, dto.StatusDownloading)
		assert.Contains(t, statuses, dto.StatusInstalling)
		assert.Contains(t, statuses, dto.StatusCompleted)

		app.MockReleasePort.AssertExpectations(t)
		app.MockUpdater.AssertExpectations(t)
	})

	t.Run("Download and install - Unauthorized missing token", func(t *testing.T) {
		app := SetupReleaseTestApp(t)

		req, err := http.NewRequest(http.MethodPost, app.Server.URL+"/api/v1/releases/download-and-install", nil)
		require.NoError(t, err)

		resp, err := app.Server.Client().Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Download and install - Stream handles download error gracefully", func(t *testing.T) {
		app := SetupReleaseTestApp(t)

		sampleRelease := &port.LatestReleaseInfo{
			TagName: "v9.9.9",
			Assets: []port.ReleaseAssetInfo{
				{
					FileName:    "Devaulty_9.9.9_amd64.deb",
					DownloadUrl: "https://github.com/owner/repo/releases/download/v9.9.9/Devaulty_9.9.9_amd64.deb",
					SizeInBytes: 20000000,
				},
			},
		}

		app.MockReleasePort.On("GetLatestRelease", mock.Anything).Return(sampleRelease, nil)

		expectedUrl := sampleRelease.Assets[0].DownloadUrl
		app.MockReleasePort.On("DownloadAsset", mock.Anything, expectedUrl, mock.AnythingOfType("string"), mock.Anything).
			Run(func(args mock.Arguments) {
				destPath := args.String(2)
				t.Cleanup(func() {
					_ = os.Remove(destPath)
				})
			}).
			Return(errors.New("connection reset by peer"))

		req, err := http.NewRequest(http.MethodPost, app.Server.URL+"/api/v1/releases/download-and-install", nil)
		require.NoError(t, err)
		req.Header.Set("DEVAULTY_INTERNAL_TOKEN", app.Token)
		req.Header.Set("Accept", "text/event-stream")

		resp, err := app.Server.Client().Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		reader := bufio.NewReader(resp.Body)
		var lastStatus dto.UpdateDownloadProgressView

		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			require.NoError(t, err)

			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "data:") {
				jsonPayload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
				_ = json.Unmarshal([]byte(jsonPayload), &lastStatus)
			}
		}

		assert.Equal(t, dto.StatusFailed, lastStatus.Status)
		assert.Contains(t, lastStatus.ErrorMessage, "error downloading asset")
		app.MockReleasePort.AssertExpectations(t)
	})

	t.Run("Download and install - Context cancellation on client disconnect", func(t *testing.T) {
		app := SetupReleaseTestApp(t)

		sampleRelease := &port.LatestReleaseInfo{
			TagName: "v9.9.9",
			Assets: []port.ReleaseAssetInfo{
				{
					FileName:    "Devaulty_9.9.9_amd64.deb",
					DownloadUrl: "https://github.com/owner/repo/releases/download/v9.9.9/Devaulty_9.9.9_amd64.deb",
				},
			},
		}

		app.MockReleasePort.On("GetLatestRelease", mock.Anything).Return(sampleRelease, nil)
		app.MockReleasePort.On("DownloadAsset", mock.Anything, mock.Anything, mock.AnythingOfType("string"), mock.Anything).
			Run(func(args mock.Arguments) {
				destPath := args.String(2)
				t.Cleanup(func() {
					_ = os.Remove(destPath)
				})
				time.Sleep(50 * time.Millisecond)
			}).
			Return(errors.New("context canceled"))

		req, err := http.NewRequest(http.MethodPost, app.Server.URL+"/api/v1/releases/download-and-install", nil)
		require.NoError(t, err)
		req.Header.Set("DEVAULTY_INTERNAL_TOKEN", app.Token)

		client := &http.Client{Timeout: 30 * time.Millisecond}
		resp, err := client.Do(req)
		if err == nil {
			_ = resp.Body.Close()
		}

		// Ensure no goroutines hung or blocked
		time.Sleep(50 * time.Millisecond)
		app.MockReleasePort.AssertExpectations(t)
	})
}
