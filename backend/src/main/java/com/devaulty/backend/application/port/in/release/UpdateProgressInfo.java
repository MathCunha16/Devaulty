package com.devaulty.backend.application.port.in.release;

import com.devaulty.backend.application.port.in.release.enums.UpdateStatus;

public record UpdateProgressInfo(
        UpdateStatus status,
        int percentage,           // 0 to 100
        long downloadedBytes,     // e.g: 45.200.000
        long totalBytes,          // e.g: 154.000.000
        String errorMessage
) {
}
