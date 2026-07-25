package com.devaulty.backend.adapter.in.web.release.dto;

import java.time.LocalDateTime;

public record AppUpdateInfoResponse(
        boolean updateAvailable,
        String currentVersion,
        String latestVersion,
        String releaseTitle,
        String releaseNotes,
        String downloadUrl,
        Long downloadSizeInBytes,
        LocalDateTime publishedAt
) {
}
