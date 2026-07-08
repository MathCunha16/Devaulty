package com.devaulty.backend.adapter.in.web.project.dto;

import java.time.LocalDateTime;
import java.util.UUID;

public record ProjectViewResponse(
        UUID id,
        String name,
        String description,
        String icon,
        String color,
        boolean archived,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
