package com.devaulty.backend.adapter.in.web.note.dto;

import java.time.LocalDateTime;
import java.util.UUID;

public record NoteViewResponse(
        UUID id,
        UUID projectId,
        String title,
        String content,
        boolean archived,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
