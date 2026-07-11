package com.devaulty.backend.adapter.in.web.link.dto;

import java.time.LocalDateTime;
import java.util.UUID;

public record LinkViewResponse(
        UUID id,
        UUID projectId,
        String title,
        String url,
        String description,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
