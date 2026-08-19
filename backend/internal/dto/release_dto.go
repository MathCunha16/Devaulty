package dto

type CurrentVersionView struct {
	CurrentVersion string `json:"currentVersion"`
}

type AppUpdateInfoResponse struct {
	UpdateAvailable     bool   `json:"updateAvailable"`
	CurrentVersion      string `json:"currentVersion"`
	LatestVersion       string `json:"latestVersion"`
	ReleaseTitle        string `json:"releaseTitle,omitempty"`
	ReleaseNotes        string `json:"releaseNotes,omitempty"`
	DownloadUrl         string `json:"downloadUrl,omitempty"`
	DownloadSizeInBytes int64  `json:"downloadSizeInBytes,omitempty"`
	PublishedAt         string `json:"publishedAt,omitempty"`
}

type UpdateDownloadStatus string

const (
	StatusDownloading UpdateDownloadStatus = "DOWNLOADING"
	StatusInstalling  UpdateDownloadStatus = "INSTALLING"
	StatusCompleted   UpdateDownloadStatus = "COMPLETED"
	StatusFailed      UpdateDownloadStatus = "FAILED"
)

type UpdateDownloadProgressView struct {
	Status          UpdateDownloadStatus `json:"status"`
	Percentage      int                  `json:"percentage"`
	DownloadedBytes int64                `json:"downloadedBytes"`
	TotalBytes      int64                `json:"totalBytes"`
	ErrorMessage    string               `json:"errorMessage,omitempty"`
}
