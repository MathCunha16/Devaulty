package com.devaulty.backend.adapter.out.external.github.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.List;

public record GitHubReleaseResponse(
        @JsonProperty("tag_name")
        String tagName,
        String name,
        String body,
        @JsonProperty("html_url")
        String htmlUrl,
        @JsonProperty("published_at")
        Instant publishedAt,
        @JsonProperty("prerelease")
        boolean preRelease,
        List<GitHubAssetResponse> assets
) {
}
