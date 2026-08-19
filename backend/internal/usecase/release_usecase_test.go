package usecase_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"
	"devaulty-backend/internal/dto"
	"devaulty-backend/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MOCKS ---

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

// --- HELPER FIXTURES ---

func createSampleReleaseInfo(tag string, assetName string) *port.LatestReleaseInfo {
	return &port.LatestReleaseInfo{
		TagName:      tag,
		Name:         "Release " + tag,
		Body:         "Release notes for " + tag,
		HtmlUrl:      "https://github.com/owner/repo/releases/tag/" + tag,
		PublishedAt:  "2026-03-10T12:00:00Z",
		IsPreRelease: false,
		Assets: []port.ReleaseAssetInfo{
			{
				FileName:    assetName,
				DownloadUrl: "https://github.com/owner/repo/releases/download/" + tag + "/" + assetName,
				SizeInBytes: 10485760,
				ContentType: "application/octet-stream",
			},
		},
	}
}

// --- UNIT TESTS ---

func TestReleaseUseCase_GetCurrentVersion(t *testing.T) {
	mockReleasePort := new(MockReleasePort)
	mockUpdater := new(MockAppUpdater)
	uc := usecase.NewReleaseUseCase(mockReleasePort, mockUpdater)

	result := uc.GetCurrentVersion()

	assert.Equal(t, model.AppVersion, result.CurrentVersion)
}

func TestReleaseUseCase_CheckUpdates_UpdateAvailable(t *testing.T) {
	mockReleasePort := new(MockReleasePort)
	mockUpdater := new(MockAppUpdater)
	uc := usecase.NewReleaseUseCase(mockReleasePort, mockUpdater)
	ctx := context.Background()

	sampleRelease := createSampleReleaseInfo("v99.0.0", "Devaulty_99.0.0_amd64.deb")
	mockReleasePort.On("GetLatestRelease", ctx).Return(sampleRelease, nil)

	res, err := uc.CheckUpdates(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.UpdateAvailable)
	assert.Equal(t, model.AppVersion, res.CurrentVersion)
	assert.Equal(t, "v99.0.0", res.LatestVersion)
	assert.Equal(t, "Release v99.0.0", res.ReleaseTitle)
	assert.Equal(t, "Release notes for v99.0.0", res.ReleaseNotes)
	assert.Equal(t, int64(10485760), res.DownloadSizeInBytes)
	mockReleasePort.AssertExpectations(t)
}

func TestReleaseUseCase_CheckUpdates_NoUpdateAvailable(t *testing.T) {
	mockReleasePort := new(MockReleasePort)
	mockUpdater := new(MockAppUpdater)
	uc := usecase.NewReleaseUseCase(mockReleasePort, mockUpdater)
	ctx := context.Background()

	sampleRelease := createSampleReleaseInfo(model.AppVersion, "Devaulty_amd64.deb")
	mockReleasePort.On("GetLatestRelease", ctx).Return(sampleRelease, nil)

	res, err := uc.CheckUpdates(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.False(t, res.UpdateAvailable)
	assert.Equal(t, model.AppVersion, res.CurrentVersion)
	assert.Equal(t, model.AppVersion, res.LatestVersion)
	mockReleasePort.AssertExpectations(t)
}

func TestReleaseUseCase_CheckUpdates_GetLatestReleaseError(t *testing.T) {
	mockReleasePort := new(MockReleasePort)
	mockUpdater := new(MockAppUpdater)
	uc := usecase.NewReleaseUseCase(mockReleasePort, mockUpdater)
	ctx := context.Background()

	expectedErr := errors.New("github api rate limit exceeded")
	mockReleasePort.On("GetLatestRelease", ctx).Return(nil, expectedErr)

	res, err := uc.CheckUpdates(ctx)

	assert.Nil(t, res)
	assert.ErrorIs(t, err, expectedErr)
	mockReleasePort.AssertExpectations(t)
}

func TestReleaseUseCase_DownloadAndInstall_Success(t *testing.T) {
	mockReleasePort := new(MockReleasePort)
	mockUpdater := new(MockAppUpdater)
	uc := usecase.NewReleaseUseCase(mockReleasePort, mockUpdater)
	ctx := context.Background()

	sampleRelease := createSampleReleaseInfo("v99.0.0", "Devaulty_99.0.0_amd64.deb")
	mockReleasePort.On("GetLatestRelease", ctx).Return(sampleRelease, nil)

	expectedUrl := sampleRelease.Assets[0].DownloadUrl

	mockReleasePort.On("DownloadAsset", ctx, expectedUrl, mock.AnythingOfType("string"), mock.Anything).
		Run(func(args mock.Arguments) {
			destinationPath := args.String(2)
			cb := args.Get(3).(port.ReleaseProgressCallback)

			// Ensure temporary files are cleaned up
			t.Cleanup(func() {
				_ = os.Remove(destinationPath)
			})

			// Simulate download progress
			cb(5242880, 10485760)  // 50%
			cb(10485760, 10485760) // 100%
		}).
		Return(nil)

	mockUpdater.On("InstallUpdate", ctx, mock.AnythingOfType("string")).Return(nil)
	mockUpdater.On("RestartApp").Return(nil)

	var progressEvents []dto.UpdateDownloadProgressView
	err := uc.DownloadAndInstall(ctx, func(view dto.UpdateDownloadProgressView) {
		progressEvents = append(progressEvents, view)
	})

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(progressEvents), 3)

	// Verify status transitions: DOWNLOADING -> INSTALLING -> COMPLETED
	statuses := make([]dto.UpdateDownloadStatus, len(progressEvents))
	for i, e := range progressEvents {
		statuses[i] = e.Status
	}
	assert.Contains(t, statuses, dto.StatusDownloading)
	assert.Contains(t, statuses, dto.StatusInstalling)
	assert.Contains(t, statuses, dto.StatusCompleted)

	mockReleasePort.AssertExpectations(t)
	mockUpdater.AssertExpectations(t)
}

func TestReleaseUseCase_DownloadAndInstall_GetLatestReleaseError(t *testing.T) {
	mockReleasePort := new(MockReleasePort)
	mockUpdater := new(MockAppUpdater)
	uc := usecase.NewReleaseUseCase(mockReleasePort, mockUpdater)
	ctx := context.Background()

	expectedErr := errors.New("network connection timeout")
	mockReleasePort.On("GetLatestRelease", ctx).Return(nil, expectedErr)

	var progressEvents []dto.UpdateDownloadProgressView
	err := uc.DownloadAndInstall(ctx, func(view dto.UpdateDownloadProgressView) {
		progressEvents = append(progressEvents, view)
	})

	assert.Error(t, err)
	assert.Equal(t, 1, len(progressEvents))
	assert.Equal(t, dto.StatusFailed, progressEvents[0].Status)
	assert.Contains(t, progressEvents[0].ErrorMessage, "error fetching release info")

	mockReleasePort.AssertExpectations(t)
}

func TestReleaseUseCase_DownloadAndInstall_NoMatchingAsset(t *testing.T) {
	mockReleasePort := new(MockReleasePort)
	mockUpdater := new(MockAppUpdater)
	uc := usecase.NewReleaseUseCase(mockReleasePort, mockUpdater)
	ctx := context.Background()

	// Release without matching asset
	releaseInfoNoAssets := &port.LatestReleaseInfo{
		TagName: "v99.0.0",
		Assets:  []port.ReleaseAssetInfo{},
	}
	mockReleasePort.On("GetLatestRelease", ctx).Return(releaseInfoNoAssets, nil)

	var progressEvents []dto.UpdateDownloadProgressView
	err := uc.DownloadAndInstall(ctx, func(view dto.UpdateDownloadProgressView) {
		progressEvents = append(progressEvents, view)
	})

	assert.Error(t, err)
	assert.Equal(t, "no matching asset found", err.Error())
	assert.Equal(t, 1, len(progressEvents))
	assert.Equal(t, dto.StatusFailed, progressEvents[0].Status)

	mockReleasePort.AssertExpectations(t)
}

func TestReleaseUseCase_DownloadAndInstall_DownloadAssetError(t *testing.T) {
	mockReleasePort := new(MockReleasePort)
	mockUpdater := new(MockAppUpdater)
	uc := usecase.NewReleaseUseCase(mockReleasePort, mockUpdater)
	ctx := context.Background()

	sampleRelease := createSampleReleaseInfo("v99.0.0", "Devaulty_99.0.0_amd64.deb")
	mockReleasePort.On("GetLatestRelease", ctx).Return(sampleRelease, nil)

	downloadErr := errors.New("download stream interrupted")
	mockReleasePort.On("DownloadAsset", ctx, mock.Anything, mock.AnythingOfType("string"), mock.Anything).
		Run(func(args mock.Arguments) {
			tempPath := args.String(2)
			t.Cleanup(func() {
				_ = os.Remove(tempPath)
			})
		}).
		Return(downloadErr)

	var progressEvents []dto.UpdateDownloadProgressView
	err := uc.DownloadAndInstall(ctx, func(view dto.UpdateDownloadProgressView) {
		progressEvents = append(progressEvents, view)
	})

	assert.Error(t, err)
	assert.Equal(t, dto.StatusFailed, progressEvents[len(progressEvents)-1].Status)
	assert.Contains(t, progressEvents[len(progressEvents)-1].ErrorMessage, "error downloading asset")

	mockReleasePort.AssertExpectations(t)
}

func TestReleaseUseCase_DownloadAndInstall_InstallUpdateError(t *testing.T) {
	mockReleasePort := new(MockReleasePort)
	mockUpdater := new(MockAppUpdater)
	uc := usecase.NewReleaseUseCase(mockReleasePort, mockUpdater)
	ctx := context.Background()

	sampleRelease := createSampleReleaseInfo("v99.0.0", "Devaulty_99.0.0_amd64.deb")
	mockReleasePort.On("GetLatestRelease", ctx).Return(sampleRelease, nil)

	mockReleasePort.On("DownloadAsset", ctx, mock.Anything, mock.AnythingOfType("string"), mock.Anything).
		Run(func(args mock.Arguments) {
			tempPath := args.String(2)
			t.Cleanup(func() {
				_ = os.Remove(tempPath)
			})
		}).
		Return(nil)

	installErr := errors.New("permission denied during binary copy")
	mockUpdater.On("InstallUpdate", ctx, mock.AnythingOfType("string")).Return(installErr)

	var progressEvents []dto.UpdateDownloadProgressView
	err := uc.DownloadAndInstall(ctx, func(view dto.UpdateDownloadProgressView) {
		progressEvents = append(progressEvents, view)
	})

	assert.Error(t, err)
	assert.Equal(t, dto.StatusFailed, progressEvents[len(progressEvents)-1].Status)
	assert.Contains(t, progressEvents[len(progressEvents)-1].ErrorMessage, "error applying update")

	mockReleasePort.AssertExpectations(t)
	mockUpdater.AssertExpectations(t)
}

func TestReleaseUseCase_DownloadAndInstall_RestartAppError(t *testing.T) {
	mockReleasePort := new(MockReleasePort)
	mockUpdater := new(MockAppUpdater)
	uc := usecase.NewReleaseUseCase(mockReleasePort, mockUpdater)
	ctx := context.Background()

	sampleRelease := createSampleReleaseInfo("v99.0.0", "Devaulty_99.0.0_amd64.deb")
	mockReleasePort.On("GetLatestRelease", ctx).Return(sampleRelease, nil)

	mockReleasePort.On("DownloadAsset", ctx, mock.Anything, mock.AnythingOfType("string"), mock.Anything).
		Run(func(args mock.Arguments) {
			tempPath := args.String(2)
			t.Cleanup(func() {
				_ = os.Remove(tempPath)
			})
		}).
		Return(nil)

	mockUpdater.On("InstallUpdate", ctx, mock.AnythingOfType("string")).Return(nil)

	restartErr := errors.New("failed to start new process")
	mockUpdater.On("RestartApp").Return(restartErr)

	var progressEvents []dto.UpdateDownloadProgressView
	err := uc.DownloadAndInstall(ctx, func(view dto.UpdateDownloadProgressView) {
		progressEvents = append(progressEvents, view)
	})

	assert.ErrorIs(t, err, restartErr)

	mockReleasePort.AssertExpectations(t)
	mockUpdater.AssertExpectations(t)
}
