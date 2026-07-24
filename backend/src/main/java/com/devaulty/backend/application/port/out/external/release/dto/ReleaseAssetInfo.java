package com.devaulty.backend.application.port.out.external.release.dto;

public record ReleaseAssetInfo(
        String fileName,
        String downloadUrl,
        long sizeInBytes,
        String contentType
) {
}
