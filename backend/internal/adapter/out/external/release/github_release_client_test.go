package release_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"devaulty-backend/internal/adapter/out/external/release"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- UNIT / INTEGRATION TESTS FOR GITHUB RELEASE CLIENT ---

func TestGitHubReleaseClient_DownloadAsset_Success(t *testing.T) {
	// Spin up a local HTTP test server simulating GitHub CDN binary asset download
	sampleBinaryData := []byte("binary-payload-data-for-devaulty-update-installer")

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Devaulty-Desktop-App", r.Header.Get("User-Agent"))
		assert.Empty(t, r.Header.Get("Accept"), "Download request must not contain API JSON Accept header")

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(sampleBinaryData)
	}))
	t.Cleanup(func() {
		testServer.Close()
	})

	client := release.NewGitHubReleaseClient()
	ctx := context.Background()

	tempDir := t.TempDir()
	destinationPath := filepath.Join(tempDir, "test_downloaded_release.deb")

	t.Cleanup(func() {
		_ = os.Remove(destinationPath)
	})

	var progressReports []int64
	progressCb := func(downloadedBytes, totalBytes int64) {
		progressReports = append(progressReports, downloadedBytes)
	}

	downloadUrl := testServer.URL + "/releases/download/v9.9.9/Devaulty_9.9.9_amd64.deb"
	err := client.DownloadAsset(ctx, downloadUrl, destinationPath, progressCb)

	assert.NoError(t, err)
	assert.NotEmpty(t, progressReports)

	// Verify file was written to disk correctly
	writtenBytes, err := os.ReadFile(destinationPath)
	require.NoError(t, err)
	assert.Equal(t, sampleBinaryData, writtenBytes)
}

func TestGitHubReleaseClient_DownloadAsset_HTTPError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(func() {
		testServer.Close()
	})

	client := release.NewGitHubReleaseClient()
	ctx := context.Background()

	tempDir := t.TempDir()
	destinationPath := filepath.Join(tempDir, "test_failed_download.tmp")
	t.Cleanup(func() {
		_ = os.Remove(destinationPath)
	})

	downloadUrl := testServer.URL + "/releases/download/v9.9.9/nonexistent.deb"
	err := client.DownloadAsset(ctx, downloadUrl, destinationPath, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected HTTP status")
}

func TestGitHubReleaseClient_DownloadAsset_InvalidDestinationPath(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("data"))
	}))
	t.Cleanup(func() {
		testServer.Close()
	})

	client := release.NewGitHubReleaseClient()
	ctx := context.Background()

	invalidPath := "/nonexistent_folder_12345/impossible_file.tmp"
	downloadUrl := testServer.URL + "/test.deb"

	err := client.DownloadAsset(ctx, downloadUrl, invalidPath, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create destination file")
}

func TestGitHubReleaseClient_DownloadAsset_ContextCanceled(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(func() {
		testServer.Close()
	})

	client := release.NewGitHubReleaseClient()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel context immediately

	tempDir := t.TempDir()
	destinationPath := filepath.Join(tempDir, "test_canceled_download.tmp")
	t.Cleanup(func() {
		_ = os.Remove(destinationPath)
	})

	downloadUrl := testServer.URL + "/test.deb"
	err := client.DownloadAsset(ctx, downloadUrl, destinationPath, nil)

	assert.Error(t, err)
}
