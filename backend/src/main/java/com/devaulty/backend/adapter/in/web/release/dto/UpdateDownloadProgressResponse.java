package com.devaulty.backend.adapter.in.web.release.dto;

import com.devaulty.backend.application.port.in.release.enums.UpdateStatus;

public record UpdateDownloadProgressResponse(
        UpdateStatus status,
        int percentage,
        long downloadedBytes,
        long totalBytes,
        String errorMessage
) {
}
