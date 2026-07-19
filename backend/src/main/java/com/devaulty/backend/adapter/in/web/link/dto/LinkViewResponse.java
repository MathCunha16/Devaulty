package com.devaulty.backend.adapter.in.web.link.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import com.devaulty.backend.adapter.in.web.tag.dto.TagSummaryResponse;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

public record LinkViewResponse(
        UUID id,
        UUID projectId,
        String title,
        String url,
        String description,
        @Schema(description = "List of tags associated with this link")
        List<TagSummaryResponse> tags,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
