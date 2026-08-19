package usecase

import (
	"context"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"
	"devaulty-backend/internal/dto"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

type ReleaseUseCase struct {
	releasePort port.ReleasePort
	updaterPort port.AppUpdater
}

func NewReleaseUseCase(releasePort port.ReleasePort, updaterPort port.AppUpdater) *ReleaseUseCase {
	return &ReleaseUseCase{
		releasePort: releasePort,
		updaterPort: updaterPort,
	}
}

func (uc *ReleaseUseCase) GetCurrentVersion() dto.CurrentVersionView {
	return dto.CurrentVersionView{
		CurrentVersion: model.AppVersion,
	}
}

func (uc *ReleaseUseCase) CheckUpdates(ctx context.Context) (*dto.AppUpdateInfoResponse, error) {
	currentVersion := model.AppVersion
	relInfo, err := uc.releasePort.GetLatestRelease(ctx)
	if err != nil {
		return nil, err
	}

	isAvailable := isNewerVersion(currentVersion, relInfo.TagName)
	selectedAsset := uc.selectMatchingAsset(relInfo.Assets)

	var downloadUrl string
	var downloadSize int64
	if selectedAsset != nil {
		downloadUrl = selectedAsset.DownloadUrl
		downloadSize = selectedAsset.SizeInBytes
	}

	return &dto.AppUpdateInfoResponse{
		UpdateAvailable:     isAvailable,
		CurrentVersion:      currentVersion,
		LatestVersion:       relInfo.TagName,
		ReleaseTitle:        relInfo.Name,
		ReleaseNotes:        relInfo.Body,
		DownloadUrl:         downloadUrl,
		DownloadSizeInBytes: downloadSize,
		PublishedAt:         relInfo.PublishedAt,
	}, nil
}

func (uc *ReleaseUseCase) DownloadAndInstall(ctx context.Context, progressCb func(view dto.UpdateDownloadProgressView)) error {
	relInfo, err := uc.releasePort.GetLatestRelease(ctx)
	if err != nil {
		progressCb(dto.UpdateDownloadProgressView{
			Status:       dto.StatusFailed,
			ErrorMessage: fmt.Sprintf("error fetching release info: %s", err.Error()),
		})
		return err
	}

	selectedAsset := uc.selectMatchingAsset(relInfo.Assets)
	if selectedAsset == nil {
		progressCb(dto.UpdateDownloadProgressView{
			Status:       dto.StatusFailed,
			ErrorMessage: "no matching asset found",
		})
		return fmt.Errorf("no matching asset found")
	}

	tempFile, err := os.CreateTemp("", "devaulty-update-*.tmp")
	if err != nil {
		progressCb(dto.UpdateDownloadProgressView{
			Status:       dto.StatusFailed,
			ErrorMessage: fmt.Sprintf("error creating temporary file: %s", err.Error()),
		})
		return err
	}

	tempPath := tempFile.Name()
	_ = tempFile.Close()

	err = uc.releasePort.DownloadAsset(ctx, selectedAsset.DownloadUrl, tempPath, func(downloadedBytes, totalBytes int64) {
		if totalBytes <= 0 {
			totalBytes = selectedAsset.SizeInBytes
		}

		percentage := 0
		if totalBytes > 0 {
			percentage = int((float64(downloadedBytes) / float64(totalBytes)) * 100)
		}

		progressCb(dto.UpdateDownloadProgressView{
			Status:          dto.StatusDownloading,
			Percentage:      percentage,
			DownloadedBytes: downloadedBytes,
			TotalBytes:      totalBytes,
		})
	})

	if err != nil {
		_ = os.Remove(tempPath)
		progressCb(dto.UpdateDownloadProgressView{
			Status:       dto.StatusFailed,
			ErrorMessage: fmt.Sprintf("download Error: %v", err),
		})
		return err
	}

	// Status INSTALLING
	progressCb(dto.UpdateDownloadProgressView{
		Status:          dto.StatusInstalling,
		Percentage:      100,
		DownloadedBytes: selectedAsset.SizeInBytes,
		TotalBytes:      selectedAsset.SizeInBytes,
	})

	if err := uc.updaterPort.InstallUpdate(ctx, tempPath); err != nil {
		progressCb(dto.UpdateDownloadProgressView{
			Status:       dto.StatusFailed,
			ErrorMessage: fmt.Sprintf("error trying to apply the update: %v", err),
		})
		return err
	}

	// Status COMPLETED
	progressCb(dto.UpdateDownloadProgressView{
		Status:          dto.StatusCompleted,
		Percentage:      100,
		DownloadedBytes: selectedAsset.SizeInBytes,
		TotalBytes:      selectedAsset.SizeInBytes,
	})

	return uc.updaterPort.RestartApp()
}

func (uc *ReleaseUseCase) selectMatchingAsset(assets []port.ReleaseAssetInfo) *port.ReleaseAssetInfo {
	osName := runtime.GOOS

	var extension string

	switch osName {
	case "windows":
		extension = ".msi"
	case "darwin":
		extension = ".dmg"
	case "linux":
		extension = detectLinuxExtension()
		if extension == "" {
			return nil
		}
	default:
		log.Printf("unsupported OS: %s", osName)
		return nil
	}

	for _, a := range assets {
		if strings.HasSuffix(strings.ToLower(a.FileName), extension) {
			return &a
		}
	}

	return nil
}

func detectLinuxExtension() string {
	content, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ".deb" // safe callback
	}

	contentLower := strings.ToLower(string(content))

	isRpmBased := strings.Contains(contentLower, "rhel") ||
		strings.Contains(contentLower, "fedora") ||
		strings.Contains(contentLower, "suse")

	isDebBased := strings.Contains(contentLower, "debian") ||
		strings.Contains(contentLower, "ubuntu")

	if isRpmBased {
		return ".rpm"
	}

	if isDebBased {
		return ".deb"
	}

	log.Printf("could not detect Linux distribution from /etc/os-release")
	return ".deb"
}

func isNewerVersion(currentVersion, latestVersion string) bool {
	cleanCurrent := strings.TrimPrefix(currentVersion, "v")
	cleanLatest := strings.TrimPrefix(latestVersion, "v")
	return cleanLatest != "" && cleanLatest != cleanCurrent
}
