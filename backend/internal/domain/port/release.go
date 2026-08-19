package port

import "context"

type ReleaseAssetInfo struct {
	FileName    string
	DownloadUrl string
	SizeInBytes int64
	ContentType string
}

type LatestReleaseInfo struct {
	TagName      string
	Name         string
	Body         string
	HtmlUrl      string
	PublishedAt  string
	IsPreRelease bool
	Assets       []ReleaseAssetInfo
}

type ReleaseProgressCallback func(downloadedBytes, totalBytes int64)

type ReleasePort interface {
	GetLatestRelease(ctx context.Context) (*LatestReleaseInfo, error)
	DownloadAsset(ctx context.Context, downloadUrl string, destinationPath string, progressCb ReleaseProgressCallback) error
}

type AppUpdater interface {
	InstallUpdate(ctx context.Context, tempDownloadFilePath string) error
	RestartApp() error
	CleanupResidualFiles() error
}
