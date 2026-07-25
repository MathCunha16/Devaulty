package com.devaulty.backend.application.port.out.external.release.dto;

import java.time.Instant;
import java.util.List;

public record LatestReleaseInfo(
        String tagName,
        String name,
        String body,
        String htmlUrl,
        Instant publishedAt,
        boolean isPreRelease,
        List<ReleaseAssetInfo> assets
) {
}
