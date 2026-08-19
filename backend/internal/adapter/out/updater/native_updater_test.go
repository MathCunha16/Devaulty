package updater_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"devaulty-backend/internal/adapter/out/updater"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- UNIT TESTS FOR NATIVE UPDATER ---

func TestNativeUpdater_NewNativeUpdater(t *testing.T) {
	u := updater.NewNativeUpdater()
	assert.NotNil(t, u)
}

func TestNativeUpdater_CleanupResidualFiles(t *testing.T) {
	u := updater.NewNativeUpdater()

	// Create a dummy temp update file in os.TempDir() to test cleanup
	tempFile, err := os.CreateTemp("", "devaulty-update-test-*.tmp")
	require.NoError(t, err)
	tempPath := tempFile.Name()
	_ = tempFile.Close()

	// Guarantee cleanup in case test fails early
	t.Cleanup(func() {
		_ = os.Remove(tempPath)
	})

	// Verify temp file exists before cleanup
	_, err = os.Stat(tempPath)
	require.NoError(t, err)

	// Execute residual cleanup
	err = u.CleanupResidualFiles()
	assert.NoError(t, err)

	// Verify temp file was removed by CleanupResidualFiles
	_, err = os.Stat(tempPath)
	assert.True(t, os.IsNotExist(err), "Temp update file should be removed by CleanupResidualFiles")
}

func TestNativeUpdater_InstallUpdate_NonExistentTempFile(t *testing.T) {
	u := updater.NewNativeUpdater()
	ctx := context.Background()

	nonExistentPath := filepath.Join(os.TempDir(), "nonexistent_devaulty_update_file_98765.tmp")

	err := u.InstallUpdate(ctx, nonExistentPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error trying to install new binary")
}
