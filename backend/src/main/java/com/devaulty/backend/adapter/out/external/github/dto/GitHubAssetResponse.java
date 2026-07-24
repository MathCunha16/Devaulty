package com.devaulty.backend.adapter.out.external.github.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

public record GitHubAssetResponse(
    String name,
    @JsonProperty("browser_download_url")
    String browserDownloadUrl,
    long size,
    @JsonProperty("content_type")
    String contentType,
    String digest
) {
}
