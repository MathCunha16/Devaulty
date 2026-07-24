package com.devaulty.backend.application.port.in.release;

import java.time.Instant;

public record AppUpdateInfo(
    boolean updateAvailable,
    String currentVersion,
    String latestVersion,
    String releaseTitle,
    String releaseNotes,
    String downloadUrl,
    Long downloadSizeInBytes,
    Instant publishedAt
) {
}
