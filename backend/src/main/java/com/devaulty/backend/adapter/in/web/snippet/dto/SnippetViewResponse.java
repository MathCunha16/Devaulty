package com.devaulty.backend.adapter.in.web.snippet.dto;

import com.devaulty.backend.domain.model.enums.SnippetLanguage;
import com.devaulty.backend.domain.model.enums.SnippetType;

import java.time.LocalDateTime;
import java.util.UUID;

public record SnippetViewResponse(
        UUID id,
        UUID projectId,
        String title,
        String description,
        String content,
        SnippetLanguage language,
        SnippetType snippetType,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
