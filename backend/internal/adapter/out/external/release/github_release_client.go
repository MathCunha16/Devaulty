package release

import (
	"context"
	"devaulty-backend/internal/domain/port"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	ConnectTimeout           = time.Second * 10
	ApiReadWriteTimeout      = time.Second * 15
	DownloadReadWriteTimeout = time.Minute * 10
	GitHubOwnerName          = "MathCunha16"
	GitHubRepoName           = "Devaulty"
)

type githubReleaseResponse struct {
	TagName     string        `json:"tag_name"`
	Name        string        `json:"name"`
	Body        string        `json:"body"`
	HtmlUrl     string        `json:"html_url"`
	PublishedAt string        `json:"published_at"`
	Prerelease  bool          `json:"prerelease"`
	Assets      []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadUrl string `json:"browser_download_url"`
	Size               int64  `json:"size"`
	ContentType        string `json:"content_type"`
}

type githubReleaseClient struct {
	httpClient *http.Client
}

func NewGitHubReleaseClient() port.ReleasePort {
	return &githubReleaseClient{
		httpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: ConnectTimeout,
				}).DialContext,
			},
		},
	}
}

func (g *githubReleaseClient) GetLatestRelease(ctx context.Context) (*port.LatestReleaseInfo, error) {
	apiCtx, cancel := context.WithTimeout(ctx, ApiReadWriteTimeout)
	defer cancel()

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", GitHubOwnerName, GitHubRepoName)
	req, err := http.NewRequestWithContext(apiCtx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Devaulty-Desktop-App")
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get latest release: %s", resp.Status)
	}

	var rel githubReleaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("failed to decode latest release: %w", err)
	}

	assets := make([]port.ReleaseAssetInfo, len(rel.Assets))
	for i, a := range rel.Assets {
		assets[i] = port.ReleaseAssetInfo{
			FileName:    a.Name,
			DownloadUrl: a.BrowserDownloadUrl,
			SizeInBytes: a.Size,
			ContentType: a.ContentType,
		}
	}
	return &port.LatestReleaseInfo{
		TagName:      rel.TagName,
		Name:         rel.Name,
		Body:         rel.Body,
		HtmlUrl:      rel.HtmlUrl,
		PublishedAt:  rel.PublishedAt,
		IsPreRelease: rel.Prerelease,
		Assets:       assets,
	}, nil
}

func (g *githubReleaseClient) DownloadAsset(ctx context.Context, downloadUrl string, destinationPath string, progressCb port.ReleaseProgressCallback) error {
	downloadCtx, cancel := context.WithTimeout(ctx, DownloadReadWriteTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(downloadCtx, http.MethodGet, downloadUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Devaulty-Desktop-App")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status downloading asset: %s", resp.Status)
	}

	out, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	totalBytes := resp.ContentLength
	buffer := make([]byte, 32*1024) // 32kb
	var downloadedBytes int64

	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			_, werr := out.Write(buffer[:n])
			if werr != nil {
				return fmt.Errorf("failed to write buffer to file: %w", werr)
			}
			downloadedBytes += int64(n)

			if progressCb != nil {
				progressCb(downloadedBytes, totalBytes)
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading stream: %w", err)
		}
	}

	return nil
}
