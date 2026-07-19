package com.devaulty.backend.adapter.in.web.note.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import com.devaulty.backend.adapter.in.web.tag.dto.TagSummaryResponse;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

public record NoteSummaryResponse(
        UUID id,
        UUID projectId,
        String title,
        boolean archived,
        @Schema(description = "List of tags associated with this note")
        List<TagSummaryResponse> tags,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
