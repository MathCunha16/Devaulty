package com.devaulty.backend.adapter.in.web.note.dto;

import java.time.LocalDateTime;
import java.util.UUID;

public record NoteSummaryResponse(
        UUID id,
        UUID projectId,
        String title,
        boolean archived,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
