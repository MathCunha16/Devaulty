package com.devaulty.backend.adapter.in.web.snippet.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import com.devaulty.backend.adapter.in.web.tag.dto.TagSummaryResponse;
import com.devaulty.backend.domain.model.enums.SnippetLanguage;
import com.devaulty.backend.domain.model.enums.SnippetType;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

public record SnippetSummaryResponse(
        UUID id,
        UUID projectId,
        String title,
        String description,
        SnippetLanguage language,
        SnippetType snippetType,
        @Schema(description = "List of tags associated with this snippet")
        List<TagSummaryResponse> tags,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
