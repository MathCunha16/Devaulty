package handler

import (
	"devaulty-backend/internal/dto"
	"devaulty-backend/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReleaseHandler struct {
	releaseUseCase *usecase.ReleaseUseCase
}

func NewReleaseHandler(releaseUseCase *usecase.ReleaseUseCase) *ReleaseHandler {
	return &ReleaseHandler{releaseUseCase: releaseUseCase}
}

func (h *ReleaseHandler) GetCurrentVersion(c *gin.Context) {
	c.JSON(http.StatusOK, h.releaseUseCase.GetCurrentVersion())
}

func (h *ReleaseHandler) CheckUpdates(c *gin.Context) {
	response, err := h.releaseUseCase.CheckUpdates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *ReleaseHandler) DownloadAndInstall(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	progressCh := make(chan dto.UpdateDownloadProgressView, 100)
	errCh := make(chan error, 1)

	go func() {
		errCh <- h.releaseUseCase.DownloadAndInstall(
			c.Request.Context(),
			func(view dto.UpdateDownloadProgressView) {
				progressCh <- view
			},
		)

		close(progressCh)
	}()

	for {
		select {
		case view, ok := <-progressCh:
			if !ok {
				err := <-errCh
				if err != nil {
					return
				}
				return
			}

			c.SSEvent("progress", view)
			c.Writer.Flush()

		case <-c.Request.Context().Done():
			return
		}
	}
}
