package com.devaulty.backend.adapter.in.web.tag.dto;

import java.time.LocalDateTime;
import java.util.UUID;

public record TagViewResponse(
        UUID id,
        UUID projectId,
        String name,
        String color,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
